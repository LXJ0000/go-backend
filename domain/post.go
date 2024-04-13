package domain

import (
	"context"
	"gorm.io/gorm"
)

const (
	DefaultPage = 0
	DefaultSize = 10
)

const (
	PostStatusHide uint8 = iota
	PostStatusPublish
)

type Post struct {
	gorm.Model
	PostID int64 `json:"post_id" gorm:"primaryKey"`

	//Title    string `json:"title" form:"title"`
	//Abstract string `json:"abstract" form:"abstract"`
	//Content  string `json:"content" form:"content"`
	//
	//AuthorID int64 `json:"author_id" form:"author_id"`
	//
	//Status uint8 `json:"status" form:"status"`

	Title    string `json:"title" form:"title" binding:"required"`
	Abstract string `json:"abstract" form:"abstract" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	AuthorID int64  `json:"author_id,string" form:"author_id" binding:"required"`
	Status   uint8  `json:"status" form:"status" binding:"required"`
}

func (Post) TableName() string {
	return `post`
}

type PostRepository interface {
	Create(c context.Context, post *Post) error
	GetByID(c context.Context, id int64) (Post, error)
	FindMany(c context.Context, filter *Post, page, size int) ([]Post, error)
}

type PostUsecase interface {
	Create(c context.Context, post *Post) error
	List(c context.Context, filter *Post, page, size int) ([]Post, error)
	Info(c context.Context, postID int64) (Post, error)
}

type PostListRequest struct {
	Page     int   `json:"page" form:"page"`
	Size     int   `json:"size" form:"size"`
	AuthorID int64 `json:"author_id" form:"author_id"`
	Status   uint8 `json:"status" form:"status"`
}

type PostListResponse struct {
	Count int    `json:"count"`
	Data  []Post `json:"data"`
}
