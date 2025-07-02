package domain

import (
	"context"
	"football_tgbot/internal/types"
)

type Calculator interface {
	HandleGetLeague(ctx context.Context, leagueKey string, teamID int) (string, error)
	HandleGetTeamStanding(ctx context.Context, leagueKey string, teamID int) (int, error)
	HandleGetRecentMatches(ctx context.Context, teamID int, lastN int) ([]types.Match, error)
	HandleGetTeamShortName(ctx context.Context, leagueKey string, teamName string) (string, error)
	HandleGetMatches(ctx context.Context) ([]types.Match, error)
}
