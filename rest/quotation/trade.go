package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Trade는 체결 내역 정보를 나타내는 구조체입니다.
type Trade struct {
	// Market은 종목 코드입니다.
	Market string `json:"market"`
	// TradeDate는 체결 일자(UTC 기준)입니다.
	TradeDate string `json:"trade_date_utc"`
	// TradeTime은 체결 시각(UTC 기준)입니다.
	TradeTime string `json:"trade_time_utc"`
	// Timestamp는 체결 타임스탬프입니다.
	Timestamp int64 `json:"timestamp"`
	// TradePrice는 체결 가격입니다.
	TradePrice float64 `json:"trade_price"`
	// TradeVolume은 체결량입니다.
	TradeVolume float64 `json:"trade_volume"`
	// PrevClosingPrice는 전일 종가(UTC 0시 기준)입니다.
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// ChangePrice는 변화량입니다.
	ChangePrice float64 `json:"change_price"`
	// AskBid는 매도/매수입니다.
	AskBid string `json:"ask_bid"`
	// SequentialID는 체결 번호(Unique)입니다.
	SequentialID int64 `json:"sequential_id"`
}

// GetTrades는 최근 체결 내역을 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/trades/ticks
func (q *Quotation) GetTrades(market string, to string, count int, cursor string, daysAgo int) ([]Trade, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}

	params := map[string]string{
		"market": market,
	}

	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		if count > 500 {
			count = 500 // 최대 500개로 제한
		}
		params["count"] = fmt.Sprintf("%d", count)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	if daysAgo > 0 {
		if daysAgo > 7 {
			return nil, errors.New("days_ago cannot exceed 7")
		}
		params["days_ago"] = fmt.Sprintf("%d", daysAgo)
	}

	resp, err := q.Client.Get("/trades/ticks", params)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
