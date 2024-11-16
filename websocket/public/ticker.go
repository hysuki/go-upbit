package public

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/go-upbit/websocket/common"
)

// Ticker는 현재가 정보를 담는 구조체입니다
type Ticker struct {
	Type               string               `json:"type"`                  // 타입
	Code               string               `json:"code"`                  // 마켓 코드
	OpeningPrice       float64              `json:"opening_price"`         // 시가
	HighPrice          float64              `json:"high_price"`            // 고가
	LowPrice           float64              `json:"low_price"`             // 저가
	TradePrice         float64              `json:"trade_price"`           // 현재가
	PrevClosingPrice   float64              `json:"prev_closing_price"`    // 전일 종가
	Change             common.ChangeType    `json:"change"`                // RISE, EVEN, FALL
	ChangePrice        float64              `json:"change_price"`          // 변화액의 절대값
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
	StreamType         common.StreamType    `json:"stream_type"`           // SNAPSHOT, REALTIME
}

// ParseTicker는 JSON 데이터를 Ticker 구조체로 파싱합니다
func ParseTicker(data []byte) (*Ticker, error) {
	var ticker Ticker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil, fmt.Errorf("티커 데이터 파싱 실패: %v", err)
	}
	return &ticker, nil
}

// SubscribeTicker는 지정된 마켓 코드들의 현재가 정보를 구독합니다
func (c *Client) SubscribeTicker(codes []string) error {
	return c.Subscribe("", "ticker", codes, nil)
}

// GetTicker는 수신된 메시지를 Ticker 구조체로 변환합니다
func (c *Client) GetTicker() (*Ticker, error) {
	data, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // 서버 상태 메시지인 경우
	}

	return ParseTicker(data)
}
