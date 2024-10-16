package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewPostRouter(PostUsecase domain.PostUsecase,
	InteractionUseCase domain.InteractionUseCase,
	FeedUsecase domain.FeedUsecase,
	UserUsecase domain.UserUsecase,
	group *gin.RouterGroup) {

	c := &controller.PostController{
		PostUsecase:        PostUsecase,
		InteractionUseCase: InteractionUseCase,
		FeedUsecase:        FeedUsecase,
		UserUsecase:        UserUsecase,
	}

	group.POST("/post", c.CreateOrPublish)
	group.GET("/post", c.Info)
	group.GET("/post/reader", c.ReaderList)
	group.GET("/post/writer", c.WriterList)
	group.POST("/post/like", c.Like)
	group.POST("/post/collect", c.Collect)
	group.GET("/post/rank", c.Rank)
}
