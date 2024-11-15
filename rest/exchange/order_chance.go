package exchange

import (
	"encoding/json"
	"errors"
)

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
	MinTotal float64 `json:"min_total"`
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
