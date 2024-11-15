package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenGenerator는 인증 토큰을 생성하는 기본 인터페이스입니다.
type TokenGenerator interface {
	GenerateToken() (string, error)
}

// Credentials는 API 인증 정보를 담는 구조체입니다.
type Credentials struct {
	AccessKey string
	SecretKey string
}

// RestTokenGenerator는 REST API용 토큰 생성기입니다.
type RestTokenGenerator struct {
	creds Credentials
}

// WebSocketTokenGenerator는 WebSocket용 토큰 생성기입니다.
type WebSocketTokenGenerator struct {
	creds Credentials
}

// NewRestTokenGenerator는 새로운 REST API 토큰 생성기를 생성합니다.
func NewRestTokenGenerator(creds Credentials) *RestTokenGenerator {
	return &RestTokenGenerator{creds: creds}
}

// NewWebSocketTokenGenerator는 새로운 WebSocket 토큰 생성기를 생성합니다.
func NewWebSocketTokenGenerator(creds Credentials) *WebSocketTokenGenerator {
	return &WebSocketTokenGenerator{creds: creds}
}

// GenerateToken은 REST API용 JWT 토큰을 생성합니다.
func (g *RestTokenGenerator) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"access_key": g.creds.AccessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.creds.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate REST token: %v", err)
	}

	return tokenString, nil
}

// GenerateToken은 WebSocket용 JWT 토큰을 생성합니다.
func (g *WebSocketTokenGenerator) GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"access_key": g.creds.AccessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.creds.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate WebSocket token: %v", err)
	}

	return "Bearer " + tokenString, nil
}
