package minio

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/LXJ0000/go-backend/pkg/file"
	"github.com/minio/minio-go/v7"
)

type Minio struct {
	client *minio.Client
}

func NewMinio(client *minio.Client) file.FileStorage {
	return &Minio{client: client}
}

func (m *Minio) CreateBucket(ctx context.Context, bucketName string) error { // TODO MakeBucketOptions
	return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

func (m *Minio) DeleteBucket(ctx context.Context, bucketName string) error {
	return m.client.RemoveBucket(ctx, bucketName)
}

func (m *Minio) ListBuckets(ctx context.Context) ([]string, error) {
	bucketInfos, err := m.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	var bucketNames []string
	for _, bucketInfo := range bucketInfos {
		bucketNames = append(bucketNames, bucketInfo.Name)
	}
	return bucketNames, nil
}

func (m *Minio) UploadFile(ctx context.Context, bucketName, objectName string, file []byte) error { // TODO origin_name
	m.client.PutObject(ctx,
		bucketName, objectName, bytes.NewReader(file), int64(len(file)),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	return nil
}

func (m *Minio) DownloadFile(ctx context.Context, bucketName, objectName string, filePath string) error {
	object, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, buf)
	return err
}

func (m *Minio) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *Minio) ListFiles(ctx context.Context, bucketName string) ([]string, error) {
	var objectNames []string
	objectCh := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objectNames = append(objectNames, object.Key)
	}
	return objectNames, nil
}
