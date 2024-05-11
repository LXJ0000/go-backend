package main

import (
	"encoding/json"
	"github.com/LXJ0000/go-backend/internal/event"
	"log/slog"

	"github.com/IBM/sarama"
)

// 基于sarama第三方库开发的kafka client
var addr = []string{"localhost:9094"}

func main() {
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	//config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	// 连接kafka
	client, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		slog.Error("producer closed", "err", err)
		return
	}
	defer client.Close()

	// 构造一个消息
	//msg := &sarama.ProducerMessage{}
	//msg.Topic = "post_read"
	//msg.Value = sarama.StringEncoder("this is a test log")
	e := event.ReadEvent{
		UserID: 0,
		PostID: 169724846903660544,
	}
	data, _ := json.Marshal(e)
	msg := sarama.ProducerMessage{
		Topic: "post_read",
		Value: sarama.ByteEncoder(data),
	}
	slog.Info("消息组装完毕... 准备发送噜~")
	for range 100 {
		slog.Info("", e)
		_, _, err = client.SendMessage(&msg)
	}

	// 发送消息
	//pid, offset, err := client.SendMessage(msg)
	//if err != nil {
	//	fmt.Println("send msg failed, err:", err)
	//	return
	//}
	//fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
