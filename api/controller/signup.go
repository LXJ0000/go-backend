package controller

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"net/http"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}

	_, err = sc.SignupUsecase.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, domain.ErrorResp("User already exists with the given email", err))
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Encrypted password fail", err))
		return
	}

	request.Password = string(encryptedPassword)

	//now := time.Now().UnixMicro()
	user := domain.User{
		UserID:   snowflake.GenID(),
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
	}
	//user.CreatedAt = now
	//user.UpdatedAt = now
	err = sc.SignupUsecase.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create user fail with db error", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(nil))

	// accessToken, err := sc.SignupUsecase.CreateAccessToken(user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create access token fail", err))
	// 	return
	// }

	// refreshToken, err := sc.SignupUsecase.CreateRefreshToken(user, sc.Env.RefreshTokenSecret, sc.Env.RefreshTokenExpiryHour)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create refresh token fail", err))
	// 	return
	// }

	// c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
	// 	"access_token":  accessToken,
	// 	"refresh_token": refreshToken,
	// }))
}
