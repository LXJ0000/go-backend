package controller

import (
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	domain.UserUsecase
	domain.RelationUsecase
	domain.PostUsecase
	Env *bootstrap.Env
}

func (col *UserController) Profile(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", nil))
		return
	}
	userID, err := lib.Str2Int64(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	profile, err := col.UserUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get profile by user_id fail with db error", err))
		return
	}
	stat, err := col.RelationUsecase.Stat(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get relation stat fail with db error", err))
		return
	}
	profile.RelationStat = stat
	postCnt, err := col.PostUsecase.Count(c, map[string]interface{}{
		"author_id": userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get post count fail with db error", err))
		return
	}
	profile.PostCnt = postCnt
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"profile": profile,
	}))
}

func (col *UserController) Fetch(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	profile, err := col.UserUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get profile by user_id fail with db error", err))
		return
	}
	stat, err := col.RelationUsecase.Stat(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get relation stat fail with db error", err))
		return
	}
	profile.RelationStat = stat

	postCnt, err := col.PostUsecase.Count(c, map[string]interface{}{
		"author_id": userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get post count fail with db error", err))
		return
	}
	profile.PostCnt = postCnt

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

func (col *UserController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}

	user, err := col.UserUsecase.GetUserByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("User not found with the given email", err))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Invalid credentials", err))
		return
	}
	ssid := uuid.New().String()
	accessToken, err := col.UserUsecase.CreateAccessToken(user, ssid, col.Env.AccessTokenSecret, col.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create access token fail", err))
		return
	}

	refreshToken, err := col.UserUsecase.CreateRefreshToken(user, ssid, col.Env.RefreshTokenSecret, col.Env.RefreshTokenExpiryHour)
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

func (col *UserController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}

	_, err = col.UserUsecase.GetUserByEmail(c, request.Email)
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

	user := domain.User{
		UserID:   snowflake.GenID(),
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
	}
	err = col.UserUsecase.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create user fail with db error", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}
