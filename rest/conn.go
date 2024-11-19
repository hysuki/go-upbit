package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hysuki/go-upbit/rest/exchange"
	"github.com/hysuki/go-upbit/rest/quotation"
)

const (
	BaseURL = "https://api.upbit.com/v1" // REST API 엔드포인트
)

// TokenGenerator는 인증 토큰을 생성하는 인터페이스입니다.
type TokenGenerator interface {
	// GenerateToken은 기본 인증 토큰을 생성합니다.
	GenerateToken() (string, error)
	// GenerateTokenWithQuery는 쿼리 파라미터를 포함한 인증 토큰을 생성합니다.
	GenerateTokenWithQuery(query url.Values) (string, error)
	// GenerateTokenWithBody는 JSON body를 해시화하여 인증 토큰을 생성합니다.
	GenerateTokenWithBody(body string) (string, error)
}

// Client는 REST API 클라이언트 인터페이스를 정의합니다.
type Client interface {
	// Get은 GET 요청을 수행합니다.
	Get(path string, params map[string]string) ([]byte, error)
	// Post는 POST 요청을 수행합니다.
	Post(path string, body interface{}) ([]byte, error)
	// Delete는 DELETE 요청을 수행합니다.
	Delete(path string, params map[string]string) ([]byte, error)
	// GetExchange는 거래소 API 클라이언트를 반환합니다.
	GetExchange() *exchange.Exchange
	// GetQuotation는 시세 API 클라이언트를 반환합니다.
	GetQuotation() *quotation.Quotation
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
	Quotation  *quotation.Quotation
}

// NewClient는 새로운 REST API 클라이언트를 생성합니다.
func NewClient(tokenGen TokenGenerator) *client {
	c := &client{
		httpClient: &http.Client{},
		tokenGen:   tokenGen,
		baseURL:    BaseURL,
	}
	c.Exchange = exchange.NewExchange(c)
	c.Quotation = quotation.NewQuotation(c)
	return c
}

// handleResponse는 API 응답을 처리하고 에러가 있다면 에러를 반환합니다
func (c *client) handleResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 응답이 에러인지 확인
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return nil, &errResp.Error
	}

	return body, nil
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

	// 인증 토큰 생성 시 쿼리 파라미터 전달
	token, err := c.tokenGen.GenerateTokenWithQuery(req.URL.Query())
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	// HTTP 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp)
}

func (c *client) Post(path string, body interface{}) ([]byte, error) {
	// body를 map[string]interface{}로 변환
	var bodyMap map[string]interface{}

	// 1. 이미 map[string]interface{}인 경우
	if m, ok := body.(map[string]interface{}); ok {
		bodyMap = m
	} else {
		// 2. 구조체를 JSON으로 마샬링 후 다시 map으로 언마샬링
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}

		if err := json.Unmarshal(jsonBytes, &bodyMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal body to map: %w", err)
		}
	}

	// url.Values로 변환
	values := make(url.Values)
	for k, v := range bodyMap {
		values.Add(k, fmt.Sprintf("%v", v))
	}

	// url.Values로 인코딩된 문자열 생성
	encodedBody := values.Encode()

	// 동일한 values로 토큰 생성
	token, err := c.tokenGen.GenerateTokenWithQuery(values)
	if err != nil {
		return nil, err
	}

	// 요청 생성 (body는 인코딩된 form 데이터)
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(encodedBody))
	if err != nil {
		return nil, err
	}

	// 헤더 설정 (Content-Type을 form 데이터로 변경)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", token)

	// HTTP 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp)
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

	return c.handleResponse(resp)
}

func (c *client) GetExchange() *exchange.Exchange {
	return c.Exchange
}

func (c *client) GetQuotation() *quotation.Quotation {
	return c.Quotation
}
