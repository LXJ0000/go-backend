package repository

import (
	"github.com/LXJ0000/go-backend/script"

	"errors"
	"fmt"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type interactionRepository struct {
	dao   orm.Database
	cache redis.Cache
}

func NewInteractionRepository(dao orm.Database, cache redis.Cache) domain.InteractionRepository {
	return &interactionRepository{
		dao:   dao,
		cache: cache,
	}
}

// BatchIncrReadCount 批量增加read_cnt 需保证 len(biz) == len(id)
func (repo *interactionRepository) BatchIncrReadCount(c context.Context, biz []string, id []int64) error {
	fn := func(tx *gorm.DB) error {
		update := map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
		}

		for i := 0; i < len(biz); i++ {
			i := i // 1.22 可不写
			create := &domain.Interaction{
				BizID:   id[i],
				Biz:     biz[i],
				ReadCnt: 1,
			}
			if err := repo.dao.Upsert(c, &domain.Interaction{}, update, create); err != nil {
				slog.Error("IncrReadCount Fail", "Error", err.Error(), "biz", biz[i], "biz_id", id[i])
			}
			go func() { // TODO new lua script or pipeline
				if err := repo.cacheIncrCnt(context.Background(), biz[i], id[i], "read_cnt"); err != nil {
					slog.Warn("Redis Op Fail With CacheIncrReadCnt", "Error", err.Error(), "biz", biz[i], "bizID", id[i])
				}
			}()
		}
		return nil
	}
	_ = repo.dao.Transaction(c, fn)
	return nil
}

func (repo *interactionRepository) IncrReadCount(c context.Context, biz string, id int64) error {
	update := map[string]interface{}{
		"read_cnt": gorm.Expr("`read_cnt` + 1"),
	}
	create := &domain.Interaction{
		BizID:   id,
		Biz:     biz,
		ReadCnt: 1,
	}
	if err := repo.dao.Upsert(c, &domain.Interaction{}, update, create); err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, id, "read_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrReadCnt", "Error", err.Error(), "biz", biz, "bizID", id)
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) Like(c context.Context, biz string, bizID, userID int64) error {
	updateInteraction := map[string]interface{}{
		"like_cnt": gorm.Expr("`like_cnt` + 1"),
	}
	createInteraction := &domain.Interaction{
		BizID:   bizID,
		Biz:     biz,
		LikeCnt: 1,
	}
	updateUserLike := map[string]interface{}{
		"status": true,
	}
	createUserLike := &domain.UserLike{
		BizID:  bizID,
		Biz:    biz,
		UserID: userID,
		Status: true,
	}
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		if err := dao.Upsert(c, &domain.Interaction{}, updateInteraction, createInteraction); err != nil {
			return err
		}
		return dao.Upsert(c, &domain.UserLike{}, updateUserLike, createUserLike)
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrLikeCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) CancelLike(c context.Context, biz string, bizID, userID int64) error {
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		//1. 更新 UserLike status = false
		if _, err := dao.UpdateOne(c,
			&domain.UserLike{},
			&domain.UserLike{
				UserID: userID,
				BizID:  bizID,
				Biz:    biz,
			},
			map[string]interface{}{
				"status": false,
			},
		); err != nil {
			return err
		}
		//2. 更新 interaction like_cnt - 1
		if _, err := dao.UpdateOne(c,
			&domain.Interaction{},
			&domain.Interaction{
				BizID: bizID,
				Biz:   biz,
			},
			map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
			},
		); err != nil {
			return err
		}
		return nil
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheDecrCnt(context.Background(), biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheDecrLikeCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) Info(c context.Context, biz string, bizID, userID int64) (domain.Interaction, domain.UserInteractionInfo, error) {
	var isLike, isCollect bool
	var err error
	var interaction domain.Interaction
	eg := errgroup.Group{}
	eg.Go(func() error {
		isLike, err = repo.isLike(c, biz, bizID, userID)
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		isCollect, err = repo.isCollect(c, biz, bizID, userID)
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})
	eg.Go(func() error {
		res, err := repo.cache.HGetAll(c, key(biz, bizID))
		if err == nil && len(res) > 0 {
			interaction.CollectCnt, _ = strconv.Atoi(res["collect_cnt"])
			interaction.ReadCnt, _ = strconv.Atoi(res["read_cnt"])
			interaction.LikeCnt, _ = strconv.Atoi(res["like_cnt"])
			return nil
		}
		item, err := repo.dao.FindOne(c, &domain.Interaction{}, &domain.Interaction{
			BizID: bizID,
			Biz:   biz,
		})
		if err != nil {
			return err
		}
		interaction = *item.(*domain.Interaction)
		go func() {
			if err := repo.cache.HSet(context.Background(),
				key(biz, bizID),
				"read_cnt", interaction.ReadCnt,
				"collect_cnt", interaction.CollectCnt,
				"like_cnt", interaction.LikeCnt,
			); err != nil {
				slog.Warn("Redis Op Fail With HSet", "Error", err.Error(), "biz", biz, "bizID", bizID, "Key", key(biz, bizID))
			}
		}()
		return nil
	})
	if err := eg.Wait(); err != nil {
		return domain.Interaction{}, domain.UserInteractionInfo{}, err
	}
	return interaction,
		domain.UserInteractionInfo{
			Liked:     isLike,
			Collected: isCollect,
		}, nil
}

