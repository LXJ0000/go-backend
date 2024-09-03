package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
)

type defaultHandler struct {
	feedRepo domain.FeedRepository
}

func NewFeedDefaultHandler(feedRepo domain.FeedRepository) domain.FeedHandler {
	return &defaultHandler{feedRepo: feedRepo}
}

func (h *defaultHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	userID := lib.Str2Int64DefaultZero(content["user_id"])
	return h.feedRepo.CreatePush(ctx, domain.Feed{
		UserID: userID, Type: t, Content: content,
	})
}

func (h *defaultHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil
}
