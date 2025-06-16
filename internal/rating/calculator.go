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
		"Premier League":   0.9,
		"LaLiga":           0.8,
		"Bundesliga":       0.8,
		"SerieA":           0.8,
		"Ligue1":           0.7,
	}

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

func CalculatePositionOfTeams(ctx context.Context, teamsService *service.TeamsService, standingsService *service.StandingsService, match types.Match, rating int) ([][]int, error) {
	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := getLeaguesForTeams(ctx, teamsService, HomeID, AwayID)
	if err != nil {
		return nil, fmt.Errorf("error getting leagues for teams: %w", err)
	}

	// Теперь HomeLeague и AwayLeague должны содержать корректные имена лиг
	posHome, err := standingsService.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return nil, fmt.Errorf("error getting home team standing (league: %s, id: %d): %w", HomeLeague, HomeID, err)
	}

	posAway, err := standingsService.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return nil, fmt.Errorf("error getting away team standing (league: %s, id: %d): %w", AwayLeague, AwayID, err)
	}

	answ := make([][]int, 2)
	answ[0] = make([]int, 2)
	answ[1] = make([]int, 2)

	answ[0][0] = posHome
	answ[1][0] = posAway

	answ[0][1] = HomeID
	answ[1][1] = AwayID
	return answ, nil

}

// func CalculateMatchRating(match types.Match) float64 {
// 	var RatingOfHomeTeam float64
// 	var RatingOfAwayTeam float64
