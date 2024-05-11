package domain

import (
	"context"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskID int64  `json:"task_id"`
	Title  string `json:"title" form:"title"`
	UserID int64  `json:"user_id" form:"user_id"`
}

func (Task) TableName() string {
	return `task`
}

type TaskRepository interface {
	Create(c context.Context, task Task) error
	Delete(c context.Context, taskID int64) error
}

type TaskUsecase interface {
	Create(c context.Context, task Task) error
	Delete(c context.Context, taskID int64) error
}
