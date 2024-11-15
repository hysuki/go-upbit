package exchange

import (
	"encoding/json"
	"errors"
	"time"
)

// ClosedOrderParams는 종료된 주문 조회를 위한 파라미터입니다.
type ClosedOrderParams struct {
	// Market은 마켓 ID입니다.
	Market string `json:"market,omitempty"`
	// State는 주문 상태입니다.
	// - done: 전체 체결 완료
	// - cancel: 주문 취소
	State string `json:"state,omitempty"`
	// States는 주문 상태의 목록입니다.
	// 기본값은 [done, cancel]이며, state와 동시에 사용할 수 없습니다.
	States []string `json:"states,omitempty"`
	// StartTime은 조회 시작 시각입니다. (주문 생성 시각 기준)
	StartTime *time.Time `json:"start_time,omitempty"`
	// EndTime은 조회 종료 시각입니다. (주문 생성 시각 기준)
	EndTime *time.Time `json:"end_time,omitempty"`
	// Limit는 요청 개수입니다. (기본값: 100, 최대: 1000)
	Limit int `json:"limit,omitempty"`
	// OrderBy는 정렬 방식입니다.
	// - asc: 오름차순
	// - desc: 내림차순 (default)
	OrderBy string `json:"order_by,omitempty"`
}

// GetClosedOrders는 종료�� 주문 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders/closed
func (e *Exchange) GetClosedOrders(params *ClosedOrderParams) ([]Order, error) {
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
		if params.StartTime != nil {
			queryParams["start_time"] = params.StartTime.Format(time.RFC3339)
		}
		if params.EndTime != nil {
			queryParams["end_time"] = params.EndTime.Format(time.RFC3339)
		}
		if params.Limit > 0 {
			if params.Limit > 1000 {
				return nil, errors.New("limit cannot exceed 1000")
			}
			queryParams["limit"] = string(params.Limit)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}

		// 시간 범위 검증
		if params.StartTime != nil && params.EndTime != nil {
			duration := params.EndTime.Sub(*params.StartTime)
			if duration > 7*24*time.Hour {
				return nil, errors.New("time range cannot exceed 7 days")
			}
		}
	}

	resp, err := e.Client.Get("/orders/closed", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
