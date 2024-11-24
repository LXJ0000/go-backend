package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewCommentRouter(env *bootstrap.Env, commentUc domain.CommentUsecase,
	userUc domain.UserUsecase,
	group *gin.RouterGroup) {
	col := &controller.CommentController{
		CommentUsecase: commentUc,
		UserUsecase:    userUc,
	}
	group.POST("/comment", col.Create)
	group.DELETE("/comment", col.Delete)
	group.POST("/comment.list", col.FindTop)
}
