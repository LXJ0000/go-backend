package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
)

type FeedLikeHandler struct {
	feedRepo domain.FeedRepository
}

func NewFeedLikeHandler(feedRepo domain.FeedRepository) domain.FeedHandler {
	return &FeedLikeHandler{feedRepo: feedRepo}
}

// CreateFeedEvent need: liker liked biz bizID
// liker 点赞者 liked 被点赞者 biz 被点赞的资源类型 bizID 被点赞的资源ID
func (h *FeedLikeHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	userID, err := lib.Str2Int64(content["liked"])
	if err != nil {
		return err
	}
	// 写到被点赞者的收件箱 userID 是被点赞者 即 userID = liked
	return h.feedRepo.CreatePush(ctx, domain.Feed{
		UserID:  userID, // 收件人 被点赞的人
		Type:    domain.FeedLikeEvent,
		Content: content,
	})
}

func (h *FeedLikeHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	return h.feedRepo.FindPushWithType(ctx, domain.FeedLikeEvent, userID, timestamp, limit)
}
