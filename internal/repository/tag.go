package repository

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
)

type tagRepository struct {
	dao orm.Database
}

func NewTagRepository(dao orm.Database) domain.TagRepository {
	return &tagRepository{dao: dao}
}

func (t *tagRepository) CreateTag(c context.Context, tag domain.Tag) error {
	return t.dao.Insert(c, &domain.Tag{}, &tag)
}

func (t *tagRepository) CreateTagBiz(c context.Context, userID int64, biz string, bizID int64, tagIDs []int64) error {
	items := make([]domain.TagBiz, 0, len(tagIDs))
	for _, id := range tagIDs {
		items = append(items, domain.TagBiz{
			TagID:  id,
			UserID: userID,
			Biz:    biz,
			BizID:  bizID,
		})
	}
	return t.dao.Insert(c, &domain.Tag{}, &items)
}

func (t *tagRepository) GetTagsByUserID(c context.Context, userID int64) ([]domain.Tag, error) {
	filter := map[string]interface{}{
		"user_id": userID,
	}
	var items []domain.Tag
	err := t.dao.FindMany(c, &domain.Tag{}, filter, &items)
	return items, err
}

func (t *tagRepository) GetTagsByBiz(c context.Context, userID int64, biz string, bizID int64) ([]domain.Tag, error) {
	filter := map[string]interface{}{
		"user_id": userID,
		"biz":     biz,
		"biz_id":  bizID,
	}
	var items []domain.Tag
	err := t.dao.FindMany(c, &domain.Tag{}, filter, &items)
	return items, err
}
