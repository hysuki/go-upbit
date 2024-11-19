// Package public은 Upbit 거래소의 공개 웹소켓 API를 제공합니다.
package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/auth"
	"github.com/hysuki/go-upbit/websocket"
	"github.com/hysuki/go-upbit/websocket/common"
)

// Client는 공개 웹소켓 클라이언트입니다.
type Client struct {
	*websocket.BaseClient
	orderBookChan chan *OrderBookResponse
	tickerChan    chan *TickerResponse
	tradeChan     chan *TradeResponse
	errChan       chan error
	done          chan struct{}
}

// MessageType은 메시지 유형을 나타냅니다.
type MessageType string

// 메시지 유형을 정의하는 상수들입니다.
const (
	Ticker    MessageType = "ticker"
	Orderbook MessageType = "orderbook"
	Trade     MessageType = "trade"
)

// Message는 웹소켓 메시지를 나타냅니다.
type Message struct {
	websocket.Message
	Type MessageType `json:"type,omitempty"`
}

// NewClient는 새로운 공개 웹소켓 클라이언트를 생성합니다.
// endpoint는 웹소켓 서버 주소, tokenGen은 토큰 생성기, pingInterval은 핑 전송 간격입니다.
func NewClient(endpoint string, tokenGen *auth.WebSocketTokenGen, pingInterval time.Duration) (*Client, error) {
	base := websocket.NewBaseClient(endpoint, tokenGen, pingInterval)

	client := &Client{
		BaseClient:    base,
		orderBookChan: make(chan *OrderBookResponse, 1000),
		tickerChan:    make(chan *TickerResponse, 1000),
		tradeChan:     make(chan *TradeResponse, 1000),
		errChan:       make(chan error, 1000),
		done:          make(chan struct{}),
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// AddSubscribe는 구독 함수를 생성합니다.
// messageType은 메시지 유형, codes는 마켓 코드 목록, options는 구독 옵션입니다.
func AddSubscribe(messageType MessageType, codes []string, options *common.SubscribeOptions) websocket.SubscribeFunc {
	if len(codes) == 0 {
		return func(c *websocket.BaseClient) error {
			return fmt.Errorf("codes는 최소 하나 이상의 마켓 코드를 포함해야 합니다")
		}
	}
	return websocket.AddSubscribe(string(messageType), codes, options)
}

// Subscribe는 지정된 구독 함수들을 사용하여 구독을 시작합니다.
// ticket은 구독 식별자, f는 구독 함수 목록입니다.
func (c *Client) Subscribe(ticket *string, f ...websocket.SubscribeFunc) error {
	if len(f) == 0 {
		return fmt.Errorf("구독 함수가 제공되지 않았습니다")
	}
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
				case string(Orderbook):
					if resp, err := ParseOrderBook(data); err != nil {
						c.errChan <- err
					} else {
						c.orderBookChan <- resp
					}
				case string(Ticker):
					if resp, err := ParseTicker(data); err != nil {
						c.errChan <- err
					} else {
						c.tickerChan <- resp
					}
				case string(Trade):
					if resp, err := ParseTrade(data); err != nil {
						c.errChan <- err
					} else {
						c.tradeChan <- resp
					}
				}
			}
		}
	}()
}

// Stop은 클라이언트를 종료합니다.
func (c *Client) Stop() {
	close(c.done)
}
