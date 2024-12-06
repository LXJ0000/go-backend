package controller

import (
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RefreshTokenController struct {
	RefreshTokenUsecase domain.RefreshTokenUsecase
	UserUsecase         domain.UserUsecase
	Env                 *bootstrap.Env
}

func (col *RefreshTokenController) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}

	userID, oldSsid, err := col.RefreshTokenUsecase.ExtractIDAndSSIDFromToken(request.RefreshToken, col.Env.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("User not found", err))
		return
	}

	tokenExpiry := time.Duration(col.Env.RefreshTokenExpiryHour) * time.Hour
	if err := col.UserUsecase.Logout(c, oldSsid, tokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("logout fail", err))
		return
	}

	user, err := col.RefreshTokenUsecase.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("User not found", err))
		return
	}

	ssid := uuid.New().String()
	accessToken, err := col.RefreshTokenUsecase.CreateAccessToken(user, ssid, col.Env.AccessTokenSecret, col.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("create access token", err))
		return
	}

	refreshToken, err := col.RefreshTokenUsecase.CreateRefreshToken(user, ssid, col.Env.RefreshTokenSecret, col.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("create refresh token", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}
