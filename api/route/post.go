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
	FileUsecase domain.FileUsecase,
	apiCache func(timeout time.Duration) gin.HandlerFunc,
	group *gin.RouterGroup) {

	c := &controller.PostController{
		PostUsecase:        PostUsecase,
		InteractionUseCase: InteractionUseCase,
		FeedUsecase:        FeedUsecase,
		UserUsecase:        UserUsecase,
		CommentUsecase:     CommentUsecase,
		FileUsecase:        FileUsecase,
	}

	group.GET("/post", c.Info) // 这里加接口缓存会导致阅读计数无法生效
	group.GET("/post/reader", apiCache(time.Second*5), c.ReaderList)
	group.GET("/post/writer", apiCache(time.Second*5), c.WriterList)
	group.GET("/post/rank", c.Rank)

	group.POST("/post", c.CreateOrPublish)
	group.POST("/post.delete", c.PostDelete)
	group.POST("/post.search", c.Search)


	group.POST("/post/like", c.Like)       // 转移到 interaction 里
	group.POST("/post/collect", c.Collect) // 转移到 interaction 里
}
