package route

import (
	"github.com/LXJ0000/go-backend/internal/domain"

	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func NewRefreshTokenRouter(env *bootstrap.Env, UserUc domain.UserUsecase, RefreshTokenUsecase domain.RefreshTokenUsecase,
	group *gin.RouterGroup) {

	rtc := &controller.RefreshTokenController{
		UserUsecase:         UserUc,
		RefreshTokenUsecase: RefreshTokenUsecase,
		Env:                 env,
	}
	group.POST("/refresh", rtc.RefreshToken)
}
