package controller

import (
	"github.com/LXJ0000/go-backend/internal/domain"
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
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Get profile by user_id fail with db error", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"profile": profile,
	}))
}