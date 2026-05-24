// Package opendart 는 DART 전자공시 OpenAPI 의 Go 클라이언트다.
package opendart

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/kenshin579/opendart/disclosure"
	"github.com/kenshin579/opendart/internal/corpcode"
	"github.com/kenshin579/opendart/internal/httpclient"
	"github.com/kenshin579/opendart/ownership"
	"github.com/kenshin579/opendart/report"
)

const defaultCorpCacheTTL = 24 * time.Hour

// Client 는 opendart 라이브러리의 단일 진입점.
type Client struct {
	http *httpclient.Client
	corp *corpcode.Cache

	Disclosure *disclosure.Client // DS001 공시정보
	Report     *report.Client     // DS002 정기보고서 주요정보
	Ownership  *ownership.Client  // DS004 지분공시 종합정보
}

// NewClient 는 API 키로 Client 를 만든다.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("opendart: apiKey is required")
	}
	cfg := clientOptions{timeout: 30 * time.Second, corpCacheTTL: defaultCorpCacheTTL}
	for _, opt := range opts {
		opt(&cfg)
	}

	hc := httpclient.New(httpclient.Config{
		APIKey:     apiKey,
		BaseURL:    cfg.baseURL,
		Timeout:    cfg.timeout,
		HTTPClient: cfg.httpClient,
	})

	c := &Client{http: hc}
	c.Disclosure = disclosure.New(hc)
	c.Report = report.New(hc)
	c.Ownership = ownership.New(hc)

	cacheDir := cfg.corpCacheDir
	if cacheDir == "" {
		cacheDir = defaultCorpCacheDir()
	}
	c.corp = corpcode.New(cacheDir, cfg.corpCacheTTL, func(ctx context.Context) ([]byte, error) {
		return hc.GetBytes(ctx, "/api/corpCode.xml", nil)
	})
	return c, nil
}

// defaultCorpCacheDir 는 OS user cache dir/opendart, 실패 시 temp/opendart.
func defaultCorpCacheDir() string {
	if d, err := os.UserCacheDir(); err == nil {
		return filepath.Join(d, "opendart")
	}
	return filepath.Join(os.TempDir(), "opendart")
}
