package service

import (
	"context"

	userRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/postgres"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// UserService предоставляет методы для работы с пользователями
// Использует PostgreSQL для хранения и получения данных о пользователях
// Реализует интерфейс UserStore для взаимодействия с данными пользователей
type UserService struct {
	userStore userRepo.UserStore
}

// Конструктор для создания нового экземпляра UserService
func NewUserService(userStore userRepo.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// Метод для получения пользователя по его Telegram ID
func (s *UserService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*types.User, error) {
	return s.userStore.GetUserByTelegramID(ctx, telegramID)
}

// Метод для сохранения пользователя в базу данных
// Если пользователь с таким Telegram ID уже существует, выполняет обновление
// Если не существует, выполняет вставку нового пользователя
func (s *UserService) SaveUser(ctx context.Context, user *types.User) error {
	return s.userStore.SaveUser(ctx, user)
}
