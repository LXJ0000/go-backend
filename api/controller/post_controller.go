package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostController struct {
	Usecase domain.PostUsecase
}

func (col *PostController) Create(c *gin.Context) {
	post := &domain.Post{
		PostID:   2,
		Title:    "222 123",
		Abstract: "222 234",
		Content:  "222 345",
	}
	if err := col.Usecase.Create(c, post); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}
