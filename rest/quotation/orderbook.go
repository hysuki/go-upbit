package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// OrderbookUnit은 호가 정보를 나타냅니다.
type OrderbookUnit struct {
	AskPrice float64 `json:"ask_price"` // 매도 호가
	BidPrice float64 `json:"bid_price"` // 매수 호가
	AskSize  float64 `json:"ask_size"`  // 매도 잔량
	BidSize  float64 `json:"bid_size"`  // 매수 잔량
}

// Orderbook은 호가 정보를 나타냅니다.
type Orderbook struct {
	Market         string          `json:"market"`          // 마켓 코드
	Timestamp      int64           `json:"timestamp"`       // 호가 생성 시각
	TotalAskSize   float64         `json:"total_ask_size"`  // 호가 매도 총 잔량
	TotalBidSize   float64         `json:"total_bid_size"`  // 호가 매수 총 잔량
	OrderbookUnits []OrderbookUnit `json:"orderbook_units"` // 호가 정보
	Level          float64         `json:"level"`           // 호가 모아보기 단위 (0: 기본 호가단위)
}

// SupportedLevel은 호가 모아보기 단위 정보를 나타냅니다.
type SupportedLevel struct {
	Market          string    `json:"market"`           // 마켓 코드
	SupportedLevels []float64 `json:"supported_levels"` // 지원하는 모아보기 단위 (0: 기본 호가단위)
}

// GetOrderbooks는 호가 정보를 조회합니다.
// markets는 마켓 코드 목록, level은 호가 모아보기 단위입니다.
func (q *Quotation) GetOrderbooks(markets []string, level float64) ([]Orderbook, error) {
	if len(markets) == 0 {
		return nil, errors.New("markets is required")
	}

	params := map[string]string{
		"markets": strings.Join(markets, ","),
	}

	if level != 0 {
		params["level"] = fmt.Sprintf("%v", level)
	}

	resp, err := q.Client.Get("/orderbook", params)
	if err != nil {
		return nil, err
	}

	var orderbooks []Orderbook
	if err := json.Unmarshal(resp, &orderbooks); err != nil {
		return nil, err
	}

	return orderbooks, nil
}

// GetSupportedLevels는 호가 모아보기 단위 정보를 조회합니다.
// 원화마켓(KRW)에서만 호가 모아보기 기능을 지원합니다.
func (q *Quotation) GetSupportedLevels() ([]SupportedLevel, error) {
	resp, err := q.Client.Get("/orderbook/supported_levels", nil)
	if err != nil {
		return nil, err
	}

	var levels []SupportedLevel
	if err := json.Unmarshal(resp, &levels); err != nil {
		return nil, err
	}

	return levels, nil
}
