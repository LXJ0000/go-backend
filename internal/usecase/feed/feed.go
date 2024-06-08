package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type feedUsecase struct {
	handlerMap map[string]domain.FeedHandler // map[type]handler
}

func NewFeedUsecase(handlerMap map[string]domain.FeedHandler) domain.FeedUsecase {
	return &feedUsecase{handlerMap: handlerMap}
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
	return nil, nil
}
