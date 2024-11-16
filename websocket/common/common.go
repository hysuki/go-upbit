package common

type SubscribeOptions struct {
	Level          *float64
	IsOnlySnapshot *bool
	IsOnlyRealtime *bool
}

type StreamType string

const (
	StreamTypeSnapshot StreamType = "SNAPSHOT" // 스냅샷 데이터
	StreamTypeRealtime StreamType = "REALTIME" // 실시간 데이터
)

type ChangeType string

const (
	ChangeTypeRise ChangeType = "RISE" // 상승
	ChangeTypeEven ChangeType = "EVEN" // 보합
	ChangeTypeFall ChangeType = "FALL" // 하락
)

type MarketState string

const (
	MarketStatePreview  MarketState = "PREVIEW"  // 입금지원
	MarketStateActive   MarketState = "ACTIVE"   // 거래지원가능
	MarketStateDelisted MarketState = "DELISTED" // 거래지원종료
)

type MarketWarning string

const (
	MarketWarningNone    MarketWarning = "NONE"    // 해당없음
	MarketWarningCaution MarketWarning = "CAUTION" // 투자유의
)

type AskBidType string

const (
	AskBidTypeAsk AskBidType = "ASK" // 매도
	AskBidTypeBid AskBidType = "BID" // 매수
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"  // 지정가 주문
	OrderTypePrice  OrderType = "price"  // 시장가 매수 주문
	OrderTypeMarket OrderType = "market" // 시장가 매도 주문
	OrderTypeBest   OrderType = "best"   // 최유리 지정가 주문
)

type OrderState string

const (
	OrderStateWait   OrderState = "wait"   // 체결 대기
	OrderStateWatch  OrderState = "watch"  // 예약 주문 대기
	OrderStateTrade  OrderState = "trade"  // 체결 발생
	OrderStateDone   OrderState = "done"   // 전체 체결 완료
	OrderStateCancel OrderState = "cancel" // 주문 취소
)

type TimeInForce string

const (
	TimeInForceIOC TimeInForce = "ioc" // IOC 주문
	TimeInForceFOK TimeInForce = "fok" // FOK 주문
)
