package route

import (
	"log"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration,
	orm orm.Database, cache cache.RedisCache, localCache cache.LocalCache,
	group *gin.RouterGroup,
	producer event.Producer, saramaClient sarama.Client) {

	repoPost := repository.NewPostRepository(orm, cache)
	repoInteraction := repository.NewInteractionRepository(orm, cache)
	repoPostRank := repository.NewPostRankRepository(localCache, cache)
	postRankUsecase := usecase.NewPostRankUsecase(repoInteraction, repoPost, repoPostRank, timeout)
	col := &controller.PostController{
		PostUsecase:        usecase.NewPostUsecase(repoPost, timeout, producer, postRankUsecase),
		InteractionUseCase: usecase.NewInteractionUsecase(repoInteraction, timeout),
	}
	consumer := event.NewBatchSyncReadEventConsumer(saramaClient, repoInteraction)
	if err := consumer.Start(); err != nil {
		slog.Error("OMG！消费者启动失败")
		log.Fatal(err)
	}
	group.POST("/post", col.CreateOrPublish)
	group.GET("/post", col.Info)
	group.GET("/post/reader", col.ReaderList)
	group.GET("/post/writer", col.WriterList)
	group.POST("/post/like", col.Like)
	group.POST("/post/collect", col.Collect)
	group.GET("/post/rank", col.Rank)
}
