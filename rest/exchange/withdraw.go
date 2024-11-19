package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Package exchange는 Upbit 거래소의 출금 관련 API를 제공합니다.

// 출금 상태를 정의하는 상수들입니다.
const (
	WithdrawStateWaiting    = "WAITING"    // 대기중
	WithdrawStateProcessing = "PROCESSING" // 진행중
	WithdrawStateDone       = "DONE"       // 완료
	WithdrawStateFailed     = "FAILED"     // 실패
	WithdrawStateCancelled  = "CANCELLED"  // 취소됨
	WithdrawStateRejected   = "REJECTED"   // 거절됨
)

// 출금 유형을 정의하는 상수들입니다.
const (
	WithdrawTransactionTypeDefault  = "default"  // 일반출금
	WithdrawTransactionTypeInternal = "internal" // 바로출금
)

// 2차 인증 수단을 정의하는 상수들입니다.
const (
	TwoFactorTypeKakao = "kakao" // 카카오 인증
	TwoFactorTypeNaver = "naver" // 네이버 인증
	TwoFactorTypeHana  = "hana"  // 하나인증서 인증
)

// WithdrawInfo는 출금 정보를 나타냅니다.
type WithdrawInfo struct {
	Type            string    `json:"type,omitempty"`             // 입출금 종류
	UUID            string    `json:"uuid,omitempty"`             // 출금의 고유 ID
	Currency        string    `json:"currency,omitempty"`         // 화폐를 의미하는 영문 대문자 코드
	NetType         string    `json:"net_type,omitempty"`         // 출금 네트워크
	TxID            string    `json:"txid,omitempty"`             // 출금의 트랜잭션 ID
	State           string    `json:"state,omitempty"`            // 출금 상태
	CreatedAt       time.Time `json:"created_at,omitempty"`       // 출금 생성 시각
	DoneAt          time.Time `json:"done_at,omitempty"`          // 출금 완료 시각
	Amount          string    `json:"amount,omitempty"`           // 출금 금액/수량
	Fee             string    `json:"fee,omitempty"`              // 출금 수수료
	TransactionType string    `json:"transaction_type,omitempty"` // 출금 유형
}

// WithdrawAddress는 출금 허용 주소 정보를 나타냅니다.
type WithdrawAddress struct {
	Currency         string `json:"currency"`                    // 화폐를 의미하는 영문 대문자 코드
	NetType          string `json:"net_type,omitempty"`          // 출금 네트워크 타입
	NetworkName      string `json:"network_name,omitempty"`      // 출금 네트워크 이름
	WithdrawAddress  string `json:"withdraw_address,omitempty"`  // 출금 주소
	SecondaryAddress string `json:"secondary_address,omitempty"` // 2차 출금 주소
}

// MemberLevel은 사용자의 보안등급 정보를 나타냅니다.
type MemberLevel struct {
	SecurityLevel         int  `json:"security_level,omitempty"`           // 사용자의 보안등급
	FeeLevel              int  `json:"fee_level,omitempty"`                // 사용자의 수수료등급
	EmailVerified         bool `json:"email_verified,omitempty"`           // 이메일 인증 여부
	IdentityAuthVerified  bool `json:"identity_auth_verified,omitempty"`   // 실명 인증 여부
	BankAccountVerified   bool `json:"bank_account_verified,omitempty"`    // 계좌 인증 여부
	TwoFactorAuthVerified bool `json:"two_factor_auth_verified,omitempty"` // 2FA 인증 수단 활성화 여부
	Locked                bool `json:"locked,omitempty"`                   // 계정 보호 상태
	WalletLocked          bool `json:"wallet_locked,omitempty"`            // 출금 보호 상태
}

// WithdrawCurrency는 화폐 정보를 나타냅니다.
type WithdrawCurrency struct {
	Code          string   `json:"code,omitempty"`           // 화폐를 의미하는 영문 대문자 코드
	WithdrawFee   string   `json:"withdraw_fee,omitempty"`   // 출금 수수료
	IsCoin        bool     `json:"is_coin,omitempty"`        // 디지털 자산 여부
	WalletState   string   `json:"wallet_state,omitempty"`   // 지갑 상태
	WalletSupport []string `json:"wallet_support,omitempty"` // 지원하는 입출금 정보
}

// WithdrawLimit는 출금 제약 정보를 나타냅니다.
type WithdrawLimit struct {
	Currency            string `json:"currency,omitempty"`              // 화폐를 의미하는 영문 대문자 코드
	Minimum             string `json:"minimum,omitempty"`               // 최소 출금 금액/수량
	RemainingDailyFiat  string `json:"remaining_daily_fiat,omitempty"`  // 통합 1일 잔여 출금 한도
	FiatCurrency        string `json:"fiat_currency,omitempty"`         // 구매 가능한 법정 화폐
	WithdrawDelayedFiat string `json:"withdraw_delayed_fiat,omitempty"` // 출금지연제로 인한 제한 금액
	Fixed               int    `json:"fixed,omitempty"`                 // 출금 금액/수량 소수점 자리 수
	CanWithdraw         bool   `json:"can_withdraw,omitempty"`          // 출금 지원 여부
}

