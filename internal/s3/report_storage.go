package s3

import (
	"context"
	"github.com/minio/minio-go/v7"
)

type (
	ReportStorage interface {
		UploadReport(ctx context.Context, reportName string, reportPath string) (string, error)
	}

	reportStorage struct {
		Minio *minio.Client
	}
)

const (
	bucketName = "reports"
	policy     = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":["arn:aws:s3:::` + bucketName + `"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
)

func NewReportStorage(minioClient *minio.Client) ReportStorage {
	ctx := context.Background()
	exists, _ := minioClient.BucketExists(ctx, bucketName)
	if !exists {
		_ = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		_ = minioClient.SetBucketPolicy(ctx, bucketName, policy)
	}
	return &reportStorage{Minio: minioClient}
}

func (r *reportStorage) UploadReport(ctx context.Context, reportName string, reportPath string) (string, error) {
	_, err := r.Minio.FPutObject(ctx, bucketName, reportName, reportPath, minio.PutObjectOptions{ContentType: "text/csv"})
	return bucketName + "/" + reportName, err
}
