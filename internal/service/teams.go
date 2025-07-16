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

func (s *TeamsService) HandleSaveTeams(ctx context.Context, collectionName string, teams []types.Team) error {
	err := s.teamsStore.SaveTeamsToMongoDB(ctx, collectionName, teams)
	if err != nil {
		return err
	}
	return nil
}

func (s *TeamsService) HandleUpsertMatch(ctx context.Context, collectionName string, team types.Team) error {
	err := s.teamsStore.UpsertTeamToMongoDB(ctx, collectionName, team)
	if err != nil {
		return err
	}
	return nil
}
