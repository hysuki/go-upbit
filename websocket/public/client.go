package public

import (
	"time"

	"github.com/hysuki/go-upbit/auth"
	"github.com/hysuki/go-upbit/websocket"
)

// Client는 공개 WebSocket API 클라이언트입니다.
type Client struct {
	*websocket.BaseClient
}

// NewClient는 새로운 공개 WebSocket 클라이언트를 생성합니다.
// endpoint는 WebSocket 서버 주소, tokenGen은 토큰 생성기, pingInterval은 핑 전송 간격입니다.
func NewClient(endpoint string, tokenGen *auth.WebSocketTokenGen, pingInterval time.Duration) (*Client, error) {
	base := websocket.NewBaseClient(endpoint, tokenGen, pingInterval)

	client := &Client{
		BaseClient: base,
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// Message는 WebSocket 메시지를 처리합니다.
func (c *Client) Message() error {
	return nil
}
