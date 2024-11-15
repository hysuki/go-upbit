package exchange

import (
	"encoding/json"
)

// DepositListParams는 입금 리스트 조회를 위한 파라미터입니다.
type DepositListParams struct {
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
	// State는 입금 상태입니다.
	State string `json:"state,omitempty"`
	// UUIDs는 입금 UUID의 목록입니다.
	UUIDs []string `json:"uuids,omitempty"`
	// TxIDs는 입금 TXID의 목록입니다.
	TxIDs []string `json:"txids,omitempty"`
	// Limit는 개수 제한입니다. (default: 100, max: 100)
	Limit int `json:"limit,omitempty"`
	// Page는 페이지 수입니다. (default: 1)
	Page int `json:"page,omitempty"`
	// OrderBy는 정렬 방식입니다.
	OrderBy string `json:"order_by,omitempty"`
}

// GetDeposits는 입금 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits
func (e *Exchange) GetDeposits(params *DepositListParams) ([]DepositInfo, error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Currency != "" {
			queryParams["currency"] = params.Currency
		}
		if params.State != "" {
			queryParams["state"] = params.State
		}
		if len(params.UUIDs) > 0 {
			uuidsBytes, err := json.Marshal(params.UUIDs)
			if err != nil {
				return nil, err
			}
			queryParams["uuids"] = string(uuidsBytes)
		}
		if len(params.TxIDs) > 0 {
			txidsBytes, err := json.Marshal(params.TxIDs)
			if err != nil {
				return nil, err
			}
			queryParams["txids"] = string(txidsBytes)
		}
		if params.Limit > 0 {
			if params.Limit > 100 {
				params.Limit = 100
			}
			queryParams["limit"] = string(params.Limit)
		}
		if params.Page > 0 {
			queryParams["page"] = string(params.Page)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}
	}

	resp, err := e.Client.Get("/deposits", queryParams)
	if err != nil {
		return nil, err
	}

	var deposits []DepositInfo
	if err := json.Unmarshal(resp, &deposits); err != nil {
		return nil, err
	}

	return deposits, nil
}
