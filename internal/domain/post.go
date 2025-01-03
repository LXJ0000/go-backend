package domain

import (
	"context"
	"time"
)

const (
	PostTopNKey = "post_topN"

	PostStatusHide    string = "hide"
	PostStatusPublish string = "publish"

	PromptOfPostAbstract = "请你根据以下内容,生成文章摘要,不多于20字.你的输出将会直接被用于生成文章摘要.因此除了文章摘要,请不要包含其他内容.\\n"
)

type Post struct {
	Model
	PostID int64 `json:"post_id,string" gorm:"primaryKey"`

	Title    string `json:"title" form:"title"`
	Abstract string `json:"abstract" form:"abstract"`
	Content  string `json:"content" form:"content"`
	AuthorID int64  `json:"author_id,string" form:"author_id"`
	Status   string `json:"status" form:"status" binding:"required"`

	Images string `json:"images" form:"images"`
}

func (Post) TableName() string {
	return `post`
}

//go:generate mockgen -source=./post.go -destination=./mock/post.go -package=domain_mock
type PostRepository interface {
	Create(c context.Context, post *Post) error
	Delete(c context.Context, postID int64) error
	GetByID(c context.Context, id int64) (Post, error)
	FindMany(c context.Context, filter interface{}) ([]Post, error) // Modify
	List(c context.Context, filter interface{}, page, size int) ([]Post, error)
	ListByLastID(c context.Context, filter interface{}, size int, lastID int64) ([]Post, error)
	FindTopNPage(c context.Context, page, size int, begin time.Time) ([]Post, error)
	Count(c context.Context, filter interface{}) (int64, error)
	Update(c context.Context, id int64, post *Post) error
	Search(c context.Context, keyword string, page, size int, rule []SortRule) ([]Post, int, error)
}

type PostUsecase interface {
	Create(c context.Context, post *Post) error
	Delete(c context.Context, postID int64) error
	List(c context.Context, filter interface{}, page, size int) ([]Post, int64, error)
	ListByLastID(c context.Context, filter interface{}, size int, lastID int64) ([]Post, int64, error)
	Info(c context.Context, postID int64) (Post, error)
	TopN(c context.Context) ([]Post, error)
	Count(c context.Context, filter interface{}) (int64, error)
	Search(c context.Context, keyword string, page, size int, rule []SortRule) ([]Post, int, error)
	//ReplaceTopN(c context.Context, items []Post, expiration time.Duration) error
}

type PostListRequest struct {
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
	AuthorID int64  `json:"author_id,string" form:"author_id"`
	Status   string `json:"status" form:"status"`
	Last     int64  `json:"last,string" form:"last"`
}

//type PostListResponse struct {
//	Count int    `json:"count"`
//	Data  []Post `json:"data"`
//}

type PostResponse struct {
	Model
	PostID int64 `json:"post_id,string"`

	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content" `
	// AuthorID int64  `json:"author_id,string" form:"author_id"`
	Status string   `json:"status"`
	Author Profile  `json:"author"`
	Images []string `json:"images"`
}

type PostInfoResponse struct {
	Post         PostResponse        `json:"post"`
	Interaction  Interaction         `json:"interaction"`
	Stat         UserInteractionStat `json:"stat"`
	CommentCount int                 `json:"comment_count"`
}

type PostDeleteRequest struct {
	PostID int64 `json:"post_id,string" form:"post_id"`
}

type PostSearchRequest struct {
	Keyword string     `json:"keyword" form:"keyword"`
	Page    int        `json:"page" form:"page"`
	Size    int        `json:"size" form:"size"`
	Sort    []SortRule `json:"sort" form:"sort"`
}
