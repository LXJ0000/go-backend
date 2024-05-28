package domain

import (
	"context"
)

type SignupRequest struct {
	UserName string `form:"user_name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupUsecase interface {
	Create(c context.Context, user User) error
	GetUserByEmail(c context.Context, email string) (User, error)
	// CreateAccessToken(user User, ssid string, secret string, expiry int) (accessToken string, err error)
	// CreateRefreshToken(user User, ssid string, secret string, expiry int) (refreshToken string, err error)
}
