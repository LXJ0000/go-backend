package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	domain.UserUsecase
	domain.RelationUsecase
	domain.PostUsecase
	domain.CodeUsecase
	domain.Sync2OpenIMUsecase
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
	if err := col.UserUsecase.UpdateProfile(c, userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *UserController) Login(c *gin.Context) {
	var request domain.LoginReq

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
	
	// token
	accessToken, refreshToken, imToken, err := col.genToken(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_detail":   user,
		"im_token":      imToken,
	}))

}

func (col *UserController) LoginByPhone(c *gin.Context) {
	var request domain.LoginByPhoneReq

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}

	user, err := col.UserUsecase.GetUserByPhone(c, request.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("User not found with the given email", err))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Invalid credentials", err))
		return
	}

	// token
	accessToken, refreshToken, imToken, err := col.genToken(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_detail":   user,
		"im_token":      imToken,
	}))

}

func (col *UserController) Signup(c *gin.Context) {
	var request domain.SignupReq

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}

	if _, err := col.UserUsecase.GetUserByEmail(c, request.Email); err == nil {
		c.JSON(http.StatusConflict, domain.ErrorResp("User already exists with the given email", err))
		return
	}

	var err error
	if request.Password, err = genPassword(request.Password); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}

	user := domain.User{
		UserID:   snowflake.GenID(),
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
	}
	err = col.UserUsecase.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create user fail with db error", err))
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := col.Sync2OpenIMUsecase.SyncUser(ctx, user, domain.Sync2OpenIMOpRegister); err != nil {
			slog.Error("Sync user to openIM fail", "error", err.Error())
		}
	}()

	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *UserController) LoginBySms(c *gin.Context) {
	var req domain.LoginBySmsReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	ok, err := col.CodeUsecase.Verify(c, domain.BizUserLogin, req.Phone, req.Code)
	if err != nil {
		if errors.Is(err, domain.ErrCodeSendTooFrequently) {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrCodeInvalid.Error(), nil))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}
	if !ok {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrCodeInvalid.Error(), nil))
		return
	}

	// 用户是否存在 TODO 抽象 findOrCreate
	user, err := col.UserUsecase.GetUserByPhone(c, req.Phone)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
			return
		}

		user = domain.User{
			UserID:   snowflake.GenID(),
			UserName: fmt.Sprintf("%s_%d%s", domain.DefaultUserNamePrefix, rand.Intn(1000), req.Phone[7:]),
			Phone:    req.Phone,
		}
		user.Password, err = genPassword(domain.DefaultUserPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
			return
		}
		if err := col.UserUsecase.Create(c, &user); err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create user fail with db error", err))
			return
		}
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			if err := col.Sync2OpenIMUsecase.SyncUser(ctx, user, domain.Sync2OpenIMOpRegister); err != nil {
				slog.Error("Sync user to openIM fail", "error", err.Error())
			}
		}()
	}

	// token
	accessToken, refreshToken, imToken, err := col.genToken(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_detail":   user,
		"im_token":      imToken,
	}))
}

func (col *UserController) SendSMSCode(c *gin.Context) {
	var req domain.SendSMSCodeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	if err := col.CodeUsecase.Send(c, domain.BizUserLogin, req.Phone); err != nil {
		if errors.Is(err, domain.ErrCodeSendTooFrequently) {
			c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrCodeSendTooFrequently.Error(), err))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrSystemError.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func genPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(encryptedPassword), nil
}

func (col *UserController) genToken(ctx context.Context, user domain.User) (string, string, string, error) {
	ssid := uuid.New().String()
	accessToken, err := col.UserUsecase.CreateAccessToken(user, ssid, col.Env.AccessTokenSecret, col.Env.AccessTokenExpiryHour)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := col.UserUsecase.CreateRefreshToken(user, ssid, col.Env.RefreshTokenSecret, col.Env.RefreshTokenExpiryHour)
	if err != nil {
		return "", "", "", err
	}

	imToken, err := col.Sync2OpenIMUsecase.GetUserToken(ctx, 1, user.UserID)
	if err != nil {
		slog.Error("Get user token from openIM fail", "error", err.Error())
	}
	return accessToken, refreshToken, imToken, nil
}
