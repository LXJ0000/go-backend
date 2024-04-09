package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/repository"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, group *gin.RouterGroup, orm orm.Database) {
	repo := repository.NewPostRepository(orm)
	col := &controller.PostController{
		Usecase: usecase.NewPostUsecase(repo, timeout),
	}
	group.POST("/post", col.Create)
}
