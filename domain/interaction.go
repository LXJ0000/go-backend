package domain

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

const (
	BizPost = "post"
)

const (
	LuaInteractionIncrCnt = "script/redis/interaction_incr_cnt.lua"
)

type Interaction struct {
	gorm.Model
	// idx_biz select * from . where biz ==
	// idx_bizID_biz 联合索引 (bizID区分度高)
	BizID int64  `gorm:"uniqueIndex:idx_bizID_biz"`
	Biz   string `gorm:"uniqueIndex:idx_bizID_biz"`
	//Biz     string `gorm:"type:varchar(255);uniqueIndex:idx_bizID_biz"` // MYSQL 写法
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64 // 3个cnt 相比较 type+cnt 在读性能友好, 每次只需要读一行
}

func (Interaction) TableName() string {
	return `interaction`
}

type InteractionUseCase interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	IncrLikeCount(c context.Context, biz string, id int64) error
	IncrCollectCount(c context.Context, biz string, id int64) error
}

type InteractionRepository interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	IncrLikeCount(c context.Context, biz string, id int64) error
	IncrCollectCount(c context.Context, biz string, id int64) error

	CacheIncrCnt(c context.Context, biz string, id int64, cntType string) error
}
