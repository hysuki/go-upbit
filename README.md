# go-upbit

**Go-Upbit**은 **업비트 API**를 쉽게 사용하기 위한 **Go 클라이언트 라이브러리**입니다.  
REST API 및 WebSocket API를 지원하며, 간단한 인터페이스로 다양한 기능을 제공합니다.

> 충분히 테스트 되지 않음, 테스트 코드 작성 지원 환영

## 요청 수 제한

업비트 API는 초당/분당 요청 수 제한이 있습니다. API 사용 시 아래 제한 사항을 반드시 확인하시기 바랍니다.

자세한 내용은 [업비트 공식 문서](https://docs.upbit.com/docs/user-request-guide)를 참고하세요.

---

## 기능

### REST API
- **거래소 API**
  - 자산 조회
  - 주문 생성 및 관리
  - 입출금 관리
- **시세 API**
  - 마켓 코드 조회
  - 캔들 데이터 조회 (분/일/주/월 단위)
  - 현재가 조회
  - 호가 정보 조회
  - 최근 체결 내역

### WebSocket API
- **Public WebSocket**
  - 현재가 정보
  - 체결 내역
  - 호가 정보
- **Private WebSocket**  
  - 내 자산 실시간 조회
  - 내 주문 실시간 조회

---

## 설치
```bash
go get github.com/hysuki/go-upbit
```

---

## 사용 예시

### 클라이언트 초기화
```go
import (
	"log"
	"time"
	"github.com/hysuki/go-upbit"
)

func main() {
	client, err := upbit.NewUpbitClient(
		upbit.WithKeys("ACCESS_KEY", "SECRET_KEY"),
		upbit.WithPingInterval(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```

### REST API 사용 예시
```go
// 마켓 코드 조회
markets, err := client.RestAPI.GetQuotation().GetMarkets(false)
if err != nil {
	log.Printf("에러: %v", err)
	return
}
log.Printf("마켓 목록: %+v", markets)

// 계좌 조회
accounts, err := client.RestAPI.GetExchange().GetAccounts()
if err != nil {
	log.Printf("에러: %v", err)
	return
}
log.Printf("계좌 정보: %+v", accounts)

// 주문하기
order, err := client.RestAPI.GetExchange().CreateOrder(&exchange.CreateOrderRequest{
	Market:    "KRW-BTC",
	Side:      exchange.OrderSideBid,
	Volume:    "0.01",
	Price:     "30000000",
	OrderType: exchange.OrderTypeLimit,
})
if err != nil {
	log.Printf("에러: %v", err)
	return
}
log.Printf("주문 결과: %+v", order)
```

### WebSocket API 사용 예시
```go
// 원화 마켓 코드 필터링
var codes []string
markets, err := client.RestAPI.GetQuotation().GetMarkets(false)
if err != nil {
    log.Fatal(err)
}

for _, market := range markets {
    if strings.HasPrefix(market.Market, "KRW-") {
        codes = append(codes, market.Market)
    }
}

// WebSocket 구독 설정
isOnlyRealtime := true
opts := common.SubscribeOptions{
    IsOnlyRealtime: &isOnlyRealtime,
}

client.PublicWS.Subscribe(nil,
    public.AddSubscribe(public.MessageTypeOrderbook, codes, &opts),
    public.AddSubscribe(public.MessageTypeTicker, codes, &opts),
    public.AddSubscribe(public.MessageTypeTrade, codes, &opts),
)

// 메시지 핸들러 시작
client.PublicWS.StartMessageHandler()

// 현재가 정보 처리
go func() {
    onTicker(ctx, client, func(ticker *public.Ticker) {
        log.Printf("현재가: %+v", ticker)
    })
}()

// 체결 정보 처리
tradeChan := make(chan *public.Trade, 1000)
go func() {
    onTrade(ctx, client, func(trade *public.Trade) {
        log.Printf("체결: %+v", trade)
        tradeChan <- trade
    })
}()

// 호가 정보 처리
go func() {
    onOrderBook(ctx, client, func(orderBook *public.Orderbook) {
        log.Printf("호가: %+v", orderBook)
    })
}()
```

## 참고 문서
- [업비트 API 문서](https://docs.upbit.com/)
- [REST API 레퍼런스](https://docs.upbit.com/reference)
- [WebSocket API 레퍼런스](https://docs.upbit.com/reference/websocket-시세-유의사항)

## 라이선스
MIT License
