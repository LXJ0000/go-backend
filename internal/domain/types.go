package domain

type Model struct {
	ID        uint `gorm:"primarykey" json:"-"`
	CreatedAt int64
	UpdatedAt int64
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	DefaultPage = 0
	DefaultSize = 10
)

type BasePageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}
