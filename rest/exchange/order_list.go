package exchange

import (
	"encoding/json"
)

// OrderList는 주문 리스트 조회를 위한 파라미터를 정의하는 구조체입니다.
type OrderListParams struct {
	// Market은 마켓 아이디입니다.
	Market string `json:"market,omitempty"`
	// UUIDs는 주문 UUID의 목록입니다.
	UUIDs []string `json:"uuids,omitempty"`
	// Identifiers는 주문 identifier의 목록입니다.
	Identifiers []string `json:"identifiers,omitempty"`
	// State는 주문 상태입니다.
	State string `json:"state,omitempty"`
	// States는 주문 상태의 목록입니다.
	States []string `json:"states,omitempty"`
	// Page는 페이지 수입니다. (기본값: 1)
	Page int `json:"page,omitempty"`
	// Limit는 요청 개수입니다. (기본값: 100)
	Limit int `json:"limit,omitempty"`
	// OrderBy는 정렬 방식입니다. (asc, desc)
	OrderBy string `json:"order_by,omitempty"`
}

// GetOrders는 주문 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders
// Deprecated: 이 API는 2024년 6월부터 더 이상 사용되지 않습니다.
// 대신 GetOrdersByUUIDs, GetOpenOrders, GetClosedOrders를 사용하세요.
func (e *Exchange) GetOrders(params *OrderListParams) ([]Order, error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Market != "" {
			queryParams["market"] = params.Market
		}
		if len(params.UUIDs) > 0 {
			uuidsBytes, err := json.Marshal(params.UUIDs)
			if err != nil {
				return nil, err
			}
			queryParams["uuids"] = string(uuidsBytes)
		}
		if len(params.Identifiers) > 0 {
			identifiersBytes, err := json.Marshal(params.Identifiers)
			if err != nil {
				return nil, err
			}
			queryParams["identifiers"] = string(identifiersBytes)
		}
		if params.State != "" {
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
			queryParams["limit"] = string(params.Limit)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}
	}

	resp, err := e.Client.Get("/orders", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
