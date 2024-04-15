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

func (uc *interactionUsecase) Like(c context.Context, biz string, bizID, userID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Like(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) CancelLike(c context.Context, biz string, bizID, userID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.CancelLike(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) Info(c context.Context, biz string, bizID, userID int64) (domain.Interaction, domain.UserInteractionInfo, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Info(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) Collect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Collect(ctx, biz, bizID, userID, collectionID)
}

func (uc *interactionUsecase) CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.CancelCollect(ctx, biz, bizID, userID, collectionID)
}
