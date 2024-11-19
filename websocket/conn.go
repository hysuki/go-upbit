package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/hysuki/go-upbit/auth"
)

// BaseClient는 웹소켓 기본 클라이언트입니다.
type BaseClient struct {
	Conn         *websocket.Conn              // 웹소켓 연결
	Ctx          context.Context              // 컨텍스트
	Cancel       context.CancelFunc           // 컨텍스트 취소 함수
	IsRunning    bool                         // 실행 상태
	Mu           sync.Mutex                   // 뮤텍스
	Messages     []Message                    // 메시지 목록
	Endpoint     string                       // 웹소켓 서버 주소
	TokenGen     auth.WebSocketTokenGenerator // 토큰 생성기
	PingTicker   *time.Ticker                 // 핑 전송 타이머
	PingInterval time.Duration                // 핑 전송 간격
}

// NewBaseClient는 새로운 웹소켓 기본 클라이언트를 생성합니다.
// endpoint는 웹소켓 서버 주소, tokenGen은 토큰 생성기, pingInterval은 핑 전송 간격입니다.
func NewBaseClient(endpoint string, tokenGen auth.WebSocketTokenGenerator, pingInterval time.Duration) *BaseClient {
	return &BaseClient{
		Endpoint:     endpoint,
		TokenGen:     tokenGen,
		PingInterval: pingInterval,
	}
}

// Connect는 웹소켓 서버에 연결합니다.
// 연결에 실패하면 에러를 반환합니다.
func (c *BaseClient) Connect() error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.IsRunning {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	token, err := c.TokenGen.GenerateToken()
	if err != nil {
		cancel()
		return fmt.Errorf("토큰 생성 실패: %w", err)
	}

	conn, _, err := websocket.Dial(ctx, c.Endpoint, &websocket.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {token},
		},
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		cancel()
		return fmt.Errorf("웹소켓 연결 실패: %w", err)
	}

	c.Conn = conn
	c.Ctx = ctx
	c.Cancel = cancel
	c.IsRunning = true

	c.startPingLoop()
	return nil
}

// Ping은 웹소켓 서버에 핑을 전송합니다.
// 연결이 없으면 에러를 반환합니다.
func (c *BaseClient) Ping() error {
	c.Mu.Lock()
	conn := c.Conn
	ctx := c.Ctx
	c.Mu.Unlock()

	if conn == nil {
		return fmt.Errorf("웹소켓 연결이 없습니다")
	}
	return conn.Ping(ctx)
}

// Close는 웹소켓 연결을 종료합니다.
// 연결 종료에 실패하면 에러를 반환합니다.
func (c *BaseClient) Close() error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if !c.IsRunning {
		return nil
	}

	c.IsRunning = false
	if c.PingTicker != nil {
		c.PingTicker.Stop()
		c.PingTicker = nil
	}

	if c.Conn != nil {
		err := c.Conn.Close(websocket.StatusNormalClosure, "정상 종료")
		if c.Cancel != nil {
			c.Cancel()
			c.Cancel = nil
		}
		c.Conn = nil
		return err
	}
	return nil
}

// Reconnect는 웹소켓 연결을 재시도합니다.
// 재연결에 실패하면 에러를 반환합니다.
func (c *BaseClient) Reconnect() error {
	if err := c.Close(); err != nil {
		return fmt.Errorf("재연결 중 종료 실패: %w", err)
	}
	return c.Connect()
}

// startPingLoop는 주기적으로 핑을 전송하는 루프를 시작합니다.
func (c *BaseClient) startPingLoop() {
	if c.PingInterval == 0 {
		return
	}

	c.PingTicker = time.NewTicker(c.PingInterval)

	go func() {
		for {
			select {
			case <-c.Ctx.Done():
				return
			case <-c.PingTicker.C:
				if err := c.Ping(); err != nil {
					fmt.Printf("핑 전송 실패: %v\n", err)
					if err := c.Reconnect(); err != nil {
						fmt.Printf("재연결 실패: %v\n", err)
					}
				}
			}
		}
	}()
}
