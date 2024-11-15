package exchange

import (
	"encoding/json"
	"errors"
)

// GetDepositParams는 개별 입금 조회를 위한 파라미터입니다.
type GetDepositParams struct {
	// UUID는 입금 UUID입니다.
	UUID string `json:"uuid,omitempty"`
	// TxID는 입금 TXID입니다.
	TxID string `json:"txid,omitempty"`
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
}

// GetDeposit는 개별 입금을 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/deposit
func (e *Exchange) GetDeposit(params *GetDepositParams) (*DepositInfo, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	// UUID, TxID, Currency 중 하나는 반드시 포함되어야 합니다.
	if params.UUID == "" && params.TxID == "" && params.Currency == "" {
		return nil, errors.New("either uuid, txid, or currency must be provided")
	}

	queryParams := make(map[string]string)
	if params.UUID != "" {
		queryParams["uuid"] = params.UUID
	}
	if params.TxID != "" {
		queryParams["txid"] = params.TxID
	}
	if params.Currency != "" {
		queryParams["currency"] = params.Currency
	}

	resp, err := e.Client.Get("/deposit", queryParams)
	if err != nil {
		return nil, err
	}

	var deposit DepositInfo
	if err := json.Unmarshal(resp, &deposit); err != nil {
		return nil, err
	}

	return &deposit, nil
}
