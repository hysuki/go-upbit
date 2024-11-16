package private

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/go-upbit/websocket/common"
)

// Asset은 개별 자산 정보를 담는 구조체입니다
type Asset struct {
	Currency string  `json:"currency"` // 화폐를 의미하는 영문 대문자 코드
	Balance  float64 `json:"balance"`  // 주문가능 수량
	Locked   float64 `json:"locked"`   // 주문 중 묶여있는 수량
}

// MyAsset은 내 자산 정보를 담는 구조체입니다
type MyAsset struct {
	Type           string            `json:"type"`            // 타입 (myAsset)
	AssetUUID      string            `json:"asset_uuid"`      // 자산 고유 아이디
	Assets         []Asset           `json:"assets"`          // 자산 리스트
	AssetTimestamp int64             `json:"asset_timestamp"` // 자산 타임스탬프
	Timestamp      int64             `json:"timestamp"`       // 타임스탬프
	StreamType     common.StreamType `json:"stream_type"`     // 스트림 타입
}

// ParseMyAsset은 JSON 데이터를 MyAsset 구조체로 파싱합니다
func ParseMyAsset(data []byte) (*MyAsset, error) {
	var myAsset MyAsset
	if err := json.Unmarshal(data, &myAsset); err != nil {
		return nil, fmt.Errorf("내 자산 데이터 파싱 실패: %v", err)
	}
	return &myAsset, nil
}

// SubscribeMyAsset은 내 자산 정보를 구독합니다
func (c *Client) SubscribeMyAsset() error {
	// MyAsset은 codes 필드를 사용하지 않으므로 nil을 전달합니다
	return c.Subscribe("", "myAsset", nil, nil)
}

// GetMyAsset은 수신된 메시지를 MyAsset 구조체로 변환합니다
func (c *Client) GetMyAsset() (*MyAsset, error) {
	data, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // 서버 상태 메시지인 경우
	}

	return ParseMyAsset(data)
}
