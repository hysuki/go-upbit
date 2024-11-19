package quotation

import (
	"encoding/json"
	"errors"
	"strings"
)

// Package quotation은 Upbit 거래소의 시세 조회 관련 API를 제공합니다.
type Ticker struct {
	Market             string  `json:"market"`                // 종목 구분 코드
	TradeDate          string  `json:"trade_date"`            // 최근 거래 일자(UTC) (yyyyMMdd)
	TradeTime          string  `json:"trade_time"`            // 최근 거래 시각(UTC) (HHmmss)
	TradeDateKst       string  `json:"trade_date_kst"`        // 최근 거래 일자(KST) (yyyyMMdd)
	TradeTimeKst       string  `json:"trade_time_kst"`        // 최근 거래 시각(KST) (HHmmss)
	TradeTimestamp     int64   `json:"trade_timestamp"`       // 최근 거래 일시(UTC) (Unix Timestamp)
	OpeningPrice       float64 `json:"opening_price"`         // 시가
	HighPrice          float64 `json:"high_price"`            // 고가
	LowPrice           float64 `json:"low_price"`             // 저가
	TradePrice         float64 `json:"trade_price"`           // 종가(현재가)
	PrevClosingPrice   float64 `json:"prev_closing_price"`    // 전일 종가
	Change             string  `json:"change"`                // 전일 대비 (EVEN: 보합, RISE: 상승, FALL: 하락)
	ChangePrice        float64 `json:"change_price"`          // 변화액의 절대값
	ChangeRate         float64 `json:"change_rate"`           // 변화율의 절대값
	SignedChangePrice  float64 `json:"signed_change_price"`   // 부호가 있는 변화액
	SignedChangeRate   float64 `json:"signed_change_rate"`    // 부호가 있는 변화율
	TradeVolume        float64 `json:"trade_volume"`          // 가장 최근 거래량
	AccTradePrice      float64 `json:"acc_trade_price"`       // 누적 거래대금(UTC 0시 기준)
	AccTradePrice24h   float64 `json:"acc_trade_price_24h"`   // 24시간 누적 거래대금
	AccTradeVolume     float64 `json:"acc_trade_volume"`      // 누적 거래량(UTC 0시 기준)
	AccTradeVolume24h  float64 `json:"acc_trade_volume_24h"`  // 24시간 누적 거래량
	Highest52WeekPrice float64 `json:"highest_52_week_price"` // 52주 신고가
	Highest52WeekDate  string  `json:"highest_52_week_date"`  // 52주 신고가 달성일 (yyyy-MM-dd)
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`  // 52주 신저가
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`   // 52주 신저가 달성일 (yyyy-MM-dd)
	Timestamp          int64   `json:"timestamp"`             // 타임스탬프
}

// GetTicker는 요청 당시 종목의 스냅샷을 조회합니다.
// markets는 조회할 마켓 코드 목록입니다.
func (q *Quotation) GetTicker(markets []string) ([]Ticker, error) {
	if len(markets) == 0 {
		return nil, errors.New("markets is required")
	}

	params := map[string]string{
		"markets": strings.Join(markets, ","),
	}

	resp, err := q.Client.Get("/ticker", params)
	if err != nil {
		return nil, err
	}

	var tickers []Ticker
	if err := json.Unmarshal(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetTickersByQuote는 마켓 단위 종목들의 스냅샷을 조회합니다.
// quoteCurrencies는 기준 화폐 목록입니다. 미지정 시 모든 종목을 조회합니다.
func (q *Quotation) GetTickersByQuote(quoteCurrencies []string) ([]Ticker, error) {
	params := make(map[string]string)
	if len(quoteCurrencies) > 0 {
		params["quote_currencies"] = strings.Join(quoteCurrencies, ",")
	}

	resp, err := q.Client.Get("/ticker/all", params)
	if err != nil {
		return nil, err
	}

	var tickers []Ticker
	if err := json.Unmarshal(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
