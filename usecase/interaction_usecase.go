package usecase

import (
	"github.com/LXJ0000/go-backend/domain"
	"golang.org/x/net/context"
	"time"
)

type interactionUsecase struct {
	repo           domain.InteractionRepository
	contextTimeout time.Duration
}

func NewInteractionUsecase(repo domain.InteractionRepository, timeout time.Duration) domain.InteractionUseCase {
	return &interactionUsecase{repo: repo, contextTimeout: timeout}
}

func (uc *interactionUsecase) IncrReadCount(c context.Context, biz string, id int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.IncrReadCount(ctx, biz, id)
}
