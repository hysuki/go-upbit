package exchange

import (
	"encoding/json"
	"errors"
)

// DepositAddress는 입금 주소 정보를 나타내는 구조체입니다.
type DepositAddress struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 입금 네트워크입니다.
	NetType string `json:"net_type"`
	// DepositAddress는 입금 주소입니다.
	DepositAddress string `json:"deposit_address"`
	// SecondaryAddress는 2차 입금 주소입니다.
	SecondaryAddress string `json:"secondary_address"`
}

// GenerateCoinAddressParams는 입금 주소 생성을 위한 파라미터입니다.
type GenerateCoinAddressParams struct {
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency"`
	// NetType은 입금 네트워크입니다.
	NetType string `json:"net_type"`
}

// GenerateCoinAddressResponse는 입금 주소 생성 요청에 대한 응답입니다.
type GenerateCoinAddressResponse struct {
	// Success는 요청 성공 여부입니다.
	Success bool `json:"success"`
	// Message는 요청 결과에 대한 메세지입니다.
	Message string `json:"message"`
}

// GenerateCoinAddress는 입금 주소 생성을 요청합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits/generate_coin_address
func (e *Exchange) GenerateCoinAddress(params *GenerateCoinAddressParams) (*GenerateCoinAddressResponse, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	if params.Currency == "" {
		return nil, errors.New("currency is required")
	}
	if params.NetType == "" {
		return nil, errors.New("net_type is required")
	}

	resp, err := e.Client.Post("/deposits/generate_coin_address", params)
	if err != nil {
		return nil, err
	}

	var response GenerateCoinAddressResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCoinAddresses는 전체 입금 주소를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits/coin_addresses
func (e *Exchange) GetCoinAddresses() ([]DepositAddress, error) {
	resp, err := e.Client.Get("/deposits/coin_addresses", nil)
	if err != nil {
		return nil, err
	}

	var addresses []DepositAddress
	if err := json.Unmarshal(resp, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

// GetCoinAddress는 개별 입금 주소를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits/coin_address
func (e *Exchange) GetCoinAddress(currency string, netType string) (*DepositAddress, error) {
	if currency == "" {
		return nil, errors.New("currency is required")
	}

	params := map[string]string{
		"currency": currency,
	}
	if netType != "" {
		params["net_type"] = netType
	}

	resp, err := e.Client.Get("/deposits/coin_address", params)
	if err != nil {
		return nil, err
	}

	var address DepositAddress
	if err := json.Unmarshal(resp, &address); err != nil {
		return nil, err
	}

	return &address, nil
}
