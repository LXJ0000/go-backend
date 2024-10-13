package route

import (
	"time"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *bootstrap.Env,
	timeout time.Duration, db orm.Database, redisCache cache.RedisCache,
	userUc domain.UserUsecase,
	relationUc domain.RelationUsecase,
	postUc domain.PostUsecase,
	group *gin.RouterGroup) {
	col := &controller.UserController{
		UserUsecase:     userUc,
		Env:             env,
		PostUsecase:     postUc,
		RelationUsecase: relationUc,
	}
	group.POST("/logout", col.Logout)
	group.GET("/user/profile", col.Fetch)
	group.POST("/user/edit", col.Update)
	group.GET("/user", col.Profile)
}
