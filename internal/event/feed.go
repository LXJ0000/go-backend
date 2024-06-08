package event

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
)

type FeedEvent struct {
	Type     string             // 内部定义，分发给不同业务方
	MetaData domain.FeedContent // 业务方具体的数据
}

type FeedEventConsumer struct {
	client      sarama.Client
	feedUsecase domain.FeedUsecase
}

func NewFeedEventConsumer(client sarama.Client, feedUsecase domain.FeedUsecase) *FeedEventConsumer {
	return &FeedEventConsumer{client: client, feedUsecase: feedUsecase}
}

func (f *FeedEventConsumer) Consume(msg *sarama.ConsumerMessage, event FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return f.feedUsecase.CreateFeedEvent(ctx, domain.Feed{
		Type:    event.Type,
		Content: event.MetaData,
	})
}
