package cos

import (
	"context"
)

type FileStorage interface {
	CreateBucket(ctx context.Context, bucketName string) error
	DeleteBucket(ctx context.Context, bucketName string) error
	ListBuckets(ctx context.Context) ([]string, error)
	UploadFile(ctx context.Context, bucketName, objectName string, file []byte) error
	DownloadFile(ctx context.Context, bucketName, objectName string, filePath string) error
	DeleteFile(ctx context.Context, bucketName, objectName string) error
	ListFiles(ctx context.Context, bucketName string) ([]string, error)
}
