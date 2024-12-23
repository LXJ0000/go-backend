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
	localCodeUc domain.CodeUsecase,
	sync2OpenIMUc domain.Sync2OpenIMUsecase,
	fileUsecase domain.FileUsecase,
	publicRouter *gin.RouterGroup,
	group *gin.RouterGroup) {
	col := &controller.UserController{
		UserUsecase:        userUc,
		Env:                env,
		PostUsecase:        postUc,
		RelationUsecase:    relationUc,
		CodeUsecase:        codeUc,
		Sync2OpenIMUsecase: sync2OpenIMUc,
		FileUsecase:        fileUsecase,
	}
	group.GET("/user/profile", col.Fetch)
	group.GET("/user", col.Profile)
	group.POST("/logout", col.Logout)
	group.POST("/user/edit", col.Update)
	group.POST("/user/search", col.Search)
	group.POST("/user/avatar", col.Avatar)
	group.POST("/user.batch", col.BatchProfile)
	group.POST("/user.reset_password", col.ResetPassword)

	publicRouter.POST("/login", col.Login)
	publicRouter.POST("/signup", col.Signup)
	publicRouter.POST("/send_sms_code", col.SendSMSCode)
	publicRouter.POST("/login/sms", col.LoginBySms)
	publicRouter.POST("/login/phone", col.LoginByPhone)
}
