package event

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/kafka"
	"golang.org/x/net/context"
	"time"
)

type ReadEvent struct {
	UserID int64
	PostID int64
}

// Producer 生产者
type Producer interface {
	ProduceReadEvent(c context.Context, event ReadEvent) error
}

// SyncReadEventProducer 生产者
type SyncReadEventProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(producer sarama.SyncProducer) *SyncReadEventProducer {
	return &SyncReadEventProducer{producer: producer}
}

func (s SyncReadEventProducer) ProduceReadEvent(c context.Context, event ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: "post_read",
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = s.producer.SendMessage(&msg)
	return err
}

// SyncReadEventConsumer 消费者
type SyncReadEventConsumer struct {
	client sarama.Client
	repo   domain.InteractionRepository
}

func NewSyncReadEventConsumer(client sarama.Client, repo domain.InteractionRepository) *SyncReadEventConsumer {
	return &SyncReadEventConsumer{client: client, repo: repo}
}

func (c *SyncReadEventConsumer) Start() error {
	config, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		//	TODO slog
		return err
	}
	go func() {
		if err := config.Consume(context.Background(), []string{"post_read"}, kafka.NewConsumerHandler[ReadEvent](c.Consumer)); err != nil {
			// TODO slog
		}
	}()
	return nil
}

func (c *SyncReadEventConsumer) Consumer(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrReadCount(ctx, domain.BizPost, event.PostID)
}
