package db

import (
	"context"
	"database/sql"
	"football_tgbot/internal/types"

	sq "github.com/Masterminds/squirrel"
)

type UserStore interface {
	CreateOrUpdate(ctx context.Context, user *types.User) error
}

type PGUserStore struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewPGUserStore(db *sql.DB) *PGUserStore {
	return &PGUserStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
