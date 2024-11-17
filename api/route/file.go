package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewFileRouter(env *bootstrap.Env, fileUc domain.FileUsecase, group *gin.RouterGroup) {
	col := &controller.FileController{FileUsecase: fileUc}
	group.POST("/file/upload", col.Upload)
	group.POST("/file/uploads", col.Uploads)
	group.GET("/file/list", col.FileList)
}
