package controller

import (
	"net/http"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/utils/lib"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

func (col *TaskController) Create(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	now := time.Now().UnixMicro()
	task.TaskID = snowflake.GenID()
	task.UserID = c.MustGet(domain.XUserID).(int64)
	task.CreatedAt = now
	task.UpdatedAt = now

	err = col.TaskUsecase.Create(c, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Create task fail with db error", err))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResp(
		map[string]interface{}{
			"task_detail": task,
		}))
}

func (col *TaskController) Delete(c *gin.Context) {
	taskID, err := lib.Str2Int64(c.Query("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	if err = col.TaskUsecase.Delete(c, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Delete task fail with db error", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(nil))
}
