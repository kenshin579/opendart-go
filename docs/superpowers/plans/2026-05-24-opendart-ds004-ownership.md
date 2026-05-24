# OpenDART DS004 지분공시 종합정보 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** OpenDART DS004 지분공시 종합정보 2개 API를 신규 `ownership` 패키지로 추가하고, 공용 list 헬퍼를 `internal/httpclient` 로 추출한다.

**Architecture:** `report` 의 unexported 제네릭 list 로직을 `internal/httpclient.GetList[T]` 로 승격(공용화)하고 `report` 는 위임(공개 API·동작 불변). 신규 sub-package `ownership` 은 `corp_code` 단일 파라미터로 `httpclient.GetList` 를 호출하는 2개 메서드(`MajorStockReports`/`ExecutiveStockReports`). root `client.Ownership` 와이어링.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds004-ownership-design.md`

**검증된 사실 (실 API):** `majorstock.json?corp_code=00126380` → status 000, 40행. `elestock.json?corp_code=00126380` → 2619행. 둘 다 `corp_code` 단일 파라미터, JSON list. 숫자/비율 문자열, 빈 값 "-".

**기존 재사용 심볼:** `httpclient.Client`(GetJSON)/`httpclient.Envelope`/`httpclient.APIError`/`httpclient.ErrNoData`, `report` 의 기존 테스트 39개, root `Client`/`NewClient`.

---

## File Structure

```
internal/httpclient/
  list.go            # GetList[T] + listEnvelope[T] (신규)
  list_test.go       # (신규)
report/
  client.go          # (수정) getListParams 위임, listResponse 제거
ownership/
  client.go          # Client + New (신규)
  ownership.go       # 2개 메서드 + item struct (신규)
  client_test.go     # newTestClient 헬퍼 (신규)
  ownership_test.go  # 2개 fixture 테스트 (신규)
  testdata/          # majorstock.json, elestock.json (신규)
client.go            # (수정) Ownership 필드 + 와이어링
README.md            # (수정) DS004 커버리지
integration_test.go  # (수정) MajorStockReports 통합 케이스
```

---

### Task 1: 공용 list 헬퍼 추출 (httpclient.GetList) + report 위임

**Files:**
- Create: `internal/httpclient/list.go`, `internal/httpclient/list_test.go`
- Modify: `report/client.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `internal/httpclient/list_test.go`:
```go
package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type listItem struct {
	Name string `json:"name"`
}

func TestGetList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "KEY", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "x", r.URL.Query().Get("p"))
		w.Write([]byte(`{"status":"000","message":"정상","list":[{"name":"가"},{"name":"나"}]}`))
	}))
	t.Cleanup(srv.Close)
	c := New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})

	items, err := GetList[listItem](context.Background(), c, "/api/x.json", map[string]string{"p": "x"})
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "가", items[0].Name)
}

func TestGetList_NoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})

	_, err := GetList[listItem](context.Background(), c, "/api/x.json", nil)
	assert.ErrorIs(t, err, ErrNoData)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./internal/httpclient/ -run TestGetList -v`
Expected: FAIL — `undefined: GetList`.

- [ ] **Step 3: 구현** — `internal/httpclient/list.go`:
```go
package httpclient

import "context"

// listEnvelope 는 OpenDART list 응답 공통 형태 (status/message + list).
type listEnvelope[T any] struct {
	Envelope
	List []T `json:"list"`
}

// GetList 는 JSON list 응답을 디코드해 list 만 반환한다.
// status 검사(013→ErrNoData, 그 외→*APIError)는 GetJSON 이 수행한다.
func GetList[T any](ctx context.Context, c *Client, path string, params map[string]string) ([]T, error) {
	var resp listEnvelope[T]
	if err := c.GetJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./internal/httpclient/ -v`
Expected: PASS (기존 + 신규 2). `go vet ./internal/httpclient/` clean.

