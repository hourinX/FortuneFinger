package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient struct {
	client *minio.Client
	bucket string
}

func newMinIOClient(ctx context.Context, conf Config) (Client, error) {
	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: 创建MinIO客户端失败: %w", err)
	}

	exists, err := client.BucketExists(ctx, conf.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: 检查MinIO存储桶「%s」失败: %w", conf.Bucket, err)
	}
	if !exists {
		return nil, fmt.Errorf("storage: MinIO存储桶「%s」不存在", conf.Bucket)
	}
	return &minioClient{
		client: client,
		bucket: conf.Bucket,
	}, nil
}

func (c *minioClient) PutObject(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	_, err := c.client.PutObject(ctx, c.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("storage: 上传MinIO对象「%s」失败: %w", objectName, err)
	}
	return nil
}

func (c *minioClient) DeleteObject(ctx context.Context, objectName string) error {
	if err := c.client.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("storage: 删除MinIO对象「%s」失败: %w", objectName, err)
	}
	return nil
}

func (c *minioClient) PresignGetObject(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
) (string, error) {
	presignedURL, err := c.client.PresignedGetObject(ctx, c.bucket, objectName, expiry, url.Values{})
	if err != nil {
		return "", fmt.Errorf("storage: 生成MinIO对象「%s」预签名地址失败: %w", objectName, err)
	}
	return presignedURL.String(), nil
}
