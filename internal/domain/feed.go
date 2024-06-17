package domain

import (
	"context"
	"time"
)

const (
	FEEDLIKEEVENT   string = "feed-like"
	FEEDPOSTEVENT   string = "feed-post"
	FEEDFOLLOWEVENT string = "feed-follow"

	THRESHOLD int = 1000 // 读写扩散阈值
)

type FeedContent map[string]string

// FeedPush 推模型 - 写扩散 - 收件箱
type FeedPush struct {
	Model
	UserID    int64  `gorm:"index;not null"` // 收件人
	Type      string // 标记事件类型 决定Content解读方式
	Content   string
	CreatedAt time.Time `gorm:"index"`
	// 理论上来说没有 Update 操作，也没有 Delete 操作，但是考虑到文章可能有撤回操作
	// 可归档
}

// FeedPull 拉模型 - 读扩散
type FeedPull struct {
	Model
	UserID    int64 `gorm:"index;not null"` // 发件人
	Type      string
	Content   string
	CreatedAt time.Time `gorm:"index"`
}

type Feed struct {
	Model
	UserID    int64
	Type      string
	Content   FeedContent
	CreatedAt time.Time
}

// FeedUsecase 处理业务公共部分 并且负责找出 Handler 来处理业务的个性部分
type FeedUsecase interface {
	CreateFeedEvent(c context.Context, feed Feed) error
	GetFeedEventList(c context.Context, userID, timestamp, limit int64) ([]Feed, error)
}

// Handler 具体业务的处理逻辑 按照 type 类型来分，因为 type 天然的标记业务
type FeedHandler interface {
	CreateFeedEvent(c context.Context, t string, content FeedContent) error
	FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]Feed, error)
}

type FeedRepository interface {
	CreatePush(c context.Context, feed ...Feed) error
	CreatePull(c context.Context, feed ...Feed) error
	FindPush(c context.Context, userIDs []int64, timestamp, limit int64) ([]Feed, error)
	FindPull(c context.Context, userID, timestamp, limit int64) ([]Feed, error)
}
