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
func (h *FeedLikeHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	userID, err := lib.Str2Int64(content["liked"])
	if err != nil {
		return err
	}
	return h.feedRepo.CreatePush(ctx, domain.Feed{
		UserID:  userID, // 收件人 被点赞的人
		Type:    domain.FeedLikeEvent,
		Content: content,
	})
}
func (h *FeedLikeHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil
}
