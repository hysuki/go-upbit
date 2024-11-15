package upbit

import (
	"encoding/json"
	"fmt"
)

// OrderType은 주문 타입을 정의합니다
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"  // 지정가 주문
	OrderTypePrice  OrderType = "price"  // 시장가 매수 주문
	OrderTypeMarket OrderType = "market" // 시장가 매도 주문
	OrderTypeBest   OrderType = "best"   // 최유리 지정가 주문
)

// OrderState는 주문 상태를 정의합니다
type OrderState string

const (
	OrderStateWait   OrderState = "wait"   // 체결 대기
	OrderStateWatch  OrderState = "watch"  // 예약 주문 대기
	OrderStateTrade  OrderState = "trade"  // 체결 발생
	OrderStateDone   OrderState = "done"   // 전체 체결 완료
	OrderStateCancel OrderState = "cancel" // 주문 취소
)

// TimeInForce는 IOC, FOK 설정을 정의합니다
type TimeInForce string

const (
	TimeInForceIOC TimeInForce = "ioc"
	TimeInForceFOK TimeInForce = "fok"
)

// MyOrder는 내 주문 정보를 담는 구조체입니다
type MyOrder struct {
	Type            string      `json:"type"`             // 타입 (myOrder)
	Code            string      `json:"code"`             // 마켓 코드
	UUID            string      `json:"uuid"`             // 주문 고유 아이디
	AskBid          AskBidType  `json:"ask_bid"`          // 매수/매도 구분
	OrderType       OrderType   `json:"order_type"`       // 주문 타입
	State           OrderState  `json:"state"`            // 주문 상태
	TradeUUID       string      `json:"trade_uuid"`       // 체결의 고유 아이디
	Price           float64     `json:"price"`            // 주문 가격
	AvgPrice        float64     `json:"avg_price"`        // 평균 체결 가격
	Volume          float64     `json:"volume"`           // 주문량
	RemainingVolume float64     `json:"remaining_volume"` // 체결 후 남은 주문 양
	ExecutedVolume  float64     `json:"executed_volume"`  // 체결된 양
	TradesCount     int         `json:"trades_count"`     // 해당 주문에 걸린 체결 수
	ReservedFee     float64     `json:"reserved_fee"`     // 수수료로 예약된 비용
	RemainingFee    float64     `json:"remaining_fee"`    // 남은 수수료
	PaidFee         float64     `json:"paid_fee"`         // 사용된 수수료
	Locked          float64     `json:"locked"`           // 거래에 사용중인 비용
	ExecutedFunds   float64     `json:"executed_funds"`   // 체결된 금액
	TimeInForce     TimeInForce `json:"time_in_force"`    // IOC, FOK 설정
	TradeTimestamp  int64       `json:"trade_timestamp"`  // 체결 타임스탬프
	OrderTimestamp  int64       `json:"order_timestamp"`  // 주문 타임스탬프
	Timestamp       int64       `json:"timestamp"`        // 타임스탬프
	StreamType      string      `json:"stream_type"`      // 스트림 타입
}

// ParseMyOrder는 JSON 데이터를 MyOrder 구조체로 파싱합니다
func ParseMyOrder(data []byte) (*MyOrder, error) {
	var myOrder MyOrder
	if err := json.Unmarshal(data, &myOrder); err != nil {
		return nil, fmt.Errorf("내 주문 데이터 파싱 실패: %v", err)
	}
	return &myOrder, nil
}

// SubscribeMyOrder는 지정된 마켓 코드들의 내 주문 정보를 구독합니다
func (c *Client) SubscribeMyOrder(codes []string) error {
	return c.Subscribe("", "myorder", codes, nil)
}

// GetMyOrder는 수신된 메시지를 MyOrder 구조체로 변환합니다
func (c *Client) GetMyOrder() (*MyOrder, error) {
	data, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // 서버 상태 메시지인 경우
	}

	return ParseMyOrder(data)
}
