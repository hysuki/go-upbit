package private

import (
	"time"

	"github.com/hysuki/go-upbit/auth"
	"github.com/hysuki/go-upbit/websocket"
)

type Client struct {
	*websocket.BaseClient
}

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
