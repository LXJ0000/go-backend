package cache

import (
	"context"
	"errors"
	"github.com/LXJ0000/go-backend/domain"
	"sync/atomic"
	"time"
)

type LocalCache struct {
	topN       *atomic.Value
	ddl        *atomic.Value
	expiration time.Duration
}

func NewRankLocalCache() *LocalCache {
	return &LocalCache{
		topN:       &atomic.Value{},
		ddl:        &atomic.Value{},
		expiration: time.Hour, // TODO 考虑过期时间 可以超长的 也可以不过期
	}
}

func (r *LocalCache) Set(c context.Context, posts []domain.Post) error {
	r.topN.Store(posts) // TODO 两个原子操作 并发安全问题 用一个 struct 解决
	ddl := time.Now().Add(r.expiration)
	r.ddl.Store(ddl)
	return nil
}

func (r *LocalCache) Get(c context.Context) ([]domain.Post, error) {
	ddl := r.ddl.Load().(time.Time)
	posts := r.topN.Load().([]domain.Post)
	if len(posts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("") // TODO
	}
	return posts, nil
}
