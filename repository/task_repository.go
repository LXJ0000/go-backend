package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/orm"

	"github.com/LXJ0000/go-backend/domain"
)

type taskRepository struct {
	db orm.Database
	//collection string
}

func NewTaskRepository(db orm.Database) domain.TaskRepository {
	return &taskRepository{
		db: db,
		//collection: collection,
	}
}

func (repo *taskRepository) Create(c context.Context, task *domain.Task) error {
	_, err := repo.db.InsertOne(c, &domain.Task{}, task)
	return err
}

func (repo *taskRepository) Delete(c context.Context, taskID int64) error {
	_, err := repo.db.DeleteOne(c, &domain.Task{}, &domain.Task{TaskID: taskID})
	return err
}
