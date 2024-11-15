package exchange

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
