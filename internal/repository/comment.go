package repository

import (
	"math"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type commentRepository struct {
	dao orm.Database
}

func NewCommentRepository(dao orm.Database) domain.CommentRepository {
	return &commentRepository{dao: dao}
}

func (repo *commentRepository) Create(c context.Context, comment *domain.Comment) error {
	dao := repo.dao.Raw(c)
	return dao.Transaction(func(tx *gorm.DB) error {
		// 一级评论
		if comment.ParentCommentID == 0 {
			return tx.Create(comment).Error
		}
		// 二级评论
		// 查询父评论
		var parent domain.Comment
		if err := tx.Model(&domain.Comment{}).Where("comment_id = ?", comment.ParentCommentID).First(&parent).Error; err != nil {
			return err
		}
		// 更新父评论的回复数
		if err := tx.Model(&domain.Comment{}).Where("comment_id = ?", comment.ParentCommentID).Update("reply_count", gorm.Expr("reply_count + 1")).Error; err != nil {
			return err
		}
		// 创建二级评论
		return tx.Create(comment).Error
	})
}

func (repo *commentRepository) Delete(c context.Context, id int64) error {
	dao := repo.dao.Raw(c)
	return dao.Transaction(func(tx *gorm.DB) error {
		// 查询所有子节点
		var ids []int64
		if err := tx.Model(&domain.Comment{}).
			Where("comment_id = ? OR parent_comment_id = ?", id, id).Pluck("comment_id", &ids).Error; err != nil {
			return err
		}
		// 删除所有子节点和本节点
		return tx.Where("comment_id IN (?)", ids).Delete(&domain.Comment{}).Error
	})
}

func (repo *commentRepository) Find(c context.Context, biz string, bizID, parentID, minID int64, limit int) ([]domain.Comment, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if minID <= 0 {
		minID = math.MaxInt
	}
	db := repo.dao.Raw(c)
	var (
		res   []domain.Comment
		count int64
	)
	query := db.Model(&domain.Comment{}).Where("biz = ? AND biz_id = ? AND parent_comment_id = ?", biz, bizID, parentID)
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Where("comment_id < ?", minID).Limit(limit).Order("id DESC").Find(&res).Error // 一级评论 则 parent_id = 0
	return res, int(count), err
}

func (repo *commentRepository) Count(c context.Context, biz string, bizID int64) (int, error) {
	db := repo.dao.Raw(c)
	var count int64
	err := db.Model(&domain.Comment{}).Where("biz = ? AND biz_id = ? AND parent_comment_id = 0", biz, bizID).Count(&count).Error
	return int(count), err
}
