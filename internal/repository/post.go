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
	redisCache *cache.RedisCache
}

func NewPostRepository(dao orm.Database, redisCache *cache.RedisCache) domain.PostRepository {
	return &postRepository{dao: dao, redisCache: redisCache}
}

func (u *postRepository) Search(c context.Context, keyword string, page, size int) ([]domain.Post, int, error) {
	var (
		items []domain.Post
		count int64
	)
	db := u.dao.Raw(c)
	q := db.Model(&domain.Post{}).Where("title LIKE ? or content LIKE ? or abstract LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, int(count), nil
}

func (repo *postRepository) ListByLastID(c context.Context, filter interface{}, size int, last int64) ([]domain.Post, error) {
	var items []domain.Post
	err := repo.dao.FindByLastIDRev(c, &domain.Post{}, filter, size, last, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *postRepository) Update(c context.Context, id int64, post *domain.Post) error {
	return repo.dao.UpdateOne(c, &domain.Post{}, map[string]interface{}{"post_id": id}, post)
}

func (repo *postRepository) Delete(c context.Context, postID int64) error {
	return repo.dao.DeleteOne(c, &domain.Post{}, map[string]interface{}{"post_id": postID})
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
	err := repo.dao.FindPageRev(c, &domain.Post{}, filter, page, size, &items)
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
