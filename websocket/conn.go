package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/hysuki/go-upbit/auth"
)

type BaseClient struct {
	conn         *websocket.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	isRunning    bool
	mu           sync.Mutex
	endpoint     string
	tokenGen     auth.WebSocketTokenGenerator
	pingTicker   *time.Ticker
	pingInterval time.Duration
}

func NewBaseClient(endpoint string, tokenGen auth.WebSocketTokenGenerator, pingInterval time.Duration) *BaseClient {
	return &BaseClient{
		endpoint:     endpoint,
		tokenGen:     tokenGen,
		pingInterval: pingInterval,
	}
}

func (c *BaseClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	token, err := c.tokenGen.GenerateToken()
	if err != nil {
		cancel()
		return fmt.Errorf("토큰 생성 실패: %w", err)
	}

	conn, _, err := websocket.Dial(ctx, c.endpoint, &websocket.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {token},
		},
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		cancel()
		return fmt.Errorf("웹소켓 연결 실패: %w", err)
	}

	c.conn = conn
	c.ctx = ctx
	c.cancel = cancel
	c.isRunning = true

	c.startPingLoop()
	return nil
}

func (c *BaseClient) Ping() error {
	c.mu.Lock()
	conn := c.conn
	ctx := c.ctx
	c.mu.Unlock()

	if conn == nil {
		return fmt.Errorf("웹소켓 연결이 없습니다")
	}
	return conn.Ping(ctx)
}

func (c *BaseClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isRunning {
		return nil
	}

	c.isRunning = false
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	if c.conn != nil {
		err := c.conn.Close(websocket.StatusNormalClosure, "정상 종료")
		if c.cancel != nil {
			c.cancel()
			c.cancel = nil
		}
		c.conn = nil
		return err
	}
	return nil
}

func (c *BaseClient) Reconnect() error {
	if err := c.Close(); err != nil {
		return fmt.Errorf("재연결 중 종료 실패: %w", err)
	}
	return c.Connect()
}

func (c *BaseClient) startPingLoop() {
	if c.pingInterval == 0 {
		return
	}

	c.pingTicker = time.NewTicker(c.pingInterval)

	go func() {
		for {
			select {
			case <-c.ctx.Done():
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
