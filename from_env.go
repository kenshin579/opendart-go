package opendart

import (
	"errors"
	"os"
)

// NewClientFromEnv 는 OPENDART_API_KEY 환경변수로 Client 를 만든다.
func NewClientFromEnv(opts ...Option) (*Client, error) {
	key := os.Getenv("OPENDART_API_KEY")
	if key == "" {
		return nil, errors.New("opendart: OPENDART_API_KEY is not set")
	}
	return NewClient(key, opts...)
}
