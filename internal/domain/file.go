package domain

import (
	"mime/multipart"

	"golang.org/x/net/context"
)

const (
	FileTypeLocal = "image"
)

const (
	FileSourceLocal = "local"
)

type File struct {
	Model
	FileID int64 `gorm:"primaryKey"`
	Name   string
	Path   string
	Type   string
	Source string
}

type FileUsecase interface {
	Upload(c context.Context, file *multipart.FileHeader) (string, error) // 文件上传
	// 多文件上传
	// 文件展示
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, error)
}

type FileReposirory interface {
	Upload(c context.Context, file File) error
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, error)
}
