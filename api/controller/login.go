package controller

import (
	"net/http"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

func (col *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}

	user, err := col.LoginUsecase.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("User not found with the given email", err))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Invalid credentials", err))
		return
	}
	ssid := uuid.New().String()
	accessToken, err := col.LoginUsecase.CreateAccessToken(user, ssid, col.Env.AccessTokenSecret, col.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create access token fail", err))
		return
	}

	refreshToken, err := col.LoginUsecase.CreateRefreshToken(user, ssid, col.Env.RefreshTokenSecret, col.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create refresh token fail", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_detail":   user,
	}))

}
