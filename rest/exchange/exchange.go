package exchange

import (
	"github.com/hysuki/go-upbit/rest/client"
)

// Exchange는 거래소 API 기능을 제공하는 구조체입니다.
type Exchange struct {
	Client client.RestClient
}

// NewExchange는 새로운 Exchange 인스턴스를 생성합니다.
// client는 REST API 호출을 처리할 클라이언트입니다.
func NewExchange(client client.RestClient) *Exchange {
	return &Exchange{
		Client: client,
	}
}
