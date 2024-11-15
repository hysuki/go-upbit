package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// DepositState는 입금 상태를 정의합니다.
const (
	// DepositStateProcessing은 입금 진행중 상태입니다.
	DepositStateProcessing = "PROCESSING"
	// DepositStateAccepted는 완료 상태입니다.
	DepositStateAccepted = "ACCEPTED"
	// DepositStateCancelled는 취소됨 상태입니다.
	DepositStateCancelled = "CANCELLED"
	// DepositStateRejected는 거절됨 상태입니다.
	DepositStateRejected = "REJECTED"
	// DepositStateTravelRuleSuspected는 트래블룰 추가 인증 대기중 상태입니다.
	DepositStateTravelRuleSuspected = "TRAVEL_RULE_SUSPECTED"
	// DepositStateRefunding는 반환절차 진행중 상태입니다.
	DepositStateRefunding = "REFUNDING"
	// DepositStateRefunded는 반환됨 상태입니다.
	DepositStateRefunded = "REFUNDED"
)

// DepositTransactionType은 입금 유형을 정의합니다.
const (
	// DepositTransactionTypeDefault는 일반입금입니다.
	DepositTransactionTypeDefault = "default"
	// DepositTransactionTypeInternal는 바로입금입니다.
	DepositTransactionTypeInternal = "internal"
)

// DepositInfo는 입금 정보를 나타내는 구조체입니다.
type DepositInfo struct {
	// Type은 입출금 종류입니다.
	Type string `json:"type"`
	// UUID는 입금에 대한 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 입금 네트워크입니다.
	NetType string `json:"net_type"`
	// TxID는 입금의 트랜잭션 아이디입니다.
	TxID string `json:"txid"`
	// State는 입금 상태입니다.
	State string `json:"state"`
	// CreatedAt은 입금 생성 시간입니다.
	CreatedAt time.Time `json:"created_at"`
	// DoneAt은 입금 완료 시간입니다.
	DoneAt time.Time `json:"done_at"`
	// Amount는 입금 수량입니다.
	Amount string `json:"amount"`
	// Fee는 입금 수수료입니다.
	Fee string `json:"fee"`
	// TransactionType은 입금 유형입니다.
	TransactionType string `json:"transaction_type"`
}

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

// GetDepositParams는 개별 입금 조회를 위한 파라미터입니다.
type GetDepositParams struct {
	// UUID는 입금 UUID입니다.
	UUID string `json:"uuid,omitempty"`
	// TxID는 입금 TXID입니다.
	TxID string `json:"txid,omitempty"`
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
}

// DepositListParams는 입금 리스트 조회를 위한 파라미터입니다.
type DepositListParams struct {
	// Currency는 Currency 코드입니다.
	Currency string `json:"currency,omitempty"`
	// State는 입금 상태입니다.
	State string `json:"state,omitempty"`
	// UUIDs는 입금 UUID의 목록입니다.
	UUIDs []string `json:"uuids,omitempty"`
	// TxIDs는 입금 TXID의 목록입니다.
	TxIDs []string `json:"txids,omitempty"`
	// Limit는 개수 제한입니다. (default: 100, max: 100)
	Limit int `json:"limit,omitempty"`
	// Page는 페이지 수입니다. (default: 1)
	Page int `json:"page,omitempty"`
	// OrderBy는 정렬 방식입니다.
	OrderBy string `json:"order_by,omitempty"`
}

// GenerateCoinAddress는 입금 주소 생성을 요청합니다.
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

// GetDepositCoinChance는 디지털 자산 입금 정보를 조회합니다.
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

// DepositKRW는 원화 입금을 요청합니다.
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

// GetDeposit는 개별 입금을 조회합니다.
func (e *Exchange) GetDeposit(params *GetDepositParams) (*DepositInfo, error) {
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

	resp, err := e.Client.Get("/deposit", queryParams)
	if err != nil {
		return nil, err
	}

	var deposit DepositInfo
	if err := json.Unmarshal(resp, &deposit); err != nil {
		return nil, err
	}

	return &deposit, nil
}

// GetDeposits는 입금 리스트를 조회합니다.
func (e *Exchange) GetDeposits(params *DepositListParams) ([]DepositInfo, error) {
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

	resp, err := e.Client.Get("/deposits", queryParams)
	if err != nil {
		return nil, err
	}

	var deposits []DepositInfo
	if err := json.Unmarshal(resp, &deposits); err != nil {
		return nil, err
	}

	return deposits, nil
}
