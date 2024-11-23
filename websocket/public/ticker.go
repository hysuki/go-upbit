package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/websocket/common"
)

// UpbitTicker는 현재가 정보를 나타냅니다.
type UpbitTicker struct {
	Type               string               `json:"type"`                  // 타입
	Code               string               `json:"code"`                  // 마켓 코드
	OpeningPrice       float64              `json:"opening_price"`         // 시가
	HighPrice          float64              `json:"high_price"`            // 고가
	LowPrice           float64              `json:"low_price"`             // 저가
	TradePrice         float64              `json:"trade_price"`           // 현재가
	PrevClosingPrice   float64              `json:"prev_closing_price"`    // 전일 종가
	Change             common.ChangeType    `json:"change"`                // 전일 대비
	ChangePrice        float64              `json:"change_price"`          // 변화액의 절대값
	SignedChangePrice  float64              `json:"signed_change_price"`   // 전일 대비 값
	ChangeRate         float64              `json:"change_rate"`           // 부호 없는 전일 대비 등락율
	SignedChangeRate   float64              `json:"signed_change_rate"`    // 전일 대비 등락율
	TradeVolume        float64              `json:"trade_volume"`          // 가장 최근 거래량
	AccTradeVolume     float64              `json:"acc_trade_volume"`      // 누적 거래량
	AccTradeVolume24h  float64              `json:"acc_trade_volume_24h"`  // 24시간 누적 거래량
	AccTradePrice      float64              `json:"acc_trade_price"`       // 누적 거래대금
	AccTradePrice24h   float64              `json:"acc_trade_price_24h"`   // 24시간 누적 거래대금
	TradeDate          string               `json:"trade_date"`            // 최근 거래 일자(UTC)
	TradeTime          string               `json:"trade_time"`            // 최근 거래 시각(UTC)
	TradeTimestamp     int64                `json:"trade_timestamp"`       // 체결 타임스탬프
	AskBid             common.AskBidType    `json:"ask_bid"`               // 매수/매도 구분
	AccAskVolume       float64              `json:"acc_ask_volume"`        // 누적 매도량
	AccBidVolume       float64              `json:"acc_bid_volume"`        // 누적 매수량
	Highest52WeekPrice float64              `json:"highest_52_week_price"` // 52주 신고가
	Highest52WeekDate  string               `json:"highest_52_week_date"`  // 52주 신고가 달성일
	Lowest52WeekPrice  float64              `json:"lowest_52_week_price"`  // 52주 신저가
	Lowest52WeekDate   string               `json:"lowest_52_week_date"`   // 52주 신저가 달성일
	MarketState        common.MarketState   `json:"market_state"`          // 거래상태
	MarketWarning      common.MarketWarning `json:"market_warning"`        // 거래경고
	Timestamp          int64                `json:"timestamp"`             // 타임스탬프
	StreamType         common.StreamType    `json:"stream_type"`           // 스트림 타입
}

