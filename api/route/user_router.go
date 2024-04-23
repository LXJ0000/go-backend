package route

import (
	"fmt"
	"github.com/LXJ0000/go-backend/api/controller"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/repository"
	"github.com/LXJ0000/go-backend/usecase"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db orm.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	pc := &controller.UserController{
		UserUsecase: usecase.NewProfileUsecase(ur, timeout),
	}
	user := group.Group("/user")
	user.GET("/profile", pc.Fetch)
	user.POST("/avatar", func(c *gin.Context) {
		// 单个文件
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Println(file.Filename)
		file.Filename = fmt.Sprintf("%d%s", time.Now().UnixMilli(), filepath.Ext(file.Filename))
		dst := filepath.Join("assets/avatar", file.Filename)
		// 上传文件到指定的目录
		if err := c.SaveUploadedFile(file, dst); err != nil {
			slog.Error("save file err:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})
	user.POST("/multi_upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for index, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf("C:/tmp/%s_%d", file.Filename, index)
			// 上传文件到指定的目录
			c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})
}
