package service

import (
	"context"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// Интerface Calculator определяет методы, необходимые для расчёта рейтинга матча
type Calculator interface {
	HandleGetLeague(ctx context.Context, collectionName string, teamID int) (string, error)
	HandleGetTeamStanding(ctx context.Context, leagueKey string, teamID int) (int, error)
	HandleGetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error)
	HandleGetTeamShortName(ctx context.Context, leagueKey string, teamName string) (string, error)
}
