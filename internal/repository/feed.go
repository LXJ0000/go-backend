package repository

import (
	"encoding/json"

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
	return r.dao.Insert(c, &domain.FeedPush{}, &feed)
}

func (r *feedRepository) CreatePull(c context.Context, feed ...domain.Feed) error {
	return r.dao.Insert(c, &domain.FeedPull{}, &feed)
}

func (r *feedRepository) FindPush(c context.Context, userIDs []int64, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil // TODO
}

func (r *feedRepository) FindPull(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	return nil, nil // TODO
}

func converToPush(feed domain.Feed) domain.FeedPush {
	content, _ := json.Marshal(feed.Content)
	return domain.FeedPush{
		UserID: feed.UserID,
		Type: feed.Type,
		Con
	}
}

func converToPull(feed domain.Feed) domain.FeedPull {
	return domain.FeedPull{}
}

func converToFeedFromPull(pull domain.FeedPull) domain.Feed {
	return domain.Feed{}
}

func converToFeedFromPush(push domain.FeedPush) domain.Feed {
	return domain.Feed{}
}
