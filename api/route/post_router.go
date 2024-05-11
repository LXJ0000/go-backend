package route

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/event"
	repository2 "github.com/LXJ0000/go-backend/internal/repository"
	usecase2 "github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"time"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, orm orm.Database, cache cache.RedisCache, group *gin.RouterGroup,
	producer event.Producer, saramaClient sarama.Client) {

	repoPost := repository2.NewPostRepository(orm, cache)
	repoInteraction := repository2.NewInteractionRepository(orm, cache)
	col := &controller.PostController{
		PostUsecase:        usecase2.NewPostUsecase(repoPost, timeout, producer),
		InteractionUseCase: usecase2.NewInteractionUsecase(repoInteraction, timeout),
	}
	consumer := event.NewBatchSyncReadEventConsumer(saramaClient, repoInteraction)
	if err := consumer.Start(); err != nil {
		slog.Error("OMG！消费者启动失败")
		log.Fatal(err)
	}
	group.POST("/post", col.CreateOrPublish)
	group.GET("/post", col.Info)
	group.GET("/post/publish", col.ReaderList)
	group.GET("/post/private", col.WriterList)
	group.POST("/post/like", col.Like)
	group.POST("/post/collect", col.Collect)
}
