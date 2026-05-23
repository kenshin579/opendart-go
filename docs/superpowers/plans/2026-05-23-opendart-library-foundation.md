# OpenDART 라이브러리 Foundation + DS001 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `github.com/kenshin579/opendart` 라이브러리의 기반(Client/인증/HTTP/에러/corp_code 인프라)과 DS001 공시정보 카테고리(기업개황·공시검색·원본파일)를 구현한다.

**Architecture:** KIS 구조를 단순화. `internal/httpclient` 가 `crtfc_key` 주입 + status envelope 검사 + TLS RSA cipher 를 담당하는 단일 GET 통로. `internal/corpcode` 는 corpCode.xml ZIP 을 디스크 캐시(TTL+stale fallback)하고 종목코드↔corp_code 인덱스를 제공. root `opendart` 패키지가 이들을 와이어링하고 `disclosure` sub-client 를 노출. internal→root import 순환은 타입 별칭(`type APIError = httpclient.APIError`)으로 회피.

**Tech Stack:** Go 1.25+, 표준 `net/http`/`encoding/json`/`encoding/xml`/`archive/zip`, `github.com/stretchr/testify` (test). resty 미사용.

**Spec:** `docs/superpowers/specs/2026-05-23-opendart-library-foundation-design.md`

**검증된 사실 (실 API 호출로 확인):**
- `company.json` 응답: `status`/`message` + `corp_code` 포함, 전 필드 string.
- `list.json` 응답: `page_no`/`page_count`/`total_count`/`total_page` 는 JSON **int**, `list[]` 는 객체 배열.
- `corpCode.xml` ZIP: 단일 `CORPCODE.xml`, 루트 `<result>` 아래 `<list>` 반복, 비상장사 `stock_code` 는 공백(" ").
- 에러: `{"status":"013","message":"조회된 데이타가 없습니다."}` (데이터 없음), `{"status":"010","message":"등록되지 않은 인증키입니다."}` (키 오류).
- 서버는 TLS1.2 RSA 키교환 cipher 전용 → 기본 클라이언트에 cipher 명시 필요.

---

## File Structure

```
go.mod                         # 기존 (module github.com/kenshin579/opendart, go 1.25)
internal/httpclient/
  client.go                    # Envelope, StatusReader, APIError, ErrNoData, Config, New, GetJSON, GetBytes
  client_test.go
internal/corpcode/
  cache.go                     # Entry, FetchFunc, Cache, parseZip, 인덱스/TTL/fallback
  cache_test.go
disclosure/
  client.go                    # disclosure.Client + New
  company.go                   # Company + GetCompany
  company_test.go
  search.go                    # SearchParams + SearchResult + DisclosureItem + SearchDisclosures
  search_test.go
  document.go                  # DownloadDocument
  document_test.go
  testdata/company_samsung.json, list_samsung.json
errors.go                      # APIError/ErrNoData 별칭 + ErrCorpCodeNotFound
config.go                      # clientOptions + Option + With*
client.go                      # Client + NewClient + 와이어링
client_test.go
from_env.go                    # NewClientFromEnv
corpcode.go                    # CorpCodeEntry 별칭 + ResolveCorpCode/LookupCorpCode/CorpCodes/RefreshCorpCodes
corpcode_test.go
examples/disclosure/main.go
README.md
integration_test.go            # //go:build integration
```

---

### Task 1: internal/httpclient (HTTP/envelope 계층)

**Files:**
- Create: `internal/httpclient/client.go`
- Test: `internal/httpclient/client_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `internal/httpclient/client_test.go`:
```go
package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResp struct {
	Envelope
	CorpName string `json:"corp_name"`
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
}

