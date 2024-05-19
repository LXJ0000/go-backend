package usecase

import (
	"context"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type FeedLikeHandler struct {
	feedRepo domain.FeedRepository
}

func NewFeedLikeHandler(feedRepo domain.FeedRepository) *FeedLikeHandler {
	return &FeedLikeHandler{feedRepo: feedRepo}
}

func CreateFeedEvent(c context.Context, content domain.Content) error {
	return nil
}

func FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil
}

type FeedFollowHandler struct {
}

type FeedPostHandler struct {
}
