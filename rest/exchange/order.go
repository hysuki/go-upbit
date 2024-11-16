package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// OrderType은 주문 유형을 나타냅니다
type OrderType string

const (
	// OrderTypeLimit는 지정가 주문입니다
	OrderTypeLimit OrderType = "limit"
	// OrderTypePrice는 시장가 매수 주문입니다
	OrderTypePrice OrderType = "price"
	// OrderTypeMarket는 시장가 매도 주문입니다
	OrderTypeMarket OrderType = "market"
	// OrderTypeBest는 최유리 지정가 주문입니다
	OrderTypeBest OrderType = "best"
)

// OrderSide는 주문 종류(매수/매도)를 나타냅니다
type OrderSide string

const (
	// OrderSideBid는 매수 주문입니다
	OrderSideBid OrderSide = "bid"
	// OrderSideAsk는 매도 주문입니다
	OrderSideAsk OrderSide = "ask"
)

// TimeInForce는 주문의 체결 조건을 나타냅니다
type TimeInForce string

const (
	// TimeInForceIOC는 "Immediate or Cancel" 조건입니다
	TimeInForceIOC TimeInForce = "ioc"
	// TimeInForceFOK는 "Fill or Kill" 조건입니다
	TimeInForceFOK TimeInForce = "fok"
)

// OrderState 상수들은 주문의 상태를 정의합니다.
const (
	// OrderStateWait는 체결 대기 상태입니다.
	OrderStateWait = "wait"
	// OrderStateWatch는 예약주문 대기 상태입니다.
	OrderStateWatch = "watch"
	// OrderStateDone은 전체 체결 완료 상태입니다.
	OrderStateDone = "done"
	// OrderStateCancel은 주문 취소 상태입니다.
	OrderStateCancel = "cancel"
)

// OrderBy 상수들은 정렬 방식을 정의합니다.
const (
	// OrderByAsc는 오름차순 정렬입니다.
	OrderByAsc = "asc"
	// OrderByDesc는 내림차순 정렬입니다.
	OrderByDesc = "desc"
)

// CreateOrderRequest는 주문 생성 요청 파라미터를 정의합니다
type CreateOrderRequest struct {
	Market      string      `json:"market,omitempty"`        // 마켓 ID (필수)
	Side        OrderSide   `json:"side,omitempty"`          // 주문 종류 (필수)
	Volume      string      `json:"volume,omitempty"`        // 주문량 (지정가, 시장가 매도 시 필수)
	Price       string      `json:"price,omitempty"`         // 주문 가격 (지정가, 시장가 매수 시 필수)
	OrderType   OrderType   `json:"ord_type,omitempty"`      // 주문 타입 (필수)
	Identifier  string      `json:"identifier,omitempty"`    // 조회용 사용자 지정값 (선택)
	TimeInForce TimeInForce `json:"time_in_force,omitempty"` // IOC, FOK 주문 설정
}

// Order는 주문 정보를 나타냅니다
type Order struct {
	UUID            string    `json:"uuid,omitempty"`             // 주문의 고유 아이디
	Side            OrderSide `json:"side,omitempty"`             // 주문 종류
	OrderType       OrderType `json:"ord_type,omitempty"`         // 주문 방식
	Price           string    `json:"price,omitempty"`            // 주문 당시 화폐 가격
	State           string    `json:"state,omitempty"`            // 주문 상태
	Market          string    `json:"market,omitempty"`           // 마켓의 유일키
	CreatedAt       string    `json:"created_at,omitempty"`       // 주문 생성 시간
	Volume          string    `json:"volume,omitempty"`           // 사용자가 입력한 주문 양
	RemainingVolume string    `json:"remaining_volume,omitempty"` // 체결 후 남은 주문 양
	ReservedFee     string    `json:"reserved_fee,omitempty"`     // 수수료로 예약된 비용
	RemainingFee    string    `json:"remaining_fee,omitempty"`    // 남은 수수료
	PaidFee         string    `json:"paid_fee,omitempty"`         // 사용된 수수료
	Locked          string    `json:"locked,omitempty"`           // 거래에 사용중인 비용
	ExecutedVolume  string    `json:"executed_volume,omitempty"`  // 체결된 양
	TradesCount     int       `json:"trades_count,omitempty"`     // 해당 주문에 걸린 체결 수
	TimeInForce     string    `json:"time_in_force,omitempty"`    // IOC, FOK 설정
}

// CancelOrderParams는 주문 취소를 위한 파라미터입니다.
type CancelOrderParams struct {
	UUID       string `json:"uuid,omitempty"`       // 취소할 주문의 UUID
	Identifier string `json:"identifier,omitempty"` // 조회용 사용자 지정값
}

