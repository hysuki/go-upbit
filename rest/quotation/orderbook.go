package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// OrderbookUnit은 호가 정보를 나타내는 구조체입니다.
type OrderbookUnit struct {
	// AskPrice는 매도 호가입니다.
	AskPrice float64 `json:"ask_price"`
	// BidPrice는 매수 호가입니다.
	BidPrice float64 `json:"bid_price"`
	// AskSize는 매도 잔량입니다.
	AskSize float64 `json:"ask_size"`
	// BidSize는 매수 잔량입니다.
	BidSize float64 `json:"bid_size"`
}

// Orderbook은 호가 정보를 나타내는 구조체입니다.
type Orderbook struct {
	// Market은 마켓 코드입니다.
	Market string `json:"market"`
	// Timestamp는 호가 생성 시각입니다.
	Timestamp int64 `json:"timestamp"`
	// TotalAskSize는 호가 매도 총 잔량입니다.
	TotalAskSize float64 `json:"total_ask_size"`
	// TotalBidSize는 호가 매수 총 잔량입니다.
	TotalBidSize float64 `json:"total_bid_size"`
	// OrderbookUnits는 호가 정보입니다.
	OrderbookUnits []OrderbookUnit `json:"orderbook_units"`
	// Level은 호가 모아보기 단위입니다. (0: 기본 호가단위)
	Level float64 `json:"level"`
}

// GetOrderbooks는 호가 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orderbook
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

// SupportedLevel은 호가 모아보기 단위 정보를 나타내는 구조체입니다.
type SupportedLevel struct {
	// Market은 마켓 코드입니다.
	Market string `json:"market"`
	// SupportedLevels는 해당 종목에서 지원하는 모아보기 단위입니다.
	// 0: 기본 호가단위
	// 호가 모아보기 기능은 원화마켓(KRW)에서만 지원하므로 BTC, USDT 마켓의 경우 0만 존재합니다.
	SupportedLevels []float64 `json:"supported_levels"`
}

// GetSupportedLevels는 호가 모아보기 단위 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orderbook/supported_levels
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
