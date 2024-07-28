package controller

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TagController struct {
	TagUsecase domain.TagUsecase
}

func (col *TagController) CreateTag(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	tagName := c.Request.FormValue("name")
	if tagName == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params: name must not be empty", nil))
		return
	}
	//now := time.Now().UnixMicro()
	tag := domain.Tag{
		TagID:   snowflakeutil.GenID(),
		UserID:  userID,
		TagName: tagName,
	}
	//tag.CreatedAt = now
	//tag.UpdatedAt = now
	if err := col.TagUsecase.CreateTag(c, tag); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create Tag Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(tag))
}

func (col *TagController) CreateTagBiz(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	var req domain.CreateTagBizRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	if err := col.TagUsecase.CreateTagBiz(c, userID, req.Biz, req.BizID, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create TagBiz Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *TagController) GetTagsByBiz(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	var req domain.GetTagsByBizRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	tags, err := col.TagUsecase.GetTagsByBiz(c, userID, req.Biz, req.BizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get Tags By Biz Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(tags))
}

func (col *TagController) GetTagsByUserID(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	tags, err := col.TagUsecase.GetTagsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get Tags By UserID Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(tags))
}
