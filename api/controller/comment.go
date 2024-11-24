package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	domain.CommentUsecase
	domain.UserUsecase
	domain.InteractionUseCase
}

func (col *CommentController) Create(c *gin.Context) {
	var req domain.CommentCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}

	comment := domain.Comment{
		CommentID:       snowflakeutil.GenID(),
		UserID:          c.MustGet(domain.XUserID).(int64),
		Biz:             req.Biz,
		BizID:           lib.Str2Int64DefaultZero(req.BizID),
		Content:         req.Content,
		ParentCommentID: lib.Str2Int64DefaultZero(req.ParentID),
		ToUserID:        lib.Str2Int64DefaultZero(req.ToUserID),
	}
	var toUser domain.Profile
	if err := col.CommentUsecase.Create(c.Request.Context(), &comment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create comment fail", err))
		return
	}
	if comment.ToUserID != 0 {
		var err error
		toUser, err = col.UserUsecase.GetProfileByID(context.Background(), comment.ToUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get user profile fail", err))
			return
		}
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"comment_detail": comment,
		"to_user":        toUser,
	}))
}

func (col *CommentController) Delete(c *gin.Context) {
	commentIDRaw := c.Query("comment_id")
	commentID, err := lib.Str2Int64(commentIDRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	if err := col.CommentUsecase.Delete(c, commentID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Delete comment fail", err))
	}
	c.JSON(http.StatusOK, domain.SuccessResp("Delete comment success"))
}

func (col *CommentController) FindTop(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	var req domain.CommentListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	// 获取评论列表
	comments, count, err := col.CommentUsecase.Find(c.Request.Context(), req.Biz, lib.Str2Int64DefaultZero(req.BizID), lib.Str2Int64DefaultZero(req.ParentID), lib.Str2Int64DefaultZero(req.MinID), req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("FindTop comment fail", err))
		return
	}
	// 获取用户信息
	userIDs := make([]int64, 0, len(comments)*2)
	for _, v := range comments {
		userIDs = append(userIDs, v.UserID)
		if v.ToUserID != 0 {
			userIDs = append(userIDs, v.ToUserID)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	userID2Info := make(map[int64]domain.Profile, len(userIDs))
	users, err := col.UserUsecase.BatchGetProfileByID(ctx, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Find user fail", err))
		return
	}
	for _, u := range users {
		userID2Info[u.UserID] = u
	}
	commentInfos := make([]domain.CommentInfo, 0, len(comments))
	for _, v := range comments {
		interaction, userInteractionInfo, err := col.InteractionUseCase.Info(c, domain.BizComment, v.CommentID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResp("InteractionUseCase Info Error", err))
			slog.Warn("InteractionUseCase Info Error", "error", err.Error())
			return
		}
		item := domain.CommentInfo{
			Comment:     v,
			UserProfile: userID2Info[v.UserID],
			Liked:       userInteractionInfo.Liked,
			LikeCount:   interaction.LikeCnt,
		}
		if v.ToUserID != 0 {
			item.ToUserProfile = userID2Info[v.ToUserID]
		}
		commentInfos = append(commentInfos, item)
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"comment_list": commentInfos,
		"count":        count,
	}))
}
