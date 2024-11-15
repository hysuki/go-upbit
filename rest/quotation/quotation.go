package quotation

import (
	"upbit.yougcha.bot/pkg/upbit/rest/client"
)

// Quotation은 시세 API 기능을 제공합니다.
type Quotation struct {
	Client client.RestClient
}

// NewQuotation은 새로운 Quotation 인스턴스를 생성합니다.
func NewQuotation(client client.RestClient) *Quotation {
	return &Quotation{
		Client: client,
	}
}
