package public

import (
	"encoding/json"
	"fmt"
)

// OrderBookResponse는 호가 정보를 담는 구조체입니다
type OrderBookResponse struct {
	Type           string          `json:"type"`            // 타입 (orderbook)
	Code           string          `json:"code"`            // 마켓 코드
	TotalAskSize   float64         `json:"total_ask_size"`  // 호가 매도 총 잔량
	TotalBidSize   float64         `json:"total_bid_size"`  // 호가 매수 총 잔량
	OrderBookUnits []OrderBookUnit `json:"orderbook_units"` // 호가 정보 목록
	Timestamp      int64           `json:"timestamp"`       // 타임스탬프
	Level          float64         `json:"level"`           // 호가 모아보기 단위
}

// OrderBookUnit은 개별 호가 정보를 담는 구조체입니다
type OrderBookUnit struct {
	AskPrice float64 `json:"ask_price"` // 매도 호가
	BidPrice float64 `json:"bid_price"` // 매수 호가
	AskSize  float64 `json:"ask_size"`  // 매도 잔량
	BidSize  float64 `json:"bid_size"`  // 매수 잔량
}

// ParseOrderBook는 JSON 데이터를 OrderBook 구조체로 파싱합니다
func ParseOrderBook(data []byte) (*OrderBookResponse, error) {
	var orderbook OrderBookResponse
	if err := json.Unmarshal(data, &orderbook); err != nil {
		return nil, fmt.Errorf("호가 데이터 파싱 실패: %v", err)
	}
	return &orderbook, nil
}

// SubscribeOrderBook는 지정된 마켓 코드들의 호가 정보를 구독합니다
// func (c *Client) SubscribeOrderBook(codes []string, options *common.SubscribeOptions) error {
// 	if len(codes) == 0 {
// 		return fmt.Errorf("마켓 코드는 필수입니다")
// 	}
// 	return c.Subscribe("", "orderbook", codes, options)
// }

// GetOrderBook는 수신된 메시지를 OrderBookResponse 구조체로 변환합니다
// func (c *Client) GetOrderBook(data []byte) (*OrderBookResponse, error) {

// 	// 타입 확인
// 	readMessage := websocket.ReadMessage{}
// 	if err := json.Unmarshal(data, &readMessage); err != nil {
// 		return nil, fmt.Errorf("타입 확인 실패: %v", err)
// 	}

// 	if readMessage.Type != string(Orderbook) {
// 		return nil, nil
// 	}

// 	return ParseOrderBook(data)
// }

// GetOrderBook waits for the next order book message
func (c *Client) GetOrderBook() (*OrderBookResponse, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.orderBookChan:
		return resp, nil
	}
}
