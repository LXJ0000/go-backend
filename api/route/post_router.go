package route

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/event"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"github.com/LXJ0000/go-backend/repository"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, orm orm.Database, cache redis.Cache, group *gin.RouterGroup,
	producer event.Producer, saramaClient sarama.Client) {

	repoPost := repository.NewPostRepository(orm, cache)
	repoInteraction := repository.NewInteractionRepository(orm, cache)
	col := &controller.PostController{
		PostUsecase:        usecase.NewPostUsecase(repoPost, timeout, producer),
		InteractionUseCase: usecase.NewInteractionUsecase(repoInteraction, timeout),
	}
	consumer := event.NewSyncReadEventConsumer(saramaClient, repoInteraction)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
	group.POST("/post", col.CreateOrPublish)
	group.GET("/post", col.Info)
	group.GET("/post/publish", col.ReaderList)
	group.GET("/post/private", col.WriterList)
	group.POST("/post/like", col.Like)
	group.POST("/post/collect", col.Collect)
}
