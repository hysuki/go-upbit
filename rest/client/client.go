// 순환참조 방지를 위해 별도의 패키지로 분리
package client

// RestClient는 REST API 클라이언트 인터페이스를 정의합니다.
type RestClient interface {
	Get(path string, params map[string]string) ([]byte, error)
	Post(path string, body interface{}) ([]byte, error)
	Delete(path string, params map[string]string) ([]byte, error)
}
