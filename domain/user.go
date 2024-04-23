package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	UserName  string    `json:"user_name" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password" gorm:"size:256"`
	Avatar    string    `json:"avatar" gorm:"size:1024"`
	Birthday  time.Time `json:"birthday"`
	Telephone string    `json:"telephone" gorm:"size:20"`
	LoginType LoginType `json:"login_type" gorm:"size:20"`
	Role      Role      `json:"role" gorm:"default:2"` //
}

func (User) TableName() string {
	return `user`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id int64) (User, error)
	//UpsertAvatar(c context.Context, avatar string) error
}

type UserUsecase interface {
	GetProfileByID(c context.Context, userID int64) (*Profile, error)
}

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
