package domain

import (
	"mime/multipart"

	"golang.org/x/net/context"
)

const (
	FileTypeUnknown = "unknown"
	FileTypeImage   = "image"
)

const (
	FileSourceKnown = "unknown"
	FileSourceLocal = "local"
)

type File struct {
	Model
	FileID int64 `gorm:"primaryKey"`
	Name   string
	Path   string
	Type   string
	Source string
	Hash   string `gorm:"unique"`
}

type FileUsecase interface {
	Upload(c context.Context, file *multipart.FileHeader) (File, error) // 文件上传
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, int, error)
	// TODO 多文件上传
}

type FileRepository interface {
	Upload(c context.Context, file File) error
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, int, error)
	FindByHash(c context.Context, hash string) (File, error)
}

type FileListRequest struct {
	BasePageRequest
	Type   string `json:"type" from:"type"`
	Source string `json:"source" from:"source"`
}
