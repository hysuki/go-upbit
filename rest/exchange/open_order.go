package exchange

import (
	"encoding/json"
	"errors"
)

// OpenOrderParams는 체결 대기 주문 조회를 위한 파라미터입니다.
type OpenOrderParams struct {
	// Market은 마켓 ID입니다.
	Market string `json:"market,omitempty"`
	// State는 주문 상태입니다.
	// - wait: 체결 대기 (default)
	// - watch: 예약주문 대기
	State string `json:"state,omitempty"`
	// States는 주문 상태의 목록입니다.
	// 기본값은 wait이며, state와 동시에 사용할 수 없습니다.
	States []string `json:"states,omitempty"`
	// Page는 페이지 수입니다. (기본값: 1)
	Page int `json:"page,omitempty"`
	// Limit는 요청 개수입니다. (기본값: 100, 최대: 100)
	Limit int `json:"limit,omitempty"`
	// OrderBy는 정렬 방식입니다.
	// - asc: 오름차순
	// - desc: 내림차순 (default)
	OrderBy string `json:"order_by,omitempty"`
}

// OpenOrder는 체결 대기 주문 정보를 나타내는 구조체입니다.
type OpenOrder struct {
	// UUID는 주문의 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Side는 주문 종류입니다.
	Side string `json:"side"`
	// OrdType은 주문 방식입니다.
	OrdType string `json:"ord_type"`
	// Price는 주문 당시 화폐 가격입니다.
	Price string `json:"price"`
	// State는 주문 상태입니다.
	State string `json:"state"`
	// Market은 마켓 ID입니다.
	Market string `json:"market"`
	// CreatedAt은 주문 생성 시간입니다.
	CreatedAt string `json:"created_at"`
	// Volume은 사용자가 입력한 주문 양입니다.
	Volume string `json:"volume"`
	// RemainingVolume은 체결 후 남은 주문 양입니다.
	RemainingVolume string `json:"remaining_volume"`
	// ReservedFee는 수수료로 예약된 비용입니다.
	ReservedFee string `json:"reserved_fee"`
	// RemainingFee는 남은 수수료입니다.
	RemainingFee string `json:"remaining_fee"`
	// PaidFee는 사용된 수수료입니다.
	PaidFee string `json:"paid_fee"`
	// Locked는 거래에 사용중인 비용입니다.
	Locked string `json:"locked"`
	// ExecutedVolume은 체결된 양입니다.
	ExecutedVolume string `json:"executed_volume"`
	// ExecutedFunds는 현재까지 체결된 금액입니다.
	ExecutedFunds string `json:"executed_funds,omitempty"`
	// TradesCount는 해당 주문에 걸린 체결 수입니다.
	TradesCount int `json:"trades_count"`
	// TimeInForce는 IOC, FOK 설정입니다.
	TimeInForce string `json:"time_in_force,omitempty"`
}

// GetOpenOrders는 체결 대기 주문 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders/open
func (e *Exchange) GetOpenOrders(params *OpenOrderParams) ([]OpenOrder, error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Market != "" {
			queryParams["market"] = params.Market
		}
		if params.State != "" {
			if len(params.States) > 0 {
				return nil, errors.New("state and states cannot be used together")
			}
			queryParams["state"] = params.State
		}
		if len(params.States) > 0 {
			statesBytes, err := json.Marshal(params.States)
			if err != nil {
				return nil, err
			}
			queryParams["states"] = string(statesBytes)
		}
		if params.Page > 0 {
			queryParams["page"] = string(params.Page)
		}
		if params.Limit > 0 {
			if params.Limit > 100 {
				return nil, errors.New("limit cannot exceed 100")
			}
			queryParams["limit"] = string(params.Limit)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}
	}

	resp, err := e.Client.Get("/orders/open", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []OpenOrder
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
