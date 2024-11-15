package exchange

import (
	"encoding/json"
	"errors"
)

// DepositCoinChance는 디지털 자산 입금 정보를 나타내는 구조체입니다.
type DepositCoinChance struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 입금 네트워크입니다.
	NetType string `json:"net_type"`
	// IsDepositPossible는 입금 가능여부입니다.
	IsDepositPossible bool `json:"is_deposit_possible"`
	// DepositImpossibleReason은 입금 불가사유입니다.
	DepositImpossibleReason string `json:"deposit_impossible_reason"`
	// MinimumDepositAmount는 최소 입금 수량입니다.
	MinimumDepositAmount string `json:"minimum_deposit_amount"`
	// MinimumDepositConfirmations는 최소 입금 컨펌 수입니다.
	MinimumDepositConfirmations int `json:"minimum_deposit_confirmations"`
	// DecimalPrecision는 입금 소수점 자릿수입니다.
	DecimalPrecision int `json:"decimal_precision"`
}

// GetDepositCoinChance는 디지털 자산 입금 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits/chance/coin
func (e *Exchange) GetDepositCoinChance(currency string) (*DepositCoinChance, error) {
	if currency == "" {
		return nil, errors.New("currency is required")
	}

	params := map[string]string{
		"currency": currency,
	}

	resp, err := e.Client.Get("/deposits/chance/coin", params)
	if err != nil {
		return nil, err
	}

	var chance DepositCoinChance
	if err := json.Unmarshal(resp, &chance); err != nil {
		return nil, err
	}

	return &chance, nil
}
