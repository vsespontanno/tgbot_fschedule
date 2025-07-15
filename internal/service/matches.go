package service

import (
	"context"

	api "github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	db "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

type MatchesService struct {
	matchesStore db.MatchesStore
	apiClient    api.MatchApiClient
}

func NewMatchesService(matchesStore db.MatchesStore, apiClient api.MatchApiClient) *MatchesService {
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
