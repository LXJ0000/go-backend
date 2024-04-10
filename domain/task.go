package domain

import (
	"context"
	"gorm.io/gorm"
)

//const (
//	CollectionTask = "tasks"
//)
//
//type Task struct {
//	ID     primitive.ObjectID `bson:"_id" json:"-"`
//	Title  string             `bson:"title" form:"title" binding:"required" json:"title"`
//	UserID primitive.ObjectID `bson:"userID" json:"-"`
//}
//func (Task) TableName() string {
//	return `task`
//}
//type TaskRepository interface {
//	Create(c context.Context, task *Task) error
//	FetchByUserID(c context.Context, userID string) ([]Task, error)
//}
//
//type TaskUsecase interface {
//	Create(c context.Context, task *Task) error
//	FetchByUserID(c context.Context, userID string) ([]Task, error)
//}

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
	Create(c context.Context, task *Task) error
	Delete(c context.Context, taskID int64) error
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	Delete(c context.Context, taskID int64) error
}
