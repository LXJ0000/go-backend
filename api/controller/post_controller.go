package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	snowflake "github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"strconv"
)

type PostController struct {
	PostUsecase        domain.PostUsecase
	InteractionUseCase domain.InteractionUseCase
}

func (col *PostController) CreateOrPublish(c *gin.Context) {
	userID := c.MustGet("x-user-id").(int64)
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	post.AuthorID = userID
	post.PostID = snowflake.GenID()
	if err := col.PostUsecase.Create(c, &post); err != nil {
		c.JSON(http.StatusOK, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Post created successfully"})
}

func (col *PostController) ReaderList(c *gin.Context) {
	//读者查看列表 只能查看已发布的文章
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if req.Page == 0 || req.Size == 0 {
		req.Page = domain.DefaultPage
		req.Size = domain.DefaultSize
	}

	filter := &domain.Post{
		Status: domain.PostStatusPublish,
	}
	if req.AuthorID != 0 {
		filter.AuthorID = req.AuthorID
	}
	posts, err := col.PostUsecase.List(
		c, filter,
		req.Page, req.Size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.PostListResponse{
		Count: len(posts), // TODO fix count
		Data:  posts,
	})
}

func (col *PostController) WriterList(c *gin.Context) {
	//创作者查看列表 可以查看所有自己的帖子
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if req.Page == 0 || req.Size == 0 {
		req.Page = domain.DefaultPage
		req.Size = domain.DefaultSize
	}
	userID := c.MustGet("x-user-id").(int64)
	posts, err := col.PostUsecase.List(c, &domain.Post{AuthorID: userID}, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.PostListResponse{
		Count: len(posts),
		Data:  posts,
	})
}

func (col *PostController) Info(c *gin.Context) {
	postIDRaw := c.Query("post_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	userID := c.MustGet("x-user-id").(int64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	var post domain.Post
	var interaction domain.Interaction
	var userInteractionInfo domain.UserInteractionInfo
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
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	go func() {
		if err := col.InteractionUseCase.IncrReadCount(c, domain.BizPost, post.PostID); err != nil {
			slog.Warn("Add post read count fail", "post_id", post.PostID)
		}
	}() // 添加文件阅读数
	c.JSON(http.StatusOK, domain.PostV0{
		Post:       post,
		ReadCnt:    interaction.ReadCnt,
		LikeCnt:    interaction.LikeCnt,
		CollectCnt: interaction.CollectCnt,
		Collected:  userInteractionInfo.Collected,
		Liked:      userInteractionInfo.Liked,
	})
}

func (col *PostController) Like(c *gin.Context) {
	isLikeRaw := c.Request.FormValue("is_like")
	postIDRaw := c.Request.FormValue("post_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	isLike, err := strconv.ParseBool(isLikeRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	userID := c.MustGet("x-user-id").(int64)

	if isLike {
		err = col.InteractionUseCase.Like(c, domain.BizPost, postID, userID)
	} else {
		err = col.InteractionUseCase.CancelLike(c, domain.BizPost, postID, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Success"})
}

func (col *PostController) Collect(c *gin.Context) {
	isLikeRaw := c.Request.FormValue("is_collect")
	postIDRaw := c.Request.FormValue("post_id")
	collectionIDRaw := c.Request.FormValue("collection_id")
	postID, err := strconv.ParseInt(postIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	isLike, err := strconv.ParseBool(isLikeRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	collectID, err := strconv.ParseInt(collectionIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	userID := c.MustGet("x-user-id").(int64)

	if isLike {
		err = col.InteractionUseCase.Collect(c, domain.BizPost, postID, userID, collectID)
	} else {
		err = col.InteractionUseCase.CancelCollect(c, domain.BizPost, postID, userID, collectID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Success"})
}
