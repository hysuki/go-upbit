package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/hysuki/go-upbit/websocket/common"
)

// Message는 웹소켓 메시지 구조체입니다
type Message struct {
	Ticket         string   `json:"ticket,omitempty"`           // 식별용 티켓
	Type           string   `json:"type,omitempty"`             // 메시지 타입
	Codes          []string `json:"codes,omitempty"`            // 구독할 마켓 코드 목록
	Level          *float64 `json:"level,omitempty"`            // 모아보기 단위
	IsOnlySnapshot *bool    `json:"is_only_snapshot,omitempty"` // 스냅샷 시세만 제공
	IsOnlyRealtime *bool    `json:"is_only_realtime,omitempty"` // 실시간 시세만 제공
}

// StatusResponse는 서버 상태 응답을 위한 구조체입니다
type StatusResponse struct {
	Status string `json:"status"`
}

type SubscribeFunc func(*BaseClient) error

func AddSubscribe(messageType string, codes []string, options *common.SubscribeOptions) SubscribeFunc {
	return func(c *BaseClient) error {
		var message Message
		if options == nil {
			message = Message{
				Type:  messageType,
				Codes: codes,
			}
		} else {
			message = Message{
				Type:           messageType,
				Codes:          codes,
				Level:          options.Level,
				IsOnlySnapshot: options.IsOnlySnapshot,
				IsOnlyRealtime: options.IsOnlyRealtime,
			}
		}

		c.Messages = append(c.Messages, message)
		return nil
	}
}

func (c *BaseClient) Subscribe(ticket *string, f ...SubscribeFunc) error {
	for _, fn := range f {
		if err := fn(c); err != nil {
			return err
		}
	}
	return c.request(ticket)
}

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

// WriteJSON에 mutex 추가
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

// ReadMessage도 mutex로 보호
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
