package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type userUsecase struct {
	repo           domain.UserRepository
	relationRepo   domain.RelationRepository
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

	return &domain.Profile{
		UserName: user.UserName, Email: user.Email,
		AboutMe: user.AboutMe, Birthday: user.Birthday,
		NickName: user.NickName, Avatar: user.Avatar,
	}, nil
}

func (uc *userUsecase) Logout(c context.Context, ssid string, tokenExpiry time.Duration) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.InvalidToken(ctx, ssid, tokenExpiry)
}

func (uc *userUsecase) UpdateProfile(c context.Context, id int64, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Update(ctx, id, user)
}
