package usecase

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/LXJ0000/go-utils/net/httpx"
)

type sync2OpenIMUsecase struct {
	domain string
	client *http.Client
}

func NewSync2OpenIMUsecase() domain.Sync2OpenIMUsecase {
	return &sync2OpenIMUsecase{client: http.DefaultClient, domain: domain.DefaultOpenIMDomain}
}

func (uc *sync2OpenIMUsecase) SyncUser(ctx context.Context, user domain.User, op string) error {
	token, err := uc.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	switch op {
	case "register":
		// https://docs.openim.io/zh-Hans/restapi/apis/userManagement/userRegister
		req := domain.SyncUserRequest{
			Users: convertUser(user),
		}
		var resp domain.SyncUserResponse
		err := httpx.NewRequest(ctx, http.MethodPost, uc.domain+"/user/user_register").
			Client(uc.client).
			BodyWithJSON(req).
			Header("operationID", lib.Int642Str(time.Now().UnixNano())).
			Header("token", token).
			Do().ScanJSON(&resp)
		if err != nil || resp.ErrCode != 0 {
			slog.Error("sync user to openIM fail", "error", err.Error(), "response", resp)
		}
		return nil
	default:
		return domain.ErrSync2OpenIMOpNotFound
	}
}

func (uc *sync2OpenIMUsecase) GetAdminToken(ctx context.Context) (string, error) {
	req := domain.GetAdminTokenRequest{
		Secret: domain.DefaultOpenIMSecret,
		UserID: domain.DefaultOpenIMUserID,
	}
	var resp domain.GetAdminTokenResponse
	err := httpx.NewRequest(ctx, http.MethodPost, uc.domain+"/auth/get_admin_token").
		Client(uc.client).
		BodyWithJSON(req).
		Header("operationID", lib.Int642Str(time.Now().UnixNano())).
		Do().ScanJSON(&resp)
	if err != nil || resp.ErrCode != 0 {
		return "", err
	}
	return resp.Data.Token, nil
}

func convertUser(users ...domain.User) []domain.Sync2OpenIMUser {
	var result []domain.Sync2OpenIMUser
	for _, u := range users {
		result = append(result, domain.Sync2OpenIMUser{
			UserID:   lib.Int642Str(u.UserID),
			NickName: u.UserName,
			FaceUrl:  u.Avatar,
		})
	}
	return result
}
