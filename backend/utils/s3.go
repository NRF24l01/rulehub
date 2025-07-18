package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/google/uuid"
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

func FullGeneratePresignedURL(client *minio.Client, bucketName string, expires time.Duration) (string, string, error) {
	var uniqueID string
	for {
		uuidObj := uuid.New()
		uniqueID = uuidObj.String()
		_, err := client.StatObject(context.Background(), bucketName, uniqueID, minio.StatObjectOptions{})
		if minio.ToErrorResponse(err).Code == "NoSuchKey" || minio.ToErrorResponse(err).Code == "NotFound" {
			break
		}
		if err != nil && minio.ToErrorResponse(err).Code != "NoSuchKey" && minio.ToErrorResponse(err).Code != "NotFound" {
			return "", "", fmt.Errorf("error checking object existence: %w", err)
		}
	}

	// Generate presigned URL for the unique UUID object
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, uniqueID, expires, nil)
	if err != nil {
		return "", "", fmt.Errorf("error generating presigned URL: %w", err)
	}
	return uniqueID, presignedURL.String(), nil
}