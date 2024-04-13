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
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Create(ctx, post)
}

func (uc *postUsecase) List(c context.Context, filter *domain.Post, page, size int) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.FindMany(ctx, filter, page, size)
}

func (uc *postUsecase) Info(c context.Context, postID int64) (domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.GetByID(ctx, postID)
}
