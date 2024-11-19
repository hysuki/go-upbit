package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Package exchange는 Upbit 거래소의 주문 관련 API를 제공합니다.

// OrderType은 주문 유형을 나타냅니다.
type OrderType string

// 주문 유형을 정의하는 상수들입니다.
const (
	OrderTypeLimit  OrderType = "limit"  // 지정가 주문
	OrderTypePrice  OrderType = "price"  // 시장가 매수 주문
	OrderTypeMarket OrderType = "market" // 시장가 매도 주문
	OrderTypeBest   OrderType = "best"   // 최유리 지정가 주문
)

// OrderSide는 주문 종류를 나타냅니다.
type OrderSide string

// 주문 종류를 정의하는 상수들입니다.
const (
	OrderSideBid OrderSide = "bid" // 매수 주문
	OrderSideAsk OrderSide = "ask" // 매도 주문
)

// TimeInForce는 주문의 체결 조건을 나타냅니다.
type TimeInForce string

// 주문 체결 조건을 정의하는 상수들입니다.
const (
	TimeInForceIOC TimeInForce = "ioc" // Immediate or Cancel 조건
	TimeInForceFOK TimeInForce = "fok" // Fill or Kill 조건
)

// 주문 상태를 정의하는 상수들입니다.
const (
	OrderStateWait   = "wait"   // 체결 대기
	OrderStateWatch  = "watch"  // 예약주문 대기
	OrderStateDone   = "done"   // 전체 체결 완료
	OrderStateCancel = "cancel" // 주문 취소
)

// 정렬 방식을 정의하는 상수들입니다.
const (
	OrderByAsc  = "asc"  // 오름차순 정렬
	OrderByDesc = "desc" // 내림차순 정렬
)

// CreateOrderRequest는 주문 생성에 필요한 파라미터입니다.
type CreateOrderRequest struct {
	Market      string      `json:"market,omitempty"`        // 마켓 ID
	Side        OrderSide   `json:"side,omitempty"`          // 주문 종류
	Volume      string      `json:"volume,omitempty"`        // 주문량
	Price       string      `json:"price,omitempty"`         // 주문 가격
	OrderType   OrderType   `json:"ord_type,omitempty"`      // 주문 타입
	Identifier  string      `json:"identifier,omitempty"`    // 조회용 사용자 지정값
	TimeInForce TimeInForce `json:"time_in_force,omitempty"` // 체결 조건
}

// Order는 주문 정보를 나타냅니다.
type Order struct {
	UUID            string    `json:"uuid,omitempty"`             // 주문의 고유 ID
	Side            OrderSide `json:"side,omitempty"`             // 주문 종류
	OrderType       OrderType `json:"ord_type,omitempty"`         // 주문 방식
	Price           string    `json:"price,omitempty"`            // 주문 가격
	State           string    `json:"state,omitempty"`            // 주문 상태
	Market          string    `json:"market,omitempty"`           // 마켓 ID
	CreatedAt       string    `json:"created_at,omitempty"`       // 주문 생성 시각
	Volume          string    `json:"volume,omitempty"`           // 주문량
	RemainingVolume string    `json:"remaining_volume,omitempty"` // 잔여 주문량
	ReservedFee     string    `json:"reserved_fee,omitempty"`     // 예약된 수수료
	RemainingFee    string    `json:"remaining_fee,omitempty"`    // 잔여 수수료
	PaidFee         string    `json:"paid_fee,omitempty"`         // 사용된 수수료
	Locked          string    `json:"locked,omitempty"`           // 거래에 사용된 비용
	ExecutedVolume  string    `json:"executed_volume,omitempty"`  // 체결된 양
	TradesCount     int       `json:"trades_count,omitempty"`     // 체결 수
	TimeInForce     string    `json:"time_in_force,omitempty"`    // 체결 조건
}

// CancelOrderParams는 주문 취소에 필요한 파라미터입니다.
type CancelOrderParams struct {
	UUID       string `json:"uuid,omitempty"`       // 취소할 주문의 UUID
	Identifier string `json:"identifier,omitempty"` // 조회용 사용자 지정값
}

// ClosedOrderParams는 완료된 주문 조회에 필요한 파라미터입니다.
type ClosedOrderParams struct {
	Market    string     `json:"market,omitempty"`     // 마켓 ID
	State     string     `json:"state,omitempty"`      // 주문 상태
	States    []string   `json:"states,omitempty"`     // 주문 상태의 목록
	StartTime *time.Time `json:"start_time,omitempty"` // 조회 시작 시각
	EndTime   *time.Time `json:"end_time,omitempty"`   // 조회 종료 시각
	Limit     int        `json:"limit,omitempty"`      // 요청 개수
	OrderBy   string     `json:"order_by,omitempty"`   // 정렬 방식
}

