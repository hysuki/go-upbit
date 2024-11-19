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
	Conn         *websocket.Conn
	Ctx          context.Context
	Cancel       context.CancelFunc
	IsRunning    bool
	Mu           sync.Mutex
	Messages     []Message
	Endpoint     string
	TokenGen     auth.WebSocketTokenGenerator
	PingTicker   *time.Ticker
	PingInterval time.Duration
}

func NewBaseClient(endpoint string, tokenGen auth.WebSocketTokenGenerator, pingInterval time.Duration) *BaseClient {
	return &BaseClient{
		Endpoint:     endpoint,
		TokenGen:     tokenGen,
		PingInterval: pingInterval,
	}
}

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

func (c *BaseClient) Reconnect() error {
	if err := c.Close(); err != nil {
		return fmt.Errorf("재연결 중 종료 실패: %w", err)
	}
	return c.Connect()
}

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
