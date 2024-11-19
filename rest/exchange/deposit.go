package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Package exchange는 Upbit 거래소의 입출금 관련 API를 제공합니다.
const (
	DepositStateProcessing          = "PROCESSING"            // 입금 진행중
	DepositStateAccepted            = "ACCEPTED"              // 입금 완료
	DepositStateCancelled           = "CANCELLED"             // 입금 취소
	DepositStateRejected            = "REJECTED"              // 입금 거절
	DepositStateTravelRuleSuspected = "TRAVEL_RULE_SUSPECTED" // 트래블룰 추가 인증 필요
	DepositStateRefunding           = "REFUNDING"             // 입금액 반환 진행중
	DepositStateRefunded            = "REFUNDED"              // 입금액 반환 완료
)

// 입금 유형을 나타내는 상수들
const (
	DepositTransactionTypeDefault  = "default"  // 일반 입금
	DepositTransactionTypeInternal = "internal" // 바로 입금
)

// DepositInfo는 입금 정보를 나타냅니다.
type DepositInfo struct {
	Type            string    `json:"type"`             // 입금 종류
	UUID            string    `json:"uuid"`             // 입금 고유 식별자
	Currency        string    `json:"currency"`         // 화폐를 의미하는 영문 대문자 코드
	NetType         string    `json:"net_type"`         // 입금 네트워크 종류
	TxID            string    `json:"txid"`             // 블록체인 트랜잭션 ID
	State           string    `json:"state"`            // 입금 상태
	CreatedAt       time.Time `json:"created_at"`       // 입금 요청 시각
	DoneAt          time.Time `json:"done_at"`          // 입금 완료 시각
	Amount          string    `json:"amount"`           // 입금 수량
	Fee             string    `json:"fee"`              // 입금 수수료
	TransactionType string    `json:"transaction_type"` // 입금 유형
}

// DepositAddress는 입금 주소 정보를 나타내는 구조체입니다.
type DepositAddress struct {
	Currency         string `json:"currency"`          // 화폐를 의미하는 영문 대문자 코드
	NetType          string `json:"net_type"`          // 입금 네트워크 종류
	DepositAddress   string `json:"deposit_address"`   // 입금 주소
	SecondaryAddress string `json:"secondary_address"` // 2차 입금 주소 (일부 화폐에서 사용)
}

// GenerateCoinAddressParams는 입금 주소 생성 요청에 필요한 파라미터입니다.
type GenerateCoinAddressParams struct {
	Currency string `json:"currency"` // 화폐를 의미하는 영문 대문자 코드
	NetType  string `json:"net_type"` // 입금 네트워크 종류
}

// GenerateCoinAddressResponse는 입금 주소 생성 요청에 대한 응답입니다.
type GenerateCoinAddressResponse struct {
	Success bool   `json:"success"` // 요청 성공 여부
	Message string `json:"message"` // 요청 결과 메시지
}

// DepositCoinChance는 암호화폐 입금 관련 정보를 나타내는 구조체입니다.
type DepositCoinChance struct {
	Currency                    string `json:"currency"`                      // 화폐를 의미하는 영문 대문자 코드
	NetType                     string `json:"net_type"`                      // 입금 네트워크 종류
	IsDepositPossible           bool   `json:"is_deposit_possible"`           // 입금 가능 여부
	DepositImpossibleReason     string `json:"deposit_impossible_reason"`     // 입금이 불가능한 경우 그 사유
	MinimumDepositAmount        string `json:"minimum_deposit_amount"`        // 최소 입금 가능 수량
	MinimumDepositConfirmations int    `json:"minimum_deposit_confirmations"` // 입금 확인에 필요한 최소 블록 확인 수
	DecimalPrecision            int    `json:"decimal_precision"`             // 입금 수량의 소수점 자릿수
}

// DepositKRWParams는 원화 입금 요청에 필요한 파라미터입니다.
type DepositKRWParams struct {
	Amount        string `json:"amount"`          // 입금 금액
	TwoFactorType string `json:"two_factor_type"` // 2차 인증 방식 (kakao, naver, hana)
}

// GetDepositParams는 개별 입금 조회 시 필요한 파라미터입니다.
type GetDepositParams struct {
	UUID     string `json:"uuid,omitempty"`     // 입금 UUID
	TxID     string `json:"txid,omitempty"`     // 입금 트랜잭션 ID
	Currency string `json:"currency,omitempty"` // 화폐를 의미하는 영문 대문자 코드
}

// DepositListParams는 입금 목록 조회 시 필요한 파라미터입니다.
type DepositListParams struct {
	Currency string   `json:"currency,omitempty"` // 화폐를 의미하는 영문 대문자 코드
	State    string   `json:"state,omitempty"`    // 입금 상태
	UUIDs    []string `json:"uuids,omitempty"`    // 조회하고자 하는 입금 UUID 목록
	TxIDs    []string `json:"txids,omitempty"`    // 조회하고자 하는 트랜잭션 ID 목록
	Limit    int      `json:"limit,omitempty"`    // 반환되는 항목의 최대 개수 (최대 100)
	Page     int      `json:"page,omitempty"`     // 조회할 페이지
	OrderBy  string   `json:"order_by,omitempty"` // 정렬 기준
}

// GenerateCoinAddress는 입금 주소를 생성합니다.
// 생성된 주소 정보와 함께 성공 여부를 반환합니다.
// 파라미터로 화폐 코드와 네트워크 유형이 필요합니다.
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

// GetCoinAddresses는 사용자의 전체 입금 주소 목록을 조회합니다.
// 등록된 모든 화폐의 입금 주소 정보를 반환합니다.
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

// GetCoinAddress는 특정 화폐의 입금 주소를 조회합니다.
// currency는 화폐 코드, netType은 네트워크 유형을 지정합니다.
// 해당 화폐의 입금 주소가 없는 경우 오류를 반환합니다.
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

// GetDepositCoinChance는 암호화폐 입금 관련 정보를 조회합니다.
// currency로 지정한 화폐의 입금 가능 여부와 최소 입금액 등의 정보를 반환합니다.
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
// 입금액과 2차 인증 방식을 지정해야 하며, 입금 성공 시 입금 정보를 반환합니다.
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

// GetDeposit는 특정 입금 내역을 조회합니다.
// UUID, TxID, Currency 중 최소 하나의 식별자가 필요합니다.
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

// GetDeposits는 입금 목록을 조회합니다.
// 조회 조건을 params로 지정할 수 있으며, 미지정 시 기본값이 적용됩니다.
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
