package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// WithdrawState는 출금 상태를 정의합니다.
const (
	// WithdrawStateWaiting은 대기중 상태입니다.
	WithdrawStateWaiting = "WAITING"
	// WithdrawStateProcessing은 진행중 상태입니다.
	WithdrawStateProcessing = "PROCESSING"
	// WithdrawStateDone은 완료 상태입니다.
	WithdrawStateDone = "DONE"
	// WithdrawStateFailed는 실패 상태입니다.
	WithdrawStateFailed = "FAILED"
	// WithdrawStateCancelled는 취소됨 상태입니다.
	WithdrawStateCancelled = "CANCELLED"
	// WithdrawStateRejected는 거절됨 상태입니다.
	WithdrawStateRejected = "REJECTED"
)

// WithdrawTransactionType은 출금 유형을 정의합니다.
const (
	// WithdrawTransactionTypeDefault는 일반출금입니다.
	WithdrawTransactionTypeDefault = "default"
	// WithdrawTransactionTypeInternal은 바로출금입니다.
	WithdrawTransactionTypeInternal = "internal"
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

// WithdrawInfo는 출금 정보를 나타내는 구조체입니다.
type WithdrawInfo struct {
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
	CreatedAt time.Time `json:"created_at"`
	// DoneAt은 출금 완료 시간입니다.
	DoneAt time.Time `json:"done_at"`
	// Amount는 출금 금액/수량입니다.
	Amount string `json:"amount"`
	// Fee는 출금 수수료입니다.
	Fee string `json:"fee"`
	// TransactionType은 출금 유형입니다.
	TransactionType string `json:"transaction_type"`
}

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

// MemberLevel은 사용자의 보안등급 정보를 나타내는 구조체입니다.
type MemberLevel struct {
	// SecurityLevel은 사용자의 보안등급입니다.
	SecurityLevel int `json:"security_level"`
	// FeeLevel은 사용자의 수수료등급입니다.
	FeeLevel int `json:"fee_level"`
	// EmailVerified는 사용자의 이메일 인증 여부입니다.
	EmailVerified bool `json:"email_verified"`
	// IdentityAuthVerified는 사용자의 실명 인증 여부입니다.
	IdentityAuthVerified bool `json:"identity_auth_verified"`
	// BankAccountVerified는 사용자의 계좌 인증 여부입니다.
	BankAccountVerified bool `json:"bank_account_verified"`
	// TwoFactorAuthVerified는 2FA 인증 수단의 활성화 여부입니다.
	TwoFactorAuthVerified bool `json:"two_factor_auth_verified"`
	// Locked는 사용자의 계정 보호 상태입니다.
	Locked bool `json:"locked"`
	// WalletLocked는 사용자의 출금 보호 상태입니다.
	WalletLocked bool `json:"wallet_locked"`
}

// WithdrawCurrency는 화폐 정보를 나타내는 구조체입니다.
type WithdrawCurrency struct {
	// Code는 화폐를 의미하는 영문 대문자 코드입니다.
	Code string `json:"code"`
	// WithdrawFee는 해당 화폐의 출금 수수료입니다.
	WithdrawFee string `json:"withdraw_fee"`
	// IsCoin은 화폐의 디지털 자산 여부입니다.
	IsCoin bool `json:"is_coin"`
	// WalletState는 해당 화폐의 지갑 상태입니다.
	WalletState string `json:"wallet_state"`
	// WalletSupport는 해당 화폐가 지원하는 입출금 정보입니다.
	WalletSupport []string `json:"wallet_support"`
}

// WithdrawLimit는 출금 제약 정보를 나타내는 구조체입니다.
type WithdrawLimit struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// Minimum은 출금 최소 금액/수량입니다.
	Minimum string `json:"minimum"`
	// RemainingDailyFiat은 통합 1일 잔여 출금 한도입니다.
	RemainingDailyFiat string `json:"remaining_daily_fiat"`
	// FiatCurrency는 해당 자산을 구매할 수 있는 법정 화폐입니다.
	FiatCurrency string `json:"fiat_currency"`
	// WithdrawDelayedFiat은 24시간 출금지연제에 의해 출금 제한되어 있는 금액입니다.
	WithdrawDelayedFiat string `json:"withdraw_delayed_fiat"`
	// Fixed는 출금 금액/수량 소수점 자리 수입니다.
	Fixed int `json:"fixed"`
	// CanWithdraw는 출금 지원 여부입니다.
	CanWithdraw bool `json:"can_withdraw"`
}

