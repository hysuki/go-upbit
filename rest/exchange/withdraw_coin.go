package exchange

import (
	"encoding/json"
	"errors"
)

// WithdrawCoinParams는 디지털 자산 출금을 위한 파라미터입니다.
type WithdrawCoinParams struct {
	// Currency는 Currency 코드입니다. (필수)
	Currency string `json:"currency"`
	// NetType은 출금 네트워크입니다. (필수)
	NetType string `json:"net_type"`
	// Amount는 출금 수량입니다. (필수)
	Amount string `json:"amount"`
	// Address는 출금 가능 주소에 등록된 출금 주소입니다. (필수)
	Address string `json:"address"`
	// SecondaryAddress는 2차 출금 주소입니다. (필요한 디지털 자산에 한해서)
	SecondaryAddress string `json:"secondary_address,omitempty"`
	// TransactionType은 출금 유형입니다.
	// default: 일반출금
	// internal: 바로출금
	TransactionType string `json:"transaction_type,omitempty"`
}

// WithdrawCoinResponse는 디지털 자산 출금 요청에 대한 응답입니다.
type WithdrawCoinResponse struct {
	// Type은 입출금 종류입니다.
	Type string `json:"type"`
	// UUID는 출금의 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 출금 네트워크입니다.
	NetType string `json:"net_type"`
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
	// KrwAmount는 원화 환산 가격입니다.
	KrwAmount string `json:"krw_amount"`
	// TransactionType은 출금 유형입니다.
	TransactionType string `json:"transaction_type"`
}

// WithdrawCoin은 디지털 자산 출금을 요청합니다.
// 엔드포인트: https://api.upbit.com/v1/withdraws/coin
func (e *Exchange) WithdrawCoin(params *WithdrawCoinParams) (*WithdrawCoinResponse, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	// 필수 파라미터 검증
	if params.Currency == "" {
		return nil, errors.New("currency is required")
	}
	if params.NetType == "" {
		return nil, errors.New("net_type is required")
	}
	if params.Amount == "" {
		return nil, errors.New("amount is required")
	}
	if params.Address == "" {
		return nil, errors.New("address is required")
	}

	resp, err := e.Client.Post("/withdraws/coin", params)
	if err != nil {
		return nil, err
	}

	var withdrawResp WithdrawCoinResponse
	if err := json.Unmarshal(resp, &withdrawResp); err != nil {
		return nil, err
	}

	return &withdrawResp, nil
}
