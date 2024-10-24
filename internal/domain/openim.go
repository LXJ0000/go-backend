package domain

import (
	"context"
	"errors"
)

var (
	ErrSync2OpenIMOpNotFound = errors.New("synchronous methods do not exist, only registration and editing are supported")
)

const (
	DefaultOpenIMDomain = "http://localhost:10002"
	DefaultOpenIMSecret = "openIM123"
	DefaultOpenIMUserID = "imAdmin"

	Sync2OpenIMOpRegister = "register"
	Sync2OpenIMOpEdit     = "edit"
)

type Sync2OpenIMUsecase interface {
	GetAdminToken(ctx context.Context) (string, error)

	SyncUser(ctx context.Context, user User, op string) error
}

type SyncUserRequest struct {
	Users []Sync2OpenIMUser `json:"users"`
}

type SyncUserResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}

type Sync2OpenIMUser struct {
	UserID   string `json:"userID"`
	NickName string `json:"nickname"`
	FaceUrl  string `json:"faceURL"`
}

type GetAdminTokenRequest struct {
	Secret string `json:"secret"`
	UserID string `json:"userID"`
}

type GetAdminTokenResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    struct {
		Token             string `json:"token"`
		ExpireTimeSeconds int    `json:"expireTimeSeconds"`
	} `json:"data"`
}