func TestGetJSON_SuccessAndKeyInjection(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "KEY", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
		w.Write([]byte(`{"status":"000","message":"정상","corp_name":"삼성전자(주)"}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", map[string]string{"corp_code": "00126380"}, &out)
	require.NoError(t, err)
	assert.Equal(t, "삼성전자(주)", out.CorpName)
}

func TestGetJSON_NoData(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", nil, &out)
	assert.ErrorIs(t, err, ErrNoData)
}

func TestGetJSON_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", nil, &out)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "010", apiErr.Status)
}

func TestGetBytes_Binary(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PK\x03\x04zipdata"))
	})
	b, err := c.GetBytes(context.Background(), "/api/corpCode.xml", nil)
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04zipdata", string(b))
}

func TestGetBytes_ErrorJSON(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	})
	_, err := c.GetBytes(context.Background(), "/api/document.xml", nil)
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
}

func TestGetJSON_HTTPError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/x", nil, &out)
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNoData)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./internal/httpclient/ -v`
Expected: FAIL — `undefined: Client`, `New`, `Config`, `Envelope`, `APIError`, `ErrNoData`.

- [ ] **Step 3: 구현** — `internal/httpclient/client.go`:
```go
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
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./internal/httpclient/ -v`
Expected: PASS (6 tests). Run `go vet ./internal/httpclient/` — clean.

- [ ] **Step 5: Commit**

```bash
git add internal/httpclient/
git commit -m "feat(httpclient): GET layer with crtfc_key + status envelope + TLS cipher"
```

---

### Task 2: internal/corpcode (corp_code 디스크 캐시 + 인덱스)

**Files:**
- Create: `internal/corpcode/cache.go`
- Test: `internal/corpcode/cache_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `internal/corpcode/cache_test.go`:
```go
package corpcode

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<result>
	<list><corp_code>00126380</corp_code><corp_name>삼성전자(주)</corp_name><corp_eng_name>SAMSUNG</corp_eng_name><stock_code>005930</stock_code><modify_date>20240101</modify_date></list>
	<list><corp_code>00434003</corp_code><corp_name>다코</corp_name><corp_eng_name>Daco</corp_eng_name><stock_code> </stock_code><modify_date>20170630</modify_date></list>
</result>`

func makeZip(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	require.NoError(t, err)
	_, err = w.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func TestCache_IndexAndLookup(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	c := New(t.TempDir(), time.Hour, func(ctx context.Context) ([]byte, error) { return zipBytes, nil })

	// 종목코드 → corp_code (상장사)
	e, ok, err := c.ByStockCode(context.Background(), "005930")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "00126380", e.CorpCode)

	// 비상장사(stock_code 공백)는 stockCode 인덱스에서 제외
	_, ok, err = c.ByStockCode(context.Background(), "")
	require.NoError(t, err)
	assert.False(t, ok)

	// corp_code → 엔트리 (비상장 포함)
	e2, ok, err := c.ByCorpCode(context.Background(), "00434003")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "다코", e2.CorpName)
	assert.Equal(t, "", e2.StockCode) // 공백이 trim 됨

	// 전체 엔트리
	all, err := c.Entries(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestCache_DiskCacheAvoidsRefetch(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	dir := t.TempDir()
	calls := 0
	fetch := func(ctx context.Context) ([]byte, error) { calls++; return zipBytes, nil }

	c1 := New(dir, time.Hour, fetch)
	_, _, _ = c1.ByStockCode(context.Background(), "005930") // fetch 1회 + 디스크 저장
	require.Equal(t, 1, calls)

	// 새 Cache 인스턴스(같은 dir, 신선한 캐시) → fetch 호출 안 함
	c2 := New(dir, time.Hour, fetch)
	_, _, _ = c2.ByStockCode(context.Background(), "005930")
	assert.Equal(t, 1, calls)
}

func TestCache_StaleFallbackOnFetchError(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	dir := t.TempDir()
	// 디스크에 오래된(만료된) 캐시 파일을 직접 둔다
	require.NoError(t, os.WriteFile(filepath.Join(dir, cacheFileName), zipBytes, 0o644))
	old := time.Now().Add(-48 * time.Hour)
	require.NoError(t, os.Chtimes(filepath.Join(dir, cacheFileName), old, old))

	c := New(dir, time.Hour, func(ctx context.Context) ([]byte, error) {
		return nil, errors.New("network down")
	})
	e, ok, err := c.ByCorpCode(context.Background(), "00126380")
	require.NoError(t, err) // stale fallback
	require.True(t, ok)
	assert.Equal(t, "삼성전자(주)", e.CorpName)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./internal/corpcode/ -v`
Expected: FAIL — `undefined: New`, `cacheFileName`, etc.

- [ ] **Step 3: 구현** — `internal/corpcode/cache.go`:
```go
// Package corpcode 는 OpenDART 고유번호(corpCode.xml) 의 디스크 캐시 + 인덱스다.
// 다운로드 비용이 큰 전체 회사 매핑 ZIP 을 TTL 단위로 재사용하고, 다운로드 실패 시
// 오래된 캐시로 fallback 한다 (KIS mastercache 모델).
package corpcode

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Entry 는 CORPCODE.xml 의 한 회사.
type Entry struct {
	CorpCode    string `xml:"corp_code"`     // 고유번호 (8자리)
	CorpName    string `xml:"corp_name"`     // 정식회사명칭
	CorpEngName string `xml:"corp_eng_name"` // 영문 정식명칭
	StockCode   string `xml:"stock_code"`    // 종목코드 (6자리, 비상장은 빈 값)
	ModifyDate  string `xml:"modify_date"`   // 최종변경일자 YYYYMMDD
}

// FetchFunc 는 corpCode.xml ZIP 바이트를 가져온다 (캐시 미스/갱신 시 호출).
type FetchFunc func(ctx context.Context) ([]byte, error)

const cacheFileName = "corpcode.zip"

// Cache 는 corp_code 매핑의 디스크 캐시 + 인메모리 인덱스.
type Cache struct {
	dir   string
	ttl   time.Duration
	fetch FetchFunc

	mu      sync.Mutex
	loaded  bool
	entries []Entry
	byStock map[string]Entry
	byCorp  map[string]Entry
}

// New 는 Cache 를 만든다. dir 는 캐시 디렉토리, ttl 은 재다운로드 주기.
func New(dir string, ttl time.Duration, fetch FetchFunc) *Cache {
	return &Cache{
		dir: dir, ttl: ttl, fetch: fetch,
		byStock: map[string]Entry{}, byCorp: map[string]Entry{},
	}
}

// ensure 는 인덱스가 준비되도록 한다 (lazy 1회). force=true 면 강제 재다운로드.
func (c *Cache) ensure(ctx context.Context, force bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loaded && !force {
		return nil
	}
	zipBytes, err := c.obtain(ctx, force)
	if err != nil {
		return err
	}
	entries, err := parseZip(zipBytes)
	if err != nil {
		return err
	}
	c.index(entries)
	c.loaded = true
	return nil
}

// obtain 은 신선한 캐시가 있으면 디스크에서, 아니면 fetch 후 저장한다.
// fetch 실패 시 오래된 캐시라도 있으면 fallback.
func (c *Cache) obtain(ctx context.Context, force bool) ([]byte, error) {
	path := filepath.Join(c.dir, cacheFileName)
	if !force {
		if b, ok := readFresh(path, c.ttl); ok {
			return b, nil
		}
	}
	b, err := c.fetch(ctx)
	if err != nil {
		if old, rerr := os.ReadFile(path); rerr == nil {
			return old, nil // stale fallback
		}
		return nil, err
	}
	if mkErr := os.MkdirAll(c.dir, 0o755); mkErr == nil {
		_ = os.WriteFile(path, b, 0o644)
	}
	return b, nil
}

// readFresh 는 TTL 내 캐시 파일을 읽는다.
func readFresh(path string, ttl time.Duration) ([]byte, bool) {
	fi, err := os.Stat(path)
	if err != nil || time.Since(fi.ModTime()) > ttl {
		return nil, false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return b, true
}

// index 는 엔트리로 두 맵을 구축한다. 빈 stock_code(비상장)는 stockCode 맵에서 제외.
func (c *Cache) index(entries []Entry) {
	c.entries = entries
	c.byStock = make(map[string]Entry, len(entries))
	c.byCorp = make(map[string]Entry, len(entries))
	for _, e := range entries {
		c.byCorp[e.CorpCode] = e
		if e.StockCode != "" {
			c.byStock[e.StockCode] = e
		}
	}
}

// parseZip 은 ZIP 안 첫 .xml 을 찾아 엔트리로 파싱하고 stock_code 를 trim 한다.
func parseZip(zipBytes []byte) ([]Entry, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("corpcode: open zip: %w", err)
	}
	for _, f := range zr.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		var doc struct {
			List []Entry `xml:"list"`
		}
		if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
			return nil, fmt.Errorf("corpcode: parse xml: %w", err)
		}
		for i := range doc.List {
			doc.List[i].StockCode = strings.TrimSpace(doc.List[i].StockCode)
		}
		return doc.List, nil
	}
	return nil, fmt.Errorf("corpcode: no xml in zip")
}

// Entries 는 전체 엔트리를 반환한다.
func (c *Cache) Entries(ctx context.Context) ([]Entry, error) {
	if err := c.ensure(ctx, false); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entries, nil
}

// ByStockCode 는 종목코드로 엔트리를 찾는다.
func (c *Cache) ByStockCode(ctx context.Context, stockCode string) (Entry, bool, error) {
	if err := c.ensure(ctx, false); err != nil {
		return Entry{}, false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.byStock[strings.TrimSpace(stockCode)]
	return e, ok, nil
}

// ByCorpCode 는 고유번호로 엔트리를 찾는다.
func (c *Cache) ByCorpCode(ctx context.Context, corpCode string) (Entry, bool, error) {
	if err := c.ensure(ctx, false); err != nil {
		return Entry{}, false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.byCorp[strings.TrimSpace(corpCode)]
	return e, ok, nil
}

// Refresh 는 TTL 을 무시하고 강제 재다운로드한다.
func (c *Cache) Refresh(ctx context.Context) error {
	return c.ensure(ctx, true)
}
```

> 주: `ByStockCode("")` 는 TrimSpace 후 ""→byStock 에 빈 키 없음→ok=false (테스트 기대와 일치). 빈 stock_code 가 인덱스에 안 들어가므로 안전.

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./internal/corpcode/ -v`
Expected: PASS (3 tests). `go vet ./internal/corpcode/` clean.

- [ ] **Step 5: Commit**

```bash
git add internal/corpcode/
git commit -m "feat(corpcode): corpCode.xml disk cache + stock/corp index"
```

---

### Task 3: disclosure — 기업개황 (GetCompany)

**Files:**
- Create: `disclosure/client.go`, `disclosure/company.go`
- Test: `disclosure/company_test.go`
- Create: `disclosure/testdata/company_samsung.json`

- [ ] **Step 1: fixture 작성** — `disclosure/testdata/company_samsung.json`:
```json
{
    "status": "000",
    "message": "정상",
    "corp_code": "00126380",
    "corp_name": "삼성전자(주)",
    "corp_name_eng": "SAMSUNG ELECTRONICS CO,.LTD",
    "stock_name": "삼성전자",
    "stock_code": "005930",
    "ceo_nm": "전영현, 노태문",
    "corp_cls": "Y",
    "jurir_no": "1301110006246",
    "bizr_no": "1248100998",
    "adres": "경기도 수원시 영통구  삼성로 129 (매탄동)",
    "hm_url": "www.samsung.com/sec",
    "ir_url": "",
    "phn_no": "02-2255-0114",
    "fax_no": "031-200-7538",
    "induty_code": "264",
    "est_dt": "19690113",
    "acc_mt": "12"
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `disclosure/company_test.go`:
```go
package disclosure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 testdata 의 fixture 를 path 별로 서빙하는 disclosure.Client 를 만든다.
func newTestClient(t *testing.T, routes map[string]string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, ok := routes[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		b, err := os.ReadFile(filepath.Join("testdata", fixture))
		require.NoError(t, err)
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	return New(hc)
}

func TestGetCompany(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/company.json": "company_samsung.json"})
	got, err := c.GetCompany(context.Background(), "00126380")
	require.NoError(t, err)
	assert.Equal(t, "삼성전자(주)", got.CorpName)
	assert.Equal(t, "005930", got.StockCode)
	assert.Equal(t, "Y", got.CorpCls)
	assert.Equal(t, "19690113", got.EstDate)
	assert.Equal(t, "12", got.AccMonth)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./disclosure/ -run TestGetCompany -v`
Expected: FAIL — `undefined: Client`, `New`, `GetCompany`.

- [ ] **Step 4: 구현** — `disclosure/client.go`:
```go
// Package disclosure 는 OpenDART DS001 공시정보 API sub-client 다.
// opendart.Client.Disclosure 로 접근한다.
package disclosure

import "github.com/kenshin579/opendart/internal/httpclient"

// Client 는 공시정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }
```

`disclosure/company.go`:
```go
package disclosure

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Company 는 기업개황 (company.json) 응답.
type Company struct {
	httpclient.Envelope
	CorpCode    string `json:"corp_code"`     // 고유번호 (8자리)
	CorpName    string `json:"corp_name"`     // 정식명칭
	CorpNameEng string `json:"corp_name_eng"` // 영문명칭
	StockName   string `json:"stock_name"`    // 종목명(상장사)/약식명칭(기타)
	StockCode   string `json:"stock_code"`    // 종목코드 (6자리)
	CeoName     string `json:"ceo_nm"`        // 대표자명
	CorpCls     string `json:"corp_cls"`      // 법인구분 Y(유가)/K(코스닥)/N(코넥스)/E(기타)
	JurirNo     string `json:"jurir_no"`      // 법인등록번호
	BizrNo      string `json:"bizr_no"`       // 사업자등록번호
	Address     string `json:"adres"`         // 주소
	HomeURL     string `json:"hm_url"`        // 홈페이지
	IRURL       string `json:"ir_url"`        // IR 홈페이지
	PhoneNo     string `json:"phn_no"`        // 전화번호
	FaxNo       string `json:"fax_no"`        // 팩스번호
	IndutyCode  string `json:"induty_code"`   // 업종코드
	EstDate     string `json:"est_dt"`        // 설립일 YYYYMMDD
	AccMonth    string `json:"acc_mt"`        // 결산월 MM
}

// GetCompany 는 corp_code(8자리)로 기업개황을 조회한다.
func (c *Client) GetCompany(ctx context.Context, corpCode string) (*Company, error) {
	var out Company
	if err := c.http.GetJSON(ctx, "/api/company.json", map[string]string{"corp_code": corpCode}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./disclosure/ -run TestGetCompany -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add disclosure/client.go disclosure/company.go disclosure/company_test.go disclosure/testdata/company_samsung.json
git commit -m "feat(disclosure): GetCompany (기업개황)"
```

---

### Task 4: disclosure — 공시검색 (SearchDisclosures)

**Files:**
- Create: `disclosure/search.go`
- Test: `disclosure/search_test.go`
- Create: `disclosure/testdata/list_samsung.json`

- [ ] **Step 1: fixture 작성** — `disclosure/testdata/list_samsung.json`:
```json
{
    "status": "000",
    "message": "정상",
    "page_no": 1,
    "page_count": 2,
    "total_count": 15,
    "total_page": 8,
    "list": [
        {
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "stock_code": "005930",
            "corp_cls": "Y",
            "report_nm": "특수관계인과의내부거래",
            "rcept_no": "20240131000326",
            "flr_nm": "삼성전자",
            "rcept_dt": "20240131",
            "rm": "공"
        },
        {
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "stock_code": "005930",
            "corp_cls": "Y",
            "report_nm": "수시공시의무관련사항(공정공시)",
            "rcept_no": "20240131800110",
            "flr_nm": "삼성전자",
            "rcept_dt": "20240131",
            "rm": "유"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `disclosure/search_test.go` (재사용: `newTestClient` 은 company_test.go 에 정의됨):
```go
package disclosure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchDisclosures(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/list.json": "list_samsung.json"})
	res, err := c.SearchDisclosures(context.Background(), SearchParams{CorpCode: "00126380", PageCount: 2})
	require.NoError(t, err)
	assert.Equal(t, 1, res.PageNo)
	assert.Equal(t, 15, res.TotalCount)
	assert.Equal(t, 8, res.TotalPage)
	require.Len(t, res.List, 2)
	assert.Equal(t, "20240131000326", res.List[0].RceptNo)
	assert.Equal(t, "특수관계인과의내부거래", res.List[0].ReportNm)
}

func TestSearchParams_toMap(t *testing.T) {
	m := SearchParams{CorpCode: "00126380", PblntfTy: "A", PageNo: 2, PageCount: 50}.toMap()
	assert.Equal(t, "00126380", m["corp_code"])
	assert.Equal(t, "A", m["pblntf_ty"])
	assert.Equal(t, "2", m["page_no"])
	assert.Equal(t, "50", m["page_count"])
	// 빈 값/0 은 생략
	_, hasBgn := m["bgn_de"]
	assert.False(t, hasBgn)

	empty := SearchParams{}.toMap()
	_, hasPage := empty["page_no"]
	assert.False(t, hasPage)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./disclosure/ -run 'TestSearchDisclosures|TestSearchParams' -v`
Expected: FAIL — `undefined: SearchParams`, `SearchResult`, `SearchDisclosures`.

- [ ] **Step 4: 구현** — `disclosure/search.go`:
```go
package disclosure

import (
	"context"
	"strconv"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// SearchParams 는 공시검색 (list.json) 요청 인자. 빈 값/0 은 쿼리에서 생략한다.
type SearchParams struct {
	CorpCode       string // 고유번호 (8자리)
	BgnDe          string // 검색 시작 접수일자 YYYYMMDD
	EndDe          string // 검색 종료 접수일자 YYYYMMDD
	LastReprtAt    string // 최종보고서만 검색 여부 (Y/N, 기본 N)
	PblntfTy       string // 공시유형 (A 정기/B 주요사항/C 발행/D 지분/E 기타/F 외부감사/G 펀드/H 자산유동화/I 거래소/J 공정위)
	PblntfDetailTy string // 공시상세유형 (4자리 코드)
	CorpCls        string // 법인구분 (Y/K/N/E)
	Sort           string // 정렬 (date/crp/rpt, 기본 date)
	SortMth        string // 정렬방법 (asc/desc, 기본 desc)
	PageNo         int    // 페이지 번호 (1~n, 0이면 생략)
	PageCount      int    // 페이지당 건수 (1~100, 0이면 생략)
}

// toMap 은 빈 값/0 을 제외한 쿼리 파라미터 맵을 만든다.
func (p SearchParams) toMap() map[string]string {
	m := map[string]string{}
	put := func(k, v string) {
		if v != "" {
			m[k] = v
		}
	}
	put("corp_code", p.CorpCode)
	put("bgn_de", p.BgnDe)
	put("end_de", p.EndDe)
	put("last_reprt_at", p.LastReprtAt)
	put("pblntf_ty", p.PblntfTy)
	put("pblntf_detail_ty", p.PblntfDetailTy)
	put("corp_cls", p.CorpCls)
	put("sort", p.Sort)
	put("sort_mth", p.SortMth)
	if p.PageNo > 0 {
		m["page_no"] = strconv.Itoa(p.PageNo)
	}
	if p.PageCount > 0 {
		m["page_count"] = strconv.Itoa(p.PageCount)
	}
	return m
}

// DisclosureItem 은 공시검색 결과 한 건.
type DisclosureItem struct {
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpName  string `json:"corp_name"`  // 종목명(법인명)
	CorpCode  string `json:"corp_code"`  // 고유번호 (8자리)
	StockCode string `json:"stock_code"` // 종목코드 (6자리)
	ReportNm  string `json:"report_nm"`  // 보고서명
	RceptNo   string `json:"rcept_no"`   // 접수번호 (DownloadDocument 인자)
	FlrNm     string `json:"flr_nm"`     // 공시 제출인명
	RceptDt   string `json:"rcept_dt"`   // 접수일자 YYYYMMDD
	Rm        string `json:"rm"`         // 비고
}

// SearchResult 는 공시검색 응답 (페이지네이션 + 목록).
type SearchResult struct {
	httpclient.Envelope
	PageNo     int              `json:"page_no"`     // 페이지 번호
	PageCount  int              `json:"page_count"`  // 페이지당 건수
	TotalCount int              `json:"total_count"` // 총 건수
	TotalPage  int              `json:"total_page"`  // 총 페이지 수
	List       []DisclosureItem `json:"list"`        // 공시 목록
}

// SearchDisclosures 는 조건별 공시보고서를 검색한다.
func (c *Client) SearchDisclosures(ctx context.Context, params SearchParams) (*SearchResult, error) {
	var out SearchResult
	if err := c.http.GetJSON(ctx, "/api/list.json", params.toMap(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./disclosure/ -run 'TestSearchDisclosures|TestSearchParams' -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add disclosure/search.go disclosure/search_test.go disclosure/testdata/list_samsung.json
git commit -m "feat(disclosure): SearchDisclosures (공시검색)"
```

---

### Task 5: disclosure — 공시서류원본파일 (DownloadDocument)

**Files:**
- Create: `disclosure/document.go`
- Test: `disclosure/document_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `disclosure/document_test.go`:
```go
package disclosure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

func TestDownloadDocument_Binary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "20240131000326", r.URL.Query().Get("rcept_no"))
		w.Write([]byte("PK\x03\x04docdata"))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	b, err := c.DownloadDocument(context.Background(), "20240131000326")
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04docdata", string(b))
}

func TestDownloadDocument_ErrorJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	_, err := c.DownloadDocument(context.Background(), "x")
	var apiErr *httpclient.APIError
	assert.ErrorAs(t, err, &apiErr)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./disclosure/ -run TestDownloadDocument -v`
Expected: FAIL — `undefined: DownloadDocument`.

- [ ] **Step 3: 구현** — `disclosure/document.go`:
```go
package disclosure

import "context"

// DownloadDocument 은 접수번호(rcept_no)로 공시서류 원본 ZIP 을 그대로 반환한다.
// 압축 해제·파싱은 호출자 몫 (임의 공시 원본은 형태가 다양한 바이너리).
func (c *Client) DownloadDocument(ctx context.Context, rceptNo string) ([]byte, error) {
	return c.http.GetBytes(ctx, "/api/document.xml", map[string]string{"rcept_no": rceptNo})
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./disclosure/ -v`
Expected: PASS (전체 disclosure 테스트). `go vet ./disclosure/` clean.

- [ ] **Step 5: Commit**

```bash
git add disclosure/document.go disclosure/document_test.go
git commit -m "feat(disclosure): DownloadDocument (공시서류원본파일)"
```

---

### Task 6: root — errors / config / client / from_env (생성 & 와이어링)

**Files:**
- Create: `errors.go`, `config.go`, `client.go`, `from_env.go`
- Test: `client_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `client_test.go`:
```go
package opendart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_RequiresAPIKey(t *testing.T) {
	_, err := NewClient("")
	require.Error(t, err)
}

func TestNewClient_WiresSubClients(t *testing.T) {
	c, err := NewClient("KEY", WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)
}

func TestNewClientFromEnv(t *testing.T) {
	t.Setenv("OPENDART_API_KEY", "ENVKEY")
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)

	t.Setenv("OPENDART_API_KEY", "")
	_, err = NewClientFromEnv()
	require.Error(t, err)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test . -run 'TestNewClient' -v`
Expected: FAIL — `undefined: NewClient`, `WithCorpCodeCacheDir`, `NewClientFromEnv`.

- [ ] **Step 3: 구현**

`errors.go`:
```go
package opendart

import (
	"errors"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// APIError 는 OpenDART status != "000" 응답. errors.As 로 Status/Message 접근.
type APIError = httpclient.APIError

// ErrNoData 는 status 013 (조회된 데이터 없음).
var ErrNoData = httpclient.ErrNoData

// ErrCorpCodeNotFound 는 종목코드/고유번호 매핑 실패.
var ErrCorpCodeNotFound = errors.New("opendart: corp_code not found")
```

`config.go`:
```go
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
```

`client.go`:
```go
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
)

const defaultCorpCacheTTL = 24 * time.Hour

// Client 는 opendart 라이브러리의 단일 진입점.
type Client struct {
	apiKey string
	http   *httpclient.Client
	corp   *corpcode.Cache

	Disclosure *disclosure.Client // DS001 공시정보
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

	c := &Client{apiKey: apiKey, http: hc}
	c.Disclosure = disclosure.New(hc)

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
```

`from_env.go`:
```go
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
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test . -run 'TestNewClient' -v`
Expected: PASS. `go build ./...` 성공(이제 root 패키지가 컴파일됨). `go vet ./...` clean.

- [ ] **Step 5: Commit**

```bash
git add errors.go config.go client.go from_env.go client_test.go
git commit -m "feat(opendart): Client construction + options + from_env"
```

---

### Task 7: root — corp_code 메서드 (corpcode.go)

**Files:**
- Create: `corpcode.go`
- Test: `corpcode_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `corpcode_test.go` (httptest 로 corpCode.xml ZIP 서빙):
```go
package opendart

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const corpXML = `<?xml version="1.0" encoding="UTF-8"?>
<result>
	<list><corp_code>00126380</corp_code><corp_name>삼성전자(주)</corp_name><corp_eng_name>SAMSUNG</corp_eng_name><stock_code>005930</stock_code><modify_date>20240101</modify_date></list>
	<list><corp_code>00434003</corp_code><corp_name>다코</corp_name><corp_eng_name>Daco</corp_eng_name><stock_code> </stock_code><modify_date>20170630</modify_date></list>
</result>`

func corpCodeServer(t *testing.T) *httptest.Server {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("CORPCODE.xml")
	require.NoError(t, err)
	_, err = w.Write([]byte(corpXML))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	zipBytes := buf.Bytes()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(zipBytes)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newCorpTestClient(t *testing.T) *Client {
	srv := corpCodeServer(t)
	c, err := NewClient("KEY",
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCorpCodeCacheDir(t.TempDir()),
	)
	require.NoError(t, err)
	return c
}

func TestResolveCorpCode(t *testing.T) {
	c := newCorpTestClient(t)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	assert.Equal(t, "00126380", corp)

	_, err = c.ResolveCorpCode(context.Background(), "999999")
	assert.ErrorIs(t, err, ErrCorpCodeNotFound)
}

func TestLookupAndList(t *testing.T) {
	c := newCorpTestClient(t)
	e, err := c.LookupCorpCode(context.Background(), "00434003")
	require.NoError(t, err)
	assert.Equal(t, "다코", e.CorpName)

	all, err := c.CorpCodes(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 2)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test . -run 'TestResolveCorpCode|TestLookupAndList' -v`
Expected: FAIL — `undefined: (*Client).ResolveCorpCode`, `CorpCodeEntry` 등.

- [ ] **Step 3: 구현** — `corpcode.go`:
```go
package opendart

import (
	"context"

	"github.com/kenshin579/opendart/internal/corpcode"
)

// CorpCodeEntry 는 corp_code 매핑의 한 회사. (internal 타입 별칭 — 외부에서 opendart.CorpCodeEntry 로 사용)
type CorpCodeEntry = corpcode.Entry

// ResolveCorpCode 는 종목코드(6자리)를 corp_code(8자리)로 변환한다.
// 매핑이 없으면 ErrCorpCodeNotFound.
func (c *Client) ResolveCorpCode(ctx context.Context, stockCode string) (string, error) {
	e, ok, err := c.corp.ByStockCode(ctx, stockCode)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrCorpCodeNotFound
	}
	return e.CorpCode, nil
}

// LookupCorpCode 는 corp_code(8자리)로 회사 엔트리를 조회한다.
func (c *Client) LookupCorpCode(ctx context.Context, corpCode string) (CorpCodeEntry, error) {
	e, ok, err := c.corp.ByCorpCode(ctx, corpCode)
	if err != nil {
		return CorpCodeEntry{}, err
	}
	if !ok {
		return CorpCodeEntry{}, ErrCorpCodeNotFound
	}
	return e, nil
}

// CorpCodes 는 전체 회사 엔트리를 반환한다 (사용자 직접 필터용).
func (c *Client) CorpCodes(ctx context.Context) ([]CorpCodeEntry, error) {
	return c.corp.Entries(ctx)
}

// RefreshCorpCodes 는 TTL 을 무시하고 corp_code 매핑을 강제 재다운로드한다.
func (c *Client) RefreshCorpCodes(ctx context.Context) error {
	return c.corp.Refresh(ctx)
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test . -run 'TestResolveCorpCode|TestLookupAndList' -v`
Expected: PASS. `go test ./...` 전체 PASS. `go vet ./...` clean.

- [ ] **Step 5: Commit**

```bash
git add corpcode.go corpcode_test.go
git commit -m "feat(opendart): corp_code resolution methods on root Client"
```

---

### Task 8: 예제 · README · 통합 테스트 · 최종 검증

**Files:**
- Create: `examples/disclosure/main.go`, `README.md`, `integration_test.go`

- [ ] **Step 1: 통합 테스트 작성** — `integration_test.go` (실 API, build tag 로 기본 제외):
```go
//go:build integration

package opendart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// 실행: OPENDART_API_KEY=... go test -tags integration -run TestIntegration -v
func TestIntegration_GetCompany(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	require.Equal(t, "00126380", corp)

	company, err := c.Disclosure.GetCompany(context.Background(), corp)
	require.NoError(t, err)
	require.Contains(t, company.CorpName, "삼성전자")
}
```

- [ ] **Step 2: 통합 테스트가 기본 빌드에서 제외되는지 확인**

Run: `go vet ./...` (build tag 파일은 기본 제외되어 영향 없음)
Run: `go test ./... 2>&1 | tail -5` (통합 테스트 미실행, 나머지 PASS)
Expected: 전체 PASS, integration_test.go 는 실행 안 됨.

- [ ] **Step 3: 예제 작성** — `examples/disclosure/main.go` (`SearchParams` 는 `disclosure` 패키지 타입이므로 import 함):
```go
// examples/disclosure — DS001 공시정보 사용 예제.
//
// 실행: OPENDART_API_KEY=... go run ./examples/disclosure
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kenshin579/opendart"
	"github.com/kenshin579/opendart/disclosure"
)

func main() {
	client, err := opendart.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	// 종목코드 → corp_code
	corp, err := client.ResolveCorpCode(ctx, "005930")
	if err != nil {
		log.Fatalf("ResolveCorpCode: %v", err)
	}

	// 기업개황
	company, err := client.Disclosure.GetCompany(ctx, corp)
	if err != nil {
		log.Fatalf("GetCompany: %v", err)
	}
	fmt.Printf("회사명: %s (%s) 대표: %s 설립: %s\n",
		company.CorpName, company.StockCode, company.CeoName, company.EstDate)

	// 공시검색 (최근 5건)
	res, err := client.Disclosure.SearchDisclosures(ctx, disclosure.SearchParams{CorpCode: corp, PageCount: 5})
	if err != nil {
		log.Fatalf("SearchDisclosures: %v", err)
	}
	fmt.Printf("총 %d건 중 %d건:\n", res.TotalCount, len(res.List))
	for _, d := range res.List {
		fmt.Printf("  [%s] %s (%s)\n", d.RceptDt, d.ReportNm, d.RceptNo)
	}
}
```

- [ ] **Step 4: 예제 컴파일 확인**

Run: `go build ./examples/disclosure/`
Expected: 성공 (Step 3 주석대로 disclosure.SearchParams 사용 시).

- [ ] **Step 5: README 작성** — `README.md`:
```markdown
# opendart

DART 전자공시시스템 OpenAPI 의 Go 클라이언트 라이브러리.

## 설치

```bash
go get github.com/kenshin579/opendart
```

## 사용

```go
client, _ := opendart.NewClientFromEnv() // OPENDART_API_KEY
ctx := context.Background()

corp, _ := client.ResolveCorpCode(ctx, "005930")        // 종목코드 → corp_code
company, _ := client.Disclosure.GetCompany(ctx, corp)   // 기업개황
res, _ := client.Disclosure.SearchDisclosures(ctx, disclosure.SearchParams{CorpCode: corp})
zip, _ := client.Disclosure.DownloadDocument(ctx, "20240131000326") // 원본 ZIP
```

## 인증

발급받은 API 키를 `OPENDART_API_KEY` 환경변수로 두거나 `opendart.NewClient(apiKey)` 로 전달한다.

## 커버리지

- DS001 공시정보: 기업개황 · 공시검색 · 공시서류원본파일 · 고유번호(corp_code 매핑)
- (예정) DS002~DS006

## 에러 처리

- `errors.Is(err, opendart.ErrNoData)` — 조회 데이터 없음(013)
- `errors.As(err, &apiErr)` — 그 외 OpenDART status (`*opendart.APIError`)

## 문서

API 명세: [`docs/api/`](docs/api/README.md)
```

- [ ] **Step 6: 최종 전체 검증**

Run:
```bash
go build ./...
go vet ./...
go test ./...
gofmt -l .
```
Expected: build/vet 성공, 모든 테스트 PASS, gofmt 차이 없음(출력 비어있음 — scripts/crawl 등 기존 파일 제외하고 신규 파일 정렬됨).

- [ ] **Step 7: Commit**

```bash
git add examples/ README.md integration_test.go
git commit -m "docs(opendart): example, README, integration test"
```

---

## Self-Review Notes

- **Spec coverage:** NewClient/FromEnv/options=Task6 · APIError/ErrNoData/ErrCorpCodeNotFound=Task6 · httpclient(crtfc_key/status/TLS/GetJSON/GetBytes)=Task1 · corp_code 인프라(다운로드/unzip/파싱/캐시TTL/fallback/인덱스)=Task2 · root corp_code 메서드(Resolve/Lookup/CorpCodes/Refresh)=Task7 · disclosure GetCompany/SearchDisclosures/DownloadDocument=Task3/4/5 · 필드 한글 코멘트=각 타입 · 테스트(httptest/fixture/zip/build-tag 통합)=각 Task+Task8 · examples/README=Task8. 모두 매핑됨.
- **Type consistency:** `httpclient.{Envelope,StatusReader,APIError,ErrNoData,Config,Client,GetJSON,GetBytes}` · `corpcode.{Entry,FetchFunc,Cache,New,ByStockCode,ByCorpCode,Entries,Refresh}` · `disclosure.{Client,New,Company,GetCompany,SearchParams,SearchResult,DisclosureItem,SearchDisclosures,DownloadDocument}` · root `{Client,NewClient,NewClientFromEnv,Option,With*,APIError(별칭),ErrNoData(별칭),ErrCorpCodeNotFound,CorpCodeEntry(별칭),ResolveCorpCode,LookupCorpCode,CorpCodes,RefreshCorpCodes}` — 정의·사용 시그니처 일치. 응답 타입은 `httpclient.Envelope` 임베드로 `StatusReader` 충족.
- **Import 순환 회피:** internal/httpclient 가 APIError/ErrNoData 정의 → root 가 `type APIError = httpclient.APIError` / `var ErrNoData = httpclient.ErrNoData` 별칭. internal/corpcode 는 FetchFunc 주입으로 httpclient 비의존(root 가 콜백으로 연결). CorpCodeEntry 도 internal 타입 별칭으로 외부 노출.
- **검증된 fixture:** company_samsung.json / list_samsung.json 은 실 API 응답. corpCode XML 구조(`<result><list>`, 공백 stock_code)도 실측 반영.
- **Task 8 예제:** main.go 는 `disclosure.SearchParams` 를 쓰므로 `disclosure` 패키지를 import 한다(Step 3 코드에 반영됨).
