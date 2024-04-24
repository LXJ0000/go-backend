package repository

import (
	"context"
	"encoding/json"
	"github.com/LXJ0000/go-backend/cache"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"time"
)

type postRepository struct {
	dao        orm.Database
	cache      cache.Cache
	redisCache cache.RedisCache // or cache.Cache
	localCache cache.LocalCache // TODO move rank repository
}

func NewPostRepository(dao orm.Database, cache cache.Cache) domain.PostRepository {
	return &postRepository{dao: dao, cache: cache} // TODO
}

func (repo *postRepository) Create(c context.Context, post *domain.Post) error {
	_, err := repo.dao.InsertOne(c, &domain.Post{}, post)
	return err
}

func (repo *postRepository) GetByID(c context.Context, id int64) (domain.Post, error) {
	post, err := repo.dao.FindOne(c, &domain.Post{}, &domain.Post{PostID: id})
	if err != nil {
		return domain.Post{}, err
	}
	return *post.(*domain.Post), err
}

func (repo *postRepository) FindMany(c context.Context, filter interface{}, page, size int) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.FindMany(c, &domain.Post{}, filter, page, size, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *postRepository) ReplaceTopN(c context.Context, items []domain.Post, expiration time.Duration) error {
	// ----------------------------------------------------------- local
	_ = repo.localCache.Set(c, items) //必然不会出错
	// ----------------------------------------------------------- redis
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return repo.cache.Set(c, domain.PostTopNKey, data, expiration)
}

func (repo *postRepository) GetTopN(c context.Context) ([]domain.Post, error) {
	// ----------------------------------------------------------- local
	posts, err := repo.localCache.Get(c)
	if err == nil {
		return posts, nil
	}
	// ----------------------------------------------------------- redis
	data, err := repo.redisCache.Get(c, domain.PostTopNKey)
	if err != nil {
		return nil, err
	}
	var items []domain.Post
	if err = json.Unmarshal([]byte(data), &items); err != nil {
		return nil, err
	}
	_ = repo.localCache.Set(c, items) // restore local cache
	return items, nil
}

// TODO 1. 预加载 2. 分布式环境下 通知其他机器缓存 redis 到本地 3.
// TODO redis 奔溃... 强制从本地缓存取出数据 不检查过期时间
// 考虑到新节点一开始就没有数据，可以强制要求一定要有数据
