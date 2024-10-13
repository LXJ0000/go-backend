package controller

import (
	"net/http"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	"github.com/gin-gonic/gin"
)

type RelationController struct {
	RelationUsecase domain.RelationUsecase
}

func (col *RelationController) Follow(c *gin.Context) {
	// followee, err := lib.Str2Int64(c.Request.FormValue("followee"))
	req := struct {
		Followee int64 `json:"followee" form:"followee"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)
	if req.Followee == userID {
		c.JSON(http.StatusForbidden, domain.ErrorResp("You can't follow yourself", nil))
		return
	}
	if err := col.RelationUsecase.Follow(c, userID, req.Followee); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Follow Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *RelationController) CancelFollow(c *gin.Context) {
	// followee, err := lib.Str2Int64(c.Request.FormValue("followee"))
	req := struct {
		Followee int64 `json:"followee" form:"followee"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)
	if userID == req.Followee {
		c.JSON(http.StatusForbidden, domain.ErrorResp("You can't cancel follow yourself", nil))
		return
	}
	if err := col.RelationUsecase.CancelFollow(c, userID, req.Followee); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("CancelFollow Failed", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *RelationController) Stat(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	stat, err := col.RelationUsecase.Stat(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get Stat Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(stat))
}

func (col *RelationController) FollowerList(c *gin.Context) {
	page, _ := lib.Str2Int(c.Query("page"))
	size, _ := lib.Str2Int(c.Query("size"))
	userID := c.MustGet(domain.XUserID).(int64)
	resp, cnt, err := col.RelationUsecase.GetFollower(c, userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get FollowerList Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     cnt,
		"user_list": resp,
	}))
}

func (col *RelationController) FolloweeList(c *gin.Context) {
	page, _ := lib.Str2Int(c.Query("page"))
	size, _ := lib.Str2Int(c.Query("size"))
	userID := c.MustGet(domain.XUserID).(int64)
	resp, cnt, err := col.RelationUsecase.GetFollowee(c, userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get FolloweeList Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     cnt,
		"user_list": resp,
	}))
}
