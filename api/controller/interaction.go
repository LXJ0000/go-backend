package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type IntrController struct {
	domain.InteractionUseCase
	domain.FeedUsecase
	domain.PostUsecase
}

func (col *IntrController) Like(c *gin.Context) { // TODO 抽象成资源操作而不是针对帖子的操作
	req := struct {
		IsLike bool   `json:"is_like" form:"is_like"`
		PostID int64  `json:"post_id,string" form:"post_id"`
		BizID  int64  `json:"biz_id,string" form:"biz_id"`
		Biz    string `json:"biz" form:"biz"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	bizID := req.BizID
	isLike := req.IsLike
	biz := req.Biz
	userID := c.MustGet(domain.XUserID).(int64)
	var err error
	if isLike {
		err = col.InteractionUseCase.Like(c, biz, bizID, userID)
	} else {
		err = col.InteractionUseCase.CancelLike(c, biz, bizID, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	go func() {
		// TODO 发消息
		var authorID int64
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		switch biz {
		case domain.BizPost:
			post, err := col.PostUsecase.Info(ctx, bizID)
			if err != nil {
				slog.Error("FeedUsecase CreateFeedEvent Error Because PostUsecase Info Error", "error", err.Error())
				return
			}
			authorID = post.AuthorID
		}
		feed := domain.Feed{
			UserID: userID,
			Type:   domain.FeedLikeEvent,
			Content: domain.FeedContent{
				"biz_id": fmt.Sprintf("%d", bizID),
				"biz":    domain.BizPost,
				"liker":  fmt.Sprintf("%d", userID),   // 点赞者
				"liked":  fmt.Sprintf("%d", authorID), // biz's author id 被点赞者
			}, // liker liked biz bizID
		}
		if err := col.FeedUsecase.CreateFeedEvent(context.Background(), feed); err != nil {
			slog.Warn("FeedUsecase CreateFeedEvent Error", "error", err.Error())
		}
	}()
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *IntrController) Collect(c *gin.Context) {
	req := struct {
		IsCollect bool   `json:"is_collect" form:"is_collect"`
		CollectID int64  `json:"collect_id,string" form:"collect_id"`
		BizID     int64  `json:"biz_id,string" form:"biz_id"`
		Biz       string `json:"biz" form:"biz"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	bizID := req.BizID
	biz := req.Biz
	isCollect := req.IsCollect
	collectID := req.CollectID
	userID := c.MustGet(domain.XUserID).(int64)
	var err error
	if isCollect {
		err = col.InteractionUseCase.Collect(c, biz, bizID, userID, collectID)
	} else {
		err = col.InteractionUseCase.CancelCollect(c, biz, bizID, userID, collectID)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}
