package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LXJ0000/go-backend/internal/domain"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type PostController struct {
	PostUsecase        domain.PostUsecase
	InteractionUseCase domain.InteractionUseCase
}

func (col *PostController) CreateOrPublish(c *gin.Context) {
	userID := c.MustGet("x-user-id").(int64)
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	post.AuthorID = userID
	post.PostID = snowflake.GenID()
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
	//if req.Page == 0 || req.Size == 0 {
	//	req.Page = domain.DefaultPage
	//	req.Size = domain.DefaultSize
	//}

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
	//if req.Page == 0 || req.Size == 0 { // TODO 不处理 放在 orm 中处理
	//	req.Page = domain.DefaultPage
	//	req.Size = domain.DefaultSize
	//}
	userID := c.MustGet("x-user-id").(int64)
	posts, count, err := col.PostUsecase.List(c, &domain.Post{AuthorID: userID}, req.Page, req.Size) // TODO
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
	postIDRaw := c.Query("post_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	userID := c.MustGet("x-user-id").(int64)
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
	userID := c.MustGet("x-user-id").(int64)

	if isLike {
		err = col.InteractionUseCase.Like(c, domain.BizPost, postID, userID)
	} else {
		err = col.InteractionUseCase.CancelLike(c, domain.BizPost, postID, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *PostController) Collect(c *gin.Context) {
	isLikeRaw := c.Request.FormValue("is_collect")
	postIDRaw := c.Request.FormValue("post_id")
	collectionIDRaw := c.Request.FormValue("collection_id")
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
	collectID, err := strconv.ParseInt(collectionIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	userID := c.MustGet("x-user-id").(int64)

	if isLike {
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
