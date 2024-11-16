package upbit

import (
	"fmt"
	"sync"
	"time"

	"github.com/hysuki/go-upbit/auth"
	"github.com/hysuki/go-upbit/rest"
	"github.com/hysuki/go-upbit/websocket/private"
	"github.com/hysuki/go-upbit/websocket/public"
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
// REST API, Public WebSocket, Private WebSocket 클라이언트를 포함합니다.
type UpbitClient struct {
	credentials  auth.Credentials
	pingInterval time.Duration
	PublicWS     *public.Client
	PrivateWS    *private.Client
	RestAPI      rest.Client
}

// UpbitClientOption은 UpbitClient 설정을 위한 함수 타입입니다.
type UpbitClientOption func(*UpbitClient)

// WithKeys는 API 키와 시크릿을 설정하는 옵션을 반환합니다.
// accessKey는 API 액세스 키이고, apiSecret은 API 시크릿 키입니다.
func WithKeys(accessKey, apiSecret string) UpbitClientOption {
	return func(c *UpbitClient) {
		c.credentials = auth.Credentials{
			AccessKey: accessKey,
			SecretKey: apiSecret,
		}
	}
}

// WithPingInterval은 WebSocket 연결의 핑 전송 간격을 설정하는 옵션을 반환합니다.
// interval은 핑 전송 간격입니다.
func WithPingInterval(interval time.Duration) UpbitClientOption {
	return func(c *UpbitClient) {
		c.pingInterval = interval
	}
}

// GetPingInterval은 설정된 핑 전송 간격을 반환합니다.
func (c *UpbitClient) GetPingInterval() time.Duration {
	return c.pingInterval
}

// NewUpbitClient는 새로운 UpbitClient 인스턴스를 생성하고 초기화합니다.
// opts는 클라이언트 설정을 위한 옵션들입니다.
// 에러가 발생하면 nil과 에러를 반환합니다.
func NewUpbitClient(opts ...UpbitClientOption) (client *UpbitClient, err error) {
	// nil 포인터 참조 방지를 위한 초기화
	client = &UpbitClient{
		pingInterval: 30 * time.Second, // 기본값 설정
	}

	for _, opt := range opts {
		opt(client)
	}

	// 인증 정보 검증
	if client.credentials == (auth.Credentials{}) {
		return nil, fmt.Errorf("credentials are required")
	}

	var wg sync.WaitGroup
	wg.Add(3) // REST API + Public WS + Private WS

	// 공통 에러 채널
	errCh := make(chan error, 3)

	// REST API 클라이언트 초기화
	go func() {
		defer wg.Done()
		restTokenGen := auth.NewRestTokenGen(client.credentials)
		client.RestAPI = rest.NewClient(restTokenGen)
		if client.RestAPI == nil {
			errCh <- fmt.Errorf("failed to initialize REST API client")
		}
	}()

	// Public WebSocket 클라이언트 초기화
	go func() {
		defer wg.Done()
		wsTokenGen := auth.NewWebSocketTokenGen(client.credentials)
		pub, err := public.NewClient(PublicWebsocketEndpoint, wsTokenGen, client.pingInterval)
		if err != nil {
			errCh <- fmt.Errorf("public client error: %w", err)
			return
		}
		client.PublicWS = pub
	}()

	// Private WebSocket 클라이언트 초기화
	go func() {
		defer wg.Done()
		wsTokenGen := auth.NewWebSocketTokenGen(client.credentials)
		pri, err := private.NewClient(PrivateWebsocketEndpoint, wsTokenGen, client.pingInterval)
		if err != nil {
			errCh <- fmt.Errorf("private client error: %w", err)
			return
		}
		client.PrivateWS = pri
	}()

	wg.Wait()
	close(errCh)

	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		// 에러 발생 시 WebSocket 클라이언트만 정리
		if client.PublicWS != nil {
			client.PublicWS.Close()
		}
		if client.PrivateWS != nil {
			client.PrivateWS.Close()
		}
		return nil, fmt.Errorf("client initialization errors: %v", errors)
	}

	return client, nil
}
