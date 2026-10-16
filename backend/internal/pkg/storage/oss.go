package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type ossClient struct {
	client *oss.Client
	bucket string
}

func newOSSClient(ctx context.Context, conf Config) (Client, error) {
	ossConfig := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey)).
		WithRegion(conf.Region).
		WithEndpoint(conf.Endpoint).
		WithDisableSSL(!conf.Secure)
	client := oss.NewClient(ossConfig)

	if _, err := client.GetBucketLocation(ctx, &oss.GetBucketLocationRequest{
		Bucket: oss.Ptr(conf.Bucket),
	}); err != nil {
		return nil, fmt.Errorf("storage: 检查OSS存储桶「%s」失败: %w", conf.Bucket, err)
	}
	return &ossClient{
		client: client,
		bucket: conf.Bucket,
	}, nil
}

func (c *ossClient) PutObject(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	request := &oss.PutObjectRequest{
		Bucket:      oss.Ptr(c.bucket),
		Key:         oss.Ptr(objectName),
		ContentType: oss.Ptr(contentType),
		Body:        reader,
	}
	if size >= 0 {
		request.ContentLength = oss.Ptr(size)
	}
	if _, err := c.client.PutObject(ctx, request); err != nil {
		return fmt.Errorf("storage: 上传OSS对象「%s」失败: %w", objectName, err)
	}
	return nil
}

func (c *ossClient) DeleteObject(ctx context.Context, objectName string) error {
	if _, err := c.client.DeleteObject(ctx, &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(c.bucket),
		Key:    oss.Ptr(objectName),
	}); err != nil {
		return fmt.Errorf("storage: 删除OSS对象「%s」失败: %w", objectName, err)
	}
	return nil
}

func (c *ossClient) PresignGetObject(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
) (string, error) {
	result, err := c.client.Presign(ctx, &oss.GetObjectRequest{
		Bucket: oss.Ptr(c.bucket),
		Key:    oss.Ptr(objectName),
	}, oss.PresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("storage: 生成OSS对象「%s」预签名地址失败: %w", objectName, err)
	}
	return result.URL, nil
}
