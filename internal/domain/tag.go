package domain

import "context"

type Tag struct {
	Model
	TagID   int64 `gorm:"primaryKey"`
	UserID  int64 `gorm:"index"`
	TagName string
}

func (Tag) TableName() string {
	return `tag`
}

// SELECT * FROM tag_biz WHERE biz = ? and biz_id = ?
type TagBiz struct {
	Model
	Biz    string `gorm:"index:idx_tagBiz_biz_bizId"`
	BizID  int64  `gorm:"index:idx_tagBiz_biz_bizId"`
	UserID int64  `gorm:"index"`
	TagID  int64  `gorm:"index"`
}

func (TagBiz) TableName() string {
	return `tag_biz`
}

type TagRepository interface {
	CreateTag(c context.Context, tag Tag) error                                                  // 用户创建Tag
	CreateTagBiz(c context.Context, userID int64, biz string, bizID int64, tagIDs []int64) error // 用户为某个资源加Tag
	GetTagsByUserID(c context.Context, userID int64) ([]Tag, error)                              // 查询用户所有Tag
	GetTagsByBiz(c context.Context, userID int64, biz string, bizID int64) ([]Tag, error)        // 查询用户在某个资源上打的的Tag
}

type TagUsecase interface {
	CreateTag(c context.Context, tag Tag) error
	CreateTagBiz(c context.Context, userID int64, biz string, bizID int64, tagIDs []int64) error
	GetTagsByUserID(c context.Context, userID int64) ([]Tag, error)
	GetTagsByBiz(c context.Context, userID int64, biz string, bizID int64) ([]Tag, error)
}

type CreateTagBizRequest struct {
	Biz    string  `json:"biz" form:"biz" binding:"required"`
	BizID  int64   `json:"biz_id" form:"biz_id" binding:"required"`
	TagIDs []int64 `json:"tag_ids" form:"tag_ids" binding:"required"`
}

type GetTagsByBizRequest struct {
	Biz   string `json:"biz" form:"biz" binding:"required"`
	BizID int64  `json:"biz_id" form:"biz_id" binding:"required"`
}
