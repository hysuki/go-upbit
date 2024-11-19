package common

// SubscribeOptions는 웹소켓 구독 옵션을 나타냅니다.
type SubscribeOptions struct {
	Level          *float64 // 호가 모아보기 단위
	IsOnlySnapshot *bool    // 스냅샷만 수신 여부
	IsOnlyRealtime *bool    // 실시간만 수신 여부
}

// StreamType은 데이터 스트림 유형을 나타냅니다.
type StreamType string

// 데이터 스트림 유형을 정의하는 상수들입니다.
const (
	StreamTypeSnapshot StreamType = "SNAPSHOT" // 스냅샷 데이터
	StreamTypeRealtime StreamType = "REALTIME" // 실시간 데이터
)

// ChangeType은 가격 변화 유형을 나타냅니다.
type ChangeType string

// 가격 변화 유형을 정의하는 상수들입니다.
const (
	ChangeTypeRise ChangeType = "RISE" // 상승
	ChangeTypeEven ChangeType = "EVEN" // 보합
	ChangeTypeFall ChangeType = "FALL" // 하락
)

// MarketState는 마켓 상태를 나타냅니다.
type MarketState string

// 마켓 상태를 정의하는 상수들입니다.
const (
	MarketStatePreview  MarketState = "PREVIEW"  // 입금지원
	MarketStateActive   MarketState = "ACTIVE"   // 거래지원가능
	MarketStateDelisted MarketState = "DELISTED" // 거래지원종료
)

// MarketWarning은 마켓 경고 상태를 나타냅니다.
type MarketWarning string

// 마켓 경고 상태를 정의하는 상수들입니다.
const (
	MarketWarningNone    MarketWarning = "NONE"    // 해당없음
	MarketWarningCaution MarketWarning = "CAUTION" // 투자유의
)

// AskBidType은 매도/매수 구분을 나타냅니다.
type AskBidType string

// 매도/매수 구분을 정의하는 상수들입니다.
const (
	AskBidTypeAsk AskBidType = "ASK" // 매도
	AskBidTypeBid AskBidType = "BID" // 매수
)

// OrderType은 주문 유형을 나타냅니다.
type OrderType string

// 주문 유형을 정의하는 상수들입니다.
const (
	OrderTypeLimit  OrderType = "limit"  // 지정가 주문
	OrderTypePrice  OrderType = "price"  // 시장가 매수 주문
	OrderTypeMarket OrderType = "market" // 시장가 매도 주문
	OrderTypeBest   OrderType = "best"   // 최유리 지정가 주문
)

// OrderState는 주문 상태를 나타냅니다.
type OrderState string

// 주문 상태를 정의하는 상수들입니다.
const (
	OrderStateWait   OrderState = "wait"   // 체결 대기
	OrderStateWatch  OrderState = "watch"  // 예약 주문 대기
	OrderStateTrade  OrderState = "trade"  // 체결 발생
	OrderStateDone   OrderState = "done"   // 전체 체결 완료
	OrderStateCancel OrderState = "cancel" // 주문 취소
)

// TimeInForce는 주문의 체결 조건을 나타냅니다.
type TimeInForce string

// 주문 체결 조건을 정의하는 상수들입니다.
const (
	TimeInForceIOC TimeInForce = "ioc" // IOC 주문
	TimeInForceFOK TimeInForce = "fok" // FOK 주문
)
