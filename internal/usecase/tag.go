package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type tagUsecase struct {
	repo           domain.TagRepository
	contextTimeout time.Duration
}

func NewTagUsecase(repo domain.TagRepository, contextTimeout time.Duration) domain.TagUsecase {
	return &tagUsecase{repo: repo, contextTimeout: contextTimeout}
}

func (t *tagUsecase) CreateTag(c context.Context, tag domain.Tag) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.repo.CreateTag(ctx, tag)
}

func (t *tagUsecase) CreateTagBiz(c context.Context, userID int64, biz string, bizID int64, tagIDs []int64) error {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.repo.CreateTagBiz(ctx, userID, biz, bizID, tagIDs)
	// TODO send to es
}

func (t *tagUsecase) GetTagsByUserID(c context.Context, userID int64) ([]domain.Tag, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.repo.GetTagsByUserID(ctx, userID)
}

func (t *tagUsecase) GetTagsByBiz(c context.Context, userID int64, biz string, bizID int64) ([]domain.Tag, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	return t.repo.GetTagsByBiz(ctx, userID, biz, bizID)
}
