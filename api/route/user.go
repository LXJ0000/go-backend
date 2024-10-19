package route

import (
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *bootstrap.Env,
	userUc domain.UserUsecase,
	relationUc domain.RelationUsecase,
	postUc domain.PostUsecase,
	codeUc domain.CodeUsecase,
	publicRouter *gin.RouterGroup,
	group *gin.RouterGroup) {
	col := &controller.UserController{
		UserUsecase:     userUc,
		Env:             env,
		PostUsecase:     postUc,
		RelationUsecase: relationUc,
		CodeUsecase:     codeUc,
	}
	group.POST("/logout", col.Logout)
	group.GET("/user/profile", col.Fetch)
	group.POST("/user/edit", col.Update)
	group.GET("/user", col.Profile)

	publicRouter.POST("/login", col.Login)
	publicRouter.POST("/signup", col.Signup)
	publicRouter.POST("/send_sms_code", col.SendSMSCode)
	publicRouter.POST("/login/sms", col.LoginBySms)

}
