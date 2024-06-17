package repository

import (
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"golang.org/x/net/context"
)

type relationRepository struct {
	dao orm.Database
}

func NewRelationRepository(dao orm.Database) domain.RelationRepository {
	return &relationRepository{dao: dao}
}

func (uc *relationRepository) Follow(c context.Context, follower, followee int64) error {
	return uc.doFollow(c, follower, followee, domain.Follow)
}

func (uc *relationRepository) CancelFollow(c context.Context, follower, followee int64) error {
	return uc.doFollow(c, follower, followee, domain.UnFollow)
}

func (uc *relationRepository) GetFollower(c context.Context, userID int64, page, size int) ([]domain.Relation, error) {
	filter := map[string]interface{}{
		"followee": userID,
		"status":   domain.Follow,
	}
	var items []domain.Relation
	db := uc.dao.WithPage(page, size)
	err := db.WithContext(c).Model(&domain.Relation{}).
	Where(filter).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (uc *relationRepository) GetFollowee(c context.Context, userID int64, page, size int) ([]domain.Relation, error) {
	filter := map[string]interface{}{
		"follower": userID,
		"status":   domain.Follow,
	}
	var items []domain.Relation
	db := uc.dao.WithPage(page, size)
	err := db.WithContext(c).Model(&domain.Relation{}).
	Where(filter).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (uc *relationRepository) Detail(c context.Context, follower, followee int64) (domain.Relation, error) {
	filter := map[string]interface{}{
		"follower": follower,
		"followee": followee,
		"status":   domain.Follow,
	}
	var item domain.Relation
	err := uc.dao.FindOne(c, &domain.Relation{}, filter, item)
	return item, err
}

func (uc *relationRepository) FollowerCnt(c context.Context, userID int64) (int64, error) {
	filter := map[string]interface{}{
		"followee": userID,
		"status":   domain.Follow,
	}
	return uc.dao.Count(c, &domain.Relation{}, filter)
}

func (uc *relationRepository) FolloweeCnt(c context.Context, userID int64) (int64, error) {
	filter := map[string]interface{}{
		"follower": userID,
		"status":   domain.Follow,
	}
	return uc.dao.Count(c, &domain.Relation{}, filter)
}

func (uc *relationRepository) doFollow(c context.Context, follower, followee int64, opt bool) error {
	update := map[string]interface{}{
		"status":     opt,
		"updated_at": time.Now(),
	}
	create := &domain.Relation{
		Followee:   followee,
		Follower:   follower,
		Status:     opt,
		RelationID: snowflakeutil.GenID(),
	}
	return uc.dao.UpsertOne(c, &domain.Relation{}, update, create)
}
