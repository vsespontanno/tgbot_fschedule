package service

import (
	"context"

	mongorepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// TeamsService предоставляет методы для работы с командами
// Использует MongoDB для хранения и получения данных о командах
// Реализует интерфейс TeamsStore для взаимодействия с данными команд
type TeamsService struct {
	teamsStore mongorepo.TeamsStore
}

// Конструктор для создания нового экземпляра TeamsService
func NewTeamsService(teamsStore mongorepo.TeamsStore) *TeamsService {
	return &TeamsService{
		teamsStore: teamsStore,
	}
}

// Метод для сохранения команд в базу MongoDB
func (s *TeamsService) HandleSaveTeams(ctx context.Context, collectionName string, teams []types.Team) error {
	err := s.teamsStore.SaveTeamsToMongoDB(ctx, collectionName, teams)
	if err != nil {
		return err
	}
	return nil
}

// Метод для обновления или вставки команды в базу MongoDB
// Если команда с таким ID уже существует, выполняет обновление
// Если не существует, выполняет вставку новой команды
func (s *TeamsService) HandleUpsertMatch(ctx context.Context, collectionName string, team types.Team) error {
	err := s.teamsStore.UpsertTeamToMongoDB(ctx, collectionName, team)
	if err != nil {
		return err
	}
	return nil
}
