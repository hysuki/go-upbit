package upbit

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/coder/websocket"
// )

// // PrivateWSClient는 인증이 필요한 WebSocket 연결을 관리합니다.
// type PrivateWSClient struct {
// 	client
// }

// // PublicWSClient는 공개 WebSocket 연결을 관리합니다.
// type PublicWSClient struct {
// 	client
// }

// // client는 WebSocket 연결의 기본 구현을 제공합니다.
// type client struct {
// 	conn        *websocket.Conn
// 	ctx         context.Context
// 	cancel      context.CancelFunc
// 	isRunning   bool
// 	mu          sync.Mutex
// 	endpoint    string
// 	upbitClient *UpbitClient
// 	pingTicker  *time.Ticker
// }

// // NewClient는 지정된 타입의 새로운 클라이언트 인스턴스를 생성합니다.
// func NewClient[T ClientType](c *UpbitClient) (*T, error) {
// 	return Connect[T](c, &client{})
// }

// // Ping은 WebSocket 연결의 상태를 확인합니다.
// func (c *client) Ping() error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	if c.conn == nil {
// 		return fmt.Errorf("[Ping] 웹소켓 연결이 없습니다")
// 	}

// 	return c.conn.Ping(c.ctx)
// }

// // getEndpoint는 클라이언트 타입에 따른 WebSocket 엔드포인트를 반환합니다.
// func getEndpoint[T ClientType]() string {
// 	var zero T
// 	switch any(zero).(type) {
// 	case PublicWSClient:
// 		return PublicWebsocketEndpoint
// 	case PrivateWSClient:
// 		return PrivateWebsocketEndpoint
// 	default:
// 		return ""
// 	}
// }

// // Connect는 WebSocket 연결을 생성하고 클라이언트를 초기화합니다.
// // 기존 연결이 있다면 정리하고 새로운 연결을 설정합니다.
// func Connect[T ClientType](upbitClient *UpbitClient, c *client) (*T, error) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	if c.conn != nil {
// 		c.Close()
// 	}
// 	if c.cancel != nil {
// 		c.cancel()
// 	}

// 	c.endpoint = getEndpoint[T]()
// 	c.upbitClient = upbitClient

// 	ctx, cancel := context.WithCancel(context.Background())
// 	token, err := upbitClient.generateToken()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	opts := &websocket.DialOptions{
// 		HTTPHeader: map[string][]string{
// 			"Authorization": {token},
// 		},
// 		CompressionMode: websocket.CompressionContextTakeover,
// 	}

// 	conn, _, err := websocket.Dial(ctx, c.endpoint, opts)
// 	if err != nil {
// 		cancel()
// 		return nil, fmt.Errorf("websocket dial failed: %w", err)
// 	}

// 	c.conn = conn
// 	c.isRunning = true
// 	c.ctx = ctx
// 	c.cancel = cancel

// 	// pingLoop 시작
// 	c.startPingLoop()

// 	return &T{
// 		client: *c,
// 	}, nil
// }

// // Reconnect는 저장된 정보를 사용하여 WebSocket 연결을 재설정합니다.
// func (c *client) Reconnect() error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	if c.conn != nil {
// 		c.Close()
// 	}
// 	if c.cancel != nil {
// 		c.cancel()
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	token, err := c.upbitClient.generateToken()
// 	if err != nil {
// 		return fmt.Errorf("failed to generate token: %w", err)
// 	}

// 	opts := &websocket.DialOptions{
// 		HTTPHeader: map[string][]string{
// 			"Authorization": {token},
// 		},
// 		CompressionMode: websocket.CompressionContextTakeover,
// 	}

// 	conn, _, err := websocket.Dial(ctx, c.endpoint, opts)
// 	if err != nil {
// 		cancel()
// 		return fmt.Errorf("websocket dial failed: %w", err)
// 	}

// 	c.conn = conn
// 	c.isRunning = true
// 	c.ctx = ctx
// 	c.cancel = cancel

// 	return nil
// }

// // Close는 WebSocket 연결을 안전하게 종료합니다.
// func (c *client) Close() error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	c.isRunning = false
// 	if c.pingTicker != nil {
// 		c.pingTicker.Stop()
// 	}
// 	if c.conn != nil {
// 		err := c.conn.Close(websocket.StatusNormalClosure, "정상 종료")
// 		c.cancel()
// 		return err
// 	}
// 	return nil
// }

// // startPingLoop는 주기적으로 핑을 보내는 고루틴을 시작합니다.
// // ex. 120초 연결 타임아웃을 방지하기 위해 50초마다 핑을 전송합니다.
// func (c *client) startPingLoop() {
// 	// pingInterval이 0이면 핑을 보내지 않음
// 	if c.upbitClient.pingInterval == 0 {
// 		return
// 	}

// 	if c.pingTicker != nil {
// 		c.pingTicker.Stop()
// 	}

// 	c.pingTicker = time.NewTicker(c.upbitClient.pingInterval)

// 	go func() {
// 		for {
// 			select {
// 			case <-c.ctx.Done():
// 				if c.pingTicker != nil {
// 					c.pingTicker.Stop()
// 				}
// 				return
// 			case <-c.pingTicker.C:
// 				if err := c.Ping(); err != nil {
// 					fmt.Printf("핑 전송 실패: %v\n", err)
// 					if err := c.Reconnect(); err != nil {
// 						fmt.Printf("재연결 실패: %v\n", err)
// 					}
// 				}
// 			}
// 		}
// 	}()
// }
