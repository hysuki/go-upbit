package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/websocket/common"
)

// UpbitTrade는 체결 정보를 나타냅니다.
type UpbitTrade struct {
	Type             string            `json:"type"`               // 타입
	Code             string            `json:"code"`               // 마켓 코드
	TradePrice       float64           `json:"trade_price"`        // 체결 가격
	TradeVolume      float64           `json:"trade_volume"`       // 체결량
	AskBid           common.AskBidType `json:"ask_bid"`            // 매수/매도 구분
	PrevClosingPrice float64           `json:"prev_closing_price"` // 전일 종가
	Change           string            `json:"change"`             // 전일 대비
	ChangePrice      float64           `json:"change_price"`       // 부호 없는 전일 대비 값
	TradeDate        string            `json:"trade_date"`         // 체결 일자(UTC)
	TradeTime        string            `json:"trade_time"`         // 체결 시각(UTC)
	TradeTimestamp   int64             `json:"trade_timestamp"`    // 체결 타임스탬프
	Timestamp        int64             `json:"timestamp"`          // 타임스탬프
	SequentialId     int64             `json:"sequential_id"`      // 체결 번호
	StreamType       common.StreamType `json:"stream_type"`        // 스트림 타입
}

// Trade는 내부적으로 사용하기 위한 체결 정보 구조체입니다.
type Trade struct {
	Type             string            `json:"type"`               // 타입 (trade)
	Code             string            `json:"code"`               // 마켓 코드
	TradePrice       float64           `json:"trade_price"`        // 체결 가격
	TradeVolume      float64           `json:"trade_volume"`       // 체결량
	AskBid           common.AskBidType `json:"ask_bid"`            // 매수/매도 구분
	PrevClosingPrice float64           `json:"prev_closing_price"` // 전일 종가
	Change           string            `json:"change"`             // 전일 대비
	ChangePrice      float64           `json:"change_price"`       // 부호 없는 전일 대비 값
	TradeDate        string            `json:"trade_date"`         // 체결 일자(UTC)
	TradeTime        string            `json:"trade_time"`         // 체결 시각(UTC)
	TradeTimestamp   time.Time         `json:"trade_timestamp"`    // 체결 시각 (trade_timestamp를 KST로 변환)
	Timestamp        time.Time         `json:"timestamp"`          // 타임스탬프 (KST)
	SequentialId     int64             `json:"sequential_id"`      // 체결 번호
	StreamType       common.StreamType `json:"stream_type"`        // 스트림 타입
}

// NewTrade는 UpbitTrade를 내부 Trade 구조체로 변환합니다.
func NewTrade(u *UpbitTrade, loc *time.Location) *Trade {
	// loc이 nil인 경우 UTC를 사용
	if loc == nil {
		loc = time.UTC
	}

	return &Trade{
		Type:       u.Type,
		Code:       u.Code,
		TradePrice: u.TradePrice,

		TradeVolume:      u.TradeVolume,
		AskBid:           u.AskBid,
		PrevClosingPrice: u.PrevClosingPrice,
		Change:           u.Change,
		ChangePrice:      u.ChangePrice,
		TradeDate:        u.TradeDate,
		TradeTime:        u.TradeTime,
		TradeTimestamp:   time.UnixMilli(u.TradeTimestamp).In(loc),
		Timestamp:        time.UnixMilli(u.Timestamp).In(loc),
		SequentialId:     u.SequentialId,
		StreamType:       u.StreamType,
	}
}

// ParseTrade는 JSON 데이터를 UpbitTrade 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseTrade(data []byte) (*UpbitTrade, error) {
	var trade UpbitTrade
	if err := json.Unmarshal(data, &trade); err != nil {
		return nil, fmt.Errorf("체결 데이터 파싱 실패: %v", err)
	}
	return &trade, nil
}

// GetTrade는 다음 체결 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 체결 정보를 반환합니다.
func (c *Client) GetTrade(loc *time.Location) (*Trade, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.tradeChan:
		return NewTrade(resp, loc), nil
	}
}
