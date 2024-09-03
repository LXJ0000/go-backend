package route

import (
	"time"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	feedUsecase "github.com/LXJ0000/go-backend/internal/usecase/feed"
	feedUsecaseHandler "github.com/LXJ0000/go-backend/internal/usecase/feed/handler"

	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
)

func NewFeedRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, redisCache cache.RedisCache, group *gin.RouterGroup) {
	repo := repository.NewFeedRepository(db)
	userRepo := repository.NewUserRepository(db, redisCache)
	relationRepo := repository.NewRelationRepository(db)
	relationUc := usecase.NewRelationUsecase(relationRepo, userRepo, timeout)

	feedLikeEventHandler := feedUsecaseHandler.NewFeedLikeHandler(repo)
	feedPostEventHandler := feedUsecaseHandler.NewFeedPostHandler(repo, relationUc)
	feedFollowEventHandler := feedUsecaseHandler.NewFeedFollowHandler(repo)

	handlerMap := map[string]domain.FeedHandler{
		domain.FeedLikeEvent:   feedLikeEventHandler,
		domain.FeedPostEvent:   feedPostEventHandler,
		domain.FeedFollowEvent: feedFollowEventHandler,
	}
	uc := feedUsecase.NewFeedUsecase(handlerMap, relationUc, repo)
	c := controller.FeedController{FeedUsecase: uc}

	group.POST("/feed", c.Feed)
}