// WithdrawChance는 출금 가능 정보를 나타냅니다.
type WithdrawChance struct {
	MemberLevel   MemberLevel      `json:"member_level,omitempty"`   // 사용자의 보안등급 정보
	Currency      WithdrawCurrency `json:"currency,omitempty"`       // 화폐 정보
	Account       Accounts         `json:"account,omitempty"`        // 사용자의 계좌 정보
	WithdrawLimit WithdrawLimit    `json:"withdraw_limit,omitempty"` // 출금 제약 정보
}

// WithdrawKRWParams는 원화 출금을 위한 파라미터입니다.
type WithdrawKRWParams struct {
	Amount        string `json:"amount,omitempty"`          // 출금액
	TwoFactorType string `json:"two_factor_type,omitempty"` // 2차 인증 수단
}

// WithdrawKRWResponse는 원화 출금 요청에 대한 응답입니다.
type WithdrawKRWResponse struct {
	Type            string `json:"type,omitempty"`             // 입출금 종류
	UUID            string `json:"uuid,omitempty"`             // 출금의 고유 ID
	Currency        string `json:"currency,omitempty"`         // 화폐를 의미하는 영문 대문자 코드
	TxID            string `json:"txid,omitempty"`             // 출금의 트랜잭션 ID
	State           string `json:"state,omitempty"`            // 출금 상태
	CreatedAt       string `json:"created_at,omitempty"`       // 출금 생성 시각
	DoneAt          string `json:"done_at,omitempty"`          // 출금 완료 시각
	Amount          string `json:"amount,omitempty"`           // 출금 금액/수량
	Fee             string `json:"fee,omitempty"`              // 출금 수수료
	TransactionType string `json:"transaction_type,omitempty"` // 출금 유형
}

// WithdrawCoinParams는 디지털 자산 출금을 위한 파라미터입니다.
type WithdrawCoinParams struct {
	Currency         string `json:"currency,omitempty"`          // Currency 코드
	NetType          string `json:"net_type,omitempty"`          // 출금 네트워크
	Amount           string `json:"amount,omitempty"`            // 출금 수량
	Address          string `json:"address,omitempty"`           // 출금 주소
	SecondaryAddress string `json:"secondary_address,omitempty"` // 2차 출금 주소
	TransactionType  string `json:"transaction_type,omitempty"`  // 출금 유형
}

// WithdrawCoinResponse는 디지털 자산 출금 요청에 대한 응답입니다.
type WithdrawCoinResponse struct {
	Type            string `json:"type,omitempty"`             // 입출금 종류
	UUID            string `json:"uuid,omitempty"`             // 출금의 고유 ID
	Currency        string `json:"currency,omitempty"`         // 화폐를 의미하는 영문 대문자 코드
	NetType         string `json:"net_type,omitempty"`         // 출금 네트워크
	TxID            string `json:"txid,omitempty"`             // 출금의 트랜잭션 ID
	State           string `json:"state,omitempty"`            // 출금 상태
	CreatedAt       string `json:"created_at,omitempty"`       // 출금 생성 시각
	DoneAt          string `json:"done_at,omitempty"`          // 출금 완료 시각
	Amount          string `json:"amount,omitempty"`           // 출금 금액/수량
	Fee             string `json:"fee,omitempty"`              // 출금 수수료
	KrwAmount       string `json:"krw_amount,omitempty"`       // 원화 환산 가격
	TransactionType string `json:"transaction_type,omitempty"` // 출금 유형
}

// GetWithdrawParams는 개별 출금 조회를 위한 파라미터입니다.
type GetWithdrawParams struct {
	UUID     string `json:"uuid,omitempty"`     // 출금 UUID
	TxID     string `json:"txid,omitempty"`     // 출금 TXID
	Currency string `json:"currency,omitempty"` // Currency 코드
}

// WithdrawListParams는 출금 리스트 조회를 위한 파라미터입니다.
type WithdrawListParams struct {
	Currency string   `json:"currency,omitempty"` // Currency 코드
	State    string   `json:"state,omitempty"`    // 출금 상태
	UUIDs    []string `json:"uuids,omitempty"`    // 출금 UUID 목록
	TxIDs    []string `json:"txids,omitempty"`    // 출금 TXID 목록
	Limit    int      `json:"limit,omitempty"`    // 개수 제한
	Page     int      `json:"page,omitempty"`     // 페이지 수
	OrderBy  string   `json:"order_by,omitempty"` // 정렬 방식
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

// GetWithdrawChance는 해당 통화의 출금 가능 정보를 조회합니다.
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