// WithdrawChance는 출금 가능 정보를 나타내는 구조체입니다.
type WithdrawChance struct {
	// MemberLevel은 사용자의 보안등급 정보입니다.
	MemberLevel MemberLevel `json:"member_level"`
	// Currency는 화폐 정보입니다.
	Currency WithdrawCurrency `json:"currency"`
	// Account는 사용자의 계좌 정보입니다.
	Account Accounts `json:"account"`
	// WithdrawLimit는 출금 제약 정보입니다.
	WithdrawLimit WithdrawLimit `json:"withdraw_limit"`
}

// WithdrawKRWParams는 원화 출금을 위한 파라미터입니다.
type WithdrawKRWParams struct {
	// Amount는 출금액입니다. (필수)
	Amount string `json:"amount"`
	// TwoFactorType은 2차 인증 수단입니다. (필수)
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

// GetWithdrawParams는 개별 출금 조회를 위한 파라미터입니다.
type GetWithdrawParams struct {
	// UUID는 출금 UUID입니다.
	UUID string `json:"uuid,omitempty"`
	// TxID는 출금 TXID입니다.
	TxID string `json:"txid,omitempty"`
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
}

// WithdrawListParams는 출금 리스트 조회를 위한 파라미터입니다.
type WithdrawListParams struct {
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
	// State는 출금 상태입니다.
	State string `json:"state,omitempty"`
	// UUIDs는 출금 UUID의 목록입니다.
	UUIDs []string `json:"uuids,omitempty"`
	// TxIDs는 출금 TXID의 목록입니다.
	TxIDs []string `json:"txids,omitempty"`
	// Limit는 개수 제한입니다. (default: 100, max: 100)
	Limit int `json:"limit,omitempty"`
	// Page는 페이지 수입니다. (default: 1)
	Page int `json:"page,omitempty"`
	// OrderBy는 정렬 방식입니다.
	OrderBy string `json:"order_by,omitempty"`
}

// GetWithdrawAddresses는 등록된 출금 허용 주소 목록을 조회합니다.
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

// GetWithdrawChance는 해당 통화의 가능한 출금 정보를 조회합니다.
func (e *Exchange) GetWithdrawChance(currency string) (*WithdrawChance, error) {
	if currency == "" {
		return nil, errors.New("currency is required")
	}

	params := map[string]string{
		"currency": currency,
	}

	resp, err := e.Client.Get("/withdraws/chance", params)
	if err != nil {
		return nil, err
	}

	var chance WithdrawChance
	if err := json.Unmarshal(resp, &chance); err != nil {
		return nil, err
	}

	return &chance, nil
}

// WithdrawKRW는 원화 출금을 요청합니다.
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

// WithdrawCoin은 디지털 자산 출금을 요청합니다.
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

// GetWithdraw는 출금 UUID를 통해 개별 출금 정보를 조회합니다.
func (e *Exchange) GetWithdraw(params *GetWithdrawParams) (*WithdrawInfo, error) {
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

	resp, err := e.Client.Get("/withdraw", queryParams)
	if err != nil {
		return nil, err
	}

	var withdraw WithdrawInfo
	if err := json.Unmarshal(resp, &withdraw); err != nil {
		return nil, err
	}

	return &withdraw, nil
}

// GetWithdraws는 출금 리스트를 조회합니다.
func (e *Exchange) GetWithdraws(params *WithdrawListParams) ([]WithdrawInfo, error) {
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
			queryParams["limit"] = strconv.Itoa(params.Limit)
		}
		if params.Page > 0 {
			queryParams["page"] = strconv.Itoa(params.Page)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}
	}

	resp, err := e.Client.Get("/withdraws", queryParams)
	if err != nil {
		return nil, err
	}

	var withdraws []WithdrawInfo
	if err := json.Unmarshal(resp, &withdraws); err != nil {
		return nil, err
	}

	return withdraws, nil
}