// OpenOrderParams는 미체결 주문 조회에 필요한 파라미터입니다.
type OpenOrderParams struct {
	Market  string   `json:"market,omitempty"`   // 마켓 ID
	State   string   `json:"state,omitempty"`    // 주문 상태
	States  []string `json:"states,omitempty"`   // 주문 상태의 목록
	Page    int      `json:"page,omitempty"`     // 페이지 수
	Limit   int      `json:"limit,omitempty"`    // 요청 개수
	OrderBy string   `json:"order_by,omitempty"` // 정렬 방식
}

// 에러 정의
var (
	ErrInvalidParams = errors.New("invalid parameters")                            // 잘못된 파라미터 에러
	ErrTooManyIDs    = errors.New("too many uuids or identifiers: maximum is 100") // ID 초과 에러
)

// OrderByIDParams는 주문 조회에 필요한 파라미터입니다.
type OrderByIDParams struct {
	Market      string   `json:"market,omitempty"`      // 마켓 ID
	UUIDs       []string `json:"uuids,omitempty"`       // 주문 UUID 목록
	Identifiers []string `json:"identifiers,omitempty"` // 주문 식별자 목록
	OrderBy     string   `json:"order_by,omitempty"`    // 정렬 방식
}

// OrderChance는 마켓별 주문 가능 정보를 나타냅니다.
type OrderChance struct {
	BidFee     string       `json:"bid_fee"`     // 매수 수수료 비율
	AskFee     string       `json:"ask_fee"`     // 매도 수수료 비율
	Market     OrderMarket  `json:"market"`      // 마켓 정보
	BidAccount OrderAccount `json:"bid_account"` // 매수 계좌 정보
	AskAccount OrderAccount `json:"ask_account"` // 매도 계좌 정보
}

// OrderMarket은 마켓 정보를 나타냅니다.
type OrderMarket struct {
	ID         string          `json:"id"`          // 마켓 ID
	Name       string          `json:"name"`        // 마켓 이름
	OrderTypes []string        `json:"order_types"` // 지원 주문 방식
	AskTypes   []string        `json:"ask_types"`   // 매도 주문 지원 방식
	BidTypes   []string        `json:"bid_types"`   // 매수 주문 지원 방식
	OrderSides []string        `json:"order_sides"` // 지원 주문 종류
	Bid        OrderConstraint `json:"bid"`         // 매수 제약사항
	Ask        OrderConstraint `json:"ask"`         // 매도 제약사항
	MaxTotal   string          `json:"max_total"`   // 최대 매도/매수 금액
	State      string          `json:"state"`       // 마켓 운영 상태
}

// OrderConstraint는 마켓 거래 제약사항을 나타냅니다.
type OrderConstraint struct {
	Currency  string `json:"currency"`   // 화폐를 의미하는 영문 대문자 코드
	PriceUnit string `json:"price_unit"` // 주문금액 단위
	MinTotal  string `json:"min_total"`  // 최소 매도/매수 금액
}

// OrderAccount는 매수/매도 계좌 정보를 나타냅니다.
type OrderAccount struct {
	Currency            string `json:"currency"`               // 화폐를 의미하는 영문 대문자 코드
	Balance             string `json:"balance"`                // 주문가능 금액/수량
	Locked              string `json:"locked"`                 // 주문 중 묶여있는 금액/수량
	AvgBuyPrice         string `json:"avg_buy_price"`          // 매수평균가
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"` // 매수평균가 수정 여부
	UnitCurrency        string `json:"unit_currency"`          // 평단가 기준 화폐
}

