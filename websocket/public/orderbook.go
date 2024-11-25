package public

import (
	"encoding/json"
	"fmt"
	"time"
)

// UpbitOrderbook는 호가 정보를 나타냅니다.
type UpbitOrderbook struct {
	Type           string          `json:"type"`            // 타입 (orderbook)
	Code           string          `json:"code"`            // 마켓 코드
	TotalAskSize   float64         `json:"total_ask_size"`  // 호가 매도 총 잔량
	TotalBidSize   float64         `json:"total_bid_size"`  // 호가 매수 총 잔량
	OrderbookUnits []OrderbookUnit `json:"orderbook_units"` // 호가 정보 목록
	Timestamp      int64           `json:"timestamp"`       // 타임스탬프
	Level          float64         `json:"level"`           // 호가 모아보기 단위
}

// Orderbook은 내부적으로 사용하기 위한 호가 정보 구조체입니다.
type Orderbook struct {
	Type           string          `json:"type"`            // 타입 (orderbook)
	Code           string          `json:"code"`            // 마켓 코드
	TotalAskSize   float64         `json:"total_ask_size"`  // 호가 매도 총 잔량
	TotalBidSize   float64         `json:"total_bid_size"`  // 호가 매수 총 잔량
	OrderbookUnits []OrderbookUnit `json:"orderbook_units"` // 호가 정보 목록
	Timestamp      time.Time       `json:"timestamp"`       // 타임스탬프 (KST)
	Level          float64         `json:"level"`           // 호가 모아보기 단위
}

// OrderbookUnit은 내부적으로 사용하기 위한 개별 호가 정보 구조체입니다.
type OrderbookUnit struct {
	AskPrice float64 `json:"ask_price"` // 매도 호가
	BidPrice float64 `json:"bid_price"` // 매수 호가
	AskSize  float64 `json:"ask_size"`  // 매도 잔량
	BidSize  float64 `json:"bid_size"`  // 매수 잔량
}

// NewOrderbook은 UpbitOrderbook을 내부 Orderbook 구조체로 변환합니다.
func NewOrderbook(u *UpbitOrderbook, loc *time.Location) *Orderbook {
	// loc이 nil인 경우 UTC를 사용
	if loc == nil {
		loc = time.UTC
	}

	units := make([]OrderbookUnit, len(u.OrderbookUnits))
	for i, unit := range u.OrderbookUnits {
		units[i] = OrderbookUnit{
			AskPrice: unit.AskPrice,
			BidPrice: unit.BidPrice,
			AskSize:  unit.AskSize,
			BidSize:  unit.BidSize,
		}
	}

	return &Orderbook{
		Type:           u.Type,
		Code:           u.Code,
		TotalAskSize:   u.TotalAskSize,
		TotalBidSize:   u.TotalBidSize,
		OrderbookUnits: units,
		Timestamp:      time.UnixMilli(u.Timestamp).In(loc),
		Level:          u.Level,
	}
}

// ParseOrderBook은 JSON 데이터를 UpbitOrderbook 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseOrderBook(data []byte) (*UpbitOrderbook, error) {
	var orderbook UpbitOrderbook
	if err := json.Unmarshal(data, &orderbook); err != nil {
		return nil, fmt.Errorf("호가 데이터 파싱 실패: %v", err)
	}
	return &orderbook, nil
}

// GetOrderBook은 다음 호가 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 호가 정보를 반환합니다.
func (c *Client) GetOrderBook(loc *time.Location) (*Orderbook, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.orderBookChan:
		return NewOrderbook(resp, loc), nil
	}
}
