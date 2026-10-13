package redisc

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// 池大小不设，由 go-redis 按 GOMAXPROCS 推算
// 超时固定，避免行为随依赖版本默认值漂移
const (
	dialTimeout     = 5 * time.Second
	readTimeout     = 3 * time.Second
	writeTimeout    = 3 * time.Second
	connMaxIdleTime = 30 * time.Minute
)

type Client struct {
	client *goredis.Client
}

//go:embed script/get_and_delete.lua
var getAndDeleteLua string

var getAndDeleteScript = goredis.NewScript(getAndDeleteLua)

//go:embed script/replace_if_value.lua
var replaceIfValueLua string

var replaceIfValueScript = goredis.NewScript(replaceIfValueLua)

func Init(ctx context.Context, conf Config) (*Client, error) {
	if err := conf.validate(); err != nil {
		return nil, err
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:            conf.Host,
		Password:        conf.Password,
		DB:              conf.DB,
		ConnMaxIdleTime: connMaxIdleTime,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redisc: Redis连接「%s」ping失败: %w", conf.Host, err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

// GetAndDelete 使用Lua脚本兼容不支持GETDEL命令的旧版Redis，并保证读取和删除的原子性。
func (c *Client) GetAndDelete(ctx context.Context, key string) (string, bool, error) {
	value, err := getAndDeleteScript.Run(ctx, c.client, []string{key}).Text()
	if errors.Is(err, goredis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

func (c *Client) HGetInt64(ctx context.Context, key, field string) (int64, bool, error) {
	value, err := c.client.HGet(ctx, key, field).Int64()
	if errors.Is(err, goredis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return value, true, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Client) ReplaceIfValue(
	ctx context.Context,
	oldKey, newKey string,
	expectedValue, newValue any,
	expiration time.Duration,
) (bool, error) {
	result, err := replaceIfValueScript.Run(
		ctx,
		c.client,
		[]string{oldKey, newKey},
		expectedValue,
		newValue,
		expiration.Milliseconds(),
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (c *Client) RunInt64SliceScript(
	ctx context.Context,
	source string,
	keys []string,
	args ...any,
) ([]int64, error) {
	return goredis.NewScript(source).Run(ctx, c.client, keys, args...).Int64Slice()
}

func (c *Client) Close() error {
	return c.client.Close()
}
