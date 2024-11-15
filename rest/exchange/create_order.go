package exchange

import (
	"encoding/json"
)

// CreateOrderRequest는 주문 생성 요청 파라미터를 정의합니다
type CreateOrderRequest struct {
	Market      string      `json:"market"`        // 마켓 ID (필수)
	Side        OrderSide   `json:"side"`          // 주문 종류 (필수)
	Volume      string      `json:"volume"`        // 주문량 (지정가, 시장가 매도 시 필수)
	Price       string      `json:"price"`         // 주문 가격 (지정가, 시장가 매수 시 필수)
	OrderType   OrderType   `json:"ord_type"`      // 주문 타입 (필수)
	Identifier  string      `json:"identifier"`    // 조회용 사용자 지정값 (선택)
	TimeInForce TimeInForce `json:"time_in_force"` // IOC, FOK 주문 설정
}

// OrderResponse는 주문 응답을 나타냅니다
type Order struct {
	UUID            string    `json:"uuid"`             // 주문의 고유 아이디
	Side            OrderSide `json:"side"`             // 주문 종류
	OrderType       OrderType `json:"ord_type"`         // 주문 방식
	Price           string    `json:"price"`            // 주문 당시 화폐 가격
	State           string    `json:"state"`            // 주문 상태
	Market          string    `json:"market"`           // 마켓의 유일키
	CreatedAt       string    `json:"created_at"`       // 주문 생성 시간
	Volume          string    `json:"volume"`           // 사용자가 입력한 주문 양
	RemainingVolume string    `json:"remaining_volume"` // 체결 후 남은 주문 양
	ReservedFee     string    `json:"reserved_fee"`     // 수수료로 예약된 비용
	RemainingFee    string    `json:"remaining_fee"`    // 남은 수수료
	PaidFee         string    `json:"paid_fee"`         // 사용된 수수료
	Locked          string    `json:"locked"`           // 거래에 사용중인 비용
	ExecutedVolume  string    `json:"executed_volume"`  // 체결된 양
	TradesCount     int       `json:"trades_count"`     // 해당 주문에 걸린 체결 수
	TimeInForce     string    `json:"time_in_force"`    // IOC, FOK 설정
}

// CreateOrder는 새로운 주문을 생성합니다
func (e *Exchange) CreateOrder(request CreateOrderRequest) (*Order, error) {
	resp, err := e.Client.Post("/orders", request)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
