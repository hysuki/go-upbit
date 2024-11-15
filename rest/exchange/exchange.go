package exchange

import (
	"upbit.yougcha.bot/pkg/upbit/rest/client"
)

// Exchange는 거래소 API 기능을 제공합니다.
type Exchange struct {
	Client client.RestClient
}

// NewExchange는 새로운 Exchange 인스턴스를 생성합니다.
func NewExchange(client client.RestClient) *Exchange {
	return &Exchange{
		Client: client,
	}
}
