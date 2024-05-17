package simpleim

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/IBM/sarama"
)

type Service struct {
	producer sarama.SyncProducer
}

func (s *Service) Receiver(ctx context.Context, sender int64, msg Message) error {
	members := s.findMembers()
	for _, member := range members {
		if member == sender {
			continue
		}
		msgJSON, err := json.Marshal(Event{Msg: msg, Receiver: member})
		if err != nil {
			slog.Warn(fmt.Sprintf("Send message to %d fail", member), "err", err)
			continue
		}
		if _, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: eventName,
			Key:   sarama.ByteEncoder(strconv.FormatInt(member, 10)),
			Value: sarama.ByteEncoder(msgJSON),
		}); err != nil {
			slog.Warn(fmt.Sprintf("Send message to %d fail", member), "err", err)
			continue
		}
		slog.Info("Produce message to ", "member", member, "msg", msg.Content)
	}
	return nil
}

func (s *Service) findMembers() []int64 {
	// 固定返回 1，2，3，4
	return []int64{1, 2, 3, 4}
}