- [ ] **Step 5: report 위임으로 변경** — `report/client.go` 에서 기존 `listResponse[T]` 타입 정의와 `getListParams` 함수를 찾는다. `listResponse[T]` 정의(아래)를 **삭제**한다:
```go
// listResponse 는 DS002 공통 list 응답 envelope.
type listResponse[T any] struct {
	httpclient.Envelope
	List []T `json:"list"`
}
```
그리고 `getListParams` 본문을 httpclient.GetList 위임으로 교체한다:
```go
// getListParams 는 raw 파라미터 맵으로 list 를 조회하는 코어 헬퍼.
func getListParams[T any](ctx context.Context, hc *httpclient.Client, path string, params map[string]string) ([]T, error) {
	return httpclient.GetList[T](ctx, hc, path, params)
}
```
(`getList[T](...,ReportParams)` wrapper 는 그대로 둔다.)

- [ ] **Step 6: report 회귀 없음 확인**

Run: `go test ./report/ -v`
Expected: 기존 39개 모두 PASS. `go vet ./report/` clean, `gofmt -l report/client.go` no output.

- [ ] **Step 7: Commit**

```bash
git add internal/httpclient/list.go internal/httpclient/list_test.go report/client.go
git commit -m "refactor(httpclient): extract generic GetList; report delegates"
```

---

### Task 2: ownership 패키지 (2개 메서드)

**Files:**
- Create: `ownership/client.go`, `ownership/ownership.go`, `ownership/client_test.go`, `ownership/ownership_test.go`
- Create: `ownership/testdata/majorstock.json`, `ownership/testdata/elestock.json`

- [ ] **Step 1: fixture 작성**

`ownership/testdata/majorstock.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240524000517",
            "rcept_dt": "2024-05-24",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "report_tp": "일반",
            "repror": "삼성물산",
            "stkqy": "1,199,285,813",
            "stkqy_irds": "-205,875",
            "stkrt": "20.09",
            "stkrt_irds": "-0.00",
            "ctr_stkqy": "97,526,980",
            "ctr_stkrt": "1.63",
            "report_resn": "특별관계자 및 보유주식수 변동"
        }
    ]
}
```

`ownership/testdata/elestock.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240529000354",
            "rcept_dt": "2024-05-29",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "repror": "손준호",
            "isu_exctv_rgist_at": "비등기임원",
            "isu_exctv_ofcps": "상무",
            "isu_main_shrholdr": "-",
            "sp_stock_lmp_cnt": "0",
            "sp_stock_lmp_irds_cnt": "-1,400",
            "sp_stock_lmp_rate": "0.00",
            "sp_stock_lmp_irds_rate": "0.00"
        }
    ]
}
```

- [ ] **Step 2: client_test.go (newTestClient 헬퍼) + ownership_test.go 작성**

`ownership/client_test.go`:
```go
package ownership

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 testdata fixture 를 path 별로 서빙하는 ownership.Client 를 만든다.
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
```

`ownership/ownership_test.go`:
```go
package ownership

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMajorStockReports(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/majorstock.json": "majorstock.json"})
	items, err := c.MajorStockReports(context.Background(), "00126380")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성물산", items[0].Repror)
	assert.Equal(t, "1,199,285,813", items[0].Stkqy)
	assert.Equal(t, "20.09", items[0].Stkrt)
}

func TestExecutiveStockReports(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/elestock.json": "elestock.json"})
	items, err := c.ExecutiveStockReports(context.Background(), "00126380")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "손준호", items[0].Repror)
	assert.Equal(t, "상무", items[0].IsuExctvOfcps)
	assert.Equal(t, "비등기임원", items[0].IsuExctvRgistAt)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./ownership/ -v`
Expected: FAIL — `undefined: Client`, `New`, `MajorStockReports` 등.

- [ ] **Step 4: 구현**

