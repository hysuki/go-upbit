package upbit

import (
	"fmt"
	"sync"
	"time"

	"upbit.yougcha.bot/pkg/upbit/auth"
	"upbit.yougcha.bot/pkg/upbit/rest"
	"upbit.yougcha.bot/pkg/upbit/websocket"
)

// API 엔드포인트 상수
const (
	// PublicWebsocketEndpoint는 공개 WebSocket API 엔드포인트입니다.
	PublicWebsocketEndpoint = "wss://api.upbit.com/websocket/v1"
	// PrivateWebsocketEndpoint는 인증이 필요한 WebSocket API 엔드포인트입니다.
	PrivateWebsocketEndpoint = "wss://api.upbit.com/websocket/v1/private"
	// RestAPIEndpoint는 REST API 엔드포인트입니다.
	RestAPIEndpoint = "https://api.upbit.com/v1"
)

// UpbitClient는 업비트 API와의 통신을 위한 메인 클라이언트 구조체입니다.
type UpbitClient struct {
	credentials  auth.Credentials
	pingInterval time.Duration
	PublicWS     websocket.Client
	PrivateWS    websocket.Client
	RestAPI      rest.Client
}

// UpbitClientOption은 UpbitClient 설정을 위한 함수 타입입니다.
type UpbitClientOption func(*UpbitClient)

// WithKeys는 API 키와 시크릿을 설정하는 옵션을 반환합니다.
func WithKeys(accessKey, apiSecret string) UpbitClientOption {
	return func(c *UpbitClient) {
		c.credentials = auth.Credentials{
			AccessKey: accessKey,
			SecretKey: apiSecret,
		}
	}
}

// WithPingInterval은 핑 전송 간격을 설정하는 옵션을 반환합니다.
func WithPingInterval(interval time.Duration) UpbitClientOption {
	return func(c *UpbitClient) {
		c.pingInterval = interval
	}
}

// GetPingInterval returns the configured ping interval
func (c *UpbitClient) GetPingInterval() time.Duration {
	return c.pingInterval
}

// NewUpbitClient는 새로운 UpbitClient 인스턴스를 생성하고 초기화합니다.
func NewUpbitClient(opts ...UpbitClientOption) (*UpbitClient, error) {
	client := &UpbitClient{}

	for _, opt := range opts {
		opt(client)
	}

	// REST API 클라이언트 초기화
	restTokenGen := auth.NewRestTokenGenerator(client.credentials)
	client.RestAPI = rest.NewClient(restTokenGen)

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2)
	wsTokenGen := auth.NewWebSocketTokenGenerator(client.credentials)

	go func() {
		defer wg.Done()
		pub, err := websocket.NewClient(websocket.PublicEndpoint, wsTokenGen, client.pingInterval)
		if err != nil {
			errCh <- fmt.Errorf("public client error: %w", err)
			return
		}
		client.PublicWS = pub
	}()

	go func() {
		defer wg.Done()
		private, err := websocket.NewClient(websocket.PrivateEndpoint, wsTokenGen, client.pingInterval)
		if err != nil {
			errCh <- fmt.Errorf("private client error: %w", err)
			return
		}
		client.PrivateWS = private
	}()

	wg.Wait()
	close(errCh)

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("client initialization errors: %v", errors)
	}

	return client, nil
}
