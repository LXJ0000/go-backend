package usecase

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-lib/slice"
	"golang.org/x/sync/errgroup"
)

type feedUsecase struct {
	handlerMap      map[string]domain.FeedHandler // map[type]handler
	relationUsecase domain.RelationUsecase
	repo            domain.FeedRepository
}

func NewFeedUsecase(handlerMap map[string]domain.FeedHandler, relationUsecase domain.RelationUsecase) domain.FeedUsecase {
	return &feedUsecase{handlerMap: handlerMap, relationUsecase: relationUsecase}
}

func (uc *feedUsecase) CreateFeedEvent(c context.Context, feed domain.Feed) error {
	handler, ok := uc.handlerMap[feed.Type]
	if !ok {
		slog.Error("TODO")
		return errors.New("TODO")
		// or 走兜底路径 default handler
	}
	return handler.CreateFeedEvent(c, feed.Type, feed.Content)

}

func (uc *feedUsecase) GetFeedEventList(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
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
		followees, _, err := uc.relationUsecase.GetFollowee(ctx, userID, -1, -1)
		if err != nil {
			return err
		}
		userIDs := slice.Map(followees, func(user domain.User) int64 {
			return user.UserID
		})
		pushEvent, err = uc.repo.FindPush(ctx, userIDs, timestamp, limit)
		return err
	})
	g.Go(func() error {
		var err error
		pullEvent, err = uc.repo.FindPull(ctx, userID, timestamp, limit)
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
