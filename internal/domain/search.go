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
