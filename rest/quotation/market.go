package quotation

import (
	"encoding/json"
)

// MarketInfo는 마켓 정보를 나타내는 구조체입니다.
type MarketInfo struct {
	// Market은 업비트에서 제공중인 시장 정보입니다.
	Market string `json:"market"`
	// KoreanName은 거래 대상 디지털 자산 한글명입니다.
	KoreanName string `json:"korean_name"`
	// EnglishName은 거래 대상 디지털 자산 영문명입니다.
	EnglishName string `json:"english_name"`
	// MarketWarning은 유의 종목 여부입니다. (NONE: 해당 사항 없음, CAUTION: 투자유의)
	// Deprecated: 이 필드는 더 이상 사용되지 않습니다. MarketEvent 필드를 사용하세요.
	MarketWarning string `json:"market_warning,omitempty"`
	// MarketEvent는 유의 종목 지정 여부 및 주의 종목 지정 여부입니다.
	MarketEvent *MarketEvent `json:"market_event,omitempty"`
}

// MarketEvent는 마켓 이벤트 정보를 나타내는 구조체입니다.
type MarketEvent struct {
	// Warning은 유의 종목 지정 여부입니다.
	Warning bool `json:"warning"`
	// Caution은 주의 종목 지정 여부입니다.
	Caution MarketCaution `json:"caution"`
}

// MarketCaution은 주의 종목 경보 타입을 나타내는 구조체입니다.
type MarketCaution struct {
	// PriceFluctuations는 가격 급등락 경보 발령 여부입니다.
	PriceFluctuations bool `json:"PRICE_FLUCTUATIONS"`
	// TradingVolumeSoaring은 거래량 급증 경보 발령 여부입니다.
	TradingVolumeSoaring bool `json:"TRADING_VOLUME_SOARING"`
	// DepositAmountSoaring은 입금량 급증 경보 발령 여부입니다.
	DepositAmountSoaring bool `json:"DEPOSIT_AMOUNT_SOARING"`
	// GlobalPriceDifferences는 가격 차이 경보 발령 여부입니다.
	GlobalPriceDifferences bool `json:"GLOBAL_PRICE_DIFFERENCES"`
	// ConcentrationOfSmallAccounts은 소수 계정 집중 경보 발령 여부입니다.
	ConcentrationOfSmallAccounts bool `json:"CONCENTRATION_OF_SMALL_ACCOUNTS"`
}

// GetMarkets는 업비트에서 거래 가능한 마켓 목록을 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/market/all
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
