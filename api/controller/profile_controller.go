package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.MustGet("x-user-id")
	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
