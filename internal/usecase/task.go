package usecase

import (
	"context"
	"github.com/LXJ0000/go-backend/internal/domain"
	"time"
)

type taskUsecase struct {
	repo           domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		repo:           taskRepository,
		contextTimeout: timeout,
	}
}

func (uc *taskUsecase) Create(c context.Context, task domain.Task) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Create(ctx, task)
}

func (uc *taskUsecase) Delete(c context.Context, taskID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Delete(ctx, taskID)
}
