package db

import (
	"context"
	"database/sql"
	"football_tgbot/internal/types"
	"time"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetUserByID(ctx context.Context, id int64) (*types.User, error)
}

type PostgresUserStore struct {
	DB *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{DB: db}
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, user *types.User) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO users (id, username, first_name, last_name, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         ON CONFLICT (id) DO NOTHING`,
		user.ID, user.Username, user.FirstName, user.LastName, time.Now(), time.Now())
	return err
}

func (s *PostgresUserStore) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT id, username, first_name, last_name, created_at, updated_at FROM users WHERE id = $1`, id)
	var user types.User
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}
