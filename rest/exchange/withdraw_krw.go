package exchange

import (
	"encoding/json"
	"errors"
)

// TwoFactorType은 2차 인증 수단을 정의합니다.
const (
	// TwoFactorTypeKakao는 카카오 인증입니다.
	TwoFactorTypeKakao = "kakao"
	// TwoFactorTypeNaver는 네이버 인증입니다.
	TwoFactorTypeNaver = "naver"
	// TwoFactorTypeHana는 하나인증서 인증입니다.
	TwoFactorTypeHana = "hana"
)

// WithdrawKRWParams는 원화 출금을 위한 파라미터입니다.
type WithdrawKRWParams struct {
	// Amount는 출금액입니다. (필수)
	Amount string `json:"amount"`
	// TwoFactorType은 2차 인증 수단입니다. (필수)
	// - kakao: 카카오 인증
	// - naver: 네이버 인증
	// - hana: 하나인증서 인증
	TwoFactorType string `json:"two_factor_type"`
}

// WithdrawKRWResponse는 원화 출금 요청에 대한 응답입니다.
type WithdrawKRWResponse struct {
	// Type은 입출금 종류입니다.
	Type string `json:"type"`
	// UUID는 출금의 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// TxID는 출금의 트랜잭션 아이디입니다.
	TxID string `json:"txid"`
	// State는 출금 상태입니다.
	State string `json:"state"`
	// CreatedAt은 출금 생성 시간입니다.
	CreatedAt string `json:"created_at"`
	// DoneAt은 출금 완료 시간입니다.
	DoneAt string `json:"done_at"`
	// Amount는 출금 금액/수량입니다.
	Amount string `json:"amount"`
	// Fee는 출금 수수료입니다.
	Fee string `json:"fee"`
	// TransactionType은 출금 유형입니다.
	TransactionType string `json:"transaction_type"`
}

// WithdrawKRW는 원화 출금을 요청합니다.
// 엔드포인트: https://api.upbit.com/v1/withdraws/krw
func (e *Exchange) WithdrawKRW(params *WithdrawKRWParams) (*WithdrawKRWResponse, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	// 필수 파라미터 검증
	if params.Amount == "" {
		return nil, errors.New("amount is required")
	}
	if params.TwoFactorType == "" {
		return nil, errors.New("two_factor_type is required")
	}

	// 2차 인증 수단 검증
	switch params.TwoFactorType {
	case TwoFactorTypeKakao, TwoFactorTypeNaver, TwoFactorTypeHana:
		// 유효한 2차 인증 수단
	default:
		return nil, errors.New("invalid two_factor_type")
	}

	resp, err := e.Client.Post("/withdraws/krw", params)
	if err != nil {
		return nil, err
	}

	var withdrawResp WithdrawKRWResponse
	if err := json.Unmarshal(resp, &withdrawResp); err != nil {
		return nil, err
	}

	return &withdrawResp, nil
}
