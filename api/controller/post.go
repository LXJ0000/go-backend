package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type PostController struct {
	domain.PostUsecase
	domain.InteractionUseCase
	domain.FeedUsecase
}

func (col *PostController) CreateOrPublish(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	//now := time.Now().UnixMicro()
	post.AuthorID = userID
	post.PostID = snowflake.GenID()
	//post.CreatedAt = now
	//post.UpdatedAt = now
	if err := col.PostUsecase.Create(c, post); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to create post", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(post))
}

func (col *PostController) ReaderList(c *gin.Context) {
	//读者查看列表 只能查看已发布的文章
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}

	filter := &domain.Post{
		Status: domain.PostStatusPublish,
	}
	if req.AuthorID != 0 {
		filter.AuthorID = req.AuthorID
	}
	posts, count, err := col.PostUsecase.List(c, filter, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to list posts", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": posts,
	}))
}

func (col *PostController) WriterList(c *gin.Context) {
	//创作者查看列表 可以查看所有自己的帖子
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)
	filter := &domain.Post{AuthorID: userID}
	posts, count, err := col.PostUsecase.List(c, filter, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to list posts", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": posts,
	}))
}

func (col *PostController) Info(c *gin.Context) {
	postID, err := lib.Str2Int64(c.Query("post_id"))
	userID := c.MustGet(domain.XUserID).(int64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	var post domain.Post
	var interaction domain.Interaction
	var userInteractionInfo domain.UserInteractionStat
	eg := errgroup.Group{}
	eg.Go(func() error {
		post, err = col.PostUsecase.Info(c, postID)
		return err
	})
	eg.Go(func() error {
		interaction, userInteractionInfo, err = col.InteractionUseCase.Info(c, domain.BizPost, postID, userID)
		if err != nil {
			slog.Warn("InteractionUseCase Info Error", "error", err.Error())
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"post_detail":        post,
		"interaction_detail": interaction,
		"stat":               userInteractionInfo,
	}))

}

func (col *PostController) Like(c *gin.Context) {
	isLikeRaw := c.Request.FormValue("is_like")
	postIDRaw := c.Request.FormValue("post_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	isLike, err := strconv.ParseBool(isLikeRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)

	if isLike {
		err = col.InteractionUseCase.Like(c, domain.BizPost, postID, userID)
	} else {
		err = col.InteractionUseCase.CancelLike(c, domain.BizPost, postID, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	go func() {
		// TODO 发消息
		feed := domain.Feed{
			UserID: userID,
			Type:   domain.FeedLikeEvent,
			Content: domain.FeedContent{
				"biz_id": fmt.Sprintf("%d", postID),
				"biz":    domain.BizPost,
				"liker":  fmt.Sprintf("%d", userID),
				"liked":  fmt.Sprintf("%d", postID),
			}, // liker liked biz bizID
		}
		if err := col.FeedUsecase.CreateFeedEvent(context.Background(), feed); err != nil {

		}
	}()
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *PostController) Collect(c *gin.Context) {
	isCollectRaw := c.Request.FormValue("is_collect")
	postIDRaw := c.Request.FormValue("post_id")
	collectionIDRaw := c.Request.FormValue("collection_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	isCollect, err := strconv.ParseBool(isCollectRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	collectID, err := strconv.ParseInt(collectionIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)

	if isCollect {
		err = col.InteractionUseCase.Collect(c, domain.BizPost, postID, userID, collectID)
	} else {
		err = col.InteractionUseCase.CancelCollect(c, domain.BizPost, postID, userID, collectID)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *PostController) Rank(c *gin.Context) {
	posts, err := col.PostUsecase.TopN(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("GetTopN Fail", err))
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(posts),
		"posts": posts,
	})
}
