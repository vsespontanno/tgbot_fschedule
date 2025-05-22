package rating

import (
	"context"
	"fmt"
	"football_tgbot/internal/db"
	"football_tgbot/internal/types"
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
func (s *Service) GetMatchRating(ctx context.Context, _ string, match types.Match) (float64, error) {
	// Определяем название коллекции на основе competition ID матча
	comp_name := match.Competition.Name
	var collectionName string
	switch comp_name {
	case "EPL":
		collectionName = "PremierLeague_standings"
	case "Bundesliga":
		collectionName = "Bundesliga_standings"
	case "Serie A":
		collectionName = "SerieA_standings"
	case "Ligue 1":
		collectionName = "Ligue1_standings"
	case "UCL":
		collectionName = "ChampionsLeague_standings"
	case "La Liga":
		collectionName = "LaLiga_standings"
	default:
		collectionName = fmt.Sprintf("%s_standings", match.Competition.Name)
	}

	// Получаем данные о положении команд
	homeStanding, err := s.store.GetTeamStanding(ctx, collectionName, match.HomeTeam.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to get home team standing: %w", err)
	}
	awayStanding, err := s.store.GetTeamStanding(ctx, collectionName, match.AwayTeam.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to get away team standing: %w", err)
	}

	// Если данные не найдены, возвращаем нейтральный рейтинг
	if homeStanding == nil || awayStanding == nil {
		return 0.5, nil
	}

	// Используем существующий калькулятор для вычисления рейтингов
	matches := []types.Match{match} // Передаем текущий матч для расчета формы
	homeRating := s.calculator.CalculateTeamRating(*homeStanding, matches, match.HomeTeam.ID)
	awayRating := s.calculator.CalculateTeamRating(*awayStanding, matches, match.AwayTeam.ID)

	// Вычисляем рейтинг матча
	return s.calculator.CalculateMatchRating(match, homeRating, awayRating), nil
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
