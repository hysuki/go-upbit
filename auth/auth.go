package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RestTokenGenerator는 REST API용 토큰 생성 기능을 정의하는 인터페이스입니다.
type RestTokenGenerator interface {
	GenerateToken() (string, error)
	GenerateTokenWithQuery(query url.Values) (string, error)
	GenerateTokenWithBody(body string) (string, error)
}

// WebSocketTokenGenerator는 WebSocket API용 토큰 생성 기능을 정의하는 인터페이스입니다.
type WebSocketTokenGenerator interface {
	GenerateToken() (string, error)
}

// Credentials는 API 인증 정보를 담는 구조체입니다.
type Credentials struct {
	AccessKey string
	SecretKey string
}

type WebSocketTokenGen struct {
	creds Credentials
}

type RestTokenGen struct {
	creds Credentials
}

func NewRestTokenGen(creds Credentials) *RestTokenGen {
	return &RestTokenGen{creds: creds}
}

func NewWebSocketTokenGen(creds Credentials) *WebSocketTokenGen {
	return &WebSocketTokenGen{creds: creds}
}

// generateToken은 기본 JWT 토큰 생성 로직을 담당하는 헬퍼 함수입니다.
func generateToken(creds Credentials, payload map[string]interface{}) (string, error) {
	if payload == nil {
		payload = map[string]interface{}{}
	}

	payload["access_key"] = creds.AccessKey
	payload["nonce"] = uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))
	return token.SignedString([]byte(creds.SecretKey))
}

func (g *RestTokenGen) GenerateToken() (string, error) {
	tokenString, err := generateToken(g.creds, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate REST token: %v", err)
	}
	return tokenString, nil
}

func (g *WebSocketTokenGen) GenerateToken() (string, error) {
	tokenString, err := generateToken(g.creds, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate WebSocket token: %v", err)
	}
	return "Bearer " + tokenString, nil
}

func (g *RestTokenGen) GenerateTokenWithQuery(query url.Values) (string, error) {
	payload := map[string]interface{}{}

	if len(query) > 0 {
		queryString := query.Encode()
		queryString, err := url.QueryUnescape(queryString)
		if err != nil {
			return "", err
		}

		hash := sha512.Sum512([]byte(queryString))
		payload["query_hash"] = hex.EncodeToString(hash[:])
		payload["query_hash_alg"] = "SHA512"
	}

	tokenString, err := generateToken(g.creds, payload)
	if err != nil {
		return "", fmt.Errorf("failed to generate token with query: %v", err)
	}

	return "Bearer " + tokenString, nil
}

func (g *RestTokenGen) GenerateTokenWithBody(body string) (string, error) {
	payload := map[string]interface{}{}

	if body != "" {
		hash := sha512.Sum512([]byte(body))
		payload["query_hash"] = hex.EncodeToString(hash[:])
		payload["query_hash_alg"] = "SHA512"
	}

	tokenString, err := generateToken(g.creds, payload)
	if err != nil {
		return "", fmt.Errorf("failed to generate token with body: %v", err)
	}

	return "Bearer " + tokenString, nil
}
