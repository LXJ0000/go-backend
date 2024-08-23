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
	items := make([]domain.FeedPush, 0, len(feed))
	for _, f := range feed {
		items = append(items, convertToPush(f))
	}
	return r.dao.Insert(c, &domain.FeedPush{}, &items)
}

func (r *feedRepository) CreatePull(c context.Context, feed ...domain.Feed) error {
	items := make([]domain.FeedPull, 0, len(feed))
	for _, f := range feed {
		items = append(items, convertToPull(f))
	}
	return r.dao.Insert(c, &domain.FeedPull{}, &items)
}

func (r *feedRepository) FindPull(c context.Context, userIDs []int64, timestamp, limit int64) ([]domain.Feed, error) {
	var items []domain.FeedPull
	db := r.dao.Raw(c)
	if err := db.Model(&domain.FeedPull{}).
		Where("user_id in (?) and timestamp > ?", userIDs, timestamp).
		Order("timestamp desc").
		Limit(int(limit)).Where(&items).Error; err != nil {
		return nil, err
	}
	feeds := make([]domain.Feed, 0, len(items))
	for _, item := range items {
		feeds = append(feeds, convertToFeedFromPull(item))
	}
	return feeds, nil
}

func (r *feedRepository) FindPush(c context.Context, userID, timestamp, limit int64) ([]domain.Feed, error) {
	var items []domain.FeedPush
	db := r.dao.Raw(c)
	if err := db.Model(&domain.FeedPush{}).
		Where("user_id = ? and timestamp > ?", userID, timestamp).
		Order("timestamp desc").
		Limit(int(limit)).Where(&items).Error; err != nil {
		return nil, err
	}
	feeds := make([]domain.Feed, 0, len(items))
	for _, item := range items {
		feeds = append(feeds, convertToFeedFromPush(item))
	}
	return feeds, nil
}

func (r *feedRepository) FindPullWithType(c context.Context, event string, userIDs []int64, timestamp, limit int64) ([]domain.Feed, error) {
	var items []domain.FeedPull
	db := r.dao.Raw(c)
	if err := db.Model(&domain.FeedPull{}).
		Where("user_id in (?) and timestamp > ? and type = ?", userIDs, timestamp, event).
		Order("timestamp desc").
		Limit(int(limit)).Where(&items).Error; err != nil {
		return nil, err
	}
	feeds := make([]domain.Feed, 0, len(items))
	for _, item := range items {
		feeds = append(feeds, convertToFeedFromPull(item))
	}
	return feeds, nil
}
func (r *feedRepository) FindPushWithType(c context.Context, event string, userID, timestamp, limit int64) ([]domain.Feed, error) {
	var items []domain.FeedPush
	db := r.dao.Raw(c)
	if err := db.Model(&domain.FeedPush{}).
		Where("user_id = ? and timestamp > ? and type = ?", userID, timestamp, event).
		Order("timestamp desc").
		Limit(int(limit)).Where(&items).Error; err != nil {
		return nil, err
	}
	feeds := make([]domain.Feed, 0, len(items))
	for _, item := range items {
		feeds = append(feeds, convertToFeedFromPush(item))
	}
	return feeds, nil
}

func convertToPush(feed domain.Feed) domain.FeedPush {
	content, _ := json.Marshal(feed.Content)
	return domain.FeedPush{
		UserID:  feed.UserID,
		Type:    feed.Type,
		Content: string(content),
		//CreatedAt: feed.CreatedAt.UnixMicro(),
	}
}

func convertToPull(feed domain.Feed) domain.FeedPull {
	content, _ := json.Marshal(feed.Content)
	return domain.FeedPull{
		UserID:  feed.UserID,
		Type:    feed.Type,
		Content: string(content),
		//CreatedAt: feed.CreatedAt.UnixMicro(),
	}
}

func convertToFeedFromPull(pull domain.FeedPull) domain.Feed {
	var content domain.FeedContent
	_ = json.Unmarshal([]byte(pull.Content), &content)
	return domain.Feed{
		UserID:  pull.UserID,
		Type:    pull.Type,
		Content: content,
		//CreatedAt: time.UnixMicro(pull.CreatedAt),
	}
}

func convertToFeedFromPush(push domain.FeedPush) domain.Feed {
	var content domain.FeedContent
	_ = json.Unmarshal([]byte(push.Content), &content)
	return domain.Feed{
		UserID:  push.UserID,
		Type:    push.Type,
		Content: content,
		//CreatedAt: time.UnixMicro(push.CreatedAt),
	}
}
