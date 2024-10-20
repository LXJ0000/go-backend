package domain

import (
	"database/sql"

	"golang.org/x/net/context"
)

type Comment struct {
	Model
	CommentID int64 `gorm:"primaryKey"`
	UserID    int64 `json:"user_id,string"`

	Biz   string `gorm:"index:idx_biz_biz_id" binding:"required" json:"biz"`
	BizID int64  `gorm:"index:idx_biz_biz_id" binding:"required" json:"biz_id,string"`

	RootID        sql.NullInt64 `json:"root_id,string" gorm:"index"`
	ParentID      sql.NullInt64 `json:"parent_id,string" gorm:"index"`
	ParentComment *Comment      `json:"parent_comment" gorm:"foreignKey:ParentID;AssociationForeignKey:CommentID;constraint:OnDelete:CASCADE"` // 简化删除操作
	Content       string        `json:"content" binding:"required"`
}

func (Comment) TableName() string {
	return `comment`
}

type CommentRepository interface {
	Create(c context.Context, comment Comment) error
	Delete(c context.Context, id int64) error                                                // 删除本节点和其对应的子节点
	FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]Comment, error) // 查找一级评论
}

type CommentUsecase interface {
	Create(c context.Context, comment Comment) error
	Delete(c context.Context, id int64) error
	FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]Comment, error)
}

type CommentCreateRequest struct {
	Biz      string `json:"biz" binding:"required"`
	BizID    int64  `json:"biz_id,string" binding:"required"`
	RootID   int64  `json:"root_id,string"`
	ParentID int64  `json:"parent_id,string"`
	Content  string `json:"content" binding:"required"`
}

type CommentListRequest struct {
	Biz   string `json:"biz" form:"biz"`
	BizID int64  `json:"biz_id,string" form:"biz_id"`
	MinID int64  `json:"min_id,string" form:"min_id"`
	Limit int    `json:"limit" form:"limit"`
}
