package route

import (
	"time"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewPostRouter(PostUsecase domain.PostUsecase,
	InteractionUseCase domain.InteractionUseCase,
	FeedUsecase domain.FeedUsecase,
	UserUsecase domain.UserUsecase,
	CommentUsecase domain.CommentUsecase,
	apiCache func(timeout time.Duration) gin.HandlerFunc,
	group *gin.RouterGroup) {

	c := &controller.PostController{
		PostUsecase:        PostUsecase,
		InteractionUseCase: InteractionUseCase,
		FeedUsecase:        FeedUsecase,
		UserUsecase:        UserUsecase,
		CommentUsecase:     CommentUsecase,
	}

	group.POST("/post", c.CreateOrPublish)
	group.GET("/post", c.Info)
	group.GET("/post/reader", apiCache(time.Second), c.ReaderList)
	group.GET("/post/writer", c.WriterList)
	group.POST("/post/like", c.Like)
	group.POST("/post/collect", c.Collect)
	group.GET("/post/rank", c.Rank)
}
