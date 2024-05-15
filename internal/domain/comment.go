package domain

import (
	"database/sql"

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
