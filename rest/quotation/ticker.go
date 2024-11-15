package quotation

import (
	"encoding/json"
	"errors"
	"strings"
)

// Ticker는 현재가 정보를 나타내는 구조체입니다.
type Ticker struct {
	// Market은 종목 구분 코드입니다.
	Market string `json:"market"`
	// TradeDate는 최근 거래 일자(UTC)입니다. (포맷: yyyyMMdd)
	TradeDate string `json:"trade_date"`
	// TradeTime는 최근 거래 시각(UTC)입니다. (포맷: HHmmss)
	TradeTime string `json:"trade_time"`
	// TradeDateKst는 최근 거래 일자(KST)입니다. (포맷: yyyyMMdd)
	TradeDateKst string `json:"trade_date_kst"`
	// TradeTimeKst는 최근 거래 시각(KST)입니다. (포맷: HHmmss)
	TradeTimeKst string `json:"trade_time_kst"`
	// TradeTimestamp는 최근 거래 일시(UTC)입니다. (포맷: Unix Timestamp)
	TradeTimestamp int64 `json:"trade_timestamp"`
	// OpeningPrice는 시가입니다.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice는 고가입니다.
	HighPrice float64 `json:"high_price"`
	// LowPrice는 저가입니다.
	LowPrice float64 `json:"low_price"`
	// TradePrice는 종가(현재가)입니다.
	TradePrice float64 `json:"trade_price"`
	// PrevClosingPrice는 전일 종가입니다. (UTC 0시 기준)
	PrevClosingPrice float64 `json:"prev_closing_price"`
	// Change는 EVEN(보합), RISE(상승), FALL(하락)을 나타냅니다.
	Change string `json:"change"`
	// ChangePrice는 변화액의 절대값입니다.
	ChangePrice float64 `json:"change_price"`
	// ChangeRate는 변화율의 절대값입니다.
	ChangeRate float64 `json:"change_rate"`
	// SignedChangePrice는 부호가 있는 변화액입니다.
	SignedChangePrice float64 `json:"signed_change_price"`
	// SignedChangeRate는 부호가 있는 변화율입니다.
	SignedChangeRate float64 `json:"signed_change_rate"`
	// TradeVolume는 가장 최근 거래량입니다.
	TradeVolume float64 `json:"trade_volume"`
	// AccTradePrice는 누적 거래대금(UTC 0시 기준)입니다.
	AccTradePrice float64 `json:"acc_trade_price"`
	// AccTradePrice24h는 24시간 누적 거래대금입니다.
	AccTradePrice24h float64 `json:"acc_trade_price_24h"`
	// AccTradeVolume는 누적 거래량(UTC 0시 기준)입니다.
	AccTradeVolume float64 `json:"acc_trade_volume"`
	// AccTradeVolume24h는 24시간 누적 거래량입니다.
	AccTradeVolume24h float64 `json:"acc_trade_volume_24h"`
	// Highest52WeekPrice는 52주 신고가입니다.
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	// Highest52WeekDate는 52주 신고가 달성일입니다. (포맷: yyyy-MM-dd)
	Highest52WeekDate string `json:"highest_52_week_date"`
	// Lowest52WeekPrice는 52주 신저가입니다.
	Lowest52WeekPrice float64 `json:"lowest_52_week_price"`
	// Lowest52WeekDate는 52주 신저가 달성일입니다. (포맷: yyyy-MM-dd)
	Lowest52WeekDate string `json:"lowest_52_week_date"`
	// Timestamp는 타임스탬프입니다.
	Timestamp int64 `json:"timestamp"`
}

// GetTicker는 요청 당시 종목의 스냅샷을 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/ticker
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
// 엔드포인트: https://api.upbit.com/v1/ticker/all
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
