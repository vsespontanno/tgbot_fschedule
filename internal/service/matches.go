package service

import (
	"context"
	"fmt"
	"log"

	"github.com/vsespontanno/tgbot_fschedule/internal/client"
	db "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

type MatchesService struct {
	matchesStore db.MatchesStore
	apiClient    client.MatchApiClient
}

func NewMatchesService(matchesStore db.MatchesStore, apiClient client.MatchApiClient) *MatchesService {
	return &MatchesService{
		matchesStore: matchesStore,
		apiClient:    apiClient,
	}
}

func (s *MatchesService) HandleGetMatches(ctx context.Context) ([]types.Match, error) {
	matches, err := s.matchesStore.GetMatches(ctx, "matches")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchesService) HandleSaveMatches(matches []types.Match, from, to string) error {
	err := s.matchesStore.SaveMatchesToMongoDB(matches, from, to)
	if err != nil {
		return err
	}
	return nil
}

func (s *MatchesService) HandleReqMatches(ctx context.Context, from string, to string) ([]types.Match, error) {
	return s.apiClient.FetchMatches(ctx, from, to)
}

func (s *MatchesService) HandleSaveMatchRating(ctx context.Context, match types.Match, rating float64) error {
	return s.matchesStore.UpdateMatchRatingInMongoDB(match, rating)
}

func (s *MatchesService) HandleUpsertMatch(ctx context.Context, match types.Match) error {
	return s.matchesStore.UpsertMatch(ctx, match)
}

func (s *MatchesService) CalculateRatingOfMatch(ctx context.Context, match types.Match, calculator Calculator) (float64, error) {
	// 1) Сила команд по позициям
	homeStrength, awayStrength, err := CalculatePositionOfTeams(ctx, calculator, match)
	if err != nil {
		return 0, fmt.Errorf("error calculating team strengths: %s", err)
	}

	// 2) Лиги и вес
	homeLeague, awayLeague, err := GetLeaguesForTeams(ctx, calculator, match.HomeTeam.ID, match.AwayTeam.ID)
	if err != nil || homeLeague == "" || awayLeague == "" {
		fmt.Printf("Матч %s - %s пропущен: проблема с лигами\nЛиги: %s - %s\nАйдишники: %d - %d\n", match.HomeTeam.Name, match.AwayTeam.Name, homeLeague, awayLeague, match.HomeTeam.ID, match.AwayTeam.ID)
		return 0, nil
	}
	avgLeagueWeight := (types.LeagueNorm[homeLeague] + types.LeagueNorm[awayLeague]) / 2.0

	// 3) Форма команд
	recentMatchesHome, err := calculator.HandleGetRecentMatches(ctx, match.HomeTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for home team %d: %v", match.HomeTeam.ID, err)
	}
	recentMatchesAway, err := calculator.HandleGetRecentMatches(ctx, match.AwayTeam.ID, 5)
	if err != nil {
		log.Printf("Error getting recent matches for away team %d: %v", match.AwayTeam.ID, err)
	}
	homeForm := CalculateForm(recentMatchesHome, match.HomeTeam.ID)
	awayForm := CalculateForm(recentMatchesAway, match.AwayTeam.ID)
	formFactor := (homeForm + awayForm) / 2.0

	// 4) Бонусы
	derbyBonus := GetDerbyBonus(ctx, calculator, match)
	stageBonus := 0.0
	if homeLeague == "Champions League" && match.Stage != "" {
		stageBonus = types.CLstage[match.Stage]
	}
	crossLeagueBonus := 0.0
	if homeLeague != awayLeague {
		crossLeagueBonus = 0.15 // Увеличен с 0.1 для большего эффекта
	}

	// 5) Финальный рейтинг
	baseRating := (homeStrength+awayStrength)/2.0*0.15 + // Уменьшено влияние позиций
		avgLeagueWeight*0.35 + // Снижено с 0.4
		formFactor*0.15 // Добавлено влияние формы
	rating := baseRating * (1 + derbyBonus + stageBonus + crossLeagueBonus)

	// 6) Ограничение и минимальное значение
	if rating > 1.0 {
		rating = 1.0
	}
	if rating < 0.1 { // Минимальный рейтинг, чтобы избежать нулей
		rating = 0.1
	}

	return rating, nil
}
