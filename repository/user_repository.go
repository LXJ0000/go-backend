package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/orm"

	"github.com/LXJ0000/go-backend/domain"
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

func (repo *userRepository) Create(c context.Context, user *domain.User) error {
	_, err := repo.dao.InsertOne(c, &domain.User{}, user)
	return err
}

func (repo *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindOne(c, &domain.User{}, &domain.User{Email: email})
	if err != nil {
		return domain.User{}, err
	}
	return *user.(*domain.User), nil
}

func (repo *userRepository) GetByID(c context.Context, id int64) (domain.User, error) {
	user, err := repo.dao.FindOne(c, &domain.User{}, &domain.User{UserID: id})
	if err != nil {
		return domain.User{}, err
	}
	return *user.(*domain.User), nil
}
