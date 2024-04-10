package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/domain"
)

type profileUsecase struct {
	repo           domain.UserRepository
	contextTimeout time.Duration
}

func NewProfileUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.ProfileUsecase {
	return &profileUsecase{
		repo:           userRepository,
		contextTimeout: timeout,
	}
}

func (uc *profileUsecase) GetProfileByID(c context.Context, userID int64) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.UserName, Email: user.Email}, nil
}
