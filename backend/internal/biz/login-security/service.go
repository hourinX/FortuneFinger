package loginsecurity

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"time"

	"fortunefinger/backend/internal/pkg/redisc"
)

const keyPrefix = "login:risk:"

var lockDurations = [...]time.Duration{
	5 * time.Minute,
	10 * time.Minute,
	30 * time.Minute,
	2 * time.Hour,
	10 * time.Hour,
}

//go:embed script/record_failure.lua
var recordFailureLua string

type Service struct {
	redis       *redisc.Client
	maxFailures int
	resetAfter  time.Duration
}

type State struct {
	Failures          int64
	RemainingAttempts int64
	Level             int64
	LockedUntil       time.Time
}

func NewService(redisClient *redisc.Client, maxFailures int, resetAfter time.Duration) (*Service, error) {
	if redisClient == nil {
		return nil, errors.New("login-security: Redis客户端不能为空")
	}
	if maxFailures <= 0 {
		return nil, errors.New("login-security: 最大登录失败次数必须大于0")
	}
	if resetAfter <= 0 {
		return nil, errors.New("login-security: 登录失败重置时间必须大于0")
	}
	return &Service{
		redis:       redisClient,
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
	}, nil
}

func (s *Service) LockedUntil(ctx context.Context, username, ip string, now time.Time) (time.Time, error) {
	seconds, exists, err := s.redis.HGetInt64(ctx, s.key(username, ip), "locked_until")
	if err != nil || !exists {
		return time.Time{}, err
	}
	lockedUntil := time.Unix(seconds, 0)
	if !lockedUntil.After(now) {
		return time.Time{}, nil
	}
	return lockedUntil, nil
}

func (s *Service) RequiresEnhancedVerification(ctx context.Context, username, ip string) (bool, error) {
	failures, _, err := s.redis.HGetInt64(ctx, s.key(username, ip), "failures")
	if err != nil {
		return false, err
	}
	level, _, err := s.redis.HGetInt64(ctx, s.key(username, ip), "level")
	if err != nil {
		return false, err
	}
	return failures >= 2 || level > 0, nil
}

func (s *Service) RecordFailure(ctx context.Context, username, ip string, now time.Time) (State, error) {
	args := []any{now.Unix(), s.maxFailures, int64(s.resetAfter.Seconds())}
	for _, duration := range lockDurations {
		args = append(args, int64(duration.Seconds()))
	}
	values, err := s.redis.RunInt64SliceScript(ctx, recordFailureLua, []string{s.key(username, ip)}, args...)
	if err != nil {
		return State{}, err
	}
	remainingAttempts := int64(s.maxFailures) - values[0]
	if values[2] > now.Unix() || remainingAttempts < 0 {
		remainingAttempts = 0
	}
	return State{
		Failures:          values[0],
		RemainingAttempts: remainingAttempts,
		Level:             values[1],
		LockedUntil:       time.Unix(values[2], 0),
	}, nil
}

func (s *Service) Reset(ctx context.Context, username, ip string) error {
	return s.redis.Delete(ctx, s.key(username, ip))
}

func (s *Service) key(username, ip string) string {
	digest := sha256.Sum256([]byte(username + "\x00" + ip))
	return keyPrefix + hex.EncodeToString(digest[:])
}
