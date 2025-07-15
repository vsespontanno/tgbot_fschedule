// internal/adapters/calculator_adapter.go
package adapters

import (
	"context"

	"github.com/vsespontanno/tgbot_fschedule/internal/domain"
	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

type CalculatorAdapter struct {
	teamsStore     mongoRepo.TeamsCalcStore
	standingsStore mongoRepo.StandingsCalcStore
	matchesStore   mongoRepo.MatchCalcStore
}

func NewCalculatorAdapter(teamsStore mongoRepo.TeamsCalcStore, standingsStore mongoRepo.StandingsCalcStore, matchesStore mongoRepo.MatchCalcStore) domain.Calculator {
	return &CalculatorAdapter{teamsStore, standingsStore, matchesStore}
}

func (a *CalculatorAdapter) HandleGetTeamStanding(ctx context.Context, league string, id int) (int, error) {
	standing, err := a.standingsStore.GetTeamStanding(ctx, league, id)
	if err != nil {
		return 0, err
	}
	return standing, nil
}

func (a *CalculatorAdapter) HandleGetRecentMatches(ctx context.Context, teamID, lastN int) ([]types.Match, error) {
	return a.matchesStore.GetRecentMatches(ctx, teamID, lastN)
}

func (a *CalculatorAdapter) HandleGetMatches(ctx context.Context) ([]types.Match, error) {
	return a.matchesStore.GetMatches(ctx, "matches")

}

func (a *CalculatorAdapter) HandleGetLeague(ctx context.Context, collectionName string, id int) (string, error) {
	league, err := a.teamsStore.GetTeamLeague(ctx, collectionName, id)
	if err != nil {
		return "", err
	}
	return league, nil
}

func (a *CalculatorAdapter) HandleGetTeamShortName(ctx context.Context, collectionName string, fullName string) (string, error) {
	shortName, err := a.teamsStore.GetTeamsShortName(ctx, collectionName, fullName)
	if err != nil {
		return "", err
	}
	return shortName, nil
}
