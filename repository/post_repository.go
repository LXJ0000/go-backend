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

func (pr *postRepository) Create(c context.Context, post *domain.Post) error {
	_, err := pr.db.InsertOne(c, &domain.Post{}, post)
	return err
}
func (pr *postRepository) GetByID(c context.Context, id int64) (domain.Post, error) {
	post, err := pr.db.FindOne(c, &domain.Post{}, &domain.Post{PostID: id})
	return post.(domain.Post), err
}
