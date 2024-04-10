package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
)

type postRepository struct {
	db orm.Database
}

func NewPostRepository(db orm.Database) domain.PostRepository {
	return &postRepository{db: db}
}

func (repo *postRepository) Create(c context.Context, post *domain.Post) error {
	_, err := repo.db.InsertOne(c, &domain.Post{}, post)
	return err
}
func (repo *postRepository) GetByID(c context.Context, id int64) (domain.Post, error) {
	post, err := repo.db.FindOne(c, &domain.Post{}, &domain.Post{PostID: id})
	return post.(domain.Post), err
}
