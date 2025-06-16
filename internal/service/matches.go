package service

import (
	"context"
	db "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"
)

type MatchesService struct {
	matchesStore db.MatchesStore
	// ratingService *RatingService   -- TODO
}

func NewMatchesService(matchesStore db.MatchesStore) *MatchesService {
	return &MatchesService{
		matchesStore: matchesStore,
	}

}

func (s *MatchesService) HandleGetMatches(ctx context.Context) ([]types.Match, error) {
	matches, err := s.matchesStore.GetMatches(ctx, "matches")
	if err != nil {
		return nil, err
	}
	return matches, nil
}
