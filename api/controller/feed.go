package controller

import (
	"net/http"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type FeedController struct {
	domain.FeedUsecase
}

func (f *FeedController) Feed(c *gin.Context) {
	var req domain.GetFeedEventListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)

	feed, err := f.FeedUsecase.GetFeedEventList(c, userID, req.Last, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("internal server error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(feed))
}
