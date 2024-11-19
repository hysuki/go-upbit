package quotation

import (
	"github.com/hysuki/go-upbit/rest/client"
)

// Package quotation은 Upbit 거래소의 시세 조회 관련 API를 제공합니다.

// Quotation은 시세 조회 API를 호출하기 위한 클라이언트입니다.
// REST API 요청을 처리하는 Client를 포함합니다.
type Quotation struct {
	Client client.RestClient // REST API 클라이언트
}

// NewQuotation은 새로운 Quotation 인스턴스를 생성합니다.
// REST API 클라이언트를 파라미터로 받아 Quotation 구조체를 초기화합니다.
func NewQuotation(client client.RestClient) *Quotation {
	return &Quotation{
		Client: client,
	}
}
