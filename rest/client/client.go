// Package client는 Upbit REST API 클라이언트의 기본 인터페이스를 정의하는 패키지입니다.
// 이 패키지는 순환 참조를 방지하기 위해 별도로 분리되었습니다.
package client

// RestClient는 Upbit REST API와의 기본적인 HTTP 통신을 위한 인터페이스입니다.
type RestClient interface {
	// Get은 지정된 경로로 GET 요청을 보내고 응답을 바이트 슬라이스로 반환합니다.
	// path는 요청할 API 엔드포인트 경로이며, params는 쿼리 파라미터입니다.
	Get(path string, params map[string]string) ([]byte, error)

	// Post는 지정된 경로로 POST 요청을 보내고 응답을 바이트 슬라이스로 반환합니다.
	// path는 요청할 API 엔드포인트 경로이며, body는 요청 본문입니다.
	Post(path string, body interface{}) ([]byte, error)

	// Delete는 지정된 경로로 DELETE 요청을 보내고 응답을 바이트 슬라이스로 반환합니다.
	// path는 요청할 API 엔드포인트 경로이며, params는 쿼리 파라미터입니다.
	Delete(path string, params map[string]string) ([]byte, error)
}
