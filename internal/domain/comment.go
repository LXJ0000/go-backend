package domain

import (
	"golang.org/x/net/context"
)

// 帖子
// 一级评论 某个用户回复该帖子 commentid userid content
// 二级评论 在某个一级评论下回复 commentid userid content parentid(一级评论的commentid)
// 二级评论 在某个一级评论下回复某人 commentid userid content parentid(一级评论的commentid) to_userid
type Comment struct {
	Model
	CommentID       int64  `gorm:"primaryKey" json:"comment_id,string"`
	ParentCommentID int64  `gorm:"index" json:"parent_comment_id,string"` // 一级评论的 parent_comment_id = 0 其他评论的 parent_comment_id = 一级评论的 comment_id
	Content         string `gorm:"not null" json:"content"`
	UserID          int64  `gorm:"not null" json:"user_id,string"`
	ToUserID        int64  `json:"to_user_id,string"` // 回复某人的时候，to_user_id = 被回复的用户的 user_id

	Biz   string `gorm:"index:idx_biz_biz_id" binding:"required" json:"biz"`
	BizID int64  `gorm:"index:idx_biz_biz_id" binding:"required" json:"biz_id,string"`

	LikeCount  int `json:"like_count"`  // 点赞数
	ReplyCount int `json:"reply_count"` // 回复数
}

func (Comment) TableName() string {
	return `comment`
}

type CommentRepository interface {
	Create(c context.Context, comment *Comment) error
	Delete(c context.Context, id int64) error                                                // 删除本节点和其对应的子节点
	FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]Comment, error) // 查找一级评论
}

type CommentUsecase interface {
	Create(c context.Context, comment *Comment) error
	Delete(c context.Context, id int64) error
	FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]Comment, error)
}

type CommentCreateRequest struct {
	Biz     string `json:"biz" binding:"required"`
	BizID   string `json:"biz_id" binding:"required"`
	Content string `json:"content" binding:"required"`

	ParentID string `json:"parent_id"`  // default 0 表示回复的是帖子
	ToUserID string `json:"to_user_id"` // default 0 表示回复的是帖子 or 一级评论
}

type CommentListRequest struct {
	Biz   string `json:"biz" form:"biz"`
	BizID string `json:"biz_id,string" form:"biz_id"`
	MinID string `json:"min_id,string" form:"min_id"`
	Limit int    `json:"limit" form:"limit"`
}
