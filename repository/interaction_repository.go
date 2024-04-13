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

func (repo *interactionRepository) IncrLikeCount(c context.Context, biz string, id int64) error {
	now := time.Now()
	update := map[string]interface{}{
		"like_cnt":   gorm.Expr("like_cnt + 1"),
		"updated_at": now,
	}
	create := &domain.Interaction{
		BizID:   id,
		Biz:     biz,
		LikeCnt: 1,
	}
	if err := repo.dao.Upsert(c, &domain.Interaction{}, update, create); err != nil {
		return err
	}
	go func() {
		if err := repo.CacheIncrCnt(c, biz, id, "like_cnt"); err != nil {
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", id)
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) IncrCollectCount(c context.Context, biz string, id int64) error {
	now := time.Now()
	update := map[string]interface{}{
		"collect_cnt": gorm.Expr("collect_cnt + 1"),
		"updated_at":  now,
	}
	create := &domain.Interaction{
		BizID:      id,
		Biz:        biz,
		CollectCnt: 1,
	}
	if err := repo.dao.Upsert(c, &domain.Interaction{}, update, create); err != nil {
		return err
	}
	go func() {
		if err := repo.CacheIncrCnt(c, biz, id, "collect_cnt"); err != nil {
			slog.Warn("Redis操作失败 CacheIncrCnt", "biz", biz, "bizID", id)
		}
	}()
	// 数据库操作成功即认为业务处理成功
	return nil
}

func (repo *interactionRepository) CacheIncrCnt(c context.Context, biz string, id int64, cntType string) error {
	_, err := repo.cache.Lua(c, domain.LuaInteractionIncrCnt, []string{key(biz, id)}, cntType, "1")
	return err
}

func key(biz string, bizID int64) string {
	return fmt.Sprintf("interaction:%s:%d", biz, bizID)
}
