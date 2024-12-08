package bootstrap

import (
	"os"

	"github.com/LXJ0000/go-backend/pkg/chat"
	"github.com/LXJ0000/go-backend/pkg/chat/doubao"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

func NewDoubaoChat() chat.Chat {
	client := arkruntime.NewClientWithApiKey(
		os.Getenv("ARK_API_KEY"),
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
		arkruntime.WithRegion("cn-beijing"),
	)

	endpointId := os.Getenv("VOLC_ENDPOINTID")

	c := doubao.NewDoubaoChat(client, endpointId)

	return c
}
