package private

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/websocket/common"
)

// UpbitMyOrder는 내 주문 정보를 나타냅니다.
type UpbitMyOrder struct {
	Type            string             `json:"type"`             // 타입 (myOrder)
	Code            string             `json:"code"`             // 마켓 코드
	UUID            string             `json:"uuid"`             // 주문 고유 ID
	AskBid          common.AskBidType  `json:"ask_bid"`          // 매수/매도 구분
	OrderType       common.OrderType   `json:"order_type"`       // 주문 타입
	State           common.OrderState  `json:"state"`            // 주문 상태
	TradeUUID       string             `json:"trade_uuid"`       // 체결의 고유 ID
	Price           float64            `json:"price"`            // 주문 가격
	AvgPrice        float64            `json:"avg_price"`        // 평균 체결 가격
	Volume          float64            `json:"volume"`           // 주문량
	RemainingVolume float64            `json:"remaining_volume"` // 체결 후 남은 주문량
	ExecutedVolume  float64            `json:"executed_volume"`  // 체결된 양
	TradesCount     int                `json:"trades_count"`     // 해당 주문에 걸린 체결 수
	ReservedFee     float64            `json:"reserved_fee"`     // 수수료로 예약된 비용
	RemainingFee    float64            `json:"remaining_fee"`    // 남은 수수료
	PaidFee         float64            `json:"paid_fee"`         // 사용된 수수료
	Locked          float64            `json:"locked"`           // 거래에 사용중인 비용
	ExecutedFunds   float64            `json:"executed_funds"`   // 체결된 금액
	TimeInForce     common.TimeInForce `json:"time_in_force"`    // IOC, FOK 설정
	TradeTimestamp  int64              `json:"trade_timestamp"`  // 체결 타임스탬프
	OrderTimestamp  int64              `json:"order_timestamp"`  // 주문 타임스탬프
	Timestamp       int64              `json:"timestamp"`        // 타임스탬프
	StreamType      common.StreamType  `json:"stream_type"`      // 스트림 타입
}

// MyOrder는 내부적으로 사용하기 위한 주문 정보 구조체입니다.
type MyOrder struct {
	Type            string             `json:"type"`             // 타입 (myOrder)
	Code            string             `json:"code"`             // 마켓 코드
	UUID            string             `json:"uuid"`             // 주문 고유 ID
	AskBid          common.AskBidType  `json:"ask_bid"`          // 매수/매도 구분
	OrderType       common.OrderType   `json:"order_type"`       // 주문 타입
	State           common.OrderState  `json:"state"`            // 주문 상태
	TradeUUID       string             `json:"trade_uuid"`       // 체결의 고유 ID
	Price           float64            `json:"price"`            // 주문 가격
	AvgPrice        float64            `json:"avg_price"`        // 평균 체결 가격
	Volume          float64            `json:"volume"`           // 주문량
	RemainingVolume float64            `json:"remaining_volume"` // 체결 후 남은 주문량
	ExecutedVolume  float64            `json:"executed_volume"`  // 체결된 양
	TradesCount     int                `json:"trades_count"`     // 해당 주문에 걸린 체결 수
	ReservedFee     float64            `json:"reserved_fee"`     // 수수료로 예약된 비용
	RemainingFee    float64            `json:"remaining_fee"`    // 남은 수수료
	PaidFee         float64            `json:"paid_fee"`         // 사용된 수수료
	Locked          float64            `json:"locked"`           // 거래에 사용중인 비용
	ExecutedFunds   float64            `json:"executed_funds"`   // 체결된 금액
	TimeInForce     common.TimeInForce `json:"time_in_force"`    // IOC, FOK 설정
	TradeAt         time.Time          // 체결 시각 (trade_timestamp를 KST로 변환)
	OrderAt         time.Time          // 주문 시각 (order_timestamp를 KST로 변환)
	Timestamp       time.Time          // 타임스탬프 (KST)
	StreamType      common.StreamType  `json:"stream_type"` // 스트림 타입
}

// NewMyOrder는 UpbitMyOrder를 내부 MyOrder 구조체로 변환합니다.
// UTC 시간을 KST(UTC+9)로 변환하여 저장합니다.
func NewMyOrder(u *UpbitMyOrder) *MyOrder {
	kst := time.FixedZone("KST", 9*60*60) // UTC+9
	return &MyOrder{
		Type:            u.Type,
		Code:            u.Code,
		UUID:            u.UUID,
		AskBid:          u.AskBid,
		OrderType:       u.OrderType,
		State:           u.State,
		TradeUUID:       u.TradeUUID,
		Price:           u.Price,
		AvgPrice:        u.AvgPrice,
		Volume:          u.Volume,
		RemainingVolume: u.RemainingVolume,
		ExecutedVolume:  u.ExecutedVolume,
		TradesCount:     u.TradesCount,
		ReservedFee:     u.ReservedFee,
		RemainingFee:    u.RemainingFee,
		PaidFee:         u.PaidFee,
		Locked:          u.Locked,
		ExecutedFunds:   u.ExecutedFunds,
		TimeInForce:     u.TimeInForce,
		TradeAt:         time.UnixMilli(u.TradeTimestamp).In(kst),
		OrderAt:         time.UnixMilli(u.OrderTimestamp).In(kst),
		Timestamp:       time.UnixMilli(u.Timestamp).In(kst),
		StreamType:      u.StreamType,
	}
}

// ParseMyOrder는 JSON 데이터를 UpbitMyOrder 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseMyOrder(data []byte) (*UpbitMyOrder, error) {
	var myOrder UpbitMyOrder
	if err := json.Unmarshal(data, &myOrder); err != nil {
		return nil, fmt.Errorf("내 주문 데이터 파싱 실패: %v", err)
	}
	return &myOrder, nil
}

// GetMyOrder는 다음 주문 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 주문 정보를 반환합니다.
func (c *Client) GetMyOrder() (*MyOrder, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.myOrderChan:
		return NewMyOrder(resp), nil
	}
}
