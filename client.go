package upbit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// Client는 Upbit 웹소켓 클라이언트를 구현합니다
type Client struct {
	accessKey string             // Upbit API 액세스 키
	apiSecret string             // Upbit API 시크릿 키
	conn      *websocket.Conn    // 웹소켓 연결
	ctx       context.Context    // 컨텍스트
	cancel    context.CancelFunc // 컨텍스트 취소 함수
	isRunning bool               // 클라이언트 실행 상태
	mu        sync.Mutex         // 동시성 제어를 위한 뮤텍스
}

// NewClient는 새로운 Upbit 웹소켓 클라이언트를 생성합니다
func NewClient(accessKey, apiSecret string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		accessKey: accessKey,
		apiSecret: apiSecret,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Connect는 Upbit 웹소켓 서버에 연결합니다
func (c *Client) Connect() error {
	token, err := c.generateToken()
	if err != nil {
		return err
	}

	opts := &websocket.DialOptions{
		HTTPHeader: map[string][]string{
			"Authorization": {token},
		},
		CompressionMode: websocket.CompressionContextTakeover, // 웹소켓 압축 활성화
	}

	conn, _, err := websocket.Dial(c.ctx, "wss://api.upbit.com/websocket/v1", opts)
	if err != nil {
		return fmt.Errorf("웹소켓 연결 실패: %v", err)
	}

	c.conn = conn
	c.isRunning = true

	// 연결 유지를 위한 핑 고루틴 시작
	go c.startPingLoop()

	return nil
}

// startPingLoop는 주기적으로 핑을 보내는 고루틴을 시작합니다
// 120초 연결 타임아웃을 방지하기 위해 50초마다 핑을 전송합니다
func (c *Client) startPingLoop() {
	ticker := time.NewTicker(50 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.Ping(); err != nil {
				fmt.Printf("핑 전송 실패: %v\n", err)
				c.reconnect()
			}
		}
	}
}

// Ping은 서버로 핑 메시지를 전송합니다
// 먼저 PING 프레임을 시도하고, 실패시 텍스트 메시지를 사용합니다
func (c *Client) Ping() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("[Ping] 웹소켓 연결이 없습니다")
	}

	err := c.conn.Ping(c.ctx)
	if err != nil {
		// PING 프레임 실패시 텍스트 메시지로 시도
		return c.WriteJSON("PING")
	}
	return nil
}

// reconnect는 웹소켓 서버로의 재연결을 시도합니다
// 최대 5번까지 재시도하며, 실패시 클라이언트를 종료합니다
func (c *Client) reconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isRunning {
		return
	}

	fmt.Println("웹소켓 재연결 시도 중...")

	// 기존 연결 종료
	if c.conn != nil {
		c.conn.Close(websocket.StatusNormalClosure, "reconnecting")
	}

	// 재연결 시도 (최대 5회)
	for i := 0; i < 5; i++ {
		err := c.Connect()
		if err == nil {
			fmt.Println("웹소켓 재연결 성공")
			return
		}
		fmt.Printf("재연결 시도 %d 실패: %v\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}

	fmt.Println("재연결 실패, 클라이언트를 종료합니다")
	c.Close()
}

// Close는 웹소켓 연결을 종료합니다
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isRunning = false
	if c.conn != nil {
		err := c.conn.Close(websocket.StatusNormalClosure, "정상 종료")
		c.cancel()
		return err
	}
	return nil
}
