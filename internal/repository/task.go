package repository

import (
	"context"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
)

type taskRepository struct {
	dao orm.Database
}

func NewTaskRepository(dao orm.Database) domain.TaskRepository {
	return &taskRepository{
		dao: dao,
	}
}

func (repo *taskRepository) Create(c context.Context, task domain.Task) error {
	return repo.dao.Insert(c, &domain.Task{}, &task)
}

func (repo *taskRepository) Delete(c context.Context, taskID int64) error {
	return repo.dao.DeleteOne(c, &domain.Task{}, &domain.Task{TaskID: taskID})
}
