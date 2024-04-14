package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"log/slog"
	"math"
	"time"
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

func (repo *interactionRepository) IncrReadCount(c context.Context, biz string, id int64) error {
	now := time.Now()
	update := map[string]interface{}{
		"read_cnt":   gorm.Expr("`read_cnt` + 1"),
		"updated_at": now,
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
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", id, "error", err.Error())
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) Like(c context.Context, biz string, bizID, userID int64) error {
	now := time.Now()
	updateInteraction := map[string]interface{}{
		"like_cnt":   gorm.Expr("`like_cnt` + 1"),
		"updated_at": now,
	}
	createInteraction := &domain.Interaction{
		BizID:   bizID,
		Biz:     biz,
		LikeCnt: 1,
	}
	updateUserLike := map[string]interface{}{
		"status":     true,
		"updated_at": now,
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
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", bizID, "error", err.Error())
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
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", bizID, "error", err.Error())
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
		var cacheValue domain.CacheInteractionKey
		res, err := repo.cache.Get(c, key(biz, bizID))
		if err == nil {
			_ = json.Unmarshal([]byte(res), &cacheValue)
			interaction.CollectCnt = cacheValue.CollectCnt
			interaction.ReadCnt = cacheValue.ReadCnt
			interaction.LikeCnt = cacheValue.LikeCnt
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
			cacheValue.LikeCnt = interaction.LikeCnt
			cacheValue.ReadCnt = interaction.ReadCnt
			cacheValue.CollectCnt = interaction.CollectCnt
			val, _ := json.Marshal(cacheValue)
			if err := repo.cache.Set(context.Background(), key(biz, bizID), val, time.Duration(math.MaxInt)); err != nil {
				slog.Warn("Redis 操作失败, Set", "biz", biz, "bizID", bizID, "Key", cacheValue, "error", err.Error())
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
	_, err := repo.cache.Lua(c, domain.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, 1)
	return err
}

func (repo *interactionRepository) cacheDecrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, domain.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, -1)
	return err
}

func key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}
