package route

import (
	"time"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	pc := &controller.UserController{
		UserUsecase: usecase.NewProfileUsecase(ur, timeout),
	}
	user := group.Group("/user")
	user.GET("/profile", pc.Fetch)
}
