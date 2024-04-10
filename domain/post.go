package domain

import (
	"context"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PostID int64 `json:"post_id" gorm:"primaryKey"`

	Title    string `json:"title" form:"title"`
	Abstract string `json:"abstract" form:"abstract"`
	Content  string `json:"content" form:"content"`

	UserID int64 `json:"user_id"`
}

func (Post) TableName() string {
	return `post`
}

type PostRepository interface {
	Create(c context.Context, user *Post) error
	GetByID(c context.Context, id int64) (Post, error)
}

type PostUsecase interface {
	Create(c context.Context, post *Post) error
}
