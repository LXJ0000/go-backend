package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	DefaultPage = 0
	DefaultSize = 10

	PostTopNKey = "post_topN"
)

const (
	PostStatusHide uint8 = iota
	PostStatusPublish
)

type Post struct {
	gorm.Model
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
	Create(c context.Context, post *Post) error
	GetByID(c context.Context, id int64) (Post, error)
	FindMany(c context.Context, filter interface{}, page, size int) ([]Post, error) // Modify
	ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
	GetTopN(c context.Context) ([]Post, error)
}

type PostUsecase interface {
	Create(c context.Context, post *Post) error
	List(c context.Context, filter interface{}, page, size int) ([]Post, error)
	Info(c context.Context, postID int64) (Post, error)
	ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
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

type PostV0 struct {
	Post
	ReadCnt    int
	LikeCnt    int
	CollectCnt int

	Collected bool
	Liked     bool
}
