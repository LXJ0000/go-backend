package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewTagRouter(env *bootstrap.Env, tagUc domain.TagUsecase,
	group *gin.RouterGroup) {

	col := &controller.TagController{
		TagUsecase: tagUc,
	}
	group.POST("/tag", col.CreateTag)
	group.POST("/tag/bind", col.CreateTagBiz)
	group.GET("/tag", col.GetTagsByUserID)
	group.GET("/tag/biz", col.GetTagsByBiz)
}
