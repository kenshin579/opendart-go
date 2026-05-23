# OpenDART DS002 지분·주식·배당 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS002 정기보고서 주요정보의 공통 추상화(`report` 패키지)와 지분·주식·배당 7개 API를 구현한다.

**Architecture:** 신규 sub-package `report/`. 30개 DS002 API가 동일한 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답을 가지므로, `ReportParams`/`ReportCode` + Go 제네릭 `getList[T]`/`listResponse[T]` 공통 헬퍼를 두고 각 엔드포인트는 "item struct + 한 줄 메서드"로 구현. root `Client`에 `Report` 필드를 추가해 `client.Report` 로 노출. DS001 `disclosure` 와 동일 패턴.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds002-equity-design.md`

**검증된 사실 (실 API 호출, 삼성전자 00126380 / 2023 / 11011):** 7개 엔드포인트 모두 status 000 + list 반환. 숫자는 콤마 문자열("508,157,148"), 빈 값은 "-". 응답 필드는 아래 fixture 와 일치 확인.

---

## File Structure

```
report/
  client.go        # Client, New, ReportCode 상수, ReportParams(+toMap), listResponse[T], getList[T]
  equity.go        # 7개 메서드 + item struct
  client_test.go   # newTestClient 헬퍼, toMap, 013→ErrNoData
  equity_test.go   # 7개 메서드 fixture 테스트
  testdata/        # 7개 실 응답 JSON fixture
client.go          # (수정) Report 필드 + 와이어링
examples/report/main.go
README.md          # (수정) 커버리지에 DS002 7개
integration_test.go  # (수정) Dividend 통합 케이스
```

---

### Task 1: report 패키지 공통 헬퍼 + 첫 엔드포인트 (Dividend)

**Files:**
- Create: `report/client.go`, `report/equity.go`, `report/client_test.go`, `report/equity_test.go`
- Create: `report/testdata/alotMatter.json`

- [ ] **Step 1: fixture 작성** — `report/testdata/alotMatter.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "se": "주당액면가액(원)",
            "stock_knd": "-",
            "thstrm": "100",
            "frmtrm": "100",
            "lwfr": "100",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/client_test.go` (공통 헬퍼 + 재사용 `newTestClient`):
```go
package report

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

// newTestClient 는 testdata fixture 를 path 별로 서빙하는 report.Client 를 만든다.
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

func TestReportParams_toMap(t *testing.T) {
	m := ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport}.toMap()
	assert.Equal(t, "00126380", m["corp_code"])
	assert.Equal(t, "2023", m["bsns_year"])
	assert.Equal(t, "11011", m["reprt_code"])
}

func TestGetList_NoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	c := New(hc)
	_, err := c.Dividend(context.Background(), ReportParams{CorpCode: "x", BsnsYear: "2023", ReprtCode: AnnualReport})
	assert.ErrorIs(t, err, httpclient.ErrNoData)
}
```

- [ ] **Step 3: Dividend 테스트 작성** — `report/equity_test.go`:
```go
package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDividend(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/alotMatter.json": "alotMatter.json"})
	items, err := c.Dividend(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "주당액면가액(원)", items[0].Se)
	assert.Equal(t, "100", items[0].Thstrm)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
}
```

- [ ] **Step 4: 테스트 실패 확인**

Run: `go test ./report/ -v`
Expected: FAIL — `undefined: Client`, `New`, `ReportParams`, `AnnualReport`, `Dividend`.

- [ ] **Step 5: 공통 헬퍼 구현** — `report/client.go`:
```go
// Package report 는 OpenDART DS002 정기보고서 주요정보 API sub-client 다.
// opendart.Client.Report 로 접근한다.
package report

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Client 는 정기보고서 주요정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// ReportCode 는 정기보고서 종류 코드.
type ReportCode string

const (
	Q1Report     ReportCode = "11013" // 1분기보고서
	HalfReport   ReportCode = "11012" // 반기보고서
	Q3Report     ReportCode = "11014" // 3분기보고서
	AnnualReport ReportCode = "11011" // 사업보고서
)

// ReportParams 는 DS002 공통 요청 인자 (전부 필수).
type ReportParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
}

func (p ReportParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
	}
}

