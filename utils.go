package upbit

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Package upbit는 Upbit 암호화폐 거래소 API를 위한 Go 클라이언트 라이브러리입니다.

// PrettyPrint는 수신된 JSON 데이터를 보기 좋게 출력하는 헬퍼 함수입니다.
// prefix는 출력할 데이터의 접두사이고, data는 출력할 JSON 데이터입니다.
func PrettyPrint(prefix string, data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err == nil {
		fmt.Printf("%s:\n%s\n", prefix, prettyJSON.String())
	} else {
		fmt.Printf("%s: %s\n", prefix, string(data))
	}
}
