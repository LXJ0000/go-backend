package controller

import (
	"database/sql"
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
	var req domain.CommentCreateRequset
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}

	var comment domain.Comment

	comment.CommentID = snowflakeutil.GenID()
	comment.UserID = c.MustGet(domain.UserCtxID).(int64)
	comment.Biz = req.Biz
	comment.BizID = req.BizID
	comment.Content = req.Content
	comment.RootID = sql.NullInt64{Int64: req.RootID, Valid: req.RootID != 0}
	comment.ParentID = sql.NullInt64{Int64: req.ParentID, Valid: req.ParentID != 0}

	if err := col.CommentUsecase.Create(c.Request.Context(), comment); err != nil {
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
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	if err := col.CommentUsecase.Delete(c, commentID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Delete comment fail", err))
	}
	c.JSON(http.StatusOK, domain.SuccessResp("Delete comment success"))
}

func (col *CommentController) FindTop(c *gin.Context) {
	var req domain.CommentListRequset
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	resp, err := col.CommentUsecase.FindTop(c.Request.Context(), req.Biz, req.BizID, req.MinID, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("FindTop comment fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"comment_list": resp,
	}))
}
