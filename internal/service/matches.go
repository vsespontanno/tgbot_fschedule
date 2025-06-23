package service

import (
	"context"
	db "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/types"

	"go.mongodb.org/mongo-driver/mongo"
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
