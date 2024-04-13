package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	snowflake "github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostController struct {
	PostUsecase domain.PostUsecase
}

func (col *PostController) CreateOrPublish(c *gin.Context) {
	userID := c.MustGet("x-user-id")
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	post.AuthorID = userID.(int64)
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
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	post, err := col.PostUsecase.Info(c, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	go func() {
		//	TODO cache readCnt
	}()
	c.JSON(http.StatusOK, post)
}
