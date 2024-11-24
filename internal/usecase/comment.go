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

func (uc *commentUsecase) Find(c context.Context, biz string, bizID, parentID, minID int64, limit int) ([]domain.Comment, int, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.commentRepo.Find(ctx, biz, bizID, parentID, minID, limit)
}

func (uc *commentUsecase) Count(c context.Context, biz string, bizID int64) (int, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.commentRepo.Count(ctx, biz, bizID)
}
