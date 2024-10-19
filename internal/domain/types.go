package domain

import (
	"errors"
)

var (
	ErrSystemError = errors.New("系统错误")
	ErrUnKnowError = errors.New("未知错误")
	ErrBadParams   = errors.New("参数错误")
)

type Model struct {
	ID        uint  `gorm:"primarykey" json:"-"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	DefaultPage           = 0
	DefaultSize           = 10
	DefaultUserPassword   = "root"
	DefaultUserNamePrefix = "user"
)

type BasePageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}
