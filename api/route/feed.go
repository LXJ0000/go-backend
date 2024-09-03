package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewFeedRouter(FeedUsecase domain.FeedUsecase, group *gin.RouterGroup) {
	c := controller.FeedController{FeedUsecase: FeedUsecase}

	group.POST("/feed", c.Feed)
}
