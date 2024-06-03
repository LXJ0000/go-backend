package domain

import (
	"context"
	"time"
)

type Content map[string]string

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
	Content   Content
	CreatedAt time.Time
}

// FeedUsecase 处理业务公共部分 并且负责找出 Handler 来处理业务的个性部分
type FeedUsecase interface {
	CreateFeedEvent(c context.Context, feed Feed) error
	GetFeedEventList(c context.Context, userID, timestamp, limit int64) ([]Feed, error)
}

type Handler interface {
	CreateFeedEvent(c context.Context, content Content) error
	FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]Feed, error)
}

type FeedRepository interface {
	// CreatePush(c context.Context)
	// CreatePull(c context.Context)
	// FindPush(c context.Context)
	// FindPull(c context.Context)
}
