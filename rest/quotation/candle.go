package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Package quotation은 Upbit 거래소의 시세 조회 관련 API를 제공합니다.

// 분봉 단위를 정의하는 상수들입니다.
const (
	CandleMinute1   = 1   // 1분봉
	CandleMinute3   = 3   // 3분봉
	CandleMinute5   = 5   // 5분봉
	CandleMinute10  = 10  // 10분봉
	CandleMinute15  = 15  // 15분봉
	CandleMinute30  = 30  // 30분봉
	CandleMinute60  = 60  // 60분봉
	CandleMinute240 = 240 // 240분봉
)

// Candle은 캔들 정보를 나타냅니다.
type Candle struct {
	Market               string  `json:"market"`                          // 마켓명
	CandleDateTimeUTC    string  `json:"candle_date_time_utc"`            // 캔들 기준 시각(UTC)
	CandleDateTimeKST    string  `json:"candle_date_time_kst"`            // 캔들 기준 시각(KST)
	OpeningPrice         float64 `json:"opening_price"`                   // 시가
	HighPrice            float64 `json:"high_price"`                      // 고가
	LowPrice             float64 `json:"low_price"`                       // 저가
	TradePrice           float64 `json:"trade_price"`                     // 종가
	Timestamp            int64   `json:"timestamp"`                       // 마지막 틱이 저장된 시각
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`          // 누적 거래 금액
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`         // 누적 거래량
	FirstDayOfPeriod     string  `json:"first_day_of_period,omitempty"`   // 캔들 기간의 첫 날
	Unit                 int     `json:"unit,omitempty"`                  // 분봉 단위
	ConvertedTradePrice  float64 `json:"converted_trade_price,omitempty"` // 종가 환산 화폐 단위로 환산된 가격
	PrevClosingPrice     float64 `json:"prev_closing_price,omitempty"`    // 전일 종가
	ChangePrice          float64 `json:"change_price,omitempty"`          // 전일 종가 대비 변화 금액
	ChangeRate           float64 `json:"change_rate,omitempty"`           // 전일 종가 대비 변화량
}

// GetCandlesMinute는 분(Minute) 캔들을 조회합니다.
// unit은 분봉 단위(1, 3, 5, 10, 15, 30, 60, 240)를 지정합니다.
// market은 마켓 코드, to는 마지막 캔들 시각, count는 조회할 캔들 개수입니다.
func (q *Quotation) GetCandlesMinute(unit int, market string, to string, count int) ([]Candle, error) {
	if !isValidMinuteUnit(unit) {
		return nil, fmt.Errorf("invalid minute unit: %d", unit)
	}
	if market == "" {
		return nil, errors.New("market is required")
	}
	if count > 200 {
		count = 200
	}

	params := map[string]string{
		"market": market,
	}
	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		params["count"] = fmt.Sprintf("%d", count)
	}

	resp, err := q.Client.Get(fmt.Sprintf("/candles/minutes/%d", unit), params)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

// GetCandlesDay는 일(Day) 캔들을 조회합니다.
// market은 마켓 코드, to는 마지막 캔들 시각, count는 조회할 캔들 개수입니다.
// convertingPriceUnit으로 종가 환산 화폐 단위를 지정할 수 있습니다.
func (q *Quotation) GetCandlesDay(market string, to string, count int, convertingPriceUnit string) ([]Candle, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}
	if count > 200 {
		count = 200
	}

	params := map[string]string{
		"market": market,
	}
	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		params["count"] = fmt.Sprintf("%d", count)
	}
	if convertingPriceUnit != "" {
		params["converting_price_unit"] = convertingPriceUnit
	}

	resp, err := q.Client.Get("/candles/days", params)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

// GetCandlesWeek는 주(Week) 캔들을 조회합니다.
// market은 마켓 코드, to는 마지막 캔들 시각, count는 조회할 캔들 개수입니다.
func (q *Quotation) GetCandlesWeek(market string, to string, count int) ([]Candle, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}
	if count > 200 {
		count = 200
	}

	params := map[string]string{
		"market": market,
	}
	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		params["count"] = fmt.Sprintf("%d", count)
	}

	resp, err := q.Client.Get("/candles/weeks", params)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

// GetCandlesMonth는 월(Month) 캔들을 조회합니다.
// market은 마켓 코드, to는 마지막 캔들 시각, count는 조회할 캔들 개수입니다.
func (q *Quotation) GetCandlesMonth(market string, to string, count int) ([]Candle, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}
	if count > 200 {
		count = 200
	}

	params := map[string]string{
		"market": market,
	}
	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		params["count"] = fmt.Sprintf("%d", count)
	}

	resp, err := q.Client.Get("/candles/months", params)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

// GetCandlesYear는 년(Year) 캔들을 조회합니다.
// market은 마켓 코드, to는 마지막 캔들 시각, count는 조회할 캔들 개수입니다.
func (q *Quotation) GetCandlesYear(market string, to string, count int) ([]Candle, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}
	if count > 200 {
		count = 200
	}

	params := map[string]string{
		"market": market,
	}
	if to != "" {
		params["to"] = to
	}
	if count > 0 {
		params["count"] = fmt.Sprintf("%d", count)
	}

	resp, err := q.Client.Get("/candles/years", params)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	if err := json.Unmarshal(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}

// isValidMinuteUnit은 유효한 분봉 단위인지 확인합니다.
// 지원하는 분봉 단위인 경우 true를 반환합니다.
func isValidMinuteUnit(unit int) bool {
	validUnits := map[int]bool{
		CandleMinute1:   true,
		CandleMinute3:   true,
		CandleMinute5:   true,
		CandleMinute10:  true,
		CandleMinute15:  true,
		CandleMinute30:  true,
		CandleMinute60:  true,
		CandleMinute240: true,
	}
	return validUnits[unit]
}
