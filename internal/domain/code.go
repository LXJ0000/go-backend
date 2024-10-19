package domain

import (
	"context"
	"errors"
)

var (
	ErrCodeSendTooFrequently   = errors.New("code send too frequently")
	ErrCodeVerifyTooFrequently = errors.New("code verify too frequently")
)

type CodeRepository interface {
	SetCode(ctx context.Context, biz, number, code string) error
	VerifyCode(ctx context.Context, biz, number, code string) (bool, error)
}
