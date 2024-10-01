package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

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

	userID := c.MustGet(domain.XUserID).(int64)
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

	resp := make([]domain.PostInfoResponse, 0, len(posts))
	for _, post := range posts {
		go func() {
			interaction, userInteractionInfo, err := col.InteractionUseCase.Info(c, domain.BizPost, post.PostID, userID)
			if err != nil {
				slog.Warn("InteractionUseCase Info Error", "error", err.Error())
				return
			}
			resp = append(resp, domain.PostInfoResponse{
				Post:        post,
				Interaction: interaction,
				Stat:        userInteractionInfo,
			})
		}()
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": resp,
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

	resp := make([]domain.PostInfoResponse, 0, len(posts))
	wg := sync.WaitGroup{}
	wg.Add(len(posts))
	for _, post := range posts {
		go func() {
			defer wg.Done()
			interaction, userInteractionInfo, err := col.InteractionUseCase.Info(c, domain.BizPost, post.PostID, userID)
			if err != nil {
				slog.Warn("InteractionUseCase Info Error", "error", err.Error())
				return
			}
			resp = append(resp, domain.PostInfoResponse{
				Post:        post,
				Interaction: interaction,
				Stat:        userInteractionInfo,
			})
		}()
	}
	wg.Wait()
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": resp,
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
	c.JSON(http.StatusOK, domain.SuccessResp(domain.PostInfoResponse{
		Post:        post,
		Interaction: interaction,
		Stat:        userInteractionInfo,
	}))
}

func (col *PostController) Like(c *gin.Context) {
	req := struct {
		IsLike bool  `json:"is_like" form:"is_like"`
		PostID int64 `json:"post_id,string" form:"post_id"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	postID := req.PostID
	isLike := req.IsLike
	userID := c.MustGet(domain.XUserID).(int64)
	var err error
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
			slog.Warn("FeedUsecase CreateFeedEvent Error", "error", err.Error())
		}
	}()
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *PostController) Collect(c *gin.Context) {
	req := struct {
		IsCollect bool  `json:"is_collect" form:"is_collect"`
		PostID    int64 `json:"post_id,string" form:"post_id"`
		CollectID int64 `json:"collect_id,string" form:"collect_id"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad params", err))
		return
	}
	postID := req.PostID
	isCollect := req.IsCollect
	collectID := req.CollectID
	userID := c.MustGet(domain.XUserID).(int64)
	var err error
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
