package domain

import (
	"context"
	"errors"
)

var (
	ErrCodeSendTooFrequently   = errors.New("code send too frequently")
	ErrCodeVerifyTooFrequently = errors.New("code verify too frequently")
	ErrCodeInvalid             = errors.New("code invalid")
)

type CodeUsecase interface {
	Send(ctx context.Context, biz, number string) error
	Verify(ctx context.Context, biz, number, code string) (bool, error)
}

type CodeRepository interface {
	SetCode(ctx context.Context, biz, number, code string) error
	VerifyCode(ctx context.Context, biz, number, code string) (bool, error)
}
