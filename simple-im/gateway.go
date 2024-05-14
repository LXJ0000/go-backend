package simpleim

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/pkg/kafka"

	// "golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
)

type GateWay struct {
	conn       sync.Map
	service    Service
	client     sarama.Client
	instanceID string
}

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, data)
}

func (g *GateWay) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", g.wsHandler)
	if err := g.subscribeMsg(); err != nil {
		return err
	}
	return http.ListenAndServe(addr, mux)
}

func (g *GateWay) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	uid := g.Uid(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("init websocket fail with %v", err)))
		return
	}
	c := &Conn{conn}
	g.conn.Store(uid, c)
	go func() {
		_, message, err := c.ReadMessage()
		if err != nil {
			slog.Warn("ReadMessage fail", "err", err)
			return
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Warn("json.Unmarshal fail", "err", err)
			return
		}
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := g.service.Receiver(ctx, uid, msg); err != nil {
				slog.Error("Receiver fail", "err", err)
				if err := c.Send(Message{
					Seq:     msg.Seq,
					Content: "FAILD",
					Type:    "RESULT",
				}); err != nil {
					slog.Error("Send fail msg fail", "err", err)
				}
			}
		}()
	}()
}

func (g *GateWay) subscribeMsg() error {
	group, err := sarama.NewConsumerGroupFromClient(g.instanceID, g.client)
	if err != nil {
		return err
	}
	go func() {
		if err := group.Consume(context.Background(), []string{eventName}, kafka.NewConsumerHandler[Event](g.consume)); err != nil {
			slog.Info("退出监听消息循环", "err", err)
		}
	}()
	return nil
}

func (g *GateWay) consume(msg *sarama.ConsumerMessage, event Event) error {
	conn, ok := g.conn.Load(event.Receiver)
	if !ok {
		slog.Warn("not user exists")
		return nil
	}
	c := conn.(*Conn)
	return c.Send(event.Msg)
}

func (g *GateWay) Uid(req *http.Request) int64 {
	uidStr := req.Header.Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

type Message struct {
	Seq     string
	Type    string
	Content string
	Cid     int64
}
