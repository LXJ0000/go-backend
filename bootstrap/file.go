package bootstrap

import (
	"os"

	"github.com/LXJ0000/go-backend/pkg/file"
	minio2 "github.com/LXJ0000/go-backend/pkg/file/minio"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinio() file.FileStorage {
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

	return minio2.NewMinio(minioClient)
}
