package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

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
	domain.UserUsecase
	domain.CommentUsecase
	domain.FileUsecase
}

func (col *PostController) Search(c *gin.Context) {
	//搜索文章
	var req domain.PostSearchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	resp, count, err := col.PostUsecase.Search(c, req.Keyword, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to search posts", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": resp,
	}))
}

func (col *PostController) PostDelete(c *gin.Context) {
	//删除文章
	var req domain.PostDeleteRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)
	postID := req.PostID
	post, err := col.PostUsecase.Info(c, postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Post not found", err))
		return
	}
	if post.AuthorID != userID {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Permission denied", nil))
		return
	}
	if err := col.PostUsecase.Delete(c, postID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to delete post", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

func (col *PostController) CreateOrPublish(c *gin.Context) {
	userID := c.MustGet(domain.XUserID).(int64)
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}

	// 单独处理图片
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	files := form.File["files"]
	fileResp, err := col.FileUsecase.Uploads(c, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Upload files fail", err))
		return
	}
	images := make([]string, 0, len(files))
	for _, resp := range fileResp.Data {
		if file, ok := resp.(domain.File); ok {
			images = append(images, file.Path)
		}
	}
	bytes, err := json.Marshal(images)
	if err != nil {
		slog.Error("json.Marshal images error", "error", err.Error())
	} else {
		post.Images = string(bytes)
	}

	if post.Content == "" && post.Images == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Content or images is required", nil))
		return
	}

	post.AuthorID = userID
	post.PostID = snowflake.GenID()
	if err := col.PostUsecase.Create(c, &post); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to create post", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(post))
}

func (col *PostController) ReaderList(c *gin.Context) {
	//读者查看列表 只能查看已发布的文章
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}

	userID := c.MustGet(domain.XUserID).(int64)
	filter := &domain.Post{
		Status: domain.PostStatusPublish,
	}
	if req.AuthorID != 0 {
		filter.AuthorID = req.AuthorID
	}
	count, resp := col.getPost(c, userID, filter, req)
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": resp,
	}))
}

func (col *PostController) WriterList(c *gin.Context) {
	//创作者查看列表 可以查看所有自己的帖子
	var req domain.PostListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	userID := c.MustGet(domain.XUserID).(int64)
	filter := &domain.Post{AuthorID: userID}
	count, resp := col.getPost(c, userID, filter, req)
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     count,
		"post_list": resp,
	}))
}

func (col *PostController) Info(c *gin.Context) {
	postID, err := lib.Str2Int64(c.Query("post_id"))
	userID := c.MustGet(domain.XUserID).(int64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
		return
	}
	var (
		post                domain.Post
		interaction         domain.Interaction
		userInteractionInfo domain.UserInteractionStat
		commentCount        int
	)
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
	eg.Go(func() error {
		commentCount = col.getCommentCount(c, domain.BizPost, postID)
		return nil
	})
	if err = eg.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(err.Error(), err))
		return
	}
	postResp, err := col.parsePostResponse(c, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("parsePostResponse Error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(domain.PostInfoResponse{
		Post:         postResp,
		Interaction:  interaction,
		Stat:         userInteractionInfo,
		CommentCount: commentCount,
	}))
}

// Like 废弃方法 转移到 interaction 里
func (col *PostController) Like(c *gin.Context) { // TODO 抽象成资源操作而不是针对帖子的操作
	req := struct {
		IsLike bool  `json:"is_like" form:"is_like"`
		PostID int64 `json:"post_id,string" form:"post_id"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
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
		// 获取帖子的作者
		ctx, cancel := context.WithTimeout(c, time.Second)
		defer cancel()
		post, err := col.PostUsecase.Info(ctx, postID)
		if err != nil {
			slog.Error("FeedUsecase CreateFeedEvent Error Because PostUsecase Info Error", "error", err.Error())
			return
		}
		feed := domain.Feed{
			UserID: userID,
			Type:   domain.FeedLikeEvent,
			Content: domain.FeedContent{
				"biz_id": fmt.Sprintf("%d", postID),
				"biz":    domain.BizPost,
				"liker":  fmt.Sprintf("%d", userID),
				"liked":  fmt.Sprintf("%d", post.AuthorID),
			}, // liker liked biz bizID
		}
		// liker 点赞了 biz_id
		if err := col.FeedUsecase.CreateFeedEvent(context.Background(), feed); err != nil {
			slog.Warn("FeedUsecase CreateFeedEvent Error", "error", err.Error())
		}
	}()
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}

// Collect 废弃方法 转移到 interaction 里
func (col *PostController) Collect(c *gin.Context) {
	req := struct {
		IsCollect bool  `json:"is_collect" form:"is_collect"`
		PostID    int64 `json:"post_id,string" form:"post_id"`
		CollectID int64 `json:"collect_id,string" form:"collect_id"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp(domain.ErrBadParams.Error(), err))
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
		c.JSON(http.StatusInternalServerError, domain.ErrorResp(domain.ErrBusy.Error(), err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count": len(posts),
		"posts": posts,
	}))
}

func (col *PostController) parsePostResponse(c context.Context, post domain.Post) (domain.PostResponse, error) {
	author, err := col.UserUsecase.GetProfileByID(c, post.AuthorID)
	if err != nil {
		return domain.PostResponse{}, err
	}
	postResp := domain.PostResponse{
		Author: author,
	}
	postResp.PostID = post.PostID
	postResp.Title = post.Title
	postResp.Abstract = post.Abstract
	postResp.Content = post.Content
	postResp.Status = post.Status
	postResp.CreatedAt = post.CreatedAt
	postResp.UpdatedAt = post.UpdatedAt
	images := make([]string, 0)
	if err := json.Unmarshal([]byte(post.Images), &images); err != nil {
		slog.Error("json.Unmarshal images error", "error", err.Error())
	} else {
		postResp.Images = images
	}
	return postResp, nil
}

func (col *PostController) getCommentCount(c *gin.Context, biz string, bizID int64) int {
	count, err := col.CommentUsecase.Count(c, biz, bizID)
	if err != nil {
		slog.Error("CommentUsecase Count Error", "error", err.Error())
		return 0
	}
	return count
}

func (col *PostController) getPost(c *gin.Context, userID int64, filter interface{}, req domain.PostListRequest) (int64, []domain.PostInfoResponse) {
	posts, count, err := col.PostUsecase.ListByLastID(c, filter, req.Size, req.Last)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Failed to list posts", err))
		return 0, nil
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
			postResp, err := col.parsePostResponse(c, post)
			if err != nil {
				slog.Error("parsePostResponse Error", "error", err.Error())
				return
			}
			commentCount := col.getCommentCount(c, domain.BizPost, post.PostID)
			resp = append(resp, domain.PostInfoResponse{
				Post:         postResp,
				Interaction:  interaction,
				Stat:         userInteractionInfo,
				CommentCount: commentCount,
			})
		}()
	}
	wg.Wait()
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Post.CreatedAt > resp[j].Post.CreatedAt
	})
	return count, resp
}
