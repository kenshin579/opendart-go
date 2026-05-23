package opendart

import (
	"net/http"
	"time"
)

type clientOptions struct {
	baseURL      string
	timeout      time.Duration
	httpClient   *http.Client
	corpCacheDir string
	corpCacheTTL time.Duration
}

// Option 은 functional option.
type Option func(*clientOptions)

// WithBaseURL 은 API 베이스 URL 을 지정한다 (테스트/프록시용).
func WithBaseURL(u string) Option { return func(o *clientOptions) { o.baseURL = u } }

// WithTimeout 은 HTTP 타임아웃을 지정한다 (기본 30s).
func WithTimeout(d time.Duration) Option { return func(o *clientOptions) { o.timeout = d } }

// WithHTTPClient 는 사용자 정의 *http.Client 를 주입한다 (기본: TLS RSA cipher 내장).
func WithHTTPClient(c *http.Client) Option { return func(o *clientOptions) { o.httpClient = c } }

// WithCorpCodeCacheDir 는 corp_code 캐시 디렉토리를 지정한다.
func WithCorpCodeCacheDir(dir string) Option { return func(o *clientOptions) { o.corpCacheDir = dir } }

// WithCorpCodeCacheTTL 은 corp_code 재다운로드 주기를 지정한다 (기본 24h).
func WithCorpCodeCacheTTL(d time.Duration) Option { return func(o *clientOptions) { o.corpCacheTTL = d } }
