package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewIntrRouter(
	InteractionUseCase domain.InteractionUseCase,
	FeedUsecase domain.FeedUsecase,
	PostUsecase domain.PostUsecase,
	group *gin.RouterGroup) {

	c := &controller.IntrController{
		InteractionUseCase: InteractionUseCase,
		FeedUsecase:        FeedUsecase,
		PostUsecase:        PostUsecase,
	}

	group.POST("/intr/like", c.Like)
	group.POST("/intr/collect", c.Collect)
}
