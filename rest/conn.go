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

// Package rest는 Upbit REST API와의 통신을 담당합니다.
const (
	BaseURL = "https://api.upbit.com/v1" // REST API 기본 URL
)

// TokenGenerator는 Upbit API 인증에 필요한 토큰을 생성하는 인터페이스입니다.
type TokenGenerator interface {
	GenerateToken() (string, error)                          // 기본 JWT 토큰 생성
	GenerateTokenWithQuery(query url.Values) (string, error) // 쿼리 파라미터를 포함한 JWT 토큰 생성
	GenerateTokenWithBody(body string) (string, error)       // JSON 본문을 포함한 JWT 토큰 생성
}

// Client는 Upbit REST API와 상호작용하기 위한 메서드를 정의하는 인터페이스입니다.
type Client interface {
	Get(path string, params map[string]string) ([]byte, error)    // GET 요청 수행
	Post(path string, body interface{}) ([]byte, error)           // POST 요청 수행
	Delete(path string, params map[string]string) ([]byte, error) // DELETE 요청 수행
	GetExchange() *exchange.Exchange                              // 거래소 API 객체 반환
	GetQuotation() *quotation.Quotation                           // 시세 조회 API 객체 반환
}

// APIError는 Upbit API에서 반환하는 에러 정보를 나타냅니다.
type APIError struct {
	Name    string `json:"name"`    // 에러 이름
	Message string `json:"message"` // 에러 메시지
}

// Error는 에러 메시지를 반환합니다.
func (e *APIError) Error() string {
	return e.Message
}

// ErrorResponse는 API 에러 응답을 나타냅니다.
type ErrorResponse struct {
	Error APIError `json:"error"` // API 에러 정보
}

// client는 Upbit REST API 클라이언트입니다.
type client struct {
	httpClient *http.Client         // HTTP 클라이언트
	tokenGen   TokenGenerator       // 토큰 생성기
	baseURL    string               // API 기본 URL
	Exchange   *exchange.Exchange   // 거래소 API 객체
	Quotation  *quotation.Quotation // 시세 조회 API 객체
}

// NewClient는 새로운 Upbit REST API 클라이언트를 생성합니다.
// tokenGen은 API 인증에 사용할 토큰 생성기입니다.
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

// handleResponse는 API 응답을 처리하고 에러가 있는 경우 이를 반환합니다.
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

// Get은 지정된 경로로 GET 요청을 보내고 응답을 반환합니다.
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

// Post는 지정된 경로로 POST 요청을 보내고 응답을 반환합니다.
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

// Delete는 지정된 경로로 DELETE 요청을 보내고 응답을 반환합니다.
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

// GetExchange는 거래소 API 관련 기능을 제공하는 Exchange 객체를 반환합니다.
func (c *client) GetExchange() *exchange.Exchange {
	return c.Exchange
}

// GetQuotation은 시세 조회 관련 기능을 제공하는 Quotation 객체를 반환합니다.
func (c *client) GetQuotation() *quotation.Quotation {
	return c.Quotation
}
