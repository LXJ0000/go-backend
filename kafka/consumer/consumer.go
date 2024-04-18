package main

import (
	"context"
	"github.com/IBM/sarama"
	"log/slog"
	"time"
)

// kafka consumer
var addr = []string{"localhost:9094"}

func main() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "group", config)
	if err != nil {
		slog.Error("consumer closed", "err", err)
		return
	}
	defer consumer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = consumer.Consume(ctx, []string{"web_log"}, Handler{}); err != nil {
		slog.Error("consumer error", "err", err)
		return
	}
}

type Handler struct {
}

func (c Handler) Setup(session sarama.ConsumerGroupSession) error {
	slog.Info("Consumer set up successfully")
	return nil
}

func (c Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	slog.Info("Consumer cleanup successfully")
	return nil
}

func (c Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("Consumer claim start")
	allMsg := claim.Messages()
	for msg := range allMsg {
		slog.Info("got", "msg", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
