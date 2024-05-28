package route

import (
	"time"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
)

func NewRelationRouter(env *bootstrap.Env, timeout time.Duration,
	orm orm.Database, cache cache.RedisCache,
	group *gin.RouterGroup,
) {
	repo := repository.NewRelationRepository(orm)
	userRepo := repository.NewUserRepository(orm, cache)
	col := &controller.RelationController{
		RelationUsecase: usecase.NewRelationUsecase(repo, userRepo, timeout),
	}

	group.POST("/relation/follow", col.Follow)
	group.POST("/relation/cancel_follow", col.CancelFollow)
	group.GET("/relation/list/follower", col.FollowerList)
	group.GET("/relation/list/followee", col.FolloweeList)
	group.GET("/relation/stat", col.Stat)
}
