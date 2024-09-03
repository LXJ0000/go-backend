package controller

import (
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type FeedController struct {
	domain.FeedUsecase
}

func (f *FeedController) Feed(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)

	feed, err := f.FeedUsecase.GetFeedEventList(c, userID, time.Now().Unix(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("internal server error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(feed))
}
