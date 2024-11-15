package exchange

import "time"

// WithdrawState는 출금 상태를 정의합니다.
const (
	// WithdrawStateWaiting은 대기중 상태입니다.
	WithdrawStateWaiting = "WAITING"
	// WithdrawStateProcessing은 진행중 상태입니다.
	WithdrawStateProcessing = "PROCESSING"
	// WithdrawStateDone은 완료 상태입니다.
	WithdrawStateDone = "DONE"
	// WithdrawStateFailed는 실패 상태입니다.
	WithdrawStateFailed = "FAILED"
	// WithdrawStateCancelled는 취소됨 상태입니다.
	WithdrawStateCancelled = "CANCELLED"
	// WithdrawStateRejected는 거절됨 상태입니다.
	WithdrawStateRejected = "REJECTED"
)

// WithdrawTransactionType은 출금 유형을 정의합니다.
const (
	// WithdrawTransactionTypeDefault는 일반출금입니다.
	WithdrawTransactionTypeDefault = "default"
	// WithdrawTransactionTypeInternal은 바로출금입니다.
	WithdrawTransactionTypeInternal = "internal"
)

// WithdrawInfo는 출금 정보를 나타내는 구조체입니다.
type WithdrawInfo struct {
	// Type은 입출금 종류입니다.
	Type string `json:"type"`
	// UUID는 출금의 고유 아이디입니다.
	UUID string `json:"uuid"`
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// NetType은 출금 네트워크입니다.
	NetType string `json:"net_type"`
	// TxID는 출금의 트랜잭션 아이디입니다.
	TxID string `json:"txid"`
	// State는 출금 상태입니다.
	State string `json:"state"`
	// CreatedAt은 출금 생성 시간입니다.
	CreatedAt time.Time `json:"created_at"`
	// DoneAt은 출금 완료 시간입니다.
	DoneAt time.Time `json:"done_at"`
	// Amount는 출금 금액/수량입니다.
	Amount string `json:"amount"`
	// Fee는 출금 수수료입니다.
	Fee string `json:"fee"`
	// TransactionType은 출금 유형입니다.
	TransactionType string `json:"transaction_type"`
}
