package usecase

import (
	"github.com/LXJ0000/go-backend/domain"
	"golang.org/x/net/context"
	"time"
)

type postUsecase struct {
	repo           domain.PostRepository
	contextTimeout time.Duration
}

func NewPostUsecase(repo domain.PostRepository, timeout time.Duration) domain.PostUsecase {
	return &postUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}
func (uc *postUsecase) Create(c context.Context, post *domain.Post) error {
	return uc.repo.Create(c, post)
}
