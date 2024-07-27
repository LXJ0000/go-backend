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
	userID := c.MustGet(domain.XUserID).(int64)
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
	ssid := c.MustGet(domain.UserSessionID).(string)
	tokenExpiry := time.Duration(col.Env.RefreshTokenExpiryHour) * time.Hour
	if err := col.UserUsecase.Logout(c, ssid, tokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Logout fail with db error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *UserController) Update(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	req := struct {
		NickName string `json:"nick_name" form:"nick_name"`
		Birthday string `json:"birthday" form:"birthday"`
		AboutMe  string `json:"about_me" form:"about_me"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	user := domain.User{}
	if req.Birthday != "" {
		birthday, err := time.Parse(time.DateOnly, req.Birthday)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResp("Update fail with invalid birthday", err))
			return
		}
		user.Birthday = birthday
	}
	user.AboutMe = req.AboutMe
	user.NickName = req.NickName
	if err := col.UserUsecase.UpdateProfile(c, userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}
