package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewRelationRouter(env *bootstrap.Env,
	relationUc domain.RelationUsecase, feedUc domain.FeedUsecase,
	group *gin.RouterGroup) {
	col := &controller.RelationController{
		RelationUsecase: relationUc,
		FeedUsecase:     feedUc,
	}

	group.POST("/relation/follow", col.Follow)
	group.POST("/relation/cancel_follow", col.CancelFollow)
	group.GET("/relation/list/follower", col.FollowerList)
	group.GET("/relation/list/followee", col.FolloweeList)
	group.GET("/relation/stat", col.Stat)
}
