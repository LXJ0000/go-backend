package repository

import (
	"context"
	"encoding/json"
	"time"

	domain2 "github.com/LXJ0000/go-backend/internal/domain"
	cache "github.com/LXJ0000/go-backend/pkg/cache"
)

type postRankRepository struct {
	localCache *cache.RistrettoCache
	redisCache *cache.RedisCache
}

func NewPostRankRepository(localCache *cache.RistrettoCache, redisCache *cache.RedisCache) domain2.RankRepository {
	return &postRankRepository{
		localCache: localCache,
		redisCache: redisCache,
	}
}

func (repo *postRankRepository) ReplaceTopN(c context.Context, items []domain2.Post, expiration time.Duration) error {
	// ----------------------------------------------------------- local
	//_ = repo.localCache.Set(c, items) //必然不会出错
	// ----------------------------------------------------------- redis
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return repo.redisCache.Set(c, domain2.PostTopNKey, data, 0)
}

func (repo *postRankRepository) GetTopN(c context.Context) ([]domain2.Post, error) {
	// ----------------------------------------------------------- local
	//posts, err := repo.localCache.Get(c)
	//if err == nil {
	//	return posts, nil
	//}
	// ----------------------------------------------------------- redis
	data, err := repo.redisCache.Get(c, domain2.PostTopNKey)
	if err != nil {
		return nil, err
	}
	var items []domain2.Post
	if err = json.Unmarshal([]byte(data), &items); err != nil {
		return nil, err
	}
	//_ = repo.localCache.Set(c, items) // restore local cache
	return items, nil
}

//func (repo *postRepository) ReplaceTopN(c context.Context, items []domain.Post, expiration time.Duration) error {
//
//}
//
//func (repo *postRepository) GetTopN(c context.Context) ([]domain.Post, error) {
//
//}
//
//// TODO 1. 预加载 2. 分布式环境下 通知其他机器缓存 redis 到本地 3.
//// TODO redis 奔溃... 强制从本地缓存取出数据 不检查过期时间
//// 考虑到新节点一开始就没有数据，可以强制要求一定要有数据
