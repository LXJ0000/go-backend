package domain

import (
	"database/sql"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CommentID int64 `gorm:"primaryKey"`
	UserID    int64

	Biz   string `gorm:"index:idx_biz_biz_id"`
	BizID int64  `gorm:"index:idx_biz_biz_id"`

	RootID        int64         `gorm:"index"`
	ParentID      sql.NullInt64 `gorm:"index"`
	ParentComment *Comment      `gorm:"foreignKey:ParentCommentID;AssociationForeignKey:CommentID;constraint:OnDelete:CASCADE"` // 简化删除操作
	Content       string
}

func (Comment) TableName() string {
	return `comment`
}

type CommentRepository interface {
	Create(c context.Context, comment Comment) error
	Delete(c context.Context, id int64) error                                                  // 删除本节点和其对应的子节点
	FindByBiz(c context.Context, biz string, bizID, minID int64, limit int) ([]Comment, error) // 查找一级评论
}

type CommentUsecase interface {
}
