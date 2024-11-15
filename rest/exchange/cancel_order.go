package exchange

import (
	"encoding/json"
	"errors"
)

// CancelOrderParams는 주문 취소를 위한 파라미터입니다.
type CancelOrderParams struct {
	// UUID는 취소할 주문의 UUID입니다.
	UUID string `json:"uuid,omitempty"`
	// Identifier는 조회용 사용자 지정값입니다.
	Identifier string `json:"identifier,omitempty"`
}

// CancelOrder는 주문을 취소합니다.
// 엔드포인트: https://api.upbit.com/v1/order
func (e *Exchange) CancelOrder(params *CancelOrderParams) (*Order, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	// UUID와 Identifier 중 하나는 반드시 포함되어야 합니다.
	if params.UUID == "" && params.Identifier == "" {
		return nil, errors.New("either uuid or identifier must be provided")
	}

	// UUID와 Identifier는 동시에 사용할 수 없습니다.
	if params.UUID != "" && params.Identifier != "" {
		return nil, errors.New("uuid and identifier cannot be used together")
	}

	queryParams := make(map[string]string)
	if params.UUID != "" {
		queryParams["uuid"] = params.UUID
	}
	if params.Identifier != "" {
		queryParams["identifier"] = params.Identifier
	}

	resp, err := e.Client.Delete("/order", queryParams)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
