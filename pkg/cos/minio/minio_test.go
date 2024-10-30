package minio

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestMinio(t *testing.T) {
	accessKey := "xxx"
	accessSecret := "xxx"
	endPoint := "xx.xx.xx.xx:xx"
	useSSL := false

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	storage := NewMinio(minioClient)
	// t.Error(storage.ListBuckets(context.Background()))
	// t.Error(storage.ListFiles(context.Background(), "openim"))
	// file, err := os.Open("minio.go")
	// if err != nil {
	// 	panic(err)
	// }
	// fileContent, err := io.ReadAll(file)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := storage.UploadFile(context.Background(), "openim", "minio.go", fileContent); err != nil {
	// 	panic(err)
	// }
	// list, err := storage.ListFiles(context.Background(), "openim")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Error(list)

	if err := storage.DownloadFile(context.Background(), "openim", "minio.go", "miniov1.go"); err != nil {
		panic(err)
	}
}
