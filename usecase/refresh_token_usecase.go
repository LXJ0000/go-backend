package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/internal/tokenutil"
)

type refreshTokenUsecase struct {
	repo           domain.UserRepository
	contextTimeout time.Duration
}

func NewRefreshTokenUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		repo:           userRepository,
		contextTimeout: timeout,
	}
}

func (uc *refreshTokenUsecase) GetUserByID(c context.Context, id int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.GetByID(ctx, id)
}

func (uc *refreshTokenUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (uc *refreshTokenUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func (uc *refreshTokenUsecase) ExtractIDFromToken(requestToken string, secret string) (int64, error) {
	return tokenutil.ExtractIDFromToken(requestToken, secret)
}
