package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type feedUsecase struct {
	handler map[string]domain.Handler // map[type]handler
}

func NewFeedUsecase(handler map[string]domain.Handler) domain.FeedUsecase {
	return &feedUsecase{handler: handler}
}

func (uc *feedUsecase) CreateFeedEvent(c context.Context, feed domain.Feed) error {
	handler, ok := uc.handler[feed.Type]
	if !ok {
		slog.Error("TODO")
		return errors.New("TODO")
		// or 走兜底路径 default handler
	}
	return handler.CreateFeedEvent(c, feed.Content)

}
func (uc *feedUsecase) GetFeedEventList(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil
}