// Ticker는 내부적으로 사용하기 위한 현재가 정보 구조체입니다.
type Ticker struct {
	Type               string               `json:"type"`                 // 타입
	Code               string               `json:"code"`                 // 마켓 코드
	OpeningPrice       float64              `json:"opening_price"`        // 시가
	HighPrice          float64              `json:"high_price"`           // 고가
	LowPrice           float64              `json:"low_price"`            // 저가
	TradePrice         float64              `json:"trade_price"`          // 현재가
	PrevClosingPrice   float64              `json:"prev_closing_price"`   // 전일 종가
	Change             common.ChangeType    `json:"change"`               // 전일 대비
	ChangePrice        float64              `json:"change_price"`         // 변화액의 절대값
	SignedChangePrice  float64              `json:"signed_change_price"`  // 전일 대비 값
	ChangeRate         float64              `json:"change_rate"`          // 부호 없는 전일 대비 등락율
	SignedChangeRate   float64              `json:"signed_change_rate"`   // 전일 대비 등락율
	TradeVolume        float64              `json:"trade_volume"`         // 가장 최근 거래량
	AccTradeVolume     float64              `json:"acc_trade_volume"`     // 누적 거래량
	AccTradeVolume24h  float64              `json:"acc_trade_volume_24h"` // 24시간 누적 거래량
	AccTradePrice      float64              `json:"acc_trade_price"`      // 누적 거래대금
	AccTradePrice24h   float64              `json:"acc_trade_price_24h"`  // 24시간 누적 거래대금
	TradeDate          string               `json:"trade_date"`           // 최근 거래 일자(UTC)
	TradeTime          string               `json:"trade_time"`           // 최근 거래 시각(UTC)
	TradeAt            time.Time            // 체결 시각 (trade_timestamp를 KST로 변환)
	AskBid             common.AskBidType    `json:"ask_bid"`               // 매수/매도 구분
	AccAskVolume       float64              `json:"acc_ask_volume"`        // 누적 매도량
	AccBidVolume       float64              `json:"acc_bid_volume"`        // 누적 매수량
	Highest52WeekPrice float64              `json:"highest_52_week_price"` // 52주 신고가
	Highest52WeekDate  string               `json:"highest_52_week_date"`  // 52주 신고가 달성일
	Lowest52WeekPrice  float64              `json:"lowest_52_week_price"`  // 52주 신저가
	Lowest52WeekDate   string               `json:"lowest_52_week_date"`   // 52주 신저가 달성일
	MarketState        common.MarketState   `json:"market_state"`          // 거래상태
	MarketWarning      common.MarketWarning `json:"market_warning"`        // 거래경고
	Timestamp          time.Time            // 타임스탬프 (KST)
	StreamType         common.StreamType    `json:"stream_type"` // 스트림 타입
}

// NewTicker는 UpbitTicker를 내부 Ticker 구조체로 변환합니다.
// UTC 시간을 KST(UTC+9)로 변환하여 저장합니다.
func NewTicker(u *UpbitTicker) *Ticker {
	kst := time.FixedZone("KST", 9*60*60) // UTC+9
	return &Ticker{
		Type:               u.Type,
		Code:               u.Code,
		OpeningPrice:       u.OpeningPrice,
		HighPrice:          u.HighPrice,
		LowPrice:           u.LowPrice,
		TradePrice:         u.TradePrice,
		PrevClosingPrice:   u.PrevClosingPrice,
		Change:             u.Change,
		ChangePrice:        u.ChangePrice,
		SignedChangePrice:  u.SignedChangePrice,
		ChangeRate:         u.ChangeRate,
		SignedChangeRate:   u.SignedChangeRate,
		TradeVolume:        u.TradeVolume,
		AccTradeVolume:     u.AccTradeVolume,
		AccTradeVolume24h:  u.AccTradeVolume24h,
		AccTradePrice:      u.AccTradePrice,
		AccTradePrice24h:   u.AccTradePrice24h,
		TradeDate:          u.TradeDate,
		TradeTime:          u.TradeTime,
		TradeAt:            time.UnixMilli(u.TradeTimestamp).In(kst),
		AskBid:             u.AskBid,
		AccAskVolume:       u.AccAskVolume,
		AccBidVolume:       u.AccBidVolume,
		Highest52WeekPrice: u.Highest52WeekPrice,
		Highest52WeekDate:  u.Highest52WeekDate,
		Lowest52WeekPrice:  u.Lowest52WeekPrice,
		Lowest52WeekDate:   u.Lowest52WeekDate,
		MarketState:        u.MarketState,
		MarketWarning:      u.MarketWarning,
		Timestamp:          time.UnixMilli(u.Timestamp).In(kst),
		StreamType:         u.StreamType,
	}
}

// ParseTicker는 JSON 데이터를 UpbitTicker 구조체로 파싱합니다.
// 파싱에 실패하면 에러��� 반환합니다.
func ParseTicker(data []byte) (*UpbitTicker, error) {
	var ticker UpbitTicker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil, fmt.Errorf("티커 데이터 파싱 실패: %v", err)
	}
	return &ticker, nil
}

// GetTicker는 다음 현재가 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 현재가 정보를 반환합니다.
func (c *Client) GetTicker() (*Ticker, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.tickerChan:
		return NewTicker(resp), nil
	}
}
