package domain

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   int64  `json:"user_id" gorm:"primaryKey"`
	UserName string `json:"user_name" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"size:256"`
}

func (User) TableName() string {
	return `user`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id int64) (User, error)
}