func (repo *interactionRepository) Collect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	updateInteraction := map[string]interface{}{
		"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
	}
	createInteraction := &domain.Interaction{
		BizID:      bizID,
		Biz:        biz,
		CollectCnt: 1,
	}
	updateUserCollect := map[string]interface{}{
		"status": true,
	}
	createUserCollect := &domain.UserCollect{
		BizID:        bizID,
		Biz:          biz,
		UserID:       userID,
		CollectionID: collectionID,
		Status:       true,
	}
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		if err := dao.Upsert(c, &domain.Interaction{}, updateInteraction, createInteraction); err != nil {
			return err
		}
		return dao.Upsert(c, &domain.UserCollect{}, updateUserCollect, createUserCollect)
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheIncrCnt(context.Background(), biz, bizID, "collect_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheIncrCollectCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) CancelCollect(c context.Context, biz string, bizID, userID, collectionID int64) error {
	fn := func(tx *gorm.DB) error {
		dao := orm.NewDatabase(tx)
		//1. 更新 UserLike status = false
		if _, err := dao.UpdateOne(c,
			&domain.UserCollect{},
			&domain.UserCollect{
				UserID:       userID,
				BizID:        bizID,
				CollectionID: collectionID,
				Biz:          biz,
			},
			map[string]interface{}{
				"status": false,
			},
		); err != nil {
			return err
		}
		//2. 更新 interaction like_cnt - 1
		if _, err := dao.UpdateOne(c,
			&domain.Interaction{},
			&domain.Interaction{
				BizID: bizID,
				Biz:   biz,
			},
			map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` - 1"),
			},
		); err != nil {
			return err
		}
		return nil
	}
	err := repo.dao.Transaction(c, fn)
	if err != nil {
		return err
	}
	go func() {
		if err := repo.cacheDecrCnt(context.Background(), biz, bizID, "collect_cnt"); err != nil {
			slog.Warn("Redis Op Fail With CacheDecrCollectCnt", "Error", err.Error(), "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) isLike(c context.Context, biz string, bizID, userID int64) (bool, error) {
	item, err := repo.dao.FindOne(c, &domain.UserLike{}, &domain.UserLike{
		UserID: userID,
		BizID:  bizID,
		Biz:    biz,
	})
	if err != nil {
		return false, err
	}
	status := (*item.(*domain.UserLike)).Status
	return status, nil
}

func (repo *interactionRepository) isCollect(c context.Context, biz string, bizID, userID int64) (bool, error) {
	item, err := repo.dao.FindOne(c, &domain.UserCollect{}, &domain.UserCollect{
		UserID: userID,
		BizID:  bizID,
		Biz:    biz,
	})
	if err != nil {
		return false, err
	}
	status := (*item.(*domain.UserCollect)).Status
	return status, nil
}

func (repo *interactionRepository) cacheIncrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, script.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, 1)
	return err
}

func (repo *interactionRepository) cacheDecrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, script.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, -1)
	return err
}

func key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}
