package controller

import (
	"github.com/LXJ0000/go-backend/domain"
	snowflake "github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

func (col *TaskController) Create(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	task.TaskID = snowflake.GenID()
	task.UserID = c.MustGet("x-user-id").(int64)

	err = col.TaskUsecase.Create(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
}

func (col *TaskController) Delete(c *gin.Context) {
	taskIDRaw := c.Query("task_id")
	taskID, err := strconv.ParseInt(taskIDRaw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if err = col.TaskUsecase.Delete(c, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
}
