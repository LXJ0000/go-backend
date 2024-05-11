package domain

import (
	"context"
	"time"
)

const (
	PostTopNKey = "post_topN"

	PostStatusHide uint8 = iota
	PostStatusPublish
)

type Post struct {
	Model
	PostID int64 `json:"post_id,string" gorm:"primaryKey"`

	Title    string `json:"title" form:"title" binding:"required"`
	Abstract string `json:"abstract" form:"abstract" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	AuthorID int64  `json:"author_id,string" form:"author_id" binding:"required"`
	Status   uint8  `json:"status" form:"status" binding:"required"`
}

func (Post) TableName() string {
	return `post`
}

//go:generate mockgen -source=./post.go -destination=./mock/post.go -package=domain_mock
type PostRepository interface {
	Create(c context.Context, post Post) error
	GetByID(c context.Context, id int64) (Post, error)
	FindMany(c context.Context, filter interface{}) ([]Post, error) // Modify
	List(c context.Context, filter interface{}, page, size int) ([]Post, error)
	FindTopNPage(c context.Context, page, size int, begin time.Time) ([]Post, error)
	Count(c context.Context, filter interface{}) (int64, error)
}

type PostUsecase interface {
	Create(c context.Context, post Post) error
	List(c context.Context, filter interface{}, page, size int) ([]Post, int64, error)
	Info(c context.Context, postID int64) (Post, error)
	//ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
}

type PostListRequest struct {
	Page     int   `json:"page" form:"page"`
	Size     int   `json:"size" form:"size"`
	AuthorID int64 `json:"author_id" form:"author_id"`
	Status   uint8 `json:"status" form:"status"`
}

//type PostListResponse struct {
//	Count int    `json:"count"`
//	Data  []Post `json:"data"`
//}
