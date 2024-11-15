// Package exchange는 Upbit REST API의 거래소 관련 기능을 제공합니다.
package exchange

import (
	"encoding/json"
)

// Accounts는 계좌 정보를 나타내는 구조체입니다.
// 참고: https://docs.upbit.com/reference/%EC%A0%84%EC%B2%B4-%EA%B3%84%EC%A2%8C-%EC%A1%B0%ED%9A%8C
type Accounts struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// Balance는 주문가능 금액/수량입니다.
	Balance string `json:"balance"`
	// Locked는 주문 중 묶여있는 금액/수량입니다.
	Locked string `json:"locked"`
	// AvgBuyPrice는 매수평균가입니다.
	AvgBuyPrice string `json:"avg_buy_price"`
	// AvgBuyPriceModified는 매수평균가 수정 여부입니다.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified"`
	// UnitCurrency는 평단가 기준 화폐입니다.
	UnitCurrency string `json:"unit_currency"`
}

// GetAccounts는 사용자의 전체 계좌 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/accounts
func (e *Exchange) GetAccounts() ([]Accounts, error) {
	resp, err := e.Client.Get("/accounts", nil)
	if err != nil {
		return nil, err
	}

	var accounts []Accounts
	if err := json.Unmarshal(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}
