package domain

import (
	"context"
	"time"
)

type RankUsecase interface {
	TopN(c context.Context) error
}

type RankRepository interface {
	ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
	GetTopN(c context.Context) ([]Post, error)
}
