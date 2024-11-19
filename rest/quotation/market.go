package quotation

import (
	"encoding/json"
)

// Package quotation은 Upbit 거래소의 시세 조회 관련 API를 제공합니다.

// MarketInfo는 마켓 정보를 나타냅니다.
type MarketInfo struct {
	Market        string       `json:"market"`                   // 마켓 ID
	KoreanName    string       `json:"korean_name"`              // 거래 대상 디지털 자산 한글명
	EnglishName   string       `json:"english_name"`             // 거래 대상 디지털 자산 영문명
	MarketWarning string       `json:"market_warning,omitempty"` // 유의 종목 여부 (Deprecated)
	MarketEvent   *MarketEvent `json:"market_event,omitempty"`   // 마켓 이벤트 정보
}

// MarketEvent는 마켓 이벤트 정보를 나타냅니다.
type MarketEvent struct {
	Warning bool          `json:"warning"` // 유의 종목 지정 여부
	Caution MarketCaution `json:"caution"` // 주의 종목 지정 여부
}

// MarketCaution은 주의 종목 경보 타입을 나타냅니다.
type MarketCaution struct {
	PriceFluctuations            bool `json:"PRICE_FLUCTUATIONS"`              // 가격 급등락 경보 발령 여부
	TradingVolumeSoaring         bool `json:"TRADING_VOLUME_SOARING"`          // 거래량 급증 경보 발령 여부
	DepositAmountSoaring         bool `json:"DEPOSIT_AMOUNT_SOARING"`          // 입금량 급증 경보 발령 여부
	GlobalPriceDifferences       bool `json:"GLOBAL_PRICE_DIFFERENCES"`        // 가격 차이 경보 발령 여부
	ConcentrationOfSmallAccounts bool `json:"CONCENTRATION_OF_SMALL_ACCOUNTS"` // 소수 계정 집중 경보 발령 여부
}

// GetMarkets는 업비트에서 거래 가능한 마켓 목록을 조회합니다.
// isDetails가 true인 경우 마켓 이벤트 정보를 포함하여 반환합니다.
func (q *Quotation) GetMarkets(isDetails bool) ([]MarketInfo, error) {
	params := make(map[string]string)
	if isDetails {
		params["is_details"] = "true"
	}

	resp, err := q.Client.Get("/market/all", params)
	if err != nil {
		return nil, err
	}

	var markets []MarketInfo
	if err := json.Unmarshal(resp, &markets); err != nil {
		return nil, err
	}

	return markets, nil
}
