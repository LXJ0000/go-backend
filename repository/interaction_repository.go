package repository

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
	"golang.org/x/net/context"
	"gorm.io/gorm"
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
	return repo.dao.Upsert(c, &domain.Interaction{}, update, create)
}
