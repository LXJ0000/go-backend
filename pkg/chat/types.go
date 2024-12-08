package chat

import "context"

type Chat interface {
	NormalChat(ctx context.Context, prompt, request string) (response string, err error)
	StreamChat(ctx context.Context, prompt, request string) (response <-chan string, err error)
}
