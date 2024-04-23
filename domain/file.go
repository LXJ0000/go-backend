package domain

import (
	"context"
	"gorm.io/gorm"
	"mime/multipart"
)

type File struct {
	gorm.Model
	FileID   int64    `json:"file_id" gorm:"primaryKey"`
	FileType FileType `json:"file_type" form:"file_type" gorm:"not null"`
	Path     string   `json:"path" form:"path" gorm:"not null"`
	Size     int64    `json:"size" form:"size" gorm:"not null"`
}

func (File) TableName() string {
	return `file`
}

type FileType int

const (
	FileImage FileType = iota + 1
	FileVideo
)

type FileUsecase interface {
	UploadOne(c context.Context, file *multipart.FileHeader) (File, error)
	UploadMany(c context.Context, files []*multipart.FileHeader) ([]File, error)
}

type FileRepository interface {
	UploadOne(c context.Context, file *multipart.FileHeader) (File, error)
	UploadMany(c context.Context, files []*multipart.FileHeader) ([]File, error)
}
