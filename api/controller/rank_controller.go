package controller

import (
	domain "github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RankController struct {
	postRankUsecase usecase.PostRankUsecase
}

func (col *RankController) GetTopN(c *gin.Context) {
	posts, err := col.postRankUsecase.GetTopN(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("GetTopN Fail", err))
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"count": len(posts),
		"posts": posts,
	})
}
