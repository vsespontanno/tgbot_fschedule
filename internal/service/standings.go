package service

import (
	"context"

	db "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
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

func (s *StandingsService) HandleGetTeamStanding(ctx context.Context, league string, id int) (int, error) {
	standing, err := s.standingsStore.GetTeamStanding(ctx, league, id)
	if err != nil {
		return 0, err
	}
	return standing, nil
}
