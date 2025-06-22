package rating

import (
	"context"
	"fmt"
	mongorepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
	"log"
)

var (
	CLstage = map[string]float64{
		"PLAYOFFS":       0.25, // 1/16
		"LAST_16":        0.5,  // 1/8
		"QUARTER_FINALS": 0.75, // 1/4
		"SEMI_FINALS":    0.9,  // 1/2
		"FINAL":          1.0,  // Final
	}

	leagueNorm = map[string]float64{
		"Champions League": 1.0,
		"PremierLeague":    0.9, // Снижено с 0.9, чтобы ЛЧ доминировала
		"LaLiga":           0.8, // Увеличено с 0.8
		"SerieA":           0.8,
		"Bundesliga":       0.75,
		"Ligue1":           0.7,
	}

	teamsInLeague = map[string]int{
		"PremierLeague": 20,
		"LaLiga":        20,
		"Bundesliga":    18,
		"SerieA":        18,
		"Ligue1":        20,
	}

	derbys = map[[2]string]float64{
		// Англия (PremierLeague)
		{"Manchester United", "Manchester City"}: 0.27, // Манчестерское дерби
		{"Liverpool", "Everton"}:                 0.16, // Мерсисайдское дерби
		{"Arsenal", "Tottenham"}:                 0.25, // Северолондонское дерби
		{"Chelsea", "Arsenal"}:                   0.25,
		{"Chelsea", "Tottenham"}:                 0.25,
		{"Manchester United", "Liverpool"}:       0.26,
		{"Manchester United", "Leeds United"}:    0.15,
		{"Newcastle", "Sunderland"}:              0.14, // Тайн-Уир

		// Испания (LaLiga)
		{"Real Madrid", "Barcelona"}: 0.35, // Увеличено с 0.3 для Эль Класико
		{"Atletico", "Real Madrid"}:  0.26, // Мадридское дерби
		{"Sevilla", "Real Betis"}:    0.2,  // Севильское дерби
		{"Barcelona", "Espanyol"}:    0.18, // Барселонское дерби
		{"Valencia", "Levante"}:      0.14, // Валенсийское дерби

		// Германия (Bundesliga)
		{"Borussia D.", "Bayern"}:           0.28, // Дер Классикер
		{"Schalke 04", "Borussia Dortmund"}: 0.16, // Рурское дерби
		{"Hamburger SV", "Werder Bremen"}:   0.15, // Северное дерби
		{"Bayern", "1860 Munich"}:           0.14, // Мюнхенское дерби
		{"Cologne", "Borussia M."}:          0.14,

		// Италия (SerieA)
		{"Inter", "Milan"}:     0.29, // Миланское дерби
		{"Roma", "Lazio"}:      0.28, // Римское дерби
		{"Juventus", "Torino"}: 0.2,  // Дерби делла Моле
		{"Genoa", "Sampdoria"}: 0.18, // Дерби делла Лантерна
		{"Napoli", "Roma"}:     0.15,

		// Франция (Ligue1)
		{"PSG", "Marseille"}:                0.23, // Ле Классик
		{"Olympique Lyon", "Saint-Etienne"}: 0.18, // Ронское дерби
		{"Nice", "Monaco"}:                  0.14, // Лазурное дерби
		{"Lille", "RC Lens"}:                0.14, // Северное дерби
	}
)

// calculateForm вычисляет форму команды на основе последних матчей
func calculateForm(matches []types.Match) float64 {
	if len(matches) == 0 {
		return 0.5 // Нейтральная форма, если данных нет
	}
	wins := 0
	for _, m := range matches {
		if m.Status == "FINISHED" {
			if (m.Score.Winner == "HOME_TEAM" && m.HomeTeam.ID == m.ID) ||
				(m.Score.Winner == "AWAY_TEAM" && m.AwayTeam.ID == m.ID) {
				wins++
			}
		}
	}
	return float64(wins) / float64(len(matches))
}

// getLeaguesForTeams определяет лиги для команд
func getLeaguesForTeams(ctx context.Context, teamsService *service.TeamsService, homeTeamID int, awayTeamID int) (homeLeague string, awayLeague string, err error) {
	// Логика осталась прежней
	foundHomeLeague := false
	for leagueKey := range leagueNorm {
		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, homeTeamID)
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
	for leagueKey := range leagueNorm {
		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, awayTeamID)
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

func CalculatePositionOfTeams(ctx context.Context, teamsService *service.TeamsService, standingsService *service.StandingsService, match types.Match) (homeTeam, awayTeam float64, err error) {
	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := getLeaguesForTeams(ctx, teamsService, HomeID, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	posHome, err := standingsService.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting home team standing: %w", err)
	}
	posAway, err := standingsService.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting away team standing: %w", err)
	}

	homeTeamRating := (float64(teamsInLeague[HomeLeague] - posHome)) / float64(teamsInLeague[HomeLeague]-1)
	awayTeamRating := (float64(teamsInLeague[AwayLeague] - posAway)) / float64(teamsInLeague[AwayLeague]-1)
	return homeTeamRating, awayTeamRating, nil
}

func CalculateRatingOfMatch(ctx context.Context, match types.Match, teamService *service.TeamsService, standingsService *service.StandingsService, matchesStore mongorepo.MatchesStore) (float64, error) {
	// 1) Сила команд по позициям
	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, teamService, standingsService, match)
	if err != nil {
		return 0, fmt.Errorf("error calculating team strengths: %w", err)
	}

	// 2) Лиги и вес
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil || homeLeague == "" || awayLeague == "" {
		fmt.Printf("Матч %s - %s пропущен: проблема с лигами\n", match.HomeTeam.Name, match.AwayTeam.Name)
		return 0, nil
	}
	avgLeagueWeight := (leagueNorm[homeLeague] + leagueNorm[awayLeague]) / 2.0

	// 3) Форма команд
	recentMatchesHome, err := matchesStore.GetRecentMatches(ctx, match.HomeTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for home team %d: %v", match.HomeTeam.ID, err)
	}
	recentMatchesAway, err := matchesStore.GetRecentMatches(ctx, match.AwayTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for away team %d: %v", match.AwayTeam.ID, err)
	}
	homeForm := calculateForm(recentMatchesHome)
	awayForm := calculateForm(recentMatchesAway)
	formFactor := (homeForm + awayForm) / 2.0

	// 4) Бонусы
	derbyBonus := GetDerbyBonus(ctx, teamService, match)
	stageBonus := 0.0
	if homeLeague == "Champions League" && match.Stage != "" {
		stageBonus = CLstage[match.Stage]
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

func GetDerbyBonus(ctx context.Context, teamService *service.TeamsService, match types.Match) float64 {
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		log.Printf("Error getting leagues for derby bonus: %v", err)
		return 0.0
	}
	homeShortName, err := teamService.HandleGetTeamShortName(ctx, homeLeague, match.HomeTeam.Name)
	if err != nil {
		log.Printf("Error getting short name for home team: %v", err)
		return 0.0
	}
	awayShortName, err := teamService.HandleGetTeamShortName(ctx, awayLeague, match.AwayTeam.Name)
	if err != nil {
		log.Printf("Error getting short name for away team: %v", err)
		return 0.0
	}
	key := [2]string{homeShortName, awayShortName}
	reverseKey := [2]string{awayShortName, homeShortName}

	if bonus, ok := derbys[key]; ok {
		return bonus
	}
	if bonus, ok := derbys[reverseKey]; ok {
		return bonus
	}
	return 0.0
}
