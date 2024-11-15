package upbit

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// 수신된 데이터를 예쁘게 출력하는 헬퍼 함수
func PrettyPrint(prefix string, data []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "    "); err == nil {
		fmt.Printf("%s:\n%s\n", prefix, prettyJSON.String())
	} else {
		fmt.Printf("%s: %s\n", prefix, string(data))
	}
}
