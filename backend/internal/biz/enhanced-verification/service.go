package enhancedverification

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	verificationdriver "fortunefinger/backend/internal/biz/enhanced-verification/driver"
	"fortunefinger/backend/internal/pkg/redisc"
)

const (
	verificationKeyPrefix = "enhanced:verification:data:"
	tokenKeyPrefix        = "enhanced:verification:token:"
)

type Service struct {
	redis  *redisc.Client
	driver verificationdriver.Driver
	ttl    time.Duration
}

type Verification struct {
	ID             string
	Mode           string
	Image          string
	Prompt         string
	Thumb          string
	ThumbX         int
	ThumbY         int
	ThumbWidth     int
	ThumbHeight    int
	RequiredPoints int
	ExpiresIn      int64
}

type Point = verificationdriver.Point

type verificationRecord struct {
	Mode    string          `json:"mode"`
	Scene   string          `json:"scene"`
	Subject string          `json:"subject"`
	IP      string          `json:"ip"`
	Payload json.RawMessage `json:"payload"`
}

type tokenRecord struct {
	Scene   string `json:"scene"`
	Subject string `json:"subject"`
	IP      string `json:"ip"`
}

func NewService(redisClient *redisc.Client, mode string, ttl time.Duration) (*Service, error) {
	if redisClient == nil {
		return nil, errors.New("enhanced-verification: Redis客户端不能为空")
	}
	if ttl <= 0 {
		return nil, errors.New("enhanced-verification: 有效期必须大于0")
	}
	verificationDriver, err := verificationdriver.New(mode)
	if err != nil {
		return nil, err
	}
	return &Service{
		redis:  redisClient,
		driver: verificationDriver,
		ttl:    ttl,
	}, nil
}

func (s *Service) Create(ctx context.Context, scene, subject, ip string) (*Verification, error) {
	data, err := s.driver.Generate()
	if err != nil {
		return nil, err
	}
	id, err := randomValue(24)
	if err != nil {
		return nil, err
	}
	record, err := json.Marshal(verificationRecord{
		Mode:    s.driver.Mode(),
		Scene:   scene,
		Subject: subject,
		IP:      ip,
		Payload: data.Payload,
	})
	if err != nil {
		return nil, err
	}
	if err := s.redis.Set(ctx, verificationKeyPrefix+id, record, s.ttl); err != nil {
		return nil, err
	}
	return &Verification{
		ID:             id,
		Mode:           s.driver.Mode(),
		Image:          data.Image,
		Prompt:         data.Prompt,
		Thumb:          data.Thumb,
		ThumbX:         data.ThumbX,
		ThumbY:         data.ThumbY,
		ThumbWidth:     data.ThumbWidth,
		ThumbHeight:    data.ThumbHeight,
		RequiredPoints: data.RequiredPoints,
		ExpiresIn:      int64(s.ttl.Seconds()),
	}, nil
}

func (s *Service) Verify(
	ctx context.Context,
	id, scene, subject, ip string,
	points []Point,
) (string, bool, error) {
	value, exists, err := s.redis.GetAndDelete(ctx, verificationKeyPrefix+id)
	if err != nil || !exists {
		return "", false, err
	}
	var record verificationRecord
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return "", false, err
	}
	if record.Mode != s.driver.Mode() || record.Scene != scene || record.Subject != subject || record.IP != ip {
		return "", false, nil
	}
	verified, err := s.driver.Validate(record.Payload, points)
	if err != nil || !verified {
		return "", false, err
	}
	token, err := randomValue(32)
	if err != nil {
		return "", false, err
	}
	tokenValue, err := json.Marshal(tokenRecord{Scene: scene, Subject: subject, IP: ip})
	if err != nil {
		return "", false, err
	}
	if err := s.redis.Set(ctx, tokenKey(token), tokenValue, s.ttl); err != nil {
		return "", false, err
	}
	return token, true, nil
}

func (s *Service) Consume(ctx context.Context, token, scene, subject, ip string) (bool, error) {
	if token == "" {
		return false, nil
	}
	value, exists, err := s.redis.GetAndDelete(ctx, tokenKey(token))
	if err != nil || !exists {
		return false, err
	}
	var record tokenRecord
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return false, err
	}
	return record.Scene == scene && record.Subject == subject && record.IP == ip, nil
}

func randomValue(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func tokenKey(token string) string {
	digest := sha256.Sum256([]byte(token))
	return tokenKeyPrefix + hex.EncodeToString(digest[:])
}
