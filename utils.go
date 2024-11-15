package upbit

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// generateToken generates a JWT token for authentication
func (c *Client) generateToken() (string, error) {
	claims := jwt.MapClaims{
		"access_key": c.accessKey,
		"nonce":      uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.apiSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return "Bearer " + tokenString, nil
}

// 수신된 데이터를 예쁘게 출력하는 헬퍼 함수
func PrettyPrint(prefix string, data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err == nil {
		fmt.Printf("%s:\n%s\n", prefix, prettyJSON.String())
	} else {
		fmt.Printf("%s: %s\n", prefix, string(data))
	}
}
