package controller

import (
	"github.com/LXJ0000/go-backend/bootstrap"
	domain "github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RefreshTokenController struct {
	RefreshTokenUsecase domain.RefreshTokenUsecase
	Env                 *bootstrap.Env
}

func (col *RefreshTokenController) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}

	userID, err := col.RefreshTokenUsecase.ExtractIDFromToken(request.RefreshToken, col.Env.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("User not found", err))
		return
	}

	user, err := col.RefreshTokenUsecase.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResp("User not found", err))
		return
	}

	accessToken, err := col.RefreshTokenUsecase.CreateAccessToken(user, col.Env.AccessTokenSecret, col.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("create access token", err))
		return
	}

	refreshToken, err := col.RefreshTokenUsecase.CreateRefreshToken(user, col.Env.RefreshTokenSecret, col.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("create refresh token", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}
