package repository

import (
	"context"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
)

type userRepository struct {
	dao   orm.Database
	cache cache.RedisCache
	//collection string
}

func NewUserRepository(dao orm.Database, cache cache.RedisCache) domain.UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
		//collection: collection,
	}
}

func (u *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := u.dao.FindOne(c, &domain.User{}, map[string]interface{}{"email": email}, &user)
	return user, err
}

func (u *userRepository) Create(c context.Context, user *domain.User) error {
	return u.dao.Insert(c, &domain.User{}, user)
}

func (u *userRepository) GetByID(c context.Context, id int64) (domain.User, error) {
	var user domain.User
	err := u.dao.FindOne(c, &domain.User{}, map[string]interface{}{"user_id": id}, &user)
	return user, err
}

func (u *userRepository) FindByUserIDs(c context.Context, userIDs []int64, page, size int) ([]domain.User, error) {
	var items []domain.User
	db := u.dao.WithPage(page, size)
	err := db.WithContext(c).Model(&domain.User{}).Where("user_id IN (?)", userIDs).Find(&items).Error
	return items, err
}

func (u *userRepository) InvalidToken(c context.Context, ssid string, exp time.Duration) error {
	return u.cache.Set(c, domain.UserLogoutKey(ssid), "", exp)
}

func (u *userRepository) Update(c context.Context, id int64, user *domain.User) error {
	return u.dao.UpdateOne(c, &domain.User{}, map[string]interface{}{"user_id": id}, user)
}

func (u *userRepository) GetByPhone(c context.Context, phone string) (domain.User, error) {
	var user domain.User
	err := u.dao.FindOne(c, &domain.User{}, map[string]interface{}{"phone": phone}, &user)
	return user, err
}
