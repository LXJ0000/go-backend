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

func (File) TableName() string {
	return `file`
}

type FileUsecase interface {
	Upload(c context.Context, file *multipart.FileHeader) (File, error) // 文件上传
	Uploads(c context.Context, files []*multipart.FileHeader) (FileUploadsResponse, error)
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, int, error)
}

type FileRepository interface {
	Upload(c context.Context, file File) error
	Uploads(c context.Context, files []File) error
	FileList(c context.Context, fileType, fileSource string, page, size int) ([]File, int, error)
	FindByHash(c context.Context, hash string) (File, error)
}

type FileListRequest struct {
	BasePageRequest
	Type   string `json:"type" from:"type"`
	Source string `json:"source" from:"source"`
}

type FileUploadsResponse struct {
	Data map[string]interface{} `json:"data"` // 文件上传状态
}
