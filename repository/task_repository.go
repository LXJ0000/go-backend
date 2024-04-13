package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/orm"

	"github.com/LXJ0000/go-backend/domain"
)

type taskRepository struct {
	dao orm.Database
	//collection string
}

func NewTaskRepository(dao orm.Database) domain.TaskRepository {
	return &taskRepository{
		dao: dao,
		//collection: collection,
	}
}

func (repo *taskRepository) Create(c context.Context, task *domain.Task) error {
	_, err := repo.dao.InsertOne(c, &domain.Task{}, task)
	return err
}

func (repo *taskRepository) Delete(c context.Context, taskID int64) error {
	_, err := repo.dao.DeleteOne(c, &domain.Task{}, &domain.Task{TaskID: taskID})
	return err
}
