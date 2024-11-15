package quotation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// CandleMinute는 분봉 단위를 정의합니다.
const (
	CandleMinute1   = 1
	CandleMinute3   = 3
	CandleMinute5   = 5
	CandleMinute10  = 10
	CandleMinute15  = 15
	CandleMinute30  = 30
	CandleMinute60  = 60
	CandleMinute240 = 240
)

// Candle은 캔들 정보를 나타내는 구조체입니다.
type Candle struct {
	// Market은 마켓명입니다.
	Market string `json:"market"`
	// CandleDateTimeUTC는 캔들 기준 시각(UTC)입니다.
	CandleDateTimeUTC string `json:"candle_date_time_utc"`
	// CandleDateTimeKST는 캔들 기준 시각(KST)입니다.
	CandleDateTimeKST string `json:"candle_date_time_kst"`
	// OpeningPrice는 시가입니다.
	OpeningPrice float64 `json:"opening_price"`
	// HighPrice는 고가입니다.
	HighPrice float64 `json:"high_price"`
	// LowPrice는 저가입니다.
	LowPrice float64 `json:"low_price"`
	// TradePrice는 종가입니다.
	TradePrice float64 `json:"trade_price"`
	// Timestamp는 해당 캔들에서 마지막 틱이 저장된 시각입니다.
	Timestamp int64 `json:"timestamp"`
	// CandleAccTradePrice는 누적 거래 금액입니다.
	CandleAccTradePrice float64 `json:"candle_acc_trade_price"`
	// CandleAccTradeVolume는 누적 거래량입니다.
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	// FirstDayOfPeriod는 캔들 기간의 가장 첫 날입니다. (주/월/년봉에서만 제공)
	FirstDayOfPeriod string `json:"first_day_of_period,omitempty"`
	// Unit은 분봉 단위입니다. (분봉에서만 제공)
	Unit int `json:"unit,omitempty"`
	// ConvertedTradePrice는 종가 환산 화폐 단위로 환산된 가격입니다. (일봉에서만 제공)
	ConvertedTradePrice float64 `json:"converted_trade_price,omitempty"`
	// PrevClosingPrice는 전일 종가입니다. (일봉에서만 제공)
	PrevClosingPrice float64 `json:"prev_closing_price,omitempty"`
	// ChangePrice는 전일 종가 대비 변화 금액입니다. (일봉에서만 제공)
	ChangePrice float64 `json:"change_price,omitempty"`
	// ChangeRate는 전일 종가 대비 변화량입니다. (일봉에서만 제공)
	ChangeRate float64 `json:"change_rate,omitempty"`
}

// GetCandlesMinute는 분(Minute) 캔들을 조회합니다.
// unit: 분봉 단위 (1, 3, 5, 10, 15, 30, 60, 240)
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
