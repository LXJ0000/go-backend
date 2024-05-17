package domain

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primarykey" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	DefaultPage = 0
	DefaultSize = 10
)