// listResponse 는 DS002 공통 list 응답 envelope.
type listResponse[T any] struct {
	httpclient.Envelope
	List []T `json:"list"`
}

// getList 는 공통 list 조회 헬퍼. status 검사 후 list 만 반환한다.
// 조회 데이터 없음(013)은 httpclient 가 ErrNoData 로 변환한다.
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, p.toMap(), &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
```

- [ ] **Step 6: Dividend 구현** — `report/equity.go`:
```go
package report

import "context"

// DividendItem 은 배당에 관한 사항 (alotMatter) 한 건.
type DividendItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호 (14자리)
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호 (8자리)
	CorpName string `json:"corp_name"` // 법인명
	Se       string `json:"se"`        // 구분 (주당액면가액, 주당 현금배당금 등)
	StockKnd string `json:"stock_knd"` // 주식 종류 (보통주 등)
	Thstrm   string `json:"thstrm"`    // 당기
	Frmtrm   string `json:"frmtrm"`    // 전기
	Lwfr     string `json:"lwfr"`      // 전전기
	StlmDt   string `json:"stlm_dt"`   // 결산기준일 (YYYY-MM-DD)
}

// Dividend 는 배당에 관한 사항을 조회한다.
func (c *Client) Dividend(ctx context.Context, p ReportParams) ([]DividendItem, error) {
	return getList[DividendItem](ctx, c.http, "/api/alotMatter.json", p)
}
```

- [ ] **Step 7: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: PASS (TestReportParams_toMap, TestGetList_NoData, TestDividend). `go vet ./report/` clean.

- [ ] **Step 8: Commit**

```bash
git add report/
git commit -m "feat(report): DS002 공통 헬퍼 + 배당에 관한 사항 (Dividend)"
```

---

### Task 2: 증자(감자)·자기주식·주식총수 (3 엔드포인트)

**Files:**
- Modify: `report/equity.go` (3개 메서드+struct 추가)
- Modify: `report/equity_test.go` (3개 테스트 추가)
- Create: `report/testdata/irdsSttus.json`, `report/testdata/tesstkAcqsDspsSttus.json`, `report/testdata/stockTotqySttus.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/irdsSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "isu_dcrs_de": "-",
            "isu_dcrs_stle": "-",
            "isu_dcrs_stock_knd": "-",
            "isu_dcrs_qy": "-",
            "isu_dcrs_mstvdv_fval_amount": "-",
            "isu_dcrs_mstvdv_amount": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/tesstkAcqsDspsSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "stock_knd": "-",
            "acqs_mth1": "-",
            "acqs_mth2": "-",
            "acqs_mth3": "-",
            "bsis_qy": "-",
            "change_qy_acqs": "-",
            "change_qy_dsps": "-",
            "change_qy_incnr": "-",
            "trmend_qy": "-",
            "rm": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/stockTotqySttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "se": "보통주",
            "isu_stock_totqy": "20,000,000,000",
            "now_to_isu_stock_totqy": "7,780,466,850",
            "now_to_dcrs_stock_totqy": "1,810,684,300",
            "redc": "-",
            "profit_incnr": "1,810,684,300",
            "rdmstk_repy": "-",
            "etc": "-",
            "istc_totqy": "5,969,782,550",
            "tesstk_co": "-",
            "distb_stock_co": "5,969,782,550",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/equity_test.go` 에 추가:
```go
func TestCapitalChange(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/irdsSttus.json": "irdsSttus.json"})
	items, err := c.CapitalChange(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성전자", items[0].CorpName)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
}

func TestTreasuryStock(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tesstkAcqsDspsSttus.json": "tesstkAcqsDspsSttus.json"})
	items, err := c.TreasuryStock(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestTotalStock(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stockTotqySttus.json": "stockTotqySttus.json"})
	items, err := c.TotalStock(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "보통주", items[0].Se)
	assert.Equal(t, "20,000,000,000", items[0].IsuStockTotqy)
	assert.Equal(t, "5,969,782,550", items[0].DistbStockCo)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestCapitalChange|TestTreasuryStock|TestTotalStock' -v`
Expected: FAIL — `undefined: CapitalChange` 등.

- [ ] **Step 4: 구현** — `report/equity.go` 에 추가:
```go
// CapitalChangeItem 은 증자(감자) 현황 (irdsSttus) 한 건.
type CapitalChangeItem struct {
	RceptNo                 string `json:"rcept_no"`                    // 접수번호 (14자리)
	CorpCls                 string `json:"corp_cls"`                    // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                   // 고유번호 (8자리)
	CorpName                string `json:"corp_name"`                   // 법인명
	IsuDcrsDe               string `json:"isu_dcrs_de"`                 // 주식발행 감소일자
	IsuDcrsStle             string `json:"isu_dcrs_stle"`               // 발행 감소 형태
	IsuDcrsStockKnd         string `json:"isu_dcrs_stock_knd"`          // 발행 감소 주식 종류
	IsuDcrsQy               string `json:"isu_dcrs_qy"`                 // 발행 감소 수량
	IsuDcrsMstvdvFvalAmount string `json:"isu_dcrs_mstvdv_fval_amount"` // 발행 감소 주당 액면 가액
	IsuDcrsMstvdvAmount     string `json:"isu_dcrs_mstvdv_amount"`      // 발행 감소 주당 가액
	StlmDt                  string `json:"stlm_dt"`                     // 결산기준일
}

// CapitalChange 는 증자(감자) 현황을 조회한다.
func (c *Client) CapitalChange(ctx context.Context, p ReportParams) ([]CapitalChangeItem, error) {
	return getList[CapitalChangeItem](ctx, c.http, "/api/irdsSttus.json", p)
}

// TreasuryStockItem 은 자기주식 취득 및 처분 현황 (tesstkAcqsDspsSttus) 한 건.
type TreasuryStockItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	StockKnd      string `json:"stock_knd"`       // 주식 종류
	AcqsMth1      string `json:"acqs_mth1"`       // 취득방법 대분류
	AcqsMth2      string `json:"acqs_mth2"`       // 취득방법 중분류
	AcqsMth3      string `json:"acqs_mth3"`       // 취득방법 소분류
	BsisQy        string `json:"bsis_qy"`         // 기초 수량
	ChangeQyAcqs  string `json:"change_qy_acqs"`  // 변동 수량 취득
	ChangeQyDsps  string `json:"change_qy_dsps"`  // 변동 수량 처분
	ChangeQyIncnr string `json:"change_qy_incnr"` // 변동 수량 소각
	TrmendQy      string `json:"trmend_qy"`       // 기말 수량
	Rm            string `json:"rm"`              // 비고
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// TreasuryStock 은 자기주식 취득 및 처분 현황을 조회한다.
func (c *Client) TreasuryStock(ctx context.Context, p ReportParams) ([]TreasuryStockItem, error) {
	return getList[TreasuryStockItem](ctx, c.http, "/api/tesstkAcqsDspsSttus.json", p)
}

// TotalStockItem 은 주식의 총수 현황 (stockTotqySttus) 한 건.
type TotalStockItem struct {
	RceptNo             string `json:"rcept_no"`                // 접수번호
	CorpCls             string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`               // 고유번호
	CorpName            string `json:"corp_name"`               // 회사명
	Se                  string `json:"se"`                      // 구분 (보통주 등)
	IsuStockTotqy       string `json:"isu_stock_totqy"`         // 발행할 주식의 총수
	NowToIsuStockTotqy  string `json:"now_to_isu_stock_totqy"`  // 현재까지 발행한 주식의 총수
	NowToDcrsStockTotqy string `json:"now_to_dcrs_stock_totqy"` // 현재까지 감소한 주식의 총수
	Redc                string `json:"redc"`                    // 감자
	ProfitIncnr         string `json:"profit_incnr"`            // 이익소각
	RdmstkRepy          string `json:"rdmstk_repy"`             // 상환주식의 상환
	Etc                 string `json:"etc"`                     // 기타
	IstcTotqy           string `json:"istc_totqy"`              // 발행주식의 총수
	TesstkCo            string `json:"tesstk_co"`               // 자기주식수
	DistbStockCo        string `json:"distb_stock_co"`          // 유통주식수
	StlmDt              string `json:"stlm_dt"`                 // 결산기준일
}

// TotalStock 은 주식의 총수 현황을 조회한다.
func (c *Client) TotalStock(ctx context.Context, p ReportParams) ([]TotalStockItem, error) {
	return getList[TotalStockItem](ctx, c.http, "/api/stockTotqySttus.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 1 포함, 회귀 없음). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/equity.go report/equity_test.go report/testdata/
git commit -m "feat(report): 증자(감자)·자기주식·주식총수 현황"
```

---

### Task 3: 최대주주·최대주주변동·소액주주 (3 엔드포인트)

**Files:**
- Modify: `report/equity.go` (3개 메서드+struct 추가)
- Modify: `report/equity_test.go` (3개 테스트 추가)
- Create: `report/testdata/hyslrSttus.json`, `report/testdata/hyslrChgSttus.json`, `report/testdata/mrhlSttus.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/hyslrSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "stock_knd": "보통주",
            "nm": "삼성생명보험㈜",
            "relate": "최대주주 본인",
            "bsis_posesn_stock_co": "508,157,148",
            "bsis_posesn_stock_qota_rt": "8.51",
            "trmend_posesn_stock_co": "508,157,148",
            "trmend_posesn_stock_qota_rt": "8.51",
            "rm": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/hyslrChgSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "change_on": "2021년 04월 29일",
            "mxmm_shrholdr_nm": "삼성생명보험㈜",
            "posesn_stock_co": "1,263,050,053",
            "qota_rt": "21.16%",
            "change_cause": "변동전 최대주주의 피상속",
            "rm": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/mrhlSttus.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240312000736",
            "corp_cls": "Y",
            "corp_code": "00126380",
            "corp_name": "삼성전자",
            "se": "소액주주",
            "shrholdr_co": "4,672,039",
            "shrholdr_tot_co": "4,672,130",
            "shrholdr_rate": "99.99%",
            "hold_stock_co": "4,017,892,514",
            "stock_tot_co": "5,969,782,550",
            "hold_stock_rate": "67.30%",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/equity_test.go` 에 추가:
```go
func TestMajorShareholders(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hyslrSttus.json": "hyslrSttus.json"})
	items, err := c.MajorShareholders(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성생명보험㈜", items[0].Nm)
	assert.Equal(t, "최대주주 본인", items[0].Relate)
	assert.Equal(t, "508,157,148", items[0].TrmendPosesnStockCo)
}

func TestMajorShareholderChanges(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hyslrChgSttus.json": "hyslrChgSttus.json"})
	items, err := c.MajorShareholderChanges(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성생명보험㈜", items[0].MxmmShrholdrNm)
	assert.Equal(t, "변동전 최대주주의 피상속", items[0].ChangeCause)
}

func TestMinorityShareholders(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/mrhlSttus.json": "mrhlSttus.json"})
	items, err := c.MinorityShareholders(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "소액주주", items[0].Se)
	assert.Equal(t, "4,672,039", items[0].ShrholdrCo)
	assert.Equal(t, "67.30%", items[0].HoldStockRate)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestMajorShareholders|TestMajorShareholderChanges|TestMinorityShareholders' -v`
Expected: FAIL — `undefined: MajorShareholders` 등.

- [ ] **Step 4: 구현** — `report/equity.go` 에 추가:
```go
// MajorShareholderItem 은 최대주주 현황 (hyslrSttus) 한 건.
type MajorShareholderItem struct {
	RceptNo                 string `json:"rcept_no"`                    // 접수번호
	CorpCls                 string `json:"corp_cls"`                    // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                   // 고유번호
	CorpName                string `json:"corp_name"`                   // 법인명
	Nm                      string `json:"nm"`                          // 성명
	Relate                  string `json:"relate"`                      // 관계
	StockKnd                string `json:"stock_knd"`                   // 주식 종류
	BsisPosesnStockCo       string `json:"bsis_posesn_stock_co"`        // 기초 소유 주식 수
	BsisPosesnStockQotaRt   string `json:"bsis_posesn_stock_qota_rt"`   // 기초 소유 주식 지분율
	TrmendPosesnStockCo     string `json:"trmend_posesn_stock_co"`      // 기말 소유 주식 수
	TrmendPosesnStockQotaRt string `json:"trmend_posesn_stock_qota_rt"` // 기말 소유 주식 지분율
	Rm                      string `json:"rm"`                          // 비고
	StlmDt                  string `json:"stlm_dt"`                     // 결산기준일
}

// MajorShareholders 는 최대주주 현황을 조회한다.
func (c *Client) MajorShareholders(ctx context.Context, p ReportParams) ([]MajorShareholderItem, error) {
	return getList[MajorShareholderItem](ctx, c.http, "/api/hyslrSttus.json", p)
}

// MajorShareholderChangeItem 은 최대주주 변동현황 (hyslrChgSttus) 한 건.
type MajorShareholderChangeItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 법인명
	ChangeOn       string `json:"change_on"`        // 변동 일
	MxmmShrholdrNm string `json:"mxmm_shrholdr_nm"` // 최대 주주 명
	PosesnStockCo  string `json:"posesn_stock_co"`  // 소유 주식 수
	QotaRt         string `json:"qota_rt"`          // 지분율
	ChangeCause    string `json:"change_cause"`     // 변동 원인
	Rm             string `json:"rm"`               // 비고
	StlmDt         string `json:"stlm_dt"`          // 결산기준일
}

// MajorShareholderChanges 는 최대주주 변동현황을 조회한다.
func (c *Client) MajorShareholderChanges(ctx context.Context, p ReportParams) ([]MajorShareholderChangeItem, error) {
	return getList[MajorShareholderChangeItem](ctx, c.http, "/api/hyslrChgSttus.json", p)
}

// MinorityShareholderItem 은 소액주주 현황 (mrhlSttus) 한 건.
type MinorityShareholderItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	Se            string `json:"se"`              // 구분
	ShrholdrCo    string `json:"shrholdr_co"`     // 주주수
	ShrholdrTotCo string `json:"shrholdr_tot_co"` // 전체 주주수
	ShrholdrRate  string `json:"shrholdr_rate"`   // 주주 비율
	HoldStockCo   string `json:"hold_stock_co"`   // 보유 주식수
	StockTotCo    string `json:"stock_tot_co"`    // 총발행 주식수
	HoldStockRate string `json:"hold_stock_rate"` // 보유 주식 비율
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// MinorityShareholders 는 소액주주 현황을 조회한다.
func (c *Client) MinorityShareholders(ctx context.Context, p ReportParams) ([]MinorityShareholderItem, error) {
	return getList[MinorityShareholderItem](ctx, c.http, "/api/mrhlSttus.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 1·2 포함). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/equity.go report/equity_test.go report/testdata/
git commit -m "feat(report): 최대주주·최대주주변동·소액주주 현황"
```

---

### Task 4: root 와이어링 (client.Report)

**Files:**
- Modify: `client.go` (Report 필드 + 와이어링)
- Modify: `client_test.go` (Report NotNil 검증)

- [ ] **Step 1: 실패하는 테스트 추가** — `client_test.go` 의 `TestNewClient_WiresSubClients` 를 다음으로 교체:
```go
func TestNewClient_WiresSubClients(t *testing.T) {
	c, err := NewClient("KEY", WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)
	assert.NotNil(t, c.Report)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test . -run TestNewClient_WiresSubClients -v`
Expected: FAIL — `c.Report undefined`.

- [ ] **Step 3: 구현** — `client.go` 수정.

import 블록에 추가:
```go
	"github.com/kenshin579/opendart/report"
```

`Client` struct 에 필드 추가 (Disclosure 아래):
```go
	Report *report.Client // DS002 정기보고서 주요정보
```

`NewClient` 내부, `c.Disclosure = disclosure.New(hc)` 다음 줄에 추가:
```go
	c.Report = report.New(hc)
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test . -run TestNewClient -v`
Expected: PASS. `go build ./...` + `go vet ./...` clean.

- [ ] **Step 5: Commit**

```bash
git add client.go client_test.go
git commit -m "feat(opendart): wire Report sub-client (DS002)"
```

---

### Task 5: 예제 · README · 통합 테스트 · 최종 검증

**Files:**
- Create: `examples/report/main.go`
- Modify: `README.md` (커버리지)
- Modify: `integration_test.go` (Dividend 케이스 추가)

- [ ] **Step 1: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (파일 상단 `//go:build integration` 및 기존 import 유지):
```go
func TestIntegration_Dividend(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.Dividend(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```
그리고 `integration_test.go` import 블록에 `"github.com/kenshin579/opendart/report"` 를 추가한다.

- [ ] **Step 2: 통합 테스트 컴파일 확인 (기본 빌드 제외)**

Run: `go vet -tags integration ./...`
Expected: clean (integration 빌드도 컴파일됨).
Run: `go test ./...`
Expected: 전체 PASS, integration 미실행.

- [ ] **Step 3: 예제 작성** — `examples/report/main.go`:
```go
// examples/report — DS002 정기보고서 주요정보(지분·주식·배당) 사용 예제.
//
// 실행: OPENDART_API_KEY=... go run ./examples/report
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kenshin579/opendart"
	"github.com/kenshin579/opendart/report"
)

func main() {
	client, err := opendart.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	corp, err := client.ResolveCorpCode(ctx, "005930")
	if err != nil {
		log.Fatalf("ResolveCorpCode: %v", err)
	}

	p := report.ReportParams{CorpCode: corp, BsnsYear: "2023", ReprtCode: report.AnnualReport}

	// 배당에 관한 사항
	dividends, err := client.Report.Dividend(ctx, p)
	if errors.Is(err, opendart.ErrNoData) {
		fmt.Println("배당 데이터 없음")
	} else if err != nil {
		log.Fatalf("Dividend: %v", err)
	} else {
		fmt.Printf("배당 항목 %d건:\n", len(dividends))
		for _, d := range dividends {
			fmt.Printf("  %s (%s): 당기 %s\n", d.Se, d.StockKnd, d.Thstrm)
		}
	}

	// 최대주주 현황
	majors, err := client.Report.MajorShareholders(ctx, p)
	if errors.Is(err, opendart.ErrNoData) {
		fmt.Println("최대주주 데이터 없음")
	} else if err != nil {
		log.Fatalf("MajorShareholders: %v", err)
	} else {
		fmt.Printf("최대주주 %d명:\n", len(majors))
		for _, m := range majors {
			fmt.Printf("  %s (%s): 지분 %s%%\n", m.Nm, m.Relate, m.TrmendPosesnStockQotaRt)
		}
	}
}
```

- [ ] **Step 4: 예제 컴파일 확인**

Run: `go build -o /tmp/report-example ./examples/report/ && rm -f /tmp/report-example`
Expected: 성공. (디렉토리명 충돌 회피용 `-o`.)

- [ ] **Step 5: README 커버리지 갱신** — `README.md` 의 "## 커버리지" 섹션을 다음으로 교체:
```markdown
## 커버리지

- DS001 공시정보: 기업개황 · 공시검색 · 공시서류원본파일 · 고유번호(corp_code 매핑)
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황
- (예정) DS002 나머지 · DS003~DS006
```

- [ ] **Step 6: 최종 전체 검증**

Run:
```bash
go build ./...
go vet ./...
go test ./...
gofmt -l . | grep -v '^scripts/crawl' || echo "clean"
```
Expected: build/vet 성공, 전체 PASS, gofmt 신규 파일 차이 없음("clean").

- [ ] **Step 7: Commit**

```bash
git add examples/report/ README.md integration_test.go
git commit -m "docs(report): example, README coverage, integration test"
```

---

## Self-Review Notes

- **Spec coverage:** 공통 헬퍼(ReportParams/ReportCode/getList/listResponse)=Task1 · 7개 메서드+struct=Task1(Dividend)+Task2(3)+Task3(3) · root 와이어링=Task4 · 예제/README/통합=Task5 · 테스트(toMap/013→ErrNoData/7개 fixture)=Task1~3+5. 모두 매핑됨.
- **Type consistency:** `report.{Client,New,ReportCode,Q1Report/HalfReport/Q3Report/AnnualReport,ReportParams(+toMap),listResponse[T],getList[T]}` + 7개 `XItem`/메서드(`Dividend`/`CapitalChange`/`TreasuryStock`/`TotalStock`/`MajorShareholders`/`MajorShareholderChanges`/`MinorityShareholders`). 메서드 시그니처 `(ctx, ReportParams) ([]XItem, error)` 일관. 필드명·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 7개 모두 실 API(삼성전자/2023/사업보고서) 응답에서 캡처. 숫자 콤마 string, 빈 값 "-".
- **제네릭 주의:** `getList[T]` 는 패키지 레벨 함수(메서드는 타입 파라미터 불가). 각 메서드가 `c.http` 를 넘겨 호출.
