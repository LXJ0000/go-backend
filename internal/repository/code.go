package repository

import (
	"context"
	"fmt"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/script"
	"log/slog"
)

type codeRepository struct {
	cache cache.RedisCache
}

func NewCodeRepository(cache cache.RedisCache) domain.CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}

func (r *codeRepository) SetCode(ctx context.Context, biz, number, code string) error {
	codeKey := fmt.Sprintf("code:%s:%s", biz, number)
	res, err := r.cache.LuaWithReturnInt(ctx, script.LuaSendCode, []string{codeKey}, code)
	if err != nil {
		slog.Error("set code error", "error", err.Error())
		return err
	}
	switch res {
	case 0:
		return nil
	case -2:
		slog.Error("set code error", "error", "code send too frequently")
		return domain.ErrCodeSendTooFrequently
	default:
		return domain.ErrSystemError
	}
}

func (r *codeRepository) VerifyCode(ctx context.Context, biz, number, code string) (bool, error) {
	codeKey := fmt.Sprintf("code:%s:%s", biz, number)
	res, err := r.cache.LuaWithReturnInt(ctx, script.LuaVerifyCode, []string{codeKey}, code)
	if err != nil {
		slog.Error("verify code error", "error", err.Error())
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, domain.ErrCodeVerifyTooFrequently // 告警
	case -2:
		return false, nil
	default:
		return false, domain.ErrUnKnowError
	}
}
