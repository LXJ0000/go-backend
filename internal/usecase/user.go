package usecase

import (
	"context"
	"github.com/LXJ0000/go-backend/utils/tokenutil"
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

func (uc *userUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *userUsecase) CreateAccessToken(user domain.User, ssid string, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, ssid, secret, expiry)
}

func (uc *userUsecase) CreateRefreshToken(user domain.User, ssid string, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, ssid, secret, expiry)
}

func (uc *userUsecase) Create(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Create(ctx, user)
}
