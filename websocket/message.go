package websocket

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/coder/websocket"
// 	"github.com/google/uuid"
// )

// // Message는 웹소켓 메시지 구조체입니다
// type Message struct {
// 	Ticket         string   `json:"ticket,omitempty"`           // 식별용 티켓
// 	Type           string   `json:"type,omitempty"`             // 메시지 타입
// 	Codes          []string `json:"codes,omitempty"`            // 구독할 마켓 코드 목록
// 	Level          *float64 `json:"level,omitempty"`            // 모아보기 단위
// 	IsOnlySnapshot *bool    `json:"is_only_snapshot,omitempty"` // 스냅샷 시세만 제공
// 	IsOnlyRealtime *bool    `json:"is_only_realtime,omitempty"` // 실시간 시세만 제공
// }

// // StatusResponse는 서버 상태 응답을 위한 구조체입니다
// type StatusResponse struct {
// 	Status string `json:"status"`
// }

// type SubscribeOptions struct {
// 	Level          *float64
// 	IsOnlySnapshot *bool
// 	IsOnlyRealtime *bool
// }

// // Subscribe는 특정 마켓 데이터를 구독합니다
// func (c *Client) Subscribe(ticket string, messageType string, codes []string, options *SubscribeOptions) error {
// 	if ticket == "" {
// 		ticket = uuid.New().String()
// 	}

// 	messages := []Message{
// 		{Ticket: ticket},
// 	}

// 	if options == nil {
// 		messages = append(messages, Message{Type: messageType, Codes: codes})
// 	} else {
// 		messages = append(messages, Message{
// 			Type:           messageType,
// 			Codes:          codes,
// 			Level:          options.Level,
// 			IsOnlySnapshot: options.IsOnlySnapshot,
// 			IsOnlyRealtime: options.IsOnlyRealtime,
// 		})
// 	}

// 	return c.WriteJSON(messages)
// }

// // WriteJSON은 웹소켓 연결을 통해 JSON 데이터를 전송합니다
// func (c *client) WriteJSON(v interface{}) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	if c.conn == nil {
// 		return fmt.Errorf("웹소켓 연결이 없습니다")
// 	}

// 	data, err := json.Marshal(v)
// 	if err != nil {
// 		return fmt.Errorf("JSON 인코딩 실패: %v", err)
// 	}

// 	return c.conn.Write(c.ctx, websocket.MessageText, data)
// }

// // ReadMessage는 웹소켓 연결로부터 메시지를 읽어옵니다
// // 서버 상태 메시지("UP")는 nil을 반환합니다
// func (c *client) ReadMessage() ([]byte, error) {
// 	c.mu.Lock()
// 	if c.conn == nil {
// 		c.mu.Unlock()
// 		// 연결이 없으면 재연결 시도
// 		if err := c.Connect(); err != nil {
// 			return nil, fmt.Errorf("재연결 실패: %v", err)
// 		}
// 		c.mu.Lock()
// 	}
// 	c.mu.Unlock()

// 	_, data, err := c.conn.Read(c.ctx)
// 	if err != nil {
// 		// 연결이 끊어진 경우 재연결 시도
// 		if websocket.CloseStatus(err) != -1 {
// 			if err := c.Connect(); err != nil {
// 				return nil, fmt.Errorf("재연결 실패: %v", err)
// 			}
// 			// 재연결 후 다시 읽기 시도
// 			_, data, err = c.conn.Read(c.ctx)
// 			if err != nil {
// 				return nil, fmt.Errorf("메시지 읽기 실패: %v", err)
// 			}
// 		} else {
// 			return nil, fmt.Errorf("메시지 읽기 실패: %v", err)
// 		}
// 	}

// 	// 서버 상태 응답 확인
// 	var status StatusResponse
// 	if err := json.Unmarshal(data, &status); err == nil && status.Status == "UP" {
// 		return nil, nil // 상태 메시지는 무시
// 	}

// 	return data, nil
// }
