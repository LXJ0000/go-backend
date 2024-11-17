package kodo

import (
	"context"

	"github.com/LXJ0000/go-backend/pkg/file"
)

// Kodo TODO
type Kodo struct {
}

func NewKodo() file.FileStorage {
	return &Kodo{}
}

func (m *Kodo) CreateBucket(ctx context.Context, bucketName string) error {
	return nil
}

func (m *Kodo) DeleteBucket(ctx context.Context, bucketName string) error {
	return nil

}

func (m *Kodo) ListBuckets(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (m *Kodo) UploadFile(ctx context.Context, bucketName, objectName string, file []byte) error {
	return nil
}

func (m *Kodo) DownloadFile(ctx context.Context, bucketName, objectName string, filePath string) error {
	return nil
}

func (m *Kodo) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return nil
}

func (m *Kodo) ListFiles(ctx context.Context, bucketName string) ([]string, error) {
	return nil, nil
}

func (m *Kodo) GetFilePath(ctx context.Context, bucketName, objectName string) (string, error) {
	return "", nil
}
