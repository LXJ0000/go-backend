package usecase

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"golang.org/x/sync/errgroup"
)

type relationUsecase struct {
	repo           domain.RelationRepository
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewRelationUsecase(repo domain.RelationRepository, userRepo domain.UserRepository, contextTimeout time.Duration) domain.RelationUsecase {
	return &relationUsecase{repo: repo, userRepo: userRepo, contextTimeout: contextTimeout}
}

func (uc *relationUsecase) Follow(c context.Context, follower, followee int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Follow(ctx, follower, followee)
}

func (uc *relationUsecase) CancelFollow(c context.Context, follower, followee int64) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.CancelFollow(ctx, follower, followee)
}

func (uc *relationUsecase) GetFollower(c context.Context, userID int64, page, size int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	relations, err := uc.repo.GetFollower(ctx, userID, page, size)
	if err != nil {
		//  TODO log
		return nil, err
	}
	userIDs := make([]int64, 0, len(relations))
	for _, relation := range relations {
		userIDs = append(userIDs, relation.Follower)
	}
	return uc.userRepo.FindByUserIDs(ctx, userIDs, page, size)
}

func (uc *relationUsecase) GetFollowee(c context.Context, userID int64, page, size int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	relations, err := uc.repo.GetFollowee(ctx, userID, page, size)
	if err != nil {
		//  TODO log
		return nil, err
	}
	userIDs := make([]int64, 0, len(relations))
	for _, relation := range relations {
		userIDs = append(userIDs, relation.Followee)
	}
	return uc.userRepo.FindByUserIDs(ctx, userIDs, page, size)
}

func (uc *relationUsecase) Detail(c context.Context, follower, followee int64) (domain.Relation, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Detail(ctx, follower, followee)
}

func (uc *relationUsecase) Stat(c context.Context, userID int64) (domain.RelationStat, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	var (
		err      error
		follower int64
		followee int64
	)
	eg := errgroup.Group{}
	eg.Go(func() error {
		follower, err = uc.repo.FollowerCnt(ctx, userID)
		return err
	})
	eg.Go(func() error {
		followee, err = uc.repo.FolloweeCnt(ctx, userID)
		return err
	})
	if err := eg.Wait(); err != nil {
		//  TODO log
		return domain.RelationStat{}, err
	}
	return domain.RelationStat{
		UserID:   userID,
		Follower: int(follower),
		Followee: int(followee),
	}, nil
}
