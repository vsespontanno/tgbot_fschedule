package service

import (
	"context"
	db "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"
)

type StandingsService struct {
	standingsStore db.StandingsStore
}

func NewStandingService(standingsStore db.StandingsStore) *StandingsService {
	return &StandingsService{
		standingsStore: standingsStore,
	}
}

func (s *StandingsService) HandleGetStandings(ctx context.Context, league string) ([]types.Standing, error) {
	standings, err := s.standingsStore.GetStandings(ctx, league)
	if err != nil {
		return nil, err
	}
	return standings, nil
}

func (s *StandingsService) HandleSaveStandings(ctx context.Context, league string, standings []types.Standing) error {
	err := s.standingsStore.SaveStandings(ctx, league, standings)
	if err != nil {
		return err
	}
	return nil
}

func (s *StandingsService) HandleGetTeamStanding(ctx context.Context, league string, teamID int) (*types.Standing, error) {
	standing, err := s.standingsStore.GetTeamStanding(ctx, league, teamID)
	if err != nil {
		return nil, err
	}
	return standing, nil
}
