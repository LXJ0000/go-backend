package repository

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
)

type postRepository struct {
	dao        orm.Database
	redisCache cache.RedisCache
}

func NewPostRepository(dao orm.Database, redisCache cache.RedisCache) domain.PostRepository {
	return &postRepository{dao: dao, redisCache: redisCache}
}

func (repo *postRepository) Update(c context.Context, id int64, post *domain.Post) error {
	return repo.dao.UpdateOne(c, &domain.Post{}, map[string]interface{}{"user_id": id}, post)
}

func (repo *postRepository) Create(c context.Context, post *domain.Post) error {
	return repo.dao.Insert(c, &domain.Post{}, post)
}

func (repo *postRepository) GetByID(c context.Context, id int64) (domain.Post, error) {
	var post domain.Post
	err := repo.dao.FindOne(c, &domain.Post{}, map[string]interface{}{"post_id": id}, &post)
	return post, err
}

func (repo *postRepository) List(c context.Context, filter interface{}, page, size int) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.FindPage(c, &domain.Post{}, filter, page, size, &items)
	return items, err
}

func (repo *postRepository) FindMany(c context.Context, filter interface{}) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.FindMany(c, &domain.Post{}, filter, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *postRepository) FindTopNPage(c context.Context, page, size int, begin time.Time) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.Raw(c).Model(&domain.Post{}).
		Where("created_at < ? and status = ?", begin.Unix(), domain.PostStatusPublish).
		Offset((page - 1) * size).Limit(size).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *postRepository) Count(c context.Context, filter interface{}) (int64, error) {
	return repo.dao.Count(c, &domain.Post{}, filter)
}
