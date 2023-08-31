package s3

import (
	"dynamic-user-segmentation-service/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

func Connect(minioConfig config.Minio) (*minio.Client, error) {
	client, err := minio.New(minioConfig.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv(minioConfig.EnvUser), os.Getenv(minioConfig.EnvPassword), ""),
		Secure: minioConfig.Ssl,
	})
	return client, err
}
