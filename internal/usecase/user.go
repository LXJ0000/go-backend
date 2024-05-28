package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type userUsecase struct {
	repo           domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		repo:           userRepository,
		contextTimeout: timeout,
	}
}

func (uc *userUsecase) GetProfileByID(c context.Context, userID int64) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.UserName, Email: user.Email}, nil
}

func (uc *userUsecase) Logout(c context.Context, ssid string, tokenExpiry time.Duration) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.InvalidToken(ctx, ssid, tokenExpiry)
}
