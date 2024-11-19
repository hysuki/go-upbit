package public

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/go-upbit/websocket/common"
)

// TradeResponse는 체결 정보를 담는 구조체입니다
type TradeResponse struct {
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

// ParseTrade는 JSON 데이터를 Trade 구조체로 파싱합니다
func ParseTrade(data []byte) (*TradeResponse, error) {
	var trade TradeResponse
	if err := json.Unmarshal(data, &trade); err != nil {
		return nil, fmt.Errorf("체결 데이터 파싱 실패: %v", err)
	}
	return &trade, nil
}

// SubscribeTrade는 지정된 마켓 코드들의 체결 정보를 구독합니다
// func (c *Client) SubscribeTrade(codes []string) error {
// 	if len(codes) == 0 {
// 		return fmt.Errorf("마켓 코드는 필수입니다")
// 	}
// 	return c.Subscribe("", "trade", codes, nil)
// }

// GetTrade는 수신된 메시지를 TradeResponse 구조체로 변환합니다
// func (c *Client) GetTrade(data []byte) (*TradeResponse, error) {

// 	// 타입 확인
// 	readMessage := websocket.ReadMessage{}
// 	if err := json.Unmarshal(data, &readMessage); err != nil {
// 		return nil, fmt.Errorf("타입 확인 실패: %v", err)
// 	}

// 	if readMessage.Type != string(Trade) {
// 		return nil, nil
// 	}

// 	return ParseTrade(data)
// }

// GetTrade waits for the next trade message
func (c *Client) GetTrade() (*TradeResponse, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.tradeChan:
		return resp, nil
	}
}
