package bootstrap

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/event"
)

func NewProducer(env *Env) event.Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{env.KafkaAddr}, config)
	if err != nil {
		slog.Warn("kafka client init failed", "error", err.Error())
		fmt.Println("kafka client init failed", err.Error())
		return nil
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	return event.NewSyncProducer(producer)
}

func NewSaramaClient(env *Env) sarama.Client {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{env.KafkaAddr}, config)
	if err != nil {
		slog.Warn("kafka client init failed", "error", err.Error())
		fmt.Println("kafka client init failed", err.Error())
		return nil
	}
	return client
}
