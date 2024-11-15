package exchange

import (
	"encoding/json"
	"errors"
)

// DepositKRWParams는 원화 입금을 위한 파라미터입니다.
type DepositKRWParams struct {
	// Amount는 입금액입니다.
	Amount string `json:"amount"`
	// TwoFactorType은 2차 인증 수단입니다.
	// - kakao: 카카오 인증
	// - naver: 네이버 인증
	// - hana: 하나인증서 인증
	TwoFactorType string `json:"two_factor_type"`
}

// DepositKRW는 원화 입금을 요청합니다.
// 엔드포인트: https://api.upbit.com/v1/deposits/krw
func (e *Exchange) DepositKRW(params *DepositKRWParams) (*DepositInfo, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	if params.Amount == "" {
		return nil, errors.New("amount is required")
	}
	if params.TwoFactorType == "" {
		return nil, errors.New("two_factor_type is required")
	}

	resp, err := e.Client.Post("/deposits/krw", params)
	if err != nil {
		return nil, err
	}

	var deposit DepositInfo
	if err := json.Unmarshal(resp, &deposit); err != nil {
		return nil, err
	}

	return &deposit, nil
}
