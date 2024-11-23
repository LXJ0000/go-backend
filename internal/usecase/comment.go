package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type commentUsecase struct {
	commentRepo    domain.CommentRepository
	contextTimeout time.Duration
}

func NewCommentUsecase(commentRepo domain.CommentRepository, contextTimeout time.Duration) domain.CommentUsecase {
	return &commentUsecase{commentRepo: commentRepo, contextTimeout: contextTimeout}
}

func (uc *commentUsecase) Create(c context.Context, comment *domain.Comment) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.commentRepo.Create(ctx, comment)
}

func (uc *commentUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.commentRepo.Delete(ctx, id)
}

func (uc *commentUsecase) FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.commentRepo.FindTop(ctx, biz, bizID, minID, limit)
}
