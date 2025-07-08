package service

import (
	"context"

	mongorepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
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

func (s *TeamsService) HandleGetTeamShortName(ctx context.Context, collectionName string, fullName string) (string, error) {
	shortName, err := s.teamsStore.GetTeamsShortName(ctx, collectionName, fullName)
	if err != nil {
		return "", err
	}
	return shortName, nil
}