// GetOrdersByID는 UUID 또는 식별자로 주문 목록을 조회합니다.
func (e *Exchange) GetOrdersByID(params *OrderByIDParams) ([]Order, error) {
	if params == nil {
		return nil, ErrInvalidParams
	}

	// UUID와 identifier는 동시에 사용할 수 없습니다.
	if len(params.UUIDs) > 0 && len(params.Identifiers) > 0 {
		return nil, ErrInvalidParams
	}

	// UUID 또는 identifier 중 하나는 필수입니다.
	if len(params.UUIDs) == 0 && len(params.Identifiers) == 0 {
		return nil, ErrInvalidParams
	}

	// 최대 100개까지만 조회 가능합니다.
	if len(params.UUIDs) > 100 || len(params.Identifiers) > 100 {
		return nil, ErrTooManyIDs
	}

	queryParams := make(map[string]string)

	if params.Market != "" {
		queryParams["market"] = params.Market
	}

	if len(params.UUIDs) > 0 {
		uuidsBytes, err := json.Marshal(params.UUIDs)
		if err != nil {
			return nil, err
		}
		queryParams["uuids"] = string(uuidsBytes)
	}

	if len(params.Identifiers) > 0 {
		identifiersBytes, err := json.Marshal(params.Identifiers)
		if err != nil {
			return nil, err
		}
		queryParams["identifiers"] = string(identifiersBytes)
	}

	if params.OrderBy != "" {
		queryParams["order_by"] = params.OrderBy
	}

	resp, err := e.Client.Get("/orders/uuids", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// CreateOrder는 새로운 주문을 생성합니다.
func (e *Exchange) CreateOrder(request *CreateOrderRequest) (*Order, error) {
	if request == nil {
		return nil, ErrInvalidParams
	}

	resp, err := e.Client.Post("/orders", request)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// CancelOrder는 주문을 취소합니다.
func (e *Exchange) CancelOrder(params *CancelOrderParams) (*Order, error) {
	if params == nil {
		return nil, errors.New("params cannot be nil")
	}

	if params.UUID == "" && params.Identifier == "" {
		return nil, errors.New("either uuid or identifier must be provided")
	}

	if params.UUID != "" && params.Identifier != "" {
		return nil, errors.New("uuid and identifier cannot be used together")
	}

	queryParams := make(map[string]string)
	if params.UUID != "" {
		queryParams["uuid"] = params.UUID
	}
	if params.Identifier != "" {
		queryParams["identifier"] = params.Identifier
	}

	resp, err := e.Client.Delete("/order", queryParams)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetClosedOrders는 완료된 주문 목록을 조회합니다.
func (e *Exchange) GetClosedOrders(params *ClosedOrderParams) ([]Order, error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Market != "" {
			queryParams["market"] = params.Market
		}
		if params.State != "" {
			if len(params.States) > 0 {
				return nil, errors.New("state and states cannot be used together")
			}
			queryParams["state"] = params.State
		}
		if len(params.States) > 0 {
			statesBytes, err := json.Marshal(params.States)
			if err != nil {
				return nil, err
			}
			queryParams["states"] = string(statesBytes)
		}
		if params.StartTime != nil {
			queryParams["start_time"] = params.StartTime.Format(time.RFC3339)
		}
		if params.EndTime != nil {
			queryParams["end_time"] = params.EndTime.Format(time.RFC3339)
		}
		if params.Limit > 0 {
			if params.Limit > 1000 {
				return nil, errors.New("limit cannot exceed 1000")
			}
			queryParams["limit"] = strconv.Itoa(params.Limit)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}

		// 시간 범위 검증
		if params.StartTime != nil && params.EndTime != nil {
			duration := params.EndTime.Sub(*params.StartTime)
			if duration > 7*24*time.Hour {
				return nil, errors.New("time range cannot exceed 7 days")
			}
		}
	}

	resp, err := e.Client.Get("/orders/closed", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOpenOrders는 미체결 주문 목록을 조회합니다.
func (e *Exchange) GetOpenOrders(params *OpenOrderParams) ([]Order, error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Market != "" {
			queryParams["market"] = params.Market
		}
		if params.State != "" {
			if len(params.States) > 0 {
				return nil, errors.New("state and states cannot be used together")
			}
			queryParams["state"] = params.State
		}
		if len(params.States) > 0 {
			statesBytes, err := json.Marshal(params.States)
			if err != nil {
				return nil, err
			}
			queryParams["states"] = string(statesBytes)
		}
		if params.Page > 0 {
			queryParams["page"] = strconv.Itoa(params.Page)
		}
		if params.Limit > 0 {
			if params.Limit > 100 {
				return nil, errors.New("limit cannot exceed 100")
			}
			queryParams["limit"] = strconv.Itoa(params.Limit)
		}
		if params.OrderBy != "" {
			queryParams["order_by"] = params.OrderBy
		}
	}

	resp, err := e.Client.Get("/orders/open", queryParams)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderChance는 마켓별 주문 가능 정보를 조회합니다.
func (e *Exchange) GetOrderChance(market string) (*OrderChance, error) {
	if market == "" {
		return nil, errors.New("market is required")
	}

	params := map[string]string{
		"market": market,
	}

	// GetAccounts()와 동일한 방식으로 호출
	resp, err := e.Client.Get("/orders/chance", params)
	if err != nil {
		return nil, err
	}

	var orderChance OrderChance
	if err := json.Unmarshal(resp, &orderChance); err != nil {
		return nil, err
	}

	return &orderChance, nil
}
