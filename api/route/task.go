package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewTaskRouter(env *bootstrap.Env, taskUc domain.TaskUsecase, group *gin.RouterGroup) {
	tc := &controller.TaskController{
		TaskUsecase: taskUc,
	}
	group.POST("/task", tc.Create)
	group.DELETE("/task", tc.Delete)
}
