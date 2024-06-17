package repository

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
)

type feedRepository struct {
	dao orm.Database
}

func NewFeedRepository(dao orm.Database) domain.FeedRepository {
	return &feedRepository{dao: dao}
}

func (r *feedRepository) CreatePush(c context.Context, feed ...domain.Feed) error {
	return r.dao.Insert(c, &domain.Feed{}, &feed)
}

func (r *feedRepository) CreatePull(c context.Context, feed ...domain.Feed) error {
	return r.dao.Insert(c, &domain.Feed{}, &feed)
}

func (r *feedRepository) FindPush(c context.Context, userIDs []int64, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil// TODO
}

func (r *feedRepository) FindPull(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil // TODO
}
