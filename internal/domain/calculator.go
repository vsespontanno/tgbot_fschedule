package domain

import (
	"context"
	"fmt"
	"log"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// CalculateForm вычисляет форму команды на основе последних матчей
func CalculateForm(matches []types.Match, teamID int) float64 {
	if len(matches) == 0 {
		return 0.5 // Нейтральная форма, если данных нет
	}
	wins := 0
	for _, m := range matches {
		if m.Status == "FINISHED" {
			if (m.Score.Winner == "HOME_TEAM" && m.HomeTeam.ID == teamID) ||
				(m.Score.Winner == "AWAY_TEAM" && m.AwayTeam.ID == teamID) {
				wins++
			}
		}
	}
	return float64(wins) / float64(len(matches))
}

// GetLeaguesForTeams определяет лиги для команд
func GetLeaguesForTeams(ctx context.Context, calculator Calculator, homeTeamID int, awayTeamID int) (homeLeague string, awayLeague string, err error) {
	// Логика осталась прежней
	foundHomeLeague := false
	for leagueKey := range types.LeagueNorm {
		league, getLeagueErr := calculator.HandleGetLeague(ctx, leagueKey, homeTeamID)
		if getLeagueErr != nil {
			if err == nil {
				err = fmt.Errorf("error checking league %s for home team %d: %w", leagueKey, homeTeamID, getLeagueErr)
			}
			continue
		}
		if league != "Wrong League" {
			homeLeague = league
			foundHomeLeague = true
			err = nil
			break
		}
	}
	if !foundHomeLeague {
		fmt.Printf("Команда с id=%d не найдена ни в одной лиге\n", homeTeamID)
		return "", "", nil
	}

	var awayErr error
	foundAwayLeague := false
	for leagueKey := range types.LeagueNorm {
		league, getLeagueErr := calculator.HandleGetLeague(ctx, leagueKey, awayTeamID)
		if getLeagueErr != nil {
			if awayErr == nil {
				awayErr = fmt.Errorf("error checking league %s for away team %d: %w", leagueKey, awayTeamID, getLeagueErr)
			}
			continue
		}
		if league != "Wrong League" {
			awayLeague = league
			foundAwayLeague = true
			break
		}
	}
	if !foundAwayLeague {
		return homeLeague, "", awayErr
	}
	return homeLeague, awayLeague, nil
}

func CalculatePositionOfTeams(ctx context.Context, calculator Calculator, match types.Match) (homeTeam, awayTeam float64, err error) {
	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := GetLeaguesForTeams(ctx, calculator, HomeID, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	posHome, err := calculator.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting home team standing: %w", err)
	}
	posAway, err := calculator.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting away team standing: %w", err)
	}

	homeTeamRating := (float64(types.TeamsInLeague[HomeLeague] - posHome)) / float64(types.TeamsInLeague[HomeLeague]-1)
	awayTeamRating := (float64(types.TeamsInLeague[AwayLeague] - posAway)) / float64(types.TeamsInLeague[AwayLeague]-1)
	return homeTeamRating, awayTeamRating, nil
}

func CalculateRatingOfMatch(ctx context.Context, match types.Match, calculator Calculator) (float64, error) {
	// 1) Сила команд по позициям
	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, calculator, match)
	if err != nil {
		return 0, fmt.Errorf("error calculating team strengths: %w", err)
	}

	// 2) Лиги и вес
	homeLeague, awayLeague, err := GetLeaguesForTeams(ctx, calculator, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil || homeLeague == "" || awayLeague == "" {
		fmt.Printf("Матч %s - %s пропущен: проблема с лигами\n", match.HomeTeam.Name, match.AwayTeam.Name)
		return 0, nil
	}
	avgLeagueWeight := (types.LeagueNorm[homeLeague] + types.LeagueNorm[awayLeague]) / 2.0

	// 3) Форма команд
	recentMatchesHome, err := calculator.HandleGetRecentMatches(ctx, match.HomeTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for home team %d: %v", match.HomeTeam.ID, err)
	}
	recentMatchesAway, err := calculator.HandleGetRecentMatches(ctx, match.AwayTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for away team %d: %v", match.AwayTeam.ID, err)
	}
	homeForm := CalculateForm(recentMatchesHome, match.HomeTeam.ID)
	awayForm := CalculateForm(recentMatchesAway, match.AwayTeam.ID)
	formFactor := (homeForm + awayForm) / 2.0

	// 4) Бонусы
	derbyBonus := GetDerbyBonus(ctx, calculator, match)
	stageBonus := 0.0
	if homeLeague == "Champions League" && match.Stage != "" {
		stageBonus = types.CLstage[match.Stage]
	}
	crossLeagueBonus := 0.0
	if homeLeague != awayLeague {
		crossLeagueBonus = 0.15 // Увеличен с 0.1 для большего эффекта
	}

	// 5) Финальный рейтинг
	baseRating := (homeStrength+awayStrength)/2.0*0.15 + // Уменьшено влияние позиций
		avgLeagueWeight*0.35 + // Снижено с 0.4
		formFactor*0.15 // Добавлено влияние формы
	rating := baseRating * (1 + derbyBonus + stageBonus + crossLeagueBonus)

	// 6) Ограничение и минимальное значение
	if rating > 1.0 {
		rating = 1.0
	}
	if rating < 0.1 { // Минимальный рейтинг, чтобы избежать нулей
		rating = 0.1
	}

	return rating, nil
}

func GetDerbyBonus(ctx context.Context, calculator Calculator, match types.Match) float64 {
	homeLeague, awayLeague, err := GetLeaguesForTeams(ctx, calculator, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		log.Printf("Error getting leagues for derby bonus: %v", err)
		return 0.0
	}
	homeShortName, err := calculator.HandleGetTeamShortName(ctx, homeLeague, match.HomeTeam.Name)
	if err != nil {
		log.Printf("Error getting short name for home team: %v", err)
		return 0.0
	}
	awayShortName, err := calculator.HandleGetTeamShortName(ctx, awayLeague, match.AwayTeam.Name)
	if err != nil {
		log.Printf("Error getting short name for away team: %v", err)
		return 0.0
	}
	key := [2]string{homeShortName, awayShortName}
	reverseKey := [2]string{awayShortName, homeShortName}

	if bonus, ok := types.Derbys[key]; ok {
		return bonus
	}
	if bonus, ok := types.Derbys[reverseKey]; ok {
		return bonus
	}
	return 0.0
}
