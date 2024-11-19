package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Package quotation은 Upbit 거래소의 시세 조회 관련 API를 제공합니다.

// Trade는 체결 내역 정보를 나타냅니다.
type Trade struct {
	Market           string  `json:"market"`             // 종목 코드
	TradeDate        string  `json:"trade_date_utc"`     // 체결 일자(UTC 기준)
	TradeTime        string  `json:"trade_time_utc"`     // 체결 시각(UTC 기준)
	Timestamp        int64   `json:"timestamp"`          // 체결 타임스탬프
	TradePrice       float64 `json:"trade_price"`        // 체결 가격
	TradeVolume      float64 `json:"trade_volume"`       // 체결량
	PrevClosingPrice float64 `json:"prev_closing_price"` // 전일 종가
	ChangePrice      float64 `json:"change_price"`       // 변화량
	AskBid           string  `json:"ask_bid"`            // 매도/매수
	SequentialID     int64   `json:"sequential_id"`      // 체결 번호(Unique)
}

// GetTrades는 최근 체결 내역을 조회합니다.
// market은 마켓 코드, to는 마지막 체결 시각, count는 체결 개수입니다.
// cursor는 페이지네이션 커서, daysAgo는 최근 체결 날짜 기준 7일 이내의 이전 데이터 조회를 위한 파라미터입니다.
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
