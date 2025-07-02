package service

import (
	"context"
	api "football_tgbot/internal/infrastructure/api"
	db "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"

	"go.mongodb.org/mongo-driver/mongo"
)

type MatchesService struct {
	matchesStore db.MatchesStore
	apiClient    api.FootballDataClient
}

func NewMatchesService(matchesStore db.MatchesStore, apiClient api.FootballDataClient) *MatchesService {
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

func (s *MatchesService) HandleGetMatchByID(ctx context.Context, matchID int) (types.Match, error) {
	match, err := s.matchesStore.GetMatchByID(ctx, matchID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Match{}, nil
		}
		return types.Match{}, err
	}
	return match, nil
}

func (s *MatchesService) HandleSaveMatches(matches []types.Match, from, to string) error {
	err := s.matchesStore.SaveMatchesToMongoDB(matches, from, to)
	if err != nil {
		return err
	}
	return nil
}

func (s *MatchesService) HandleGetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error) {
	matches, err := s.matchesStore.GetRecentMatches(ctx, teamID, lastN)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *MatchesService) HandleReqMatches(ctx context.Context, from string, to string) ([]types.Match, error) {
	return s.apiClient.GetMatches(ctx, from, to)
}

func (s *MatchesService) HandleSaveMatchRating(ctx context.Context, match types.Match, rating float64) error {
	return s.matchesStore.UpdateMatchRatingInMongoDB(match, rating)
}

func (s *MatchesService) HandleUpsertMatch(ctx context.Context, match types.Match) error {
	return s.matchesStore.UpsertMatch(ctx, match)
}
