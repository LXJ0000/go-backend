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

func NewTagRouter(env *bootstrap.Env, timeout time.Duration,
	orm orm.Database, cache cache.RedisCache,
	group *gin.RouterGroup) {

	repo := repository.NewTagRepository(orm)
	useCase := usecase.NewTagUsecase(repo, timeout)
	col := &controller.TagController{
		TagUsecase: useCase,
	}
	group.POST("/tag", col.CreateTag)
	group.POST("/tag/bind", col.CreateTagBiz)
	group.GET("/tag", col.GetTagsByUserID)
	group.GET("/tag/biz", col.GetTagsByBiz)
}
