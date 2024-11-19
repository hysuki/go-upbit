package upbit

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Package upbit는 Upbit 암호화폐 거래소의 API를 사용하기 위한 Go 클라이언트 라이브러리를 제공합니다.

// PrettyPrint는 JSON 데이터를 들여쓰기가 적용된 가독성 좋은 형태로 출력합니다.
//
// prefix 파라미터는 출력될 JSON 데이터 앞에 표시될 문자열입니다.
// data 파라미터는 출력하고자 하는 JSON 바이트 데이터입니다.
//
// JSON 데이터 파싱에 실패할 경우, 원본 데이터를 그대로 출력합니다.
func PrettyPrint(prefix string, data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err == nil {
		fmt.Printf("%s:\n%s\n", prefix, prettyJSON.String())
	} else {
		fmt.Printf("%s: %s\n", prefix, string(data))
	}
}
