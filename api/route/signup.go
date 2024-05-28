package route

import (
	"time"

	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, redisCache cache.RedisCache, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, redisCache)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
}
