package repository

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
)

type commentRepository struct {
	dao orm.Database
}

func NewCommentRepository(dao orm.Database) domain.CommentRepository {
	return &commentRepository{dao: dao}
}

func (repo *commentRepository) Create(c context.Context, comment domain.Comment) error {
	return repo.dao.InsertOne(c, &domain.Comment{}, &comment)
}

func (repo *commentRepository) Delete(c context.Context, id int64) error {
	return repo.dao.DeleteOne(c, &domain.Comment{}, id)
}

func (repo *commentRepository) FindByBiz(c context.Context, biz string, bizID, minID int64, limit int) ([]domain.Comment, error) {
	db := repo.dao.Raw(c)
	var res []domain.Comment
	err := db.WithContext(c).
		Where("biz = ? AND biz_id = ? AND id < ? AND parend_id IS NULL", biz, bizID, minID). // 一级评论 则 parent_id is null
		Limit(limit).
		// Order("id asc").
		Find(&res).Error
	return res, err
}
