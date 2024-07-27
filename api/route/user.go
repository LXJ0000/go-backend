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

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, redisCache cache.RedisCache,
	group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, redisCache)
	col := &controller.UserController{
		UserUsecase: usecase.NewUserUsecase(ur, timeout),
		Env:         env,
	}
	group.POST("/logout", col.Logout)
	group.GET("/user/profile", col.Fetch)
	group.POST("/user/edit", col.Update)
}
