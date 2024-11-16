package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/hysuki/go-upbit/auth"
)

// BaseClient은 WebSocket 연결을 관리하는 구조체입니다
type BaseClient struct {
	conn         *websocket.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	isRunning    bool
	mu           *sync.Mutex
	endpoint     string
	tokenGen     auth.BaseTokenGenerator
	pingTicker   *time.Ticker
	pingInterval time.Duration
}

// NewBaseClient은 새로운 BaseClient을 생성합니다
func NewBaseClient(endpoint string, tokenGen auth.BaseTokenGenerator, pingInterval time.Duration) *BaseClient {
	return &BaseClient{
		endpoint:     endpoint,
		tokenGen:     tokenGen,
		pingInterval: pingInterval,
		mu:           &sync.Mutex{},
	}
}

// Connect implements types.WSConn
func (c *BaseClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

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

	c.startPingLoop()

	return nil
}

// Ping implements types.WSConn
func (c *BaseClient) Ping() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("[Ping] 웹소켓 연결이 없습니다")
	}
	return c.conn.Ping(c.ctx)
}

// Close implements types.WSConn
func (c *BaseClient) Close() error {
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

// Reconnect는 웹소켓 연결을 재연결합니다
func (c *BaseClient) Reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.Close(); err != nil {
		return fmt.Errorf("close failed during reconnect: %w", err)
	}
	return c.Connect()
}

// startPingLoop는 주기적으로 핑을 보내는 고루틴을 시작합니다
func (c *BaseClient) startPingLoop() {
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
