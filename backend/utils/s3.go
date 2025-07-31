package utils

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func GeneratePresignedPutURL(client *minio.Client, bucketName string, expires time.Duration) (string, string, error) {
    uniqueID := uuid.New().String()

    presignedURL, err := client.PresignedPutObject(context.Background(), bucketName, uniqueID, expires)
    if err != nil {
        return "", "", fmt.Errorf("error generating presigned PUT URL: %w", err)
    }

    urlStr := presignedURL.String()
    urlStr, err = replaceHostWithBaseURL(urlStr)
    if err != nil {
        return "", "", fmt.Errorf("error replacing host with base URL: %w", err)
    }

    return uniqueID, urlStr, nil
}

func GeneratePresignedGetURL(client *minio.Client, bucketName string, objectKey string, expires time.Duration) (string, error) {
    // Проверка существования объекта
    _, err := client.StatObject(context.Background(), bucketName, objectKey, minio.StatObjectOptions{})
    if err != nil {
        errResp := minio.ToErrorResponse(err)
        if errResp.Code == "NoSuchKey" || errResp.Code == "NotFound" {
            return "", fmt.Errorf("object %s does not exist in bucket %s", objectKey, bucketName)
        }
        return "", fmt.Errorf("error checking object existence: %w", err)
    }

    presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectKey, expires, nil)
    if err != nil {
        return "", fmt.Errorf("error generating presigned GET URL: %w", err)
    }

    urlStr := presignedURL.String()
    urlStr, err = replaceHostWithBaseURL(urlStr)
    if err != nil {
        return "", fmt.Errorf("error replacing host with base URL: %w", err)
    }

    return urlStr, nil
}

// replaceHostWithBaseURL replaces the host part of a URL with the S3_BASE_URL if set
func replaceHostWithBaseURL(originalURL string) (string, error) {
    baseURL := os.Getenv("S3_BASE_URL")
    if baseURL == "" {
        return originalURL, nil // Return original URL if S3_BASE_URL is not set
    }

    u, err := url.Parse(originalURL)
    if err != nil {
        return "", fmt.Errorf("error parsing URL: %w", err)
    }

    baseU, err := url.Parse(baseURL)
    if err != nil {
        return "", fmt.Errorf("error parsing base URL: %w", err)
    }

    // Replace the scheme, host, and port with those from the base URL
    u.Scheme = baseU.Scheme
    u.Host = baseU.Host

    // Preserve the path from the base URL if it exists and append the original path
    if baseU.Path != "" && baseU.Path != "/" {
        // Ensure base path doesn't end with slash and original path doesn't start with slash
        basePath := baseU.Path
        if basePath[len(basePath)-1] == '/' {
            basePath = basePath[:len(basePath)-1]
        }
        
        origPath := u.Path
        if len(origPath) > 0 && origPath[0] == '/' {
            origPath = origPath[1:]
        }
        
        u.Path = basePath + "/" + origPath
    }

    return u.String(), nil
}

func GetPresignedLifetime() time.Duration {
	secStr := os.Getenv("S3_PRESIGNED_LIFETIME")
	if secStr == "" {
		return time.Hour // дефолт 1 час
	}
	sec, err := strconv.Atoi(secStr)
	if err != nil || sec <= 0 {
		return time.Hour
	}
	return time.Duration(sec) * time.Second
}