package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RankController struct {
	postRankUsecase usecase.PostRankUsecase
}

func (col *RankController) GetTopN(c *gin.Context) {
	posts, err := col.postRankUsecase.GetTopN(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "GetTopN Fail"})
	}
	c.JSON(http.StatusOK, domain.PostListResponse{
		Count: len(posts),
		Data:  posts,
	})
}
