package minio

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestMinio(t *testing.T) {
	godotenv.Load("../../../.env")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	accessSecret := os.Getenv("MINIO_ACCESS_SECRET")
	endPoint := os.Getenv("MINIO_EXTERNAL_ADDRESS")
	useSSL := false

	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	storage := NewMinio(minioClient)
	t.Error(storage.ListBuckets(context.Background()))

	tool := NewMinio(minioClient)
	url, err := tool.GetFilePath(context.Background(), "go-backend", "30a057aabeb4d4bd0009e1fb217df42.png")
	t.Error(url, err)
	// err = tool.CreateBucket(context.Background(), "go-backend")
	// t.Error(err)
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
	// list, err := storage.ListFiles(context.Background(), "go-backend")
	// if err != nil {
	// 	panic(err)
	// }
	// t.Error(list)

	// if err := storage.DownloadFile(context.Background(), "openim", "minio.go", "miniov1.go"); err != nil {
	// 	panic(err)
	// }

}
