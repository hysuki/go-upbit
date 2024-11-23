package private

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hysuki/go-upbit/websocket/common"
)

// Asset은 개별 자산 정보를 나타냅니다.
type Asset struct {
	Currency string  `json:"currency,omitempty"` // 화폐를 의미하는 영문 대문자 코드
	Balance  float64 `json:"balance,omitempty"`  // 주문가능 수량
	Locked   float64 `json:"locked,omitempty"`   // 주문 중 묶여있는 수량
}

// UpbitMyAsset은 내 자산 정보를 나타냅니다.
type UpbitMyAsset struct {
	Type           string            `json:"type,omitempty"`            // 타입 (myAsset)
	AssetUUID      string            `json:"asset_uuid,omitempty"`      // 자산 고유 ID
	Assets         []Asset           `json:"assets,omitempty"`          // 자산 목록
	AssetTimestamp int64             `json:"asset_timestamp,omitempty"` // 자산 타임스탬프
	Timestamp      int64             `json:"timestamp,omitempty"`       // 타임스탬프
	StreamType     common.StreamType `json:"stream_type,omitempty"`     // 스트림 타입
}

// MyAsset은 내부적으로 사용하기 위한 자산 정보 구조체입니다.
type MyAsset struct {
	Type       string            `json:"type,omitempty"`       // 타입 (myAsset)
	AssetUUID  string            `json:"asset_uuid,omitempty"` // 자산 고유 ID
	Assets     []Asset           `json:"assets,omitempty"`     // 자산 목록
	AssetAt    time.Time         // 자산 시각 (asset_timestamp를 KST로 변환)
	Timestamp  time.Time         // 타임스탬프 (KST)
	StreamType common.StreamType `json:"stream_type,omitempty"` // 스트림 타입
}

// NewMyAsset은 UpbitMyAsset을 내부 MyAsset 구조체로 변환합니다.
// UTC 시간을 KST(UTC+9)로 변환하여 저장합니다.
func NewMyAsset(u *UpbitMyAsset) *MyAsset {
	kst := time.FixedZone("KST", 9*60*60) // UTC+9
	return &MyAsset{
		Type:       u.Type,
		AssetUUID:  u.AssetUUID,
		Assets:     u.Assets, // Asset 구조체는 공유
		AssetAt:    time.UnixMilli(u.AssetTimestamp).In(kst),
		Timestamp:  time.UnixMilli(u.Timestamp).In(kst),
		StreamType: u.StreamType,
	}
}

// ParseMyAsset은 JSON 데이터를 UpbitMyAsset 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseMyAsset(data []byte) (*UpbitMyAsset, error) {
	var myAsset UpbitMyAsset
	if err := json.Unmarshal(data, &myAsset); err != nil {
		return nil, fmt.Errorf("내 자산 데이터 파싱 실패: %v", err)
	}
	return &myAsset, nil
}

// GetMyAsset은 다음 자산 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 자산 정보를 반환합니다.
func (c *Client) GetMyAsset() (*MyAsset, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.myAssetChan:
		return NewMyAsset(resp), nil
	}
}
