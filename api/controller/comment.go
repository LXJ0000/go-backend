package controller

import (
	"net/http"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentUsecase domain.CommentUsecase
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

	if err := col.CommentUsecase.Create(c.Request.Context(), &comment); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create comment fail", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"comment_detail": comment,
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
	var req domain.CommentListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	resp, err := col.CommentUsecase.FindTop(c.Request.Context(), req.Biz, lib.Str2Int64DefaultZero(req.BizID), lib.Str2Int64DefaultZero(req.MinID), req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("FindTop comment fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"comment_list": resp,
	}))
}
