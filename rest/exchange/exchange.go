package exchange

import (
	"github.com/hysuki/go-upbit/rest/client"
)

// Package exchange는 Upbit 거래소의 API 기능을 제공합니다.

// Exchange는 Upbit 거래소 API를 호출하기 위한 클라이언트입니다.
// REST API 요청을 처리하는 Client를 포함합니다.
type Exchange struct {
	Client client.RestClient // REST API 클라이언트
}

// NewExchange는 새로운 Exchange 인스턴스를 생성합니다.
// REST API 클라이언트를 파라미터로 받아 Exchange 구조체를 초기화합니다.
func NewExchange(client client.RestClient) *Exchange {
	return &Exchange{
		Client: client,
	}
}