`ownership/client.go`:
```go
// Package ownership 는 OpenDART DS004 지분공시 종합정보 API sub-client 다.
// opendart.Client.Ownership 로 접근한다.
package ownership

import "github.com/kenshin579/opendart/internal/httpclient"

// Client 는 지분공시 종합정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }
```

`ownership/ownership.go`:
```go
package ownership

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// MajorStockItem 은 대량보유 상황보고 (majorstock) 한 건.
type MajorStockItem struct {
	RceptNo    string `json:"rcept_no"`    // 접수번호
	RceptDt    string `json:"rcept_dt"`    // 접수일자
	CorpCode   string `json:"corp_code"`   // 고유번호
	CorpName   string `json:"corp_name"`   // 회사명
	ReportTp   string `json:"report_tp"`   // 보고구분
	Repror     string `json:"repror"`      // 대표보고자
	Stkqy      string `json:"stkqy"`       // 보유주식등의 수
	StkqyIrds  string `json:"stkqy_irds"`  // 보유주식등의 증감
	Stkrt      string `json:"stkrt"`       // 보유비율
	StkrtIrds  string `json:"stkrt_irds"`  // 보유비율 증감
	CtrStkqy   string `json:"ctr_stkqy"`   // 주요체결 주식등의 수
	CtrStkrt   string `json:"ctr_stkrt"`   // 주요체결 보유비율
	ReportResn string `json:"report_resn"` // 보고사유
}

// MajorStockReports 는 대량보유 상황보고(5% 룰)를 조회한다.
func (c *Client) MajorStockReports(ctx context.Context, corpCode string) ([]MajorStockItem, error) {
	return httpclient.GetList[MajorStockItem](ctx, c.http, "/api/majorstock.json", map[string]string{"corp_code": corpCode})
}

// ExecutiveStockItem 은 임원·주요주주 소유보고 (elestock) 한 건.
type ExecutiveStockItem struct {
	RceptNo            string `json:"rcept_no"`               // 접수번호
	RceptDt            string `json:"rcept_dt"`               // 접수일자
	CorpCode           string `json:"corp_code"`              // 고유번호
	CorpName           string `json:"corp_name"`              // 회사명
	Repror             string `json:"repror"`                 // 보고자
	IsuExctvRgistAt    string `json:"isu_exctv_rgist_at"`     // 발행회사 관계 임원(등기여부)
	IsuExctvOfcps      string `json:"isu_exctv_ofcps"`        // 발행회사 관계 임원 직위
	IsuMainShrholdr    string `json:"isu_main_shrholdr"`      // 발행회사 관계 주요주주
	SpStockLmpCnt      string `json:"sp_stock_lmp_cnt"`       // 특정증권등 소유 수
	SpStockLmpIrdsCnt  string `json:"sp_stock_lmp_irds_cnt"`  // 특정증권등 소유 증감 수
	SpStockLmpRate     string `json:"sp_stock_lmp_rate"`      // 특정증권등 소유 비율
	SpStockLmpIrdsRate string `json:"sp_stock_lmp_irds_rate"` // 특정증권등 소유 증감 비율
}

// ExecutiveStockReports 는 임원·주요주주 소유보고를 조회한다.
func (c *Client) ExecutiveStockReports(ctx context.Context, corpCode string) ([]ExecutiveStockItem, error) {
	return httpclient.GetList[ExecutiveStockItem](ctx, c.http, "/api/elestock.json", map[string]string{"corp_code": corpCode})
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./ownership/ -v`
Expected: PASS (2개). `go vet ./ownership/` clean, `gofmt -l ownership/*.go` no output.

- [ ] **Step 6: Commit**

```bash
git add ownership/
git commit -m "feat(ownership): DS004 대량보유 상황보고 + 임원·주요주주 소유보고"
```

---

### Task 3: root 와이어링 (client.Ownership)

**Files:**
- Modify: `client.go` (Ownership 필드 + 와이어링)
- Modify: `client_test.go` (Ownership NotNil 검증)

