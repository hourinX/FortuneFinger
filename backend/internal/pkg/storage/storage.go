package storage

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Client interface {
	PutObject(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error
	DeleteObject(ctx context.Context, objectName string) error
	PresignGetObject(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}

func Init(ctx context.Context, conf Config) (Client, error) {
	if err := conf.validate(); err != nil {
		return nil, err
	}

	switch conf.Driver {
	case DriverMinIO:
		return newMinIOClient(ctx, conf)
	case DriverOSS:
		return newOSSClient(ctx, conf)
	default:
		return nil, fmt.Errorf("storage: 不支持的存储驱动「%s」", conf.Driver)
	}
}
