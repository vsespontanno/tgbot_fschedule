package rating

import (
	"context"
	"football_tgbot/db"
	"football_tgbot/types"
	"log"
	"time"
)

// Service сервис для управления рейтингами команд
type Service struct {
	store      db.MatchesStore
	calculator *Calculator
}

// NewService создает новый сервис для управления рейтингами
func NewService(store db.MatchesStore) *Service {
	return &Service{
		store:      store,
		calculator: NewCalculator(),
	}
}

// UpdateRatings обновляет рейтинги всех команд для указанного турнира
func (s *Service) UpdateRatings(ctx context.Context, collectionName string) error {
	// Получаем текущую таблицу
	standings, err := s.store.GetStandings(ctx, collectionName)
	if err != nil {
		return err
	}

	// Получаем последние матчи
	matches, err := s.store.GetMatches(ctx, collectionName)
	if err != nil {
		return err
	}

	// Вычисляем рейтинги для каждой команды
	var ratings []types.TeamRating
	for _, standing := range standings {
		rating := s.calculator.CalculateTeamRating(standing, matches, standing.Team.ID)
		ratings = append(ratings, rating)
	}

	// Сохраняем обновленные рейтинги
	err = s.store.SaveTeamRatings(ctx, collectionName, ratings)
	if err != nil {
		return err
	}

	log.Printf("Updated ratings for %d teams in %s", len(ratings), collectionName)
	return nil
}

// GetMatchRating возвращает рейтинг матча
func (s *Service) GetMatchRating(ctx context.Context, collectionName string, match types.Match) (float64, error) {
	// Получаем рейтинги команд
	homeRating, err := s.store.GetTeamRating(ctx, collectionName, match.HomeTeam.ID)
	if err != nil {
		return 0, err
	}

	awayRating, err := s.store.GetTeamRating(ctx, collectionName, match.AwayTeam.ID)
	if err != nil {
		return 0, err
	}

	// Если рейтинги не найдены, возвращаем нейтральный рейтинг
	if homeRating == nil || awayRating == nil {
		return 0.5, nil
	}

	// Вычисляем рейтинг матча
	return s.calculator.CalculateMatchRating(match, *homeRating, *awayRating), nil
}

// StartRatingUpdater запускает периодическое обновление рейтингов
func (s *Service) StartRatingUpdater(ctx context.Context, collectionName string, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.UpdateRatings(ctx, collectionName); err != nil {
				log.Printf("Error updating ratings: %v", err)
			}
		}
	}
}
