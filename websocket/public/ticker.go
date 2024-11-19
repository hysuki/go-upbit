package public

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/go-upbit/websocket/common"
)

// TickerResponse는 현재가 정보를 나타냅니다.
type TickerResponse struct {
	Type               string               `json:"type"`                  // 타입
	Code               string               `json:"code"`                  // 마켓 코드
	OpeningPrice       float64              `json:"opening_price"`         // 시가
	HighPrice          float64              `json:"high_price"`            // 고가
	LowPrice           float64              `json:"low_price"`             // 저가
	TradePrice         float64              `json:"trade_price"`           // 현재가
	PrevClosingPrice   float64              `json:"prev_closing_price"`    // 전일 종가
	Change             common.ChangeType    `json:"change"`                // 전일 대비 (RISE: 상승, EVEN: 보합, FALL: 하락)
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
	StreamType         common.StreamType    `json:"stream_type"`           // 스트림 타입
}

// ParseTicker는 JSON 데이터를 TickerResponse 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseTicker(data []byte) (*TickerResponse, error) {
	var ticker TickerResponse
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil, fmt.Errorf("티커 데이터 파싱 실패: %v", err)
	}
	return &ticker, nil
}

// GetTicker는 다음 현재가 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 현재가 정보를 반환합니다.
func (c *Client) GetTicker() (*TickerResponse, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.tickerChan:
		return resp, nil
	}
}
