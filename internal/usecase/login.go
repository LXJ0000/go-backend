package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/LXJ0000/go-backend/utils/tokenutil"
)

type loginUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uc *loginUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.GetByEmail(ctx, email)
}

func (uc *loginUsecase) CreateAccessToken(user domain.User, ssid string, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, ssid, secret, expiry)
}

func (uc *loginUsecase) CreateRefreshToken(user domain.User, ssid string, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, ssid, secret, expiry)
}
