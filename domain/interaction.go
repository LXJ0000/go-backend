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
	BizID int64  `gorm:"uniqueIndex:idx_interaction_bizID_biz"`
	Biz   string `gorm:"uniqueIndex:idx_interaction_bizID_biz"`
	//Biz     string `gorm:"type:varchar(255);uniqueIndex:idx_bizID_biz"` // MYSQL 写法
	ReadCnt    int
	LikeCnt    int
	CollectCnt int // 3个cnt 相比较 type+cnt 在读性能友好, 每次只需要读一行
}

type CacheInteractionKey struct {
	ReadCnt    int
	LikeCnt    int
	CollectCnt int
}

func (Interaction) TableName() string {
	return `interaction`
}

type InteractionUseCase interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	Like(c context.Context, biz string, bizID, userID int64) error
	CancelLike(c context.Context, biz string, bizID, userID int64) error
	Info(c context.Context, biz string, bizID, userID int64) (Interaction, UserInteractionInfo, error)
}

type InteractionRepository interface {
	IncrReadCount(c context.Context, biz string, id int64) error
	Like(c context.Context, biz string, bizID, userID int64) error
	CancelLike(c context.Context, biz string, bizID, userID int64) error
	Info(c context.Context, biz string, bizID, userID int64) (Interaction, UserInteractionInfo, error)
}

type UserLike struct {
	gorm.Model
	UserID int64  `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz"`
	BizID  int64  `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz"`
	Biz    string `gorm:"uniqueIndex:idx_userLike_userID_bizID_biz"`
	Status bool   // true 点赞 false 取消点赞
	// `gorm:"uniqueIndex:idx_userID_bizID_biz"`
	//	具体索引顺序，需要根据业务需求规定，此外还需根据字段区分度
	//1. 查询用户喜欢的东西 select * from user_like where user_id = ? and biz = ?
	//2. 查询某个东西的点赞数 select * from user_like where bizID = ? and biz = ?
}

func (UserLike) TableName() string {
	return `user_like`
}

type UserCollect struct {
	gorm.Model
	UserID int64  `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz"`
	BizID  int64  `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz"`
	Biz    string `gorm:"uniqueIndex:idx_userCollect_userID_bizID_biz"`
	Status bool
}

func (UserCollect) TableName() string {
	return `user_collect`
}

type UserInteractionInfo struct {
	Liked     bool
	Collected bool
}
