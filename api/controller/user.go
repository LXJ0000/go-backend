package controller

import (
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (col *UserController) Fetch(c *gin.Context) {
	userID := c.MustGet(domain.USERCTXID).(int64)
	profile, err := col.UserUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get profile by user_id fail with db error", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"profile": profile,
	}))
}

func (col *UserController) Logout(c *gin.Context) {
	ssid := c.MustGet(domain.USERSESSIONID).(string)
	tokenExpiry := time.Duration(col.Env.RefreshTokenExpiryHour) * time.Hour
	if err := col.UserUsecase.Logout(c, ssid, tokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Logout fail with db error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}
