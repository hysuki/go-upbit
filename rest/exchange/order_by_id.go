package exchange

import (
	"encoding/json"
	"errors"
)

// OrderByIDParams는 UUID 또는 identifier로 주문 리스트를 조회하기 위한 파라미터입니다.
type OrderByIDParams struct {
	// Market은 마켓 ID입니다.
	Market string `json:"market,omitempty"`
	// UUIDs는 주문 UUID의 목록입니다. (최대 100개)
	UUIDs []string `json:"uuids,omitempty"`
	// Identifiers는 주문 identifier의 목록입니다. (최대 100개)
	Identifiers []string `json:"identifiers,omitempty"`
	// OrderBy는 정렬 방식입니다. (asc, desc)
	OrderBy string `json:"order_by,omitempty"`
}

// GetOrdersByID는 UUID 또는 identifier로 주문 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders/uuids
func (e *Exchange) GetOrdersByID(params *OrderByIDParams) ([]Order, error) {
	if params == nil {
		return nil, ErrInvalidParams
	}

	// UUID와 identifier는 동시에 사용할 수 없습니다.
	if len(params.UUIDs) > 0 && len(params.Identifiers) > 0 {
		return nil, ErrInvalidParams
	}

	// UUID 또는 identifier 중 하나는 필수입니다.
	if len(params.UUIDs) == 0 && len(params.Identifiers) == 0 {
		return nil, ErrInvalidParams
	}

	// 최대 100개까지만 조회 가능합니다.
	if len(params.UUIDs) > 100 || len(params.Identifiers) > 100 {
		return nil, ErrTooManyIDs
	}

	queryParams := make(map[string]string)

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

	if params.OrderBy != "" {
		queryParams["order_by"] = params.OrderBy
	}

	resp, err := e.Client.Get("/orders/uuids", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// 에러 정의
var (
	// ErrInvalidParams는 잘못된 파라미터가 전달되었을 때 발생하는 에러입니다.
	ErrInvalidParams = errors.New("invalid parameters")
	// ErrTooManyIDs는 UUID 또는 identifier가 100개를 초과할 때 발생하는 에러입니다.
	ErrTooManyIDs = errors.New("too many uuids or identifiers: maximum is 100")
)
