package event

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/domain"
)

type FeedEvent struct {
	Type     string         // 内部定义，分发给不同业务方
	MetaData domain.Content // 业务方具体的数据
}

type FeedEventConsumer struct {
	client      sarama.Client
	feedUsecase domain.FeedUsecase
	feedService domain.
}

func NewFeedEventConsumer(client sarama.Client, feedUsecase domain.FeedUsecase) *FeedEventConsumer {
	return &FeedEventConsumer{client: client, feedUsecase: feedUsecase}
}
