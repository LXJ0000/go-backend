package bootstrap

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/event"
	"log"
)

func NewProducer(env *Env) event.Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{env.KafkaAddr}, config)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	return client
}
