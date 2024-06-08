package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/LXJ0000/go-lib/slice"
)

type FeedPostHandler struct {
	feedRepo        domain.FeedRepository
	relationUsecase domain.RelationUsecase
}

func NewFeedPostHandler(feedRepo domain.FeedRepository, relationUsecase domain.RelationUsecase) domain.FeedHandler {
	return &FeedPostHandler{feedRepo: feedRepo}
}

// CreateFeedEvent need: user_id
func (h *FeedPostHandler) CreateFeedEvent(c context.Context, t string, content domain.FeedContent) error {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	follower, err := lib.Str2Int64(content["user_id"])
	if err != nil {
		return err
	}
	followers, cnt, err := h.relationUsecase.GetFollower(ctx, follower, -1, -1)
	if err != nil {
		return err
	}
	switch { // 粉丝数超过阈值，则读扩散，否则写扩散
	case cnt > domain.THRESHOLD:
		return h.feedRepo.CreatePull(ctx, domain.Feed{
			Type:   domain.FEEDPOSTEVENT,
			UserID: follower,
		})
	default:
		events := slice.Map(followers, func(user domain.User) domain.Feed {
			return domain.Feed{
				UserID:  user.UserID,
				Type:    domain.FEEDPOSTEVENT,
				Content: content,
			}
		})
		return h.feedRepo.CreatePush(c, events...)
	}
}

func (h *FeedPostHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil
}
