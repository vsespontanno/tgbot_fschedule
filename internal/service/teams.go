package service

import (
	"context"
	mongorepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"
)

type TeamsService struct {
	teamsStore mongorepo.TeamsStore
}

func NewTeamsService(teamsStore mongorepo.TeamsStore) *TeamsService {
	return &TeamsService{
		teamsStore: teamsStore,
	}
}

func (s *TeamsService) HandleGetTeams(ctx context.Context, collectionName string) ([]types.Team, error) {
	teams, err := s.teamsStore.GetAllTeams(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *TeamsService) HandleGetLeague(ctx context.Context, collectionName string, id int) (string, error) {
	league, err := s.teamsStore.GetTeamLeague(ctx, collectionName, id)
	if err != nil {
		return "", err
	}
	return league, nil
}
