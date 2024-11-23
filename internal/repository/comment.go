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
	return repo.dao.Insert(c, &domain.Comment{}, comment)
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

func (repo *commentRepository) FindTop(c context.Context, biz string, bizID, minID int64, limit int) ([]domain.Comment, error) {
	if limit <= 0 {
		limit = 10
	}
	if minID <= 0 {
		minID = math.MaxInt
	}
	db := repo.dao.Raw(c)
	var res []domain.Comment
	err := db.
		Where("biz = ? AND biz_id = ? AND id < ? AND parent_comment_id = 0", biz, bizID, minID). // 一级评论 则 parent_id = 0
		Limit(limit).
		Order("id DESC").
		Find(&res).Error
	return res, err
}
