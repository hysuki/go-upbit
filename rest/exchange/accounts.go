// Package exchange는 Upbit 거래소 API와 관련된 기능을 제공합니다.
package exchange

import (
	"encoding/json"
)

// Accounts는 Upbit 거래소의 계좌 정보를 나타내는 구조체입니다.
type Account struct {
	// Currency는 화폐를 나타내는 영문 대문자 코드입니다.
	// 예: KRW, BTC, ETH
	Currency string `json:"currency,omitempty"`

	// Balance는 해당 화폐의 주문 가능한 잔고 수량입니다.
	Balance string `json:"balance,omitempty"`

	// Locked는 해당 화폐의 주문이 걸려있는 잔고 수량입니다.
	Locked string `json:"locked,omitempty"`

	// AvgBuyPrice는 해당 화폐의 매수 평균가입니다.
	AvgBuyPrice string `json:"avg_buy_price,omitempty"`

	// AvgBuyPriceModified는 매수 평균가 수정 여부를 나타냅니다.
	// 수정되었다면 true, 아니면 false입니다.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified"`

	// UnitCurrency는 평균 매수가의 기준이 되는 화폐를 나타냅니다.
	// 일반적으로 KRW입니다.
	UnitCurrency string `json:"unit_currency,omitempty"`
}

// GetAccounts는 보유한 자산 리스트를 조회합니다.
// 이 메서드는 인증이 필요한 API를 호출하며, 전체 계좌 잔고와 관련 정보를 반환합니다.
// 에러가 발생한 경우 error를 반환합니다.
func (e *Exchange) GetAccounts() ([]Account, error) {
	resp, err := e.Client.Get("/accounts", nil)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := json.Unmarshal(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}