- [ ] **Step 1: 실패하는 테스트 추가** — `client_test.go` 의 `TestNewClient_WiresSubClients` 를 다음으로 교체:
```go
func TestNewClient_WiresSubClients(t *testing.T) {
	c, err := NewClient("KEY", WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)
	assert.NotNil(t, c.Report)
	assert.NotNil(t, c.Ownership)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test . -run TestNewClient_WiresSubClients -v`
Expected: FAIL — `c.Ownership undefined`.

- [ ] **Step 3: 구현** — `client.go` 수정.

import 블록에 추가:
```go
	"github.com/kenshin579/opendart/ownership"
```
`Client` struct 에 `Report` 필드 다음 줄 추가:
```go
	Ownership *ownership.Client // DS004 지분공시 종합정보
```
`NewClient` 내부, `c.Report = report.New(hc)` 다음 줄에 추가:
```go
	c.Ownership = ownership.New(hc)
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test . -run TestNewClient -v`
Expected: PASS. `go build ./...` + `go vet ./...` clean, `gofmt -l client.go` no output.

- [ ] **Step 5: Commit**

```bash
git add client.go client_test.go
git commit -m "feat(opendart): wire Ownership sub-client (DS004)"
```

---

### Task 4: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS003 줄 **다음에** 새 줄을 추가한다:
```markdown
- DS004 지분공시 종합정보: 대량보유 상황보고(5% 룰) · 임원·주요주주 소유보고
```
그리고 `- (예정)` 줄을 다음으로 교체:
```markdown
- (예정) DS002 개인별 보수 Ver2.0 2종 · DS005~DS006
```

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` 유지; `ownership` import 를 import 블록에 추가):
```go
func TestIntegration_MajorStockReports(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Ownership.MajorStockReports(context.Background(), corp)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```
> 주: 이 케이스는 `ownership` 패키지 타입을 직접 쓰지 않으므로 `ownership` import 가 불필요하다.
> 기존 import 블록은 `report` 만 있으면 충분하다(추가 import 없음). 컴파일 에러가 나지 않는지 Step 3 에서 확인한다.

- [ ] **Step 3: 통합 빌드 확인 (기본 빌드 제외)**

Run: `go vet -tags integration ./...`
Expected: clean.
Run: `go test ./...`
Expected: 전체 PASS, integration 미실행.

- [ ] **Step 4: 최종 전체 검증**

Run:
```bash
go build ./...
go vet ./...
go test ./...
gofmt -l . | grep -v '^scripts/crawl' || echo "clean"
```
Expected: build/vet 성공, 전체 PASS, gofmt 신규 파일 차이 없음("clean").

- [ ] **Step 5: Commit**

```bash
git add README.md integration_test.go
git commit -m "docs(opendart): DS004 지분공시 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** httpclient.GetList 추출 + report 위임 = Task1 · ownership 2개 메서드/struct = Task2 · root 와이어링 = Task3 · README/통합 = Task4. 모두 매핑됨.
- **Type consistency:** `httpclient.GetList[T]`/`listEnvelope[T]` · `report.getListParams`(위임)·`getList`(유지) · `ownership.{Client,New,MajorStockItem,MajorStockReports,ExecutiveStockItem,ExecutiveStockReports}` · root `Ownership *ownership.Client`. 시그니처 `(ctx, corpCode string) ([]XItem, error)` 일관. 필드·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** majorstock/elestock 실 API(삼성전자) 응답 첫 항목(majorstock report_resn 다중행은 단일행으로 정리).
- **report 위임 회귀:** report 공개 메서드/시그니처 불변, 기존 39개 테스트로 검증(Task1 Step6).
- **Task4 통합 테스트 import 주의:** `MajorStockReports` 는 `c.Ownership` 메서드라 `ownership` 패키지 타입을 직접 참조하지 않음 → 추가 import 불필요(기존 `report` import 유지). vet 으로 확인.
