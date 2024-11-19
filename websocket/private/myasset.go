package private

import (
	"encoding/json"
	"fmt"

	"github.com/hysuki/go-upbit/websocket/common"
)

// Asset은 개별 자산 정보를 나타냅니다.
type Asset struct {
	Currency string  `json:"currency,omitempty"` // 화폐를 의미하는 영문 대문자 코드
	Balance  float64 `json:"balance,omitempty"`  // 주문가능 수량
	Locked   float64 `json:"locked,omitempty"`   // 주문 중 묶여있는 수량
}

// MyAssetResponse는 내 자산 정보를 나타냅니다.
type MyAssetResponse struct {
	Type           string            `json:"type,omitempty"`            // 타입 (myAsset)
	AssetUUID      string            `json:"asset_uuid,omitempty"`      // 자산 고유 ID
	Assets         []Asset           `json:"assets,omitempty"`          // 자산 목록
	AssetTimestamp int64             `json:"asset_timestamp,omitempty"` // 자산 타임스탬프
	Timestamp      int64             `json:"timestamp,omitempty"`       // 타임스탬프
	StreamType     common.StreamType `json:"stream_type,omitempty"`     // 스트림 타입
}

// ParseMyAsset은 JSON 데이터를 MyAssetResponse 구조체로 파싱합니다.
// 파싱에 실패하면 에러를 반환합니다.
func ParseMyAsset(data []byte) (*MyAssetResponse, error) {
	var myAsset MyAssetResponse
	if err := json.Unmarshal(data, &myAsset); err != nil {
		return nil, fmt.Errorf("내 자산 데이터 파싱 실패: %v", err)
	}
	return &myAsset, nil
}

// GetMyAsset은 다음 자산 메시지를 기다립니다.
// 에러가 발생하면 에러를 반환하고, 성공하면 자산 정보를 반환합니다.
func (c *Client) GetMyAsset() (*MyAssetResponse, error) {
	select {
	case err := <-c.errChan:
		return nil, err
	case resp := <-c.myAssetChan:
		return resp, nil
	}
}
