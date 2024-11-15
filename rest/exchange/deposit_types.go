package exchange

import "time"

// DepositState는 입금 상태를 정의합니다.
const (
	// DepositStateProcessing은 입금 진행중 상태입니다.
	DepositStateProcessing = "PROCESSING"
	// DepositStateAccepted는 완료 상태입니다.
	DepositStateAccepted = "ACCEPTED"
	// DepositStateCancelled는 취소됨 상태입니다.
	DepositStateCancelled = "CANCELLED"
	// DepositStateRejected는 거절됨 상태입니다.
	DepositStateRejected = "REJECTED"
	// DepositStateTravelRuleSuspected는 트래블룰 추가 인증 대기중 상태입니다.
	DepositStateTravelRuleSuspected = "TRAVEL_RULE_SUSPECTED"
	// DepositStateRefunding는 반환절차 진행중 상태입니다.
	DepositStateRefunding = "REFUNDING"
	// DepositStateRefunded는 반환됨 상태입니다.
	DepositStateRefunded = "REFUNDED"
)

// DepositTransactionType은 입금 유형을 정의합니다.
const (
	// DepositTransactionTypeDefault는 일반입금입니다.
	DepositTransactionTypeDefault = "default"
	// DepositTransactionTypeInternal는 바로입금입니다.
	DepositTransactionTypeInternal = "internal"
)

// DepositInfo는 입금 정보를 나타내는 구조체입니다.
type DepositInfo struct {
	// Type은 입출금 종류입니다.
	Type string `json:"type"`
	// UUID는 입금에 대한 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 입금 네트워크입니다.
	NetType string `json:"net_type"`
	// TxID는 입금의 트랜잭션 아이디입니다.
	TxID string `json:"txid"`
	// State는 입금 상태입니다.
	State string `json:"state"`
	// CreatedAt은 입금 생성 시간입니다.
	CreatedAt time.Time `json:"created_at"`
	// DoneAt은 입금 완료 시간입니다.
	DoneAt time.Time `json:"done_at"`
	// Amount는 입금 수량입니다.
	Amount string `json:"amount"`
	// Fee는 입금 수수료입니다.
	Fee string `json:"fee"`
	// TransactionType은 입금 유형입니다.
	TransactionType string `json:"transaction_type"`
}
