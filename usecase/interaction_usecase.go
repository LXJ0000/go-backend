package usecase

import (
	"github.com/LXJ0000/go-backend/domain"
	"golang.org/x/net/context"
)

type interactionUsecase struct {
	repo domain.InteractionRepository
}

func NewInteractionUsecase(repo domain.InteractionRepository) domain.InteractionUseCase {
	return &interactionUsecase{repo: repo}
}

func (uc *interactionUsecase) IncrReadCount(c context.Context, biz string, id int64) error {
	return nil
}
