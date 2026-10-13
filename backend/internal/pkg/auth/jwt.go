package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessToken  = "access"
	RefreshToken = "refresh"
)

type Claims struct {
	TokenType string `json:"token_type"`
	Remember  bool   `json:"remember,omitempty"`
	jwt.RegisteredClaims
}

type Manager struct {
	key                   []byte
	issuer                string
	audience              string
	accessExpire          time.Duration
	refreshExpire         time.Duration
	rememberRefreshExpire time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshIdentity struct {
	UserID   uint64
	Remember bool
}

func NewManager(
	secret, issuer, audience string,
	accessExpire, refreshExpire, rememberRefreshExpire time.Duration,
) (*Manager, error) {
	if len([]byte(secret)) < 32 {
		return nil, errors.New("auth: JWT密钥不能少于32字节")
	}
	if strings.TrimSpace(issuer) == "" || strings.TrimSpace(audience) == "" {
		return nil, errors.New("auth: JWT签发方和接收方不能为空")
	}
	if accessExpire <= 0 || refreshExpire <= accessExpire {
		return nil, errors.New("auth: JWT刷新凭证有效期必须大于访问凭证有效期")
	}
	if rememberRefreshExpire <= refreshExpire {
		return nil, errors.New("auth: JWT记住登录有效期必须大于普通刷新凭证有效期")
	}
	return &Manager{
		key:                   []byte(secret),
		issuer:                issuer,
		audience:              audience,
		accessExpire:          accessExpire,
		refreshExpire:         refreshExpire,
		rememberRefreshExpire: rememberRefreshExpire,
	}, nil
}

func (m *Manager) Generate(userID uint64, remember bool) (TokenPair, error) {
	accessToken, err := m.sign(userID, AccessToken, m.accessExpire, false)
	if err != nil {
		return TokenPair{}, err
	}
	refreshExpire := m.refreshExpire
	if remember {
		refreshExpire = m.rememberRefreshExpire
	}
	refreshToken, err := m.sign(userID, RefreshToken, refreshExpire, remember)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpire.Seconds()),
	}, nil
}

func (m *Manager) VerifyAccessToken(raw string) (uint64, error) {
	identity, err := m.verify(raw, AccessToken)
	return identity.UserID, err
}

func (m *Manager) VerifyRefreshToken(raw string) (RefreshIdentity, error) {
	return m.verify(raw, RefreshToken)
}

func (m *Manager) sign(userID uint64, tokenType string, lifetime time.Duration, remember bool) (string, error) {
	if userID == 0 {
		return "", errors.New("auth: Token用户ID无效")
	}
	jti, err := newJTI()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := Claims{
		TokenType: tokenType,
		Remember:  remember,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   strconv.FormatUint(userID, 10),
			Audience:  jwt.ClaimStrings{m.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(lifetime)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.key)
}

func newJTI() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", errors.New("auth: 生成Token唯一标识失败")
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func (m *Manager) verify(raw, tokenType string) (RefreshIdentity, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(*jwt.Token) (any, error) {
			return m.key, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithIssuedAt(),
	)
	if err != nil || !token.Valid || claims.TokenType != tokenType {
		return RefreshIdentity{}, errors.New("auth: Token无效")
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || userID == 0 {
		return RefreshIdentity{}, errors.New("auth: Token用户ID无效")
	}
	return RefreshIdentity{UserID: userID, Remember: claims.Remember}, nil
}
