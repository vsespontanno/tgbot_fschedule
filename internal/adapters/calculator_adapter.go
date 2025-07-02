// internal/adapters/calculator_adapter.go
package adapters

import (
	"context"
	"football_tgbot/internal/domain"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
)

type CalculatorAdapter struct {
	teamsSvc     *service.TeamsService
	standingsSvc *service.StandingsService
	matchesSvc   *service.MatchesService
}

func NewCalculatorAdapter(
	teamsSvc *service.TeamsService,
	standingsSvc *service.StandingsService,
	matchesSvc *service.MatchesService,
) domain.Calculator {
	return &CalculatorAdapter{teamsSvc, standingsSvc, matchesSvc}
}

func (a *CalculatorAdapter) HandleGetLeague(
	ctx context.Context, leagueKey string, teamID int,
) (string, error) {
	return a.teamsSvc.HandleGetLeague(ctx, leagueKey, teamID)
}

func (a *CalculatorAdapter) HandleGetTeamStanding(
	ctx context.Context, leagueKey string, teamID int,
) (int, error) {
	return a.standingsSvc.HandleGetTeamStanding(ctx, leagueKey, teamID)
}

func (a *CalculatorAdapter) HandleGetRecentMatches(
	ctx context.Context, teamID, lastN int,
) ([]types.Match, error) {
	return a.matchesSvc.HandleGetRecentMatches(ctx, teamID, lastN)
}

func (a *CalculatorAdapter) HandleGetTeamShortName(
	ctx context.Context, leagueKey, teamName string,
) (string, error) {
	return a.teamsSvc.HandleGetTeamShortName(ctx, leagueKey, teamName)
}

func (a *CalculatorAdapter) HandleGetMatches(ctx context.Context) ([]types.Match, error) {
	return a.matchesSvc.HandleGetMatches(ctx)

}
