package doubao

import (
	"context"
	"io"
	"log/slog"

	"github.com/LXJ0000/go-backend/pkg/chat"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

const (
	defaultPrompt = "你是豆包，是由字节跳动开发的 AI 人工智能助手"
)

type DoubaoChat struct {
	client     *arkruntime.Client
	endpointId string
}

func NewDoubaoChat(client *arkruntime.Client, endpointId string) chat.Chat {
	c := &DoubaoChat{client: client, endpointId: endpointId}
	return c
}

func (d *DoubaoChat) NormalChat(ctx context.Context, prompt, request string) (response string, err error) {
	req := d.parseToken2Req(prompt, request)
	resp, err := d.client.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.Error("doubao: standard chat error", "error", err)
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil
}

func (d *DoubaoChat) StreamChat(ctx context.Context, prompt, request string) (response <-chan string, err error) {

	req := d.parseToken2Req(prompt, request)
	stream, err := d.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		slog.Error("doubao: stream chat error", "error", err)
		return nil, err
	}
	defer stream.Close()

	resp := make(chan string, 1)

	go func() {
		defer close(resp)
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				slog.Error("doubao: stream chat error", "error", err)
				return
			}

			if len(recv.Choices) > 0 {
				resp <- recv.Choices[0].Delta.Content
			}
		}
	}()

	return resp, nil

	// req := d.parseToken2Req(prompt, request)
	// ch, err := d.r.StreamChat(d.endpointId, req)
	// if err != nil {
	// 	errVal := &api.Error{}
	// 	if errors.As(err, &errVal) { // the returned error always type of *api.Error
	// 		slog.Error("meet maas error", "error", errVal)
	// 	}
	// 	return
	// }
	// resp := make(chan string, 1)
	// go func() {
	// 	for msg := range ch {
	// 		if msg.Error != nil {
	// 			slog.Error("it is possible that error occurs during response processing", "response", d.mustMarshalJson(msg.Error))
	// 			return
	// 		}
	// 		resp <- d.mustMarshalJson(msg)
	// 		if msg.Usage != nil {
	// 			resp <- d.mustMarshalJson(msg.Usage)
	// 		}
	// 	}
	// 	close(resp)
	// }()
	// return resp, nil
}

// func (d *DoubaoChat) mustMarshalJson(v interface{}) string {
// 	s, _ := json.Marshal(v)
// 	return string(s)
// }

func (d *DoubaoChat) parseToken2Req(prompt, token string) model.ChatRequest {
	req := model.ChatCompletionRequest{
		Model: d.endpointId,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: &prompt,
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: &token,
				},
			},
		},
	}
	return req
}
