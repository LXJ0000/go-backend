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

func NewCommentRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, group *gin.RouterGroup) {
	repo := repository.NewCommentRepository(db)
	col := &controller.CommentController{
		CommentUsecase: usecase.NewCommentUsecase(repo, timeout),
	}
	group.POST("/comment", col.Create)
	group.DELETE("/comment", col.Delete)
	group.GET("/comment", col.FindTop)
}
