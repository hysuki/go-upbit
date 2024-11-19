// Package websocket은 Upbit 거래소의 웹소켓 연결을 관리합니다.
package websocket

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/hysuki/go-upbit/websocket/common"
)

// Message는 웹소켓 메시지를 나타냅니다.
type Message struct {
	Ticket         string   `json:"ticket,omitempty"`           // 식별용 티켓
	Type           string   `json:"type,omitempty"`             // 메시지 타입
	Codes          []string `json:"codes,omitempty"`            // 구독할 마켓 코드 목록
	Level          *float64 `json:"level,omitempty"`            // 호가 모아보기 단위
	IsOnlySnapshot *bool    `json:"is_only_snapshot,omitempty"` // 스냅샷 시세만 제공
	IsOnlyRealtime *bool    `json:"is_only_realtime,omitempty"` // 실시간 시세만 제공
}

// StatusResponse는 서버 상태 응답을 나타냅니다.
type StatusResponse struct {
	Status string `json:"status"` // 서버 상태
}

// SubscribeFunc는 구독 함수 타입을 정의합니다.
type SubscribeFunc func(*BaseClient) error

// AddSubscribe는 구독 함수를 생성합니다.
// messageType은 메시지 유형, codes는 마켓 코드 목록, options는 구독 옵션입니다.
func AddSubscribe(messageType string, codes []string, options *common.SubscribeOptions) SubscribeFunc {
	return func(c *BaseClient) error {
		var message Message
		upperCodes := []string{}

		// 코드 대문자로 변환
		for _, code := range codes {
			// - 가 있는지 확인
			if strings.Contains(code, "-") {
				upperCodes = append(upperCodes, strings.ToUpper(code))
			} else {
				return fmt.Errorf("마켓 코드 오류: %s", code)
			}
		}

		if options == nil {
			message = Message{
				Type:  messageType,
				Codes: upperCodes,
			}
		} else {
			message = Message{
				Type:           messageType,
				Codes:          upperCodes,
				Level:          options.Level,
				IsOnlySnapshot: options.IsOnlySnapshot,
				IsOnlyRealtime: options.IsOnlyRealtime,
			}
		}

		c.Messages = append(c.Messages, message)
		return nil
	}
}

// Subscribe는 지정된 구독 함수들을 사용하여 구독을 시작합니다.
// ticket은 구독 식별자, f는 구독 함수 목록입니다.
func (c *BaseClient) Subscribe(ticket *string, f ...SubscribeFunc) error {
	for _, fn := range f {
		if err := fn(c); err != nil {
			return err
		}
	}
	return c.request(ticket)
}

// request는 구독 요청을 전송합니다.
// ticket은 구독 식별자입니다.
func (c *BaseClient) request(ticket *string) error {
	if ticket == nil {
		uuid := uuid.New().String()
		ticket = &uuid
	}

	messages := []Message{
		{Ticket: *ticket},
	}

	if len(c.Messages) > 0 {
		messages = append(messages, c.Messages...)
	}

	return c.WriteJSON(messages)
}

// WriteJSON은 JSON 데이터를 웹소켓으로 전송합니다.
// 전송에 실패하면 에러를 반환합니다.
func (c *BaseClient) WriteJSON(v interface{}) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.Conn == nil {
		return fmt.Errorf("웹소켓 연결이 없습니다")
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("JSON 인코딩 실패: %v", err)
	}

	return c.Conn.Write(c.Ctx, websocket.MessageText, data)
}

// ReadMessage는 웹소켓 메시지를 읽어옵니다.
// 메시지 읽기에 실패하면 에러를 반환합니다.
type ReadMessage struct {
	Type string `json:"type"` // 메시지 타입
	Code string `json:"code"` // 마켓 코드
}

// ReadMessage는 웹소켓으로부터 메시지를 읽어옵니다.
// 연결이 끊어진 경우 재연결을 시도합니다.
func (c *BaseClient) ReadMessage() ([]byte, error) {
	if c.Conn == nil {
		if err := c.Connect(); err != nil {
			return nil, fmt.Errorf("재연결 실패: %v", err)
		}
	}

	_, data, err := c.Conn.Read(c.Ctx)
	if err != nil {
		if websocket.CloseStatus(err) != -1 {
			if err := c.Reconnect(); err != nil {
				return nil, fmt.Errorf("재연결 실패: %v", err)
			}
			// 재연결 후 다시 읽기 시도
			_, data, err = c.Conn.Read(c.Ctx)
			if err != nil {
				return nil, fmt.Errorf("메시지 읽기 실패: %v", err)
			}
		} else {
			return nil, fmt.Errorf("메시지 읽기 실패: %v", err)
		}
	}

	// 서버 상태 응답 확인
	var status StatusResponse
	if err := json.Unmarshal(data, &status); err == nil && status.Status == "UP" {
		return nil, nil
	}

	return data, nil
}
