package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_USERNAME")
	secretKey := os.Getenv("MINIO_PASSWORD")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT, MINIO_USERNAME, and MINIO_PASSWORD must be set")
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GeneratePresignedURL(client *minio.Client, bucketName, objectName string, expires time.Duration) (string, error) {
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}