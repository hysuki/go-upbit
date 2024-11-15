package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	// PublicEndpoint는 공개 WebSocket API 엔드포인트입니다.
	PublicEndpoint = "wss://api.upbit.com/websocket/v1"
	// PrivateEndpoint는 인증이 필요한 WebSocket API 엔드포인트입니다.
	PrivateEndpoint = "wss://api.upbit.com/websocket/v1/private"
)

// TokenGenerator는 인증 토큰을 생성하는 인터페이스입니다.
type TokenGenerator interface {
	GenerateToken() (string, error)
}

// Client는 WebSocket 클라이언트 인터페이스를 정의합니다.
type Client interface {
	Ping() error
	Reconnect() error
	Close() error
}

// client는 WebSocket 연결의 기본 구현을 제공합니다.
type client struct {
	conn         *websocket.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	isRunning    bool
	mu           sync.Mutex
	endpoint     string
	tokenGen     TokenGenerator
	pingTicker   *time.Ticker
	pingInterval time.Duration
}

// NewClient는 새로운 WebSocket 클라이언트를 생성합니다.
func NewClient(endpoint string, tokenGen TokenGenerator, pingInterval time.Duration) (Client, error) {
	c := &client{
		endpoint:     endpoint,
		tokenGen:     tokenGen,
		pingInterval: pingInterval,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

// connect는 WebSocket 연결을 설정합니다.
func (c *client) connect() error {
	ctx, cancel := context.WithCancel(context.Background())

	token, err := c.tokenGen.GenerateToken()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to generate token: %w", err)
	}

	opts := &websocket.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {token},
		},
		CompressionMode: websocket.CompressionContextTakeover,
	}

	conn, _, err := websocket.Dial(ctx, c.endpoint, opts)
	if err != nil {
		cancel()
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	c.conn = conn
	c.ctx = ctx
	c.cancel = cancel
	c.isRunning = true

	// Start ping loop
	c.startPingLoop()

	return nil
}

// Ping은 WebSocket 연결의 상태를 확인합니다.
func (c *client) Ping() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("[Ping] 웹소켓 연결이 없습니다")
	}

	return c.conn.Ping(c.ctx)
}

// Reconnect는 저장된 정보를 사용하여 WebSocket 연결을 재설정합니다.
func (c *client) Reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.Close(); err != nil {
		return fmt.Errorf("close failed during reconnect: %w", err)
	}

	return c.connect()
}

// Close는 WebSocket 연결을 안전하게 종료합니다.
func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isRunning = false
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}
	if c.conn != nil {
		err := c.conn.Close(websocket.StatusNormalClosure, "정상 종료")
		c.cancel()
		return err
	}
	return nil
}

// startPingLoop는 주기적으로 핑을 보내는 고루틴을 시작합니다.
func (c *client) startPingLoop() {
	if c.pingInterval == 0 {
		return
	}

	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}

	c.pingTicker = time.NewTicker(c.pingInterval)

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				if c.pingTicker != nil {
					c.pingTicker.Stop()
				}
				return
			case <-c.pingTicker.C:
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
