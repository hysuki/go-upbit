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

// Upbit API 서버의 엔드포인트를 정의하는 상수들입니다.
const (
	PublicWebsocketEndpoint  = "wss://api.upbit.com/websocket/v1"         // 공개 웹소켓 API 주소
	PrivateWebsocketEndpoint = "wss://api.upbit.com/websocket/v1/private" // 비공개 웹소켓 API 주소
	RestAPIEndpoint          = "https://api.upbit.com/v1"                 // REST API 주소
)

// UpbitClient는 Upbit API 클라이언트입니다.
type UpbitClient struct {
	credentials  auth.Credentials // API 인증 정보
	pingInterval time.Duration    // 웹소켓 핑 전송 간격
	PublicWS     *public.Client   // 공개 웹소켓 클라이언트
	PrivateWS    *private.Client  // 비공개 웹소켓 클라이언트
	RestAPI      rest.Client      // REST API 클라이언트
}

// UpbitClientOption은 UpbitClient의 설정을 변경하는 함수 타입입니다.
type UpbitClientOption func(*UpbitClient)

// WithKeys는 API 인증 키를 설정하는 옵션을 반환합니다.
// accessKey는 액세스 키, apiSecret은 시크릿 키입니다.
func WithKeys(accessKey, apiSecret string) UpbitClientOption {
	return func(c *UpbitClient) {
		c.credentials = auth.Credentials{
			AccessKey: accessKey,
			SecretKey: apiSecret,
		}
	}
}

// WithPingInterval은 웹소켓 핑 전송 간격을 설정하는 옵션을 반환합니다.
// interval은 핑 전송 간격입니다.
func WithPingInterval(interval time.Duration) UpbitClientOption {
	return func(c *UpbitClient) {
		c.pingInterval = interval
	}
}

// GetPingInterval은 현재 설정된 웹소켓 핑 전송 간격을 반환합니다.
func (c *UpbitClient) GetPingInterval() time.Duration {
	return c.pingInterval
}

// NewUpbitClient는 새로운 Upbit API 클라이언트를 생성합니다.
// opts로 클라이언트 설정을 지정할 수 있으며, WithKeys 옵션은 필수입니다.
func NewUpbitClient(opts ...UpbitClientOption) (client *UpbitClient, err error) {
	client = &UpbitClient{
		pingInterval: 30 * time.Second, // 기본값 설정
	}

	for _, opt := range opts {
		opt(client)
	}

	// 인증 정보 검증
	if client.credentials == (auth.Credentials{}) {
		return nil, fmt.Errorf("인증 정보가 필요합니다")
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
			errCh <- fmt.Errorf("REST API 클라이언트 초기화 실패")
		}
	}()

	// Public WebSocket 클라이언트 초기화
	go func() {
		defer wg.Done()
		wsTokenGen := auth.NewWebSocketTokenGen(client.credentials)
		pub, err := public.NewClient(PublicWebsocketEndpoint, wsTokenGen, client.pingInterval)
		if err != nil {
			errCh <- fmt.Errorf("공개 웹소켓 클라이언트 에러: %w", err)
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
			errCh <- fmt.Errorf("비공개 웹소켓 클라이언트 에러: %w", err)
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
		return nil, fmt.Errorf("클라이언트 초기화 에러: %v", errors)
	}

	return client, nil
}
