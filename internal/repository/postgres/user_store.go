package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	sq "github.com/Masterminds/squirrel"
)

// Интерфейс для работы с пользователями в PostgreSQL
type UserStore interface {
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*types.User, error)
	SaveUser(ctx context.Context, user *types.User) error
}

// PGUserStore реализует интерфейс UserStore для работы с пользователями в PostgreSQL
type PGUserStore struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

// NewPGUserStore создает новый экземпляр PGUserStore
func NewPGUserStore(db *sql.DB) UserStore {
	return &PGUserStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// GetUserByTelegramID получает пользователя по его Telegram ID
// Возвращает указатель на пользователя или ошибку, если пользователь не найден
func (s *PGUserStore) GetUserByTelegramID(ctx context.Context, telegramID int64) (*types.User, error) {
	query := s.builder.Select("id", "telegram_id", "username", "created_at").
		From("users").
		Where(sq.Eq{"telegram_id": telegramID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, sqlStr, args...)

	var user types.User
	if err := row.Scan(&user.ID, &user.TelegramID, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// SaveUser сохраняет пользователя в базе данных
// Если пользователь с таким Telegram ID уже существует, ничего не делает
// Если пользователь не существует, выполняет вставку нового пользователя
// Возвращает ошибку, если не удалось выполнить вставку
func (s *PGUserStore) SaveUser(ctx context.Context, user *types.User) error {
	// Проверка существования пользователя по telegram_id
	var exists bool
	err := s.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM users WHERE telegram_id = $1
        )`, user.TelegramID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking user existence: %w", err)
	}

	if exists {
		log.Printf("User %d already exists, skipping insert", user.TelegramID)
		return nil
	}

	// Выполняем вставку, если пользователь не существует
	query := s.builder.Insert("users").
		Columns("telegram_id", "username").
		Values(user.TelegramID, user.Username)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building insert query: %w", err)
	}

	res, err := s.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("executing insert: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Printf("Inserted user %d (rows affected: %d)", user.TelegramID, rowsAffected)
	return nil
}
