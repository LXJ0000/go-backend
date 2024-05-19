package event

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
)

type FeedEvent struct {
	Type string
	Data domain.Content
}

type FeedEventConsumer struct {
	client      sarama.Client
	feedUsecase domain.FeedUsecase
}

func NewFeedEventConsumer(client sarama.Client, feedUsecase domain.FeedUsecase) *FeedEventConsumer {
	return &FeedEventConsumer{client: client, feedUsecase: feedUsecase}
}
