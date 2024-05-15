package usecase

import "github.com/LXJ0000/go-backend/internal/domain"

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{commentRepo: commentRepo}
}
