package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func (pc *UserController) Fetch(c *gin.Context) {
	userID := c.MustGet("x-user-id")
	profile, err := pc.UserUsecase.GetProfileByID(c, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
