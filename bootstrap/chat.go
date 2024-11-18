package bootstrap

import (
	"os"

	"github.com/LXJ0000/go-backend/pkg/chat"
	"github.com/LXJ0000/go-backend/pkg/chat/doubao"
	client "github.com/volcengine/volc-sdk-golang/service/maas/v2"
)

func NewDoubaoChat() chat.Chat {
	r := client.NewInstance("maas-api.ml-platform-cn-beijing.volces.com", "cn-beijing")

	r.SetAccessKey(os.Getenv("VOLC_ACCESSKEY"))
	r.SetSecretKey(os.Getenv("VOLC_SECRETKEY"))

	endpointId := os.Getenv("VOLC_ENDPOINTID")

	c := doubao.NewDoubaoChat(r, endpointId)

	return c
}
