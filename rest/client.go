package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"upbit.yougcha.bot/pkg/upbit/rest/exchange"
)

const (
	BaseURL = "https://api.upbit.com/v1" // REST API 엔드포인트
)

// TokenGenerator는 인증 토큰을 생성하는 인터페이스입니다.
type TokenGenerator interface {
	GenerateToken() (string, error)
}

// Client는 REST API 클라이언트 인터페이스를 정의합니다.
type Client interface {
	Get(path string, params map[string]string) ([]byte, error)
	Post(path string, body interface{}) ([]byte, error)
	Delete(path string, params map[string]string) ([]byte, error)
	GetExchange() *exchange.Exchange
}

// APIError는 Upbit API 에러 응답을 나타냅니다.
type APIError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

// ErrorResponse는 API 에러 응답의 전체 구조를 나타냅니다.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// client는 REST API 호출을 관리합니다.
type client struct {
	httpClient *http.Client
	tokenGen   TokenGenerator
	baseURL    string
	Exchange   *exchange.Exchange
}

// NewClient는 새로운 REST API 클라이언트를 생성합니다.
func NewClient(tokenGen TokenGenerator) *client {
	c := &client{
		httpClient: &http.Client{},
		tokenGen:   tokenGen,
		baseURL:    BaseURL,
	}
	c.Exchange = exchange.NewExchange(c)
	return c
}

func (c *client) Get(path string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	// 쿼리 파라미터 추가
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// 인증 토큰 생성 및 헤더 추가
	token, err := c.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	// HTTP 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 4XX 에러 응답 처리
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, err
		}
		return nil, &errResp.Error
	}

	return body, nil
}

func (c *client) Post(path string, body interface{}) ([]byte, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")

	// 인증 토큰 생성 및 헤더 추가
	token, err := c.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	// HTTP 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 4XX 에러 응답 처리
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, err
		}
		return nil, &errResp.Error
	}

	return respBody, nil
}

func (c *client) Delete(path string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	// 쿼리 파라미터 추가
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// 인증 토큰 생성 및 헤더 추가
	token, err := c.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	// HTTP 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 4XX 에러 응답 처리
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, err
		}
		return nil, &errResp.Error
	}

	return body, nil
}

func (c *client) GetExchange() *exchange.Exchange {
	return c.Exchange
}
