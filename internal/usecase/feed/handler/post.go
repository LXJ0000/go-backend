package usecase

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sort"
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
	static, err := h.relationUsecase.Stat(ctx, follower)
	if err != nil {
		return err
	}

	// 考虑使用读扩散还是写扩散，考虑其他情况：铁粉、活跃用户
	switch { // 粉丝数超过阈值，则读扩散，否则写扩散
	case static.Follower > domain.THRESHOLD:
		return h.feedRepo.CreatePull(ctx, domain.Feed{
			Type:   domain.FeedPostEvent,
			UserID: follower,
		})
	default:
		followers, _, err := h.relationUsecase.GetFollower(ctx, follower, -1, -1)
		if err != nil {
			return err
		}
		events := slice.Map(followers, func(user domain.User) domain.Feed {
			return domain.Feed{
				UserID:  user.UserID,
				Type:    domain.FeedPostEvent,
				Content: content,
			}
		})
		return h.feedRepo.CreatePush(c, events...)
	}
}

func (h *FeedPostHandler) FindFeedEvent(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	ctx, cancel := context.WithTimeout(c, time.Second)
	defer cancel()
	// 1. 查收件箱
	// 2. 查发件箱
	var (
		g         = errgroup.Group{}
		pushEvent []domain.Feed
		pullEvent []domain.Feed
	)
	g.Go(func() error {
		// TODO 降级策略 跳过
		var err error
		followees, _, err := h.relationUsecase.GetFollowee(ctx, userID, -1, -1)
		if err != nil {
			return err
		}
		userIDs := slice.Map(followees, func(user domain.User) int64 {
			return user.UserID
		})
		pushEvent, err = h.feedRepo.FindPullWithType(ctx, domain.FeedPostEvent, userIDs, timestamp, limit)
		return err
	})
	g.Go(func() error {
		var err error
		pullEvent, err = h.feedRepo.FindPushWithType(ctx, domain.FeedPostEvent, userID, timestamp, limit)
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	// 3. 合并、排序（按照时间戳倒叙排序）、分页
	events := append(pushEvent, pullEvent...)
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
	return events[:min(len(events), int(limit))], nil
}
