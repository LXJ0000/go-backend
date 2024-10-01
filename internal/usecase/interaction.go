package usecase

import (
	"errors"
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"

	"golang.org/x/net/context"
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
	// 判断是否点过赞
	_, stat, err := uc.repo.Stat(ctx, biz, bizID, userID)
	if err != nil {
		return err
	}
	if stat.Liked {
		slog.Error("重复点赞", "biz", biz, "bizID", bizID, "userID", userID)
		return errors.New("请勿重复点赞")
	}
	return uc.repo.Like(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) CancelLike(c context.Context, biz string, bizID, userID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	// 判断是否没点过赞
	_, stat, err := uc.repo.Stat(ctx, biz, bizID, userID)
	if err != nil {
		return err
	}
	if !stat.Liked {
		slog.Error("取消点赞失败", "biz", biz, "bizID", bizID, "userID", userID)
		return errors.New("取消点赞失败")
	}
	return uc.repo.CancelLike(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) Info(c context.Context, biz string, bizID, userID int64) (domain.Interaction, domain.UserInteractionStat, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Stat(ctx, biz, bizID, userID)
}

func (uc *interactionUsecase) Collect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	// 判断是否收藏过 TODO 收藏夹处理
	_, stat, err := uc.repo.Stat(ctx, biz, bizID, userID)
	if err != nil {
		return err
	}
	if stat.Collected {
		slog.Error("重复收藏", "biz", biz, "bizID", bizID, "userID", userID)
		return errors.New("请勿重复收藏")
	}
	return uc.repo.Collect(ctx, biz, bizID, userID, collectionID)
}

func (uc *interactionUsecase) CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	// 判断是否没收藏过 TODO 收藏夹处理
	_, stat, err := uc.repo.Stat(ctx, biz, bizID, userID)
	if err != nil {
		return err
	}
	if !stat.Collected {
		slog.Error("取消收藏失败", "biz", biz, "bizID", bizID, "userID", userID)
		return errors.New("取消收藏失败")
	}
	return uc.repo.CancelCollect(ctx, biz, bizID, userID, collectionID)
}

func (uc *interactionUsecase) GetByIDs(c context.Context, biz string, bizIDs []int64) (map[int64]domain.Interaction, error) {
	return nil, nil
}
