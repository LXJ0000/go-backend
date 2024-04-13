package domain

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

const (
	BizPost = "post"
)

type Interaction struct {
	gorm.Model
	// idx_biz select * from . where biz ==
	// idx_bizID_biz 联合索引 (bizID区分度高)
	BizID int64 `gorm:"uniqueIndex:idx_bizID_biz"`
	//Biz     string `gorm:"uniqueIndex:idx_bizID_biz"`
	Biz     string `gorm:"type:varchar(255);uniqueIndex:idx_bizID_biz"`
	ReadCnt int64
}

func (Interaction) TableName() string {
	return `interaction`
}

type InteractionUseCase interface {
	IncrReadCount(c context.Context, biz string, id int64) error
}

type InteractionRepository interface {
	IncrReadCount(c context.Context, biz string, id int64) error
}
