package service

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
	var wins int
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
	homeLeague, err = calculator.HandleGetLeague(ctx, "Teams", homeTeamID)
	if err != nil {
		return "", "", fmt.Errorf("error getting home league: %w", err)
	}
	awayLeague, err = calculator.HandleGetLeague(ctx, "Teams", awayTeamID)
	if err != nil {
		return "", "", fmt.Errorf("error getting away league: %w", err)
	}
	return homeLeague, awayLeague, nil

}

func CalculatePositionOfTeams(ctx context.Context, calculator Calculator, match types.Match) (homeTeam, awayTeam float64, err error) {
	HomeID := match.HomeTeam.ID
	AwayID := match.AwayTeam.ID

	HomeLeague, AwayLeague, err := GetLeaguesForTeams(ctx, calculator, HomeID, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting leagues for teams: %s", err)
	}
	if HomeLeague == "" || AwayLeague == "" {
		return -1, -1, err
	}
	posHome, err := calculator.HandleGetTeamStanding(ctx, HomeLeague, HomeID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting home team standing: %s", err)
	}
	posAway, err := calculator.HandleGetTeamStanding(ctx, AwayLeague, AwayID)
	if err != nil {
		return -1, -1, fmt.Errorf("error getting away team standing: %s", err)
	}

	homeTeamRating := (float64(types.TeamsInLeague[HomeLeague] - posHome)) / float64(types.TeamsInLeague[HomeLeague]-1)
	awayTeamRating := (float64(types.TeamsInLeague[AwayLeague] - posAway)) / float64(types.TeamsInLeague[AwayLeague]-1)
	return homeTeamRating, awayTeamRating, nil
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
