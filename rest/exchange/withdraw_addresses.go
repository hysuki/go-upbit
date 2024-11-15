package exchange

import (
	"encoding/json"
)

// WithdrawAddress는 출금 허용 주소 정보를 나타내는 구조체입니다.
type WithdrawAddress struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 출금 네트워크 타입입니다.
	NetType string `json:"net_type"`
	// NetworkName은 출금 네트워크 이름입니다.
	NetworkName string `json:"network_name"`
	// WithdrawAddress는 출금 주소입니다.
	WithdrawAddress string `json:"withdraw_address"`
	// SecondaryAddress는 2차 출금 주소입니다. (필요한 디지털 자산에 한해서)
	SecondaryAddress string `json:"secondary_address"`
}

// GetWithdrawAddresses는 등록된 출금 허용 주소 목록을 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/withdraws/coin_addresses
func (e *Exchange) GetWithdrawAddresses() ([]WithdrawAddress, error) {
	resp, err := e.Client.Get("/withdraws/coin_addresses", nil)
	if err != nil {
		return nil, err
	}

	var addresses []WithdrawAddress
	if err := json.Unmarshal(resp, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}
