package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

const retryNumber = 5

type Consumer interface {
	Start() error
}

// ConsumerHandler ===========================================================================================================
type ConsumerHandler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, data T) error
}

func (h ConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	allMsg := claim.Messages()
	for msg := range allMsg {
		var data T
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			slog.Error("Unmarshal Message Fail",
				"err", err.Error(), "topic", msg.Topic, "offset", msg.Offset, "partition", msg.Partition)
			continue
		}
		done := true
		if err := h.fn(msg, data); err != nil {
			// TODO 重试
			done = false
			for range retryNumber {
				if err := h.fn(msg, data); err == nil {
					done = true
					break
				}
			}
			slog.Error("Consume Message Fail",
				"err", err.Error(), "topic", msg.Topic, "offset", msg.Offset, "partition", msg.Partition)
		}
		if done {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

func NewConsumerHandler[T any](fn func(msg *sarama.ConsumerMessage, data T) error) *ConsumerHandler[T] {
	return &ConsumerHandler[T]{
		fn: fn,
	}
}
