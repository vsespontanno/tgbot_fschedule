package service

import (
	"context"

	db "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// StandingsService предоставляет методы для работы с таблицами лиг
// Использует MongoDB для хранения и получения данных о таблицах лиг
// Реализует интерфейс StandingsStore для взаимодействия с данными таблиц лиг
type StandingsService struct {
	standingsStore db.StandingsStore
}

// Конструктор для создания нового экземпляра StandingsService
func NewStandingService(standingsStore db.StandingsStore) *StandingsService {
	return &StandingsService{
		standingsStore: standingsStore,
	}
}

// Метод для получения таблицы лиги по названию лиги
func (s *StandingsService) HandleGetStandings(ctx context.Context, league string) ([]types.Standing, error) {
	standings, err := s.standingsStore.GetStandings(ctx, league)
	if err != nil {
		return nil, err
	}
	return standings, nil
}

// Метод для сохранения таблиц лиг в базу MongoDB
func (s *StandingsService) HandleSaveStandings(ctx context.Context, league string, standings []types.Standing) error {
	err := s.standingsStore.SaveStandings(ctx, league, standings)
	if err != nil {
		return err
	}
	return nil
}
