package service

import (
	"context"

	userRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/postgres"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

type UserService struct {
	userStore userRepo.UserStore
}

func NewUserService(userStore userRepo.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*types.User, error) {
	return s.userStore.GetUserByTelegramID(ctx, telegramID)
}

func (s *UserService) SaveUser(ctx context.Context, user *types.User) error {
	return s.userStore.SaveUser(ctx, user)
}
