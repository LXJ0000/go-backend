package repository

import (
	"context"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
)

type userRepository struct {
	dao orm.Database
	//collection string
}

func NewUserRepository(dao orm.Database) domain.UserRepository {
	return &userRepository{
		dao: dao,
		//collection: collection,
	}
}

func (u *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	var user domain.User
	err := u.dao.FindOne(c, &domain.User{}, map[string]interface{}{"email": email}, &user)
	return user, err
}

func (u *userRepository) Create(c context.Context, user domain.User) error {
	return u.dao.Insert(c, &domain.User{}, &user)
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
