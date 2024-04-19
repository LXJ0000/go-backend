package main

import (
	"context"
	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"time"
)

// kafka consumer
var addr = []string{"localhost:9094"}

const batchSize = 10

func main() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "group", config)
	if err != nil {
		slog.Error("consumer closed", "err", err)
		return
	}
	defer consumer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
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

func (c Handler) ConsumeSyncClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("Consumer claim start")
	allMsg := claim.Messages()
	for msg := range allMsg {
		slog.Info("got", "msg", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

// ConsumeClaim ASync
func (c Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("Consumer claim start")
	allMsg := claim.Messages()
	//for msg := range allMsg {
	//	go func() {
	//		// 生产 远大于 消费 , 很容易造成大量goroutine
	//		slog.Info("got", "msg", string(msg.Value))
	//		session.MarkMessage(msg, "")
	//	}()
	//}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		eg := errgroup.Group{}
		done := false
		var last *sarama.ConsumerMessage
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-allMsg:
				if !ok {
					cancel()
					return nil // 消费者被关闭
				}
				last = msg
				eg.Go(func() error {
					slog.Info("got", "msg", string(msg.Value))
					time.Sleep(time.Second)
					return nil
				})
			}
		}
		if err := eg.Wait(); err != nil {
			//重试 or 记录日志 人工处理
			slog.Warn(err.Error())
			continue
		}
		if last != nil {
			session.MarkMessage(last, "")
		}
		cancel()
	}
}
