package private

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/auth"
	"github.com/hysuki/go-upbit/websocket"
	"github.com/hysuki/go-upbit/websocket/common"
)

// Client는 개인 웹소켓 클라이언트입니다.
type Client struct {
	*websocket.BaseClient
	myOrderChan chan *MyOrderResponse
	myAssetChan chan *MyAssetResponse
	errChan     chan error
	done        chan struct{}
}

// MessageType은 메시지 유형을 나타냅니다.
type MessageType string

// 메시지 유형을 정의하는 상수들입니다.
const (
	MyOrder MessageType = "myOrder"
	MyAsset MessageType = "myAsset"
)

// Message는 웹소켓 메시지를 나타냅니다.
type Message struct {
	websocket.Message
	Type MessageType `json:"type,omitempty"`
}

// NewClient는 새로운 개인 웹소켓 클라이언트를 생성합니다.
// endpoint는 웹소켓 서버 주소, tokenGen은 토큰 생성기, pingInterval은 핑 전송 간격입니다.
func NewClient(endpoint string, tokenGen *auth.WebSocketTokenGen, pingInterval time.Duration) (*Client, error) {
	base := websocket.NewBaseClient(endpoint, tokenGen, pingInterval)

	client := &Client{
		BaseClient:  base,
		myOrderChan: make(chan *MyOrderResponse, 1000),
		myAssetChan: make(chan *MyAssetResponse, 1000),
		errChan:     make(chan error, 1000),
		done:        make(chan struct{}),
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// AddSubscribe는 구독 함수를 생성합니다.
// messageType은 메시지 유형, codes는 마켓 코드 목록, options는 구독 옵션입니다.
func AddSubscribe(messageType MessageType, codes []string, options *common.SubscribeOptions) websocket.SubscribeFunc {
	if messageType == MyAsset && len(codes) != 0 {
		return func(c *websocket.BaseClient) error {
			return fmt.Errorf("MyAsset 타입은 마켓 코드를 지정할 수 없습니다")
		}
	}
	return websocket.AddSubscribe(string(messageType), codes, options)
}

// Subscribe는 지정된 구독 함수들을 사용하여 구독을 시작합니다.
// ticket은 구독 식별자, f는 구독 함수 목록입니다.
func (c *Client) Subscribe(ticket *string, f ...websocket.SubscribeFunc) error {
	return c.BaseClient.Subscribe(ticket, f...)
}

// StartMessageHandler는 메시지 처리기를 시작합니다.
// 수신된 메시지를 적절한 채널로 전달합니다.
func (c *Client) StartMessageHandler() {
	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				data, err := c.ReadMessage()
				if err != nil {
					c.errChan <- err
					continue
				}

				// 타입 확인
				readMessage := websocket.ReadMessage{}
				if err := json.Unmarshal(data, &readMessage); err != nil {
					c.errChan <- fmt.Errorf("타입 확인 실패: %v", err)
					continue
				}

				// 메시지 타입에 따라 적절한 채널로 전송
				switch readMessage.Type {
				case string(MyOrder):
					if resp, err := ParseMyOrder(data); err != nil {
						c.errChan <- err
					} else {
						c.myOrderChan <- resp
					}
				case string(MyAsset):
					if resp, err := ParseMyAsset(data); err != nil {
						c.errChan <- err
					} else {
						c.myAssetChan <- resp
					}
				}
			}
		}
	}()
}
