package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
)

type FeedInterHandler struct {
	feedRepo domain.FeedRepository
}

func NewFeedInterHandler(feedRepo domain.FeedRepository) domain.FeedHandler {
	return &FeedInterHandler{feedRepo: feedRepo}
}

// CreateFeedEvent need: from to biz bizID
func (h *FeedInterHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	userID, err := lib.Str2Int64(content["to"])
	if err != nil {
		return err
	}
	return h.feedRepo.CreatePush(ctx, domain.Feed{
		UserID:  userID,
		Type:    t,
		Content: content,
	})
}

func (h *FeedInterHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	return h.feedRepo.FindPush(ctx, userID, timestamp, limit)
	// return h.feedRepo.FindPushWithType(ctx, domain.FeedLikeEvent, userID, timestamp, limit)
}
