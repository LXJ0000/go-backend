package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
)

type FeedFollowHandler struct {
	feedRepo domain.FeedRepository
}

func NewFeedFollowHandler(feedRepo domain.FeedRepository) domain.FeedHandler {
	return &defaultHandler{feedRepo: feedRepo}
}

// CreateFeedEvent need: follower followee
func (h *FeedFollowHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	followee, err := lib.Str2Int64(content["followee"])
	if err != nil {
		return err
	}
	return h.feedRepo.CreatePush(ctx, domain.Feed{
		UserID:  followee, // 收件人 被点赞的人
		Type:    domain.FeedFollowEvent,
		Content: content,
	})
}

func (h *FeedFollowHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	return h.feedRepo.FindPushWithType(ctx, domain.FeedFollowEvent, userID, timestamp, limit)
}
