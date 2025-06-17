package rating

import (
	"context"
	"fmt"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
)

var (
	CLstage = map[string]float64{
		"Round of 16":    0.25,
		"Quarter-finals": 0.50,
		"Semi-finals":    0.75,
		"Final":          1.0,
	}

	leagueNorm = map[string]float64{
		"Champions League": 1.0,
		"PremierLeague":    0.9,
		"LaLiga":           0.8,
		"SerieA":           0.8,
		"Bundesliga":       0.7,
		"Ligue1":           0.65,
	}

	teamsInLeague = map[string]int{
		"PremierLeague": 20,
		"LaLiga":        20,
		"Bundesliga":    18,
		"SerieA":        18,
		"Ligue1":        20,
	}

	// derbys = map[string]float64{

	// }

	derbyflag float64
)

// getLeaguesForTeams пытается определить лиги для домашней и гостевой команд.
func getLeaguesForTeams(ctx context.Context, teamsService *service.TeamsService, homeTeamID int, awayTeamID int) (homeLeague string, awayLeague string, err error) {
	// Найти лигу для домашней команды
	foundHomeLeague := false
	for leagueKey := range leagueNorm {
		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, homeTeamID)
		if getLeagueErr != nil {
			// Можно логировать ошибку или обработать ее иначе, пока продолжаем поиск
			// fmt.Printf("Error checking league %s for home team %d: %v\n", leagueKey, homeTeamID, getLeagueErr)
			// Сохраняем первую возникшую ошибку, если она есть
			if err == nil {
				err = fmt.Errorf("error checking league %s for home team %d: %w", leagueKey, homeTeamID, getLeagueErr)
			}
			continue
		}
		if league != "Wrong League" {
			homeLeague = league
			foundHomeLeague = true
			err = nil // Сбрасываем ошибку, если лига найдена
			break
		}
	}
	if !foundHomeLeague {
		if err == nil { // Если ошибок при поиске не было, но лига не найдена
			return "", "", fmt.Errorf("could not determine league for home team ID %d", homeTeamID)
		}
		return "", "", fmt.Errorf("could not determine league for home team ID %d (last error: %w)", homeTeamID, err)
	}

	// Найти лигу для гостевой команды
	var awayErr error // Отдельная переменная для ошибок гостевой команды
	foundAwayLeague := false
	for leagueKey := range leagueNorm {
		league, getLeagueErr := teamsService.HandleGetLeague(ctx, leagueKey, awayTeamID)
		if getLeagueErr != nil {
			// fmt.Printf("Error checking league %s for away team %d: %v\n", leagueKey, awayTeamID, getLeagueErr)
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
		if awayErr == nil {
			return homeLeague, "", fmt.Errorf("could not determine league for away team ID %d", awayTeamID)
		}
		return homeLeague, "", fmt.Errorf("could not determine league for away team ID %d (last error: %w)", awayTeamID, awayErr)
	}
	return homeLeague, awayLeague, nil // Возвращаем найденные лиги и nil в качестве ошибки, если обе найдены
}

func CalculatePositionOfTeams(ctx context.Context, teamsService *service.TeamsService, standingsService *service.StandingsService, match types.Match) (homeTeam, awayTeam float64, err error) {
	var homeTeamRating float64
	var awayTeamRating float64

	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := getLeaguesForTeams(ctx, teamsService, HomeID, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	// Теперь HomeLeague и AwayLeague должны содержать корректные имена лиг
	posHome, err := standingsService.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting home team standing (league: %s, id: %d): %w", HomeLeague, HomeID, err)
	}

	posAway, err := standingsService.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting away team standing (league: %s, id: %d): %w", AwayLeague, AwayID, err)
	}

	homeTeamRating = (float64(teamsInLeague[HomeLeague] - posHome)) / float64(teamsInLeague[HomeLeague]-1)
	awayTeamRating = (float64(teamsInLeague[AwayLeague] - posAway)) / float64(teamsInLeague[AwayLeague]-1)
	return homeTeamRating, awayTeamRating, nil
}

func CalculateRatingOfMatch(ctx context.Context, match types.Match, teamService *service.TeamsService, standingsService *service.StandingsService) (float64, error) {
	// 1) Считаем силу по позиции в таблице (от 0.0 до 1.0)
	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, teamService, standingsService, match)
	if err != nil {
		return 0, fmt.Errorf("error calculating team strengths: %w", err)
	}

	// 2) Определяем лиги
	homeLeague, awayLeague, err := getLeaguesForTeams(ctx, teamService, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil {
		return 0, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	// 3) Средний вес лиги (раньше вы брали сумму, теперь — среднее)
	avgLeagueWeight := (leagueNorm[homeLeague] + leagueNorm[awayLeague]) / 2.0

	// 4) Бонусы (пока нулевые, добавите позже дерби и стадию ЛЧ)
	var derbyBonus float64 // = 0.15 если дерби
	var stageBonus float64 // = CLstage[stage] * stageWeight

	// 5) Собираем окончательный рейтинг
	//    teamStrength — доля 0.2 (каждая команда по 0.1)
	//    leagueWeight — доля 0.4
	//    derbyBonus  — доля 0.15
	//    stageBonus  — доля 0.25
	rating := ((homeStrength+awayStrength)/2.0)*0.2 +
		avgLeagueWeight*0.4 +
		derbyBonus +
		stageBonus

	// Гарантируем, что не выйдем за 1.0
	if rating > 1.0 {
		rating = 1.0
	}
	return rating, nil
}
