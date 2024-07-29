package domain

import "context"

const (
	ESPostIndex    = "post_index"
	ESUserIndex    = "user_index"
	ESCommentIndex = "comment_index"
)

// 命令与查询分离 CQRS

type ESSync interface {
	InputAny(ctx context.Context, item interface{}) error

	InputUser(ctx context.Context, item User) error
	InputPost(ctx context.Context, item Post) error
}

type ESSearch interface {
	Search()
}

type SearchUsecase interface {
	Search(ctx context.Context, userID int64, cmd string) error // 用户画像和表达式
}

type SearchRepository interface {
	SearchUser(ctx context.Context, keywords ...string) ([]User, error)
	SearchPost(ctx context.Context, userID int64, keywords ...string) ([]Post, error)
}

type SearchResult struct {
	Users []User
	Posts []Post
}