// ClosedOrderParams는 종료된 주문 조회를 위한 파라미터입니다.
type ClosedOrderParams struct {
	Market    string     `json:"market,omitempty"`     // 마켓 ID
	State     string     `json:"state,omitempty"`      // 주문 상태
	States    []string   `json:"states,omitempty"`     // 주문 상태의 목록
	StartTime *time.Time `json:"start_time,omitempty"` // 조회 시작 시각
	EndTime   *time.Time `json:"end_time,omitempty"`   // 조회 종료 시각
	Limit     int        `json:"limit,omitempty"`      // 요청 개수
	OrderBy   string     `json:"order_by,omitempty"`   // 정렬 방식
}

// OpenOrderParams는 미체결 주문 조회를 위한 파라미터입니다.
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
	// ErrInvalidParams는 잘못된 파라미터가 전달되었을 때 발생하는 에러입니다.
	ErrInvalidParams = errors.New("invalid parameters")
	// ErrTooManyIDs는 UUID 또는 identifier가 100개를 초과할 때 발생하는 에러입니다.
	ErrTooManyIDs = errors.New("too many uuids or identifiers: maximum is 100")
)

// OrderByIDParams는 UUID 또는 identifier로 주문 리스트를 조회하기 위한 파라미터입니다.
type OrderByIDParams struct {
	// Market은 마켓 ID입니다.
	Market string `json:"market,omitempty"`
	// UUIDs는 주문 UUID의 목록입니다. (최대 100개)
	UUIDs []string `json:"uuids,omitempty"`
	// Identifiers는 주문 identifier의 목록입니다. (최대 100개)
	Identifiers []string `json:"identifiers,omitempty"`
	// OrderBy는 정렬 방식입니다. (asc, desc)
	OrderBy string `json:"order_by,omitempty"`
}

// GetOrdersByID는 UUID 또는 identifier로 주문 리스트를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders/uuids
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

// CreateOrder는 새로운 주문을 생성니다
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

// GetClosedOrders는 종료된 주문 리스트를 조회합니다.
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

// GetOpenOrders는 미체결 주문 리스트를 조회합니다.
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

// OrderChance는 마켓별 주문 가능 정보를 나타내는 구조체입니다.
type OrderChance struct {
	// BidFee는 매수 수수료 비율입니다.
	BidFee string `json:"bid_fee"`
	// AskFee는 매도 수수료 비율입니다.
	AskFee string `json:"ask_fee"`
	// Market은 마켓에 대한 정보입니다.
	Market OrderMarket `json:"market"`
	// BidAccount는 매수 시 사용하는 화폐의 계좌 상태입니다.
	BidAccount OrderAccount `json:"bid_account"`
	// AskAccount는 매도 시 사용하는 화폐의 계좌 상태입니다.
	AskAccount OrderAccount `json:"ask_account"`
}

// OrderMarket은 마켓 정보를 나타내는 구조체입니다.
type OrderMarket struct {
	// ID는 마켓의 유일 키입니다.
	ID string `json:"id"`
	// Name은 마켓 이름입니다.
	Name string `json:"name"`
	// OrderTypes는 지원 주문 방식입니다. (만료 예정)
	OrderTypes []string `json:"order_types"`
	// AskTypes는 매도 주문 지원 방식입니다.
	AskTypes []string `json:"ask_types"`
	// BidTypes는 매수 주문 지원 방식입니다.
	BidTypes []string `json:"bid_types"`
	// OrderSides는 지원 주문 종류입니다.
	OrderSides []string `json:"order_sides"`
	// Bid는 매수 시 제약사항입니다.
	Bid OrderConstraint `json:"bid"`
	// Ask는 매도 시 제약사항입니다.
	Ask OrderConstraint `json:"ask"`
	// MaxTotal은 최대 매도/매수 금액입니다.
	MaxTotal string `json:"max_total"`
	// State는 마켓 운영 상태입니다.
	State string `json:"state"`
}

// OrderConstraint는 마켓 거래 제약사항을 나타내는 구조체입니다.
type OrderConstraint struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// PriceUnit은 주문금액 단위입니다.
	PriceUnit string `json:"price_unit"`
	// MinTotal은 최소 매도/매수 금액입니다.
	MinTotal string `json:"min_total"`
}

// OrderAccount는 매수/매도 계좌 정보를 나타내는 구조체입니다.
type OrderAccount struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// Balance는 주문가능 금액/수량입니다.
	Balance string `json:"balance"`
	// Locked는 주문 중 묶여있는 금액/수량입니다.
	Locked string `json:"locked"`
	// AvgBuyPrice는 매수평균가입니다.
	AvgBuyPrice string `json:"avg_buy_price"`
	// AvgBuyPriceModified는 매수평균가 수정 여부입니다.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified"`
	// UnitCurrency는 평단가 기준 화폐입니다.
	UnitCurrency string `json:"unit_currency"`
}

// GetOrderChance는 마켓별 주문 가능 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/orders/chance
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
