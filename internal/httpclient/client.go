// Package httpclient 는 OpenDART REST 호출의 단일 GET 통로다.
// crtfc_key 주입, status envelope 검사, TLS RSA cipher 처리를 담당한다.
package httpclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultBaseURL 은 OpenDART API 베이스 URL.
const DefaultBaseURL = "https://opendart.fss.or.kr"

// Envelope 는 모든 OpenDART JSON 응답 공통 필드. 응답 타입이 임베드한다.
type Envelope struct {
	Status  string `json:"status"`  // 000=정상
	Message string `json:"message"` // 상태 메시지
}

// APIStatus 는 임베드한 타입이 StatusReader 를 충족하게 한다.
func (e Envelope) APIStatus() (status, message string) { return e.Status, e.Message }

// StatusReader 는 응답에서 status/message 를 읽는다 (Envelope 임베드로 자동 구현).
type StatusReader interface {
	APIStatus() (status, message string)
}

// APIError 는 status != "000" 인 OpenDART 응답을 나타낸다.
type APIError struct {
	Status  string // OpenDART status 코드
	Message string // OpenDART message
}

func (e *APIError) Error() string {
	return fmt.Sprintf("opendart: [%s] %s", e.Status, e.Message)
}

// ErrNoData 는 status 013 (조회된 데이터 없음).
var ErrNoData = errors.New("opendart: no data (013)")

// Config 는 Client 생성 인자.
type Config struct {
	APIKey     string
	BaseURL    string        // 빈 값이면 DefaultBaseURL
	Timeout    time.Duration // 0이면 30s
	HTTPClient *http.Client  // nil이면 TLS RSA cipher 내장 클라이언트
}

// Client 는 OpenDART HTTP 계층.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// New 는 Config 로 Client 를 만든다.
func New(cfg Config) *Client {
	base := cfg.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	hc := cfg.HTTPClient
	if hc == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		hc = newHTTPClient(timeout)
	}
	return &Client{apiKey: cfg.APIKey, baseURL: base, http: hc}
}

// newHTTPClient 는 OpenDART 서버(TLS1.2 RSA 키교환 cipher 전용)에 맞춘 클라이언트.
// Go 기본 ClientHello 는 forward secrecy 없는 RSA cipher 를 빼서 handshake 가 실패한다.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
		},
	}
}

// get 은 crtfc_key 를 주입해 GET 후 본문 바이트를 반환한다.
func (c *Client) get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("crtfc_key", c.apiKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opendart: GET %s: http status %d", path, resp.StatusCode)
	}
	return body, nil
}

// GetJSON 은 응답을 out 으로 디코드하고 envelope status 를 검사한다.
func (c *Client) GetJSON(ctx context.Context, path string, params map[string]string, out StatusReader) error {
	body, err := c.get(ctx, path, params)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("opendart: decode %s: %w", path, err)
	}
	status, message := out.APIStatus()
	return statusError(status, message)
}

// GetBytes 는 바이너리(ZIP 등) 응답을 그대로 반환한다. 에러 시 JSON envelope 가 올 수
// 있으므로 본문이 '{' 로 시작하면 status 를 검사한다 (ZIP 매직은 'PK').
func (c *Client) GetBytes(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	body, err := c.get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 && body[0] == '{' {
		var env Envelope
		if json.Unmarshal(body, &env) == nil && env.Status != "" {
			if serr := statusError(env.Status, env.Message); serr != nil {
				return nil, serr
			}
		}
	}
	return body, nil
}

// statusError 는 status 코드를 에러로 변환한다. 000/""→nil, 013→ErrNoData, 그 외→*APIError.
func statusError(status, message string) error {
	switch status {
	case "000", "":
		return nil
	case "013":
		return ErrNoData
	default:
		return &APIError{Status: status, Message: message}
	}
}
