package doubao

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/LXJ0000/go-backend/pkg/chat"
	api "github.com/volcengine/volc-sdk-golang/service/maas/models/api/v2"
	client "github.com/volcengine/volc-sdk-golang/service/maas/v2"
)

const (
	defaultPrompt = "你是豆包，是由字节跳动开发的 AI 人工智能助手"
)

type DoubaoChat struct {
	r          *client.MaaS
	endpointId string
}

func NewDoubaoChat(r *client.MaaS, endpointId string) chat.Chat {
	c := &DoubaoChat{r: r, endpointId: endpointId}
	return c
}

func (d *DoubaoChat) NormalChat(prompt, request string) (response string, err error) {
	req := d.parseToken2Req(prompt, request)
	got, status, err := d.r.Chat(d.endpointId, req)
	if err != nil {
		errVal := &api.Error{}
		if errors.As(err, &errVal) { // the returned error always type of *api.Error
			slog.Error("meet maas error", "error", errVal, "status", status)
			return
		}
		return
	}
	return d.mustMarshalJson(got), nil
}

func (d *DoubaoChat) StreamChat(prompt, request string) (response <-chan string, err error) {
	req := d.parseToken2Req(prompt, request)
	ch, err := d.r.StreamChat(d.endpointId, req)
	if err != nil {
		errVal := &api.Error{}
		if errors.As(err, &errVal) { // the returned error always type of *api.Error
			slog.Error("meet maas error", "error", errVal)
		}
		return
	}
	resp := make(chan string, 1)
	go func() {
		for msg := range ch {
			if msg.Error != nil {
				slog.Error("it is possible that error occurs during response processing", "response", d.mustMarshalJson(msg.Error))
				return
			}
			resp <- d.mustMarshalJson(msg)
			if msg.Usage != nil {
				resp <- d.mustMarshalJson(msg.Usage)
			}
		}
		close(resp)
	}()
	return resp, nil
}

func (d *DoubaoChat) mustMarshalJson(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func (d *DoubaoChat) parseToken2Req(prompt, token string) *api.ChatReq {
	if prompt == "" {
		prompt = defaultPrompt
	}
	return &api.ChatReq{
		Messages: []*api.Message{
			{
				Role:    api.ChatRoleSystem,
				Content: prompt,
			},
			{
				Role:    api.ChatRoleUser,
				Content: token,
			},
		},
	}
}
