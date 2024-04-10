package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	snowflake "github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostController struct {
	Usecase domain.PostUsecase
}

func (col *PostController) Create(c *gin.Context) {
	userID := c.MustGet("x-user-id")
	var post domain.Post
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	post.UserID = userID.(int64)
	post.PostID = snowflake.GenID()
	if err := col.Usecase.Create(c, &post); err != nil {
		c.JSON(http.StatusOK, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Post created successfully"})
}
