package websocket

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/hysuki/go-upbit/auth"
)

// BaseClient는 웹소켓 기본 클라이언트입니다.
type BaseClient struct {
	Conn              *websocket.Conn              // 웹소켓 연결
	Ctx               context.Context              // 컨텍스트
	Cancel            context.CancelFunc           // 컨텍스트 취소 함수
	IsRunning         bool                         // 실행 상태
	Mu                sync.Mutex                   // 뮤텍스
	Messages          []Message                    // 메시지 목록
	Endpoint          string                       // 웹소켓 서버 주소
	TokenGen          auth.WebSocketTokenGenerator // 토큰 생성기
	PingTicker        *time.Ticker                 // 핑 전송 타이머
	PingInterval      time.Duration                // 핑 전송 간격
	reconnectAttempts int                          // 재연결 시도 횟수
	maxReconnectTries int                          // 최대 재연결 시도 횟수
	reconnectWait     time.Duration                // 재연결 대기 시간
}

// NewBaseClient는 새로운 웹소켓 기본 클라이언트를 생성합니다.
// endpoint는 웹소켓 서버 주소, tokenGen은 토큰 생성기, pingInterval은 핑 전송 간격입니다.
func NewBaseClient(endpoint string, tokenGen auth.WebSocketTokenGenerator, pingInterval time.Duration) *BaseClient {
	return &BaseClient{
		Endpoint:          endpoint,
		TokenGen:          tokenGen,
		PingInterval:      pingInterval,
		maxReconnectTries: 5,               // 기본값 설정
		reconnectWait:     time.Second * 3, // 기본값 설정
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
	if !c.IsRunning || c.Conn == nil {
		c.Mu.Unlock()
		return fmt.Errorf("failed to ping: use of closed network connection")
	}
	conn := c.Conn
	ctx := c.Ctx
	c.Mu.Unlock()

	select {
	case <-ctx.Done():
		return fmt.Errorf("컨텍스트 취소됨")
	default:
		if err := conn.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping: %w", err)
		}
		return nil
	}
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

	if c.Cancel != nil {
		c.Cancel()
		c.Cancel = nil
	}

	if c.Conn != nil {
		// 이미 닫힌 연결인지 확인
		err := c.Conn.Close(websocket.StatusNormalClosure, "정상 종료")
		c.Conn = nil
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			return err
		}
	}

	return nil
}

// Reconnect는 웹소켓 연결을 재시도합니다.
// 재연결에 실패하면 에러를 반환합니다.
func (c *BaseClient) Reconnect() error {
	c.Mu.Lock()
	if c.reconnectAttempts >= c.maxReconnectTries {
		c.Mu.Unlock()
		return fmt.Errorf("최대 재연결 시도 횟수(%d) 초과", c.maxReconnectTries)
	}
	c.reconnectAttempts++
	currentAttempt := c.reconnectAttempts
	c.Mu.Unlock()

	// 지수 백오프로 대기
	waitTime := c.reconnectWait * time.Duration(currentAttempt)
	time.Sleep(waitTime)

	// Close 호출 전에 현재 상태 확인
	if err := c.Close(); err != nil {
		// Close 실패는 무시하고 계속 진행
		fmt.Printf("연결 종료 중 오류 발생 (무시됨): %v\n", err)
	}

	if err := c.Connect(); err != nil {
		return fmt.Errorf("재연결 실패: %w", err)
	}

	// 이전 구독 정보 복구
	if len(c.Messages) > 0 {
		if err := c.request(nil); err != nil {
			return fmt.Errorf("구독 복구 실패: %w", err)
		}
	}

	c.Mu.Lock()
	c.reconnectAttempts = 0
	c.Mu.Unlock()

	return nil
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
					if strings.Contains(err.Error(), "컨텍스트 취소됨") {
						return
					}
					fmt.Printf("핑 전송 실패: %v\n", err)

					// 연결이 닫혔거나 실패한 경우 재연결 시도
					if strings.Contains(err.Error(), "use of closed network connection") ||
						strings.Contains(err.Error(), "failed to ping") {
						for i := 0; i < c.maxReconnectTries; i++ {
							fmt.Printf("재연결 시도 %d/%d...\n", i+1, c.maxReconnectTries)
							if err := c.Reconnect(); err != nil {
								fmt.Printf("재연결 실패: %v\n", err)
								time.Sleep(c.reconnectWait * time.Duration(i+1))
								continue
							}
							fmt.Println("재연결 성공")
							break
						}
					}
				}
			}
		}
	}()
}
