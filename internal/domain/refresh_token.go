package domain

import (
	"context"
)

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUsecase interface {
	GetUserByID(c context.Context, id int64) (User, error)
	CreateAccessToken(user User, ssid string, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user User, ssid string, secret string, expiry int) (refreshToken string, err error)
	ExtractIDAndSSIDFromToken(requestToken string, secret string) (int64, string, error)
}
