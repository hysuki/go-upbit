package exchange

// OrderState 상수들은 주문의 상태를 정의합니다.
const (
	// OrderStateWait는 체결 대기 상태입니다.
	OrderStateWait = "wait"
	// OrderStateWatch는 예약주문 대기 상태입니다.
	OrderStateWatch = "watch"
	// OrderStateDone은 전체 체결 완료 상태입니다.
	OrderStateDone = "done"
	// OrderStateCancel은 주문 취소 상태입니다.
	OrderStateCancel = "cancel"
)

// OrderBy 상수들은 정렬 방식을 정의합니다.
const (
	// OrderByAsc는 오름차순 정렬입니다.
	OrderByAsc = "asc"
	// OrderByDesc는 내림차순 정렬입니다.
	OrderByDesc = "desc"
)
