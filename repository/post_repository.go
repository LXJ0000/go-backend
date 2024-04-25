package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/cache"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"time"
)

type postRepository struct {
	dao   orm.Database
	cache cache.Cache
}

func NewPostRepository(dao orm.Database, cache cache.Cache) domain.PostRepository {
	return &postRepository{dao: dao, cache: cache}
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

func (repo *postRepository) FindPage(c context.Context, filter interface{}, page, size int) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.FindPage(c, &domain.Post{}, filter, page, size, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
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
		Where("created_at < ? and status = ?", begin, domain.PostStatusPublish).
		Offset((page - 1) * size).Limit(size).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
