package repository

import (
	"fmt"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log/slog"
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
		"read_cnt":   gorm.Expr("read_cnt + 1"),
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
		if err := repo.CacheIncrCnt(c, biz, id, "read_cnt"); err != nil {
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", id)
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) Like(c context.Context, biz string, bizID, userID int64) error {
	now := time.Now()
	updateInteraction := map[string]interface{}{
		"like_cnt":   gorm.Expr("like_cnt + 1"),
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
		if err := repo.CacheIncrCnt(c, biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", bizID)
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
				"like_cnt": gorm.Expr("like_cnt - 1"),
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
		if err := repo.CacheDecrCnt(c, biz, bizID, "like_cnt"); err != nil {
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", bizID)
		}
	}()
	return nil
}

func (repo *interactionRepository) CacheIncrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, domain.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, "1")
	return err
}

func (repo *interactionRepository) CacheDecrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, domain.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, "-1")
	return err
}

func key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}
