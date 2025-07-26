// internal/adapters/calculator_adapter.go
package service

import (
	"context"

	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// CalculatorAdapter реализует интерфейс Calculator
type CalculatorAdapter struct {
	teamsStore     mongoRepo.TeamsCalcStore
	standingsStore mongoRepo.StandingsCalcStore
	matchesStore   mongoRepo.MatchCalcStore
}

// Конструктор для создания нового экземпляра CalculatorAdapter
func NewCalculatorAdapter(teamsStore mongoRepo.TeamsCalcStore, standingsStore mongoRepo.StandingsCalcStore, matchesStore mongoRepo.MatchCalcStore) Calculator {
	return &CalculatorAdapter{teamsStore, standingsStore, matchesStore}
}

// Находит место команды в турнирной таблице по её уникальному идентификатору
func (a *CalculatorAdapter) HandleGetTeamStanding(ctx context.Context, league string, id int) (int, error) {
	standing, err := a.standingsStore.GetTeamStanding(ctx, league, id)
	if err != nil {
		return 0, err
	}
	return standing, nil
}

// Получает последние N матчей команды по её уникальному идентификатору
// Используется для анализа формы команды и её текущего состояния
func (a *CalculatorAdapter) HandleGetRecentMatches(ctx context.Context, teamID, lastN int) ([]types.Match, error) {
	return a.matchesStore.GetRecentMatches(ctx, teamID, lastN)
}

// Получает лигу команды по её уникальному идентификатору
// Используется для определения, в какой лиге играет команда
func (a *CalculatorAdapter) HandleGetLeague(ctx context.Context, collectionName string, id int) (string, error) {
	league, err := a.teamsStore.GetTeamLeague(ctx, collectionName, id)
	if err != nil {
		return "", err
	}
	return league, nil
}

// Получает короткое название команды по её полному названию
// Используется для отображения более компактного названия команды в интерфейсе
func (a *CalculatorAdapter) HandleGetTeamShortName(ctx context.Context, collectionName string, fullName string) (string, error) {
	shortName, err := a.teamsStore.GetTeamsShortName(ctx, collectionName, fullName)
	if err != nil {
		return "", err
	}
	return shortName, nil
}
