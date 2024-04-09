package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func (User) TableName() string {
	return `user`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id string) (User, error)
}

//type User struct {
//	orm.Model
//	UserID   int64  `json:"user_id" orm:"primaryKey"`
//	UserName string `json:"user_name" orm:"unique"`
//	Email    string `json:"email" orm:"unique"`
//	Password string `json:"password" orm:"size:256"`
//}
//
//type UserRepository interface {
//	Create(c context.Context, user *User) error
//	GetByEmail(c context.Context, email string) (User, error)
//	GetByID(c context.Context, id int64) (User, error)
//}
