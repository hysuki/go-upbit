package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Package auth는 Upbit API 인증을 위한 기능을 제공합니다.

// BaseTokenGenerator는 기본 토큰 생성 기능을 정의하는 인터페이스입니다.
type BaseTokenGenerator interface {
	// GenerateToken은 기본 인증 토큰을 생성합니다.
	GenerateToken() (string, error)
}

// RestTokenGenerator는 REST API용 토큰 생성 기능을 정의하는 인터페이스입니다.
type RestTokenGenerator interface {
	BaseTokenGenerator
	// GenerateTokenWithQuery는 쿼리 파라미터를 포함한 인증 토큰을 생성합니다.
	GenerateTokenWithQuery(query url.Values) (string, error)
}

// Credentials는 API 인증 정보를 담는 구조체입니다.
type Credentials struct {
	// AccessKey는 API 액세스 키입니다.
	AccessKey string
	// SecretKey는 API 시크릿 키입니다.
	SecretKey string
}

// WebSocketTokenGen은 WebSocket용 토큰 생성기입니다.
type WebSocketTokenGen struct {
	creds Credentials
}

// RestTokenGen은 REST API용 토큰 생성기입니다.
type RestTokenGen struct {
	creds Credentials
}

// NewRestTokenGen은 새로운 REST API 토큰 생성기를 생성합니다.
// creds는 API 인증 정보입니다.
func NewRestTokenGen(creds Credentials) *RestTokenGen {
	return &RestTokenGen{creds: creds}
}

// NewWebSocketTokenGen은 새로운 WebSocket 토큰 생성기를 생성합니다.
// creds는 API 인증 정보입니다.
func NewWebSocketTokenGen(creds Credentials) *WebSocketTokenGen {
	return &WebSocketTokenGen{creds: creds}
}

// GenerateToken은 REST API용 JWT 토큰을 생성합니다.
// 생성된 토큰 문자열과 에러를 반환합니다.
func (g *RestTokenGen) GenerateToken() (string, error) {
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
// 생성된 토큰 문자열과 에러를 반환합니다.
func (g *WebSocketTokenGen) GenerateToken() (string, error) {
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

// GenerateTokenWithQuery는 쿼리 파라미터를 포함한 REST API용 JWT 토큰을 생성합니다.
// query는 요청 쿼리 파라미터입니다.
// 생성된 토큰 문자열과 에러를 반환합니다.
func (g *RestTokenGen) GenerateTokenWithQuery(query url.Values) (string, error) {
	payload := map[string]interface{}{
		"access_key": g.creds.AccessKey,
		"nonce":      uuid.New().String(),
	}

	if query != nil {
		queryString := query.Encode()
		hash := sha512.Sum512([]byte(queryString))
		queryHash := hex.EncodeToString(hash[:])

		payload["query_hash"] = queryHash
		payload["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))
	return token.SignedString([]byte(g.creds.SecretKey))
}
