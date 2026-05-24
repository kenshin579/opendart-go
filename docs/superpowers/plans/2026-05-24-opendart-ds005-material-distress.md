# OpenDART DS005 부실·법적 이벤트 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 주요사항보고서 주요정보의 부실·법적 이벤트 7개 API를 신규 `material` 패키지로 추가한다.

**Architecture:** 신규 sub-package `material`. 36개 DS005 API 가 모두 동일 요청(`corp_code`+`bgn_de`+`end_de`)이므로 공통 `MaterialParams`(+toMap, 빈 값 omit) + 기존 `httpclient.GetList[T]` 재사용(DS002 report 패턴과 동형). 7개 메서드 = 7개 item struct. root `client.Material` 와이어링.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-distress-design.md`

**검증된 사실 (실 API):** 7개 모두 동일 요청 파라미터, JSON list. 실 사례로 fixture 캡처: 부도(DH오토넥스 00126089), 회생(DH오토넥스), 영업정지(태광산업 00153393), 해산(코리아퍼시픽01호 00580603), 채권은행개시(태영건설 00153861), 소송(올리패스 01070149). **채권은행 중단(bnkMngtPcsp)** 은 매우 드물어 실 주요사항보고서 데이터를 찾지 못함 → fixture 는 docs 필드 스키마 + 실 sibling(bnkMngtPcbg)을 모델로 한 doc-derived(값만 합성, 키/구조는 정확). 통합 테스트는 실 데이터인 `DefaultOccurrences` 로 한다.

**기존 재사용 심볼:** `httpclient.Client`/`httpclient.GetList[T]`, root `Client`/`NewClient`, 기존 sub-client 와이어링 패턴.

---

## File Structure

```
material/
  client.go         # Client + New + MaterialParams(+toMap) (신규)
  distress.go       # 7개 메서드 + item struct (신규)
  client_test.go    # newTestClient + MaterialParams.toMap 테스트 (신규)
  distress_test.go  # 7개 fixture 테스트 (신규)
  testdata/         # 7개 fixture (신규)
client.go           # (수정) Material 필드 + 와이어링
README.md           # (수정) DS005 커버리지
integration_test.go # (수정) DefaultOccurrences 통합 케이스
```

---

### Task 1: material 패키지 + MaterialParams + 부도발생 (DefaultOccurrences)

**Files:**
- Create: `material/client.go`, `material/distress.go`, `material/client_test.go`, `material/distress_test.go`
- Create: `material/testdata/dfOcr.json`

- [ ] **Step 1: fixture 작성** — `material/testdata/dfOcr.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20231012000317",
            "corp_cls": "Y",
            "corp_code": "00126089",
            "corp_name": "DH오토넥스",
            "df_cn": "당사 김제지점 발행 만기어음 부도",
            "df_amt": "48,322,175",
            "df_bnk": "신한은행 광산금융센터",
            "dfd": "-",
            "df_rs": "부도사유 : 법적지급제한"
        }
    ]
}
```

- [ ] **Step 2: client_test.go 작성 (newTestClient + toMap 테스트)** — `material/client_test.go`:
```go
package material

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 testdata fixture 를 path 별로 서빙하는 material.Client 를 만든다.
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

func TestMaterialParams_toMap(t *testing.T) {
	m := MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"}.toMap()
	assert.Equal(t, "00126089", m["corp_code"])
	assert.Equal(t, "20230101", m["bgn_de"])
	assert.Equal(t, "20231231", m["end_de"])

	only := MaterialParams{CorpCode: "00126089"}.toMap()
	_, hasBgn := only["bgn_de"]
	assert.False(t, hasBgn)
	assert.Equal(t, "00126089", only["corp_code"])
}
```

- [ ] **Step 3: distress_test.go 작성 (DefaultOccurrences 테스트)** — `material/distress_test.go`:
```go
package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOccurrences(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dfOcr.json": "dfOcr.json"})
	items, err := c.DefaultOccurrences(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DH오토넥스", items[0].CorpName)
	assert.Equal(t, "당사 김제지점 발행 만기어음 부도", items[0].DfCn)
	assert.Equal(t, "48,322,175", items[0].DfAmt)
}
```

- [ ] **Step 4: 테스트 실패 확인**

Run: `go test ./material/ -v`
Expected: FAIL — `undefined: Client`, `New`, `MaterialParams`, `DefaultOccurrences`.

- [ ] **Step 5: client.go 구현** — `material/client.go`:
```go
// Package material 는 OpenDART DS005 주요사항보고서 주요정보 API sub-client 다.
// opendart.Client.Material 로 접근한다.
package material

import "github.com/kenshin579/opendart/internal/httpclient"

// Client 는 주요사항보고서 주요정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// MaterialParams 는 DS005 공통 요청 인자 (날짜 범위). 빈 값은 쿼리에서 생략한다(OpenDART 기본값 적용).
type MaterialParams struct {
	CorpCode string // 고유번호 (8자리)
	BgnDe    string // 시작일 YYYYMMDD
	EndDe    string // 종료일 YYYYMMDD
}

func (p MaterialParams) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" {
		m["bgn_de"] = p.BgnDe
	}
	if p.EndDe != "" {
		m["end_de"] = p.EndDe
	}
	return m
}
```

- [ ] **Step 6: distress.go 구현 (DefaultOccurrences)** — `material/distress.go`:
```go
package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// DefaultItem 은 부도발생 (dfOcr) 한 건.
type DefaultItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DfCn     string `json:"df_cn"`     // 부도내용
	DfAmt    string `json:"df_amt"`    // 부도금액
	DfBnk    string `json:"df_bnk"`    // 부도발생은행
	Dfd      string `json:"dfd"`       // 최종부도(당좌거래정지)일자
	DfRs     string `json:"df_rs"`     // 부도사유 및 경위
}

// DefaultOccurrences 는 부도발생을 조회한다.
func (c *Client) DefaultOccurrences(ctx context.Context, p MaterialParams) ([]DefaultItem, error) {
	return httpclient.GetList[DefaultItem](ctx, c.http, "/api/dfOcr.json", p.toMap())
}
```

- [ ] **Step 7: 테스트 통과 확인**

Run: `go test ./material/ -v`
Expected: PASS (TestMaterialParams_toMap, TestDefaultOccurrences). `go vet ./material/` clean, `gofmt -l material/` no output.

- [ ] **Step 8: Commit**

```bash
git add material/
git commit -m "feat(material): DS005 공통 MaterialParams + 부도발생"
```

---

### Task 2: 영업정지·회생절차·해산사유 (3 엔드포인트)

**Files:**
- Modify: `material/distress.go`, `material/distress_test.go`
- Create: `material/testdata/bsnSp.json`, `material/testdata/ctrcvsBgrq.json`, `material/testdata/dsRsOcr.json`

- [ ] **Step 1: fixture 작성**

`material/testdata/bsnSp.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20230616000582",
            "corp_cls": "Y",
            "corp_code": "00153393",
            "corp_name": "태광산업",
            "bsnsp_rm": "방적 사업",
            "bsnsp_amt": "97,254,982,693",
            "rsl": "2,703,851,289,364",
            "sl_vs": "3.6",
            "ls_atn": "해당",
            "krx_stt_atn": "O",
            "bsnsp_cn": "면사, 혼방사",
            "bsnsp_rs": "사업환경 및 사업실적의 지속적인 악화에 따른 사업 중단",
            "ft_ctp": "잔여 사업 집중과 신규 사업 추진을 통한 수익성 개선",
            "bsnsp_af": "장기적으로 수익 개선 효과 기대",
            "bsnspd": "2023년 08월 31일",
            "bddd": "2023년 06월 16일",
            "od_a_at_t": "2",
            "od_a_at_b": "1",
            "adt_a_atn": "참석"
        }
    ]
}
```

`material/testdata/ctrcvsBgrq.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20230926000754",
            "corp_cls": "Y",
            "corp_code": "00126089",
            "corp_name": "DH오토넥스",
            "apcnt": "주식회사 대유플러스",
            "cpct": "서울회생법원",
            "rq_rs": "경영정상화 및 향후 계속 기업으로의 가치 보존",
            "rqd": "2023년 09월 25일",
            "ft_ctp_sc": "2023.09.25자로 서울회생법원에 회생절차개시를 신청하였습니다."
        }
    ]
}
```

`material/testdata/dsRsOcr.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20200327000589",
            "corp_cls": "E",
            "corp_code": "00580603",
            "corp_name": "코리아퍼시픽01호선박투자회사",
            "ds_rs": "존립기간의 만료",
            "ds_rsd": "2020년 03월 27일",
            "od_a_at_t": "-",
            "od_a_at_b": "-",
            "adt_a_atn": "-"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `material/distress_test.go` 에 추가:
```go
func TestBusinessSuspensions(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnSp.json": "bsnSp.json"})
	items, err := c.BusinessSuspensions(context.Background(), MaterialParams{CorpCode: "00153393", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "방적 사업", items[0].BsnspRm)
	assert.Equal(t, "97,254,982,693", items[0].BsnspAmt)
	assert.Equal(t, "2023년 08월 31일", items[0].Bsnspd)
}

func TestRehabilitationApplications(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ctrcvsBgrq.json": "ctrcvsBgrq.json"})
	items, err := c.RehabilitationApplications(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "서울회생법원", items[0].Cpct)
	assert.Equal(t, "주식회사 대유플러스", items[0].Apcnt)
}

func TestDissolutionCauses(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dsRsOcr.json": "dsRsOcr.json"})
	items, err := c.DissolutionCauses(context.Background(), MaterialParams{CorpCode: "00580603", BgnDe: "20200101", EndDe: "20201231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "존립기간의 만료", items[0].DsRs)
	assert.Equal(t, "2020년 03월 27일", items[0].DsRsd)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./material/ -run 'TestBusinessSuspensions|TestRehabilitationApplications|TestDissolutionCauses' -v`
Expected: FAIL — `undefined: ... BusinessSuspensions` 등.

- [ ] **Step 4: 구현** — `material/distress.go` 에 추가:
```go
// BusinessSuspensionItem 은 영업정지 (bsnSp) 한 건.
type BusinessSuspensionItem struct {
	RceptNo   string `json:"rcept_no"`    // 접수번호
	CorpCls   string `json:"corp_cls"`    // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`   // 고유번호
	CorpName  string `json:"corp_name"`   // 회사명
	BsnspRm   string `json:"bsnsp_rm"`    // 영업정지 분야
	BsnspAmt  string `json:"bsnsp_amt"`   // 영업정지 내역(영업정지금액)
	Rsl       string `json:"rsl"`         // 영업정지 내역(최근매출총액)
	SlVs      string `json:"sl_vs"`       // 영업정지 내역(매출액 대비)
	LsAtn     string `json:"ls_atn"`      // 영업정지 내역(대규모법인여부)
	KrxSttAtn string `json:"krx_stt_atn"` // 영업정지 내역(거래소 의무공시 해당 여부)
	BsnspCn   string `json:"bsnsp_cn"`    // 영업정지 내용
	BsnspRs   string `json:"bsnsp_rs"`    // 영업정지사유
	FtCtp     string `json:"ft_ctp"`      // 향후대책
	BsnspAf   string `json:"bsnsp_af"`    // 영업정지영향
	Bsnspd    string `json:"bsnspd"`      // 영업정지일자
	Bddd      string `json:"bddd"`        // 이사회결의일(결정일)
	OdAAtT    string `json:"od_a_at_t"`   // 사외이사 참석여부(참석)
	OdAAtB    string `json:"od_a_at_b"`   // 사외이사 참석여부(불참)
	AdtAAtn   string `json:"adt_a_atn"`   // 감사(감사위원) 참석여부
}

// BusinessSuspensions 는 영업정지를 조회한다.
func (c *Client) BusinessSuspensions(ctx context.Context, p MaterialParams) ([]BusinessSuspensionItem, error) {
	return httpclient.GetList[BusinessSuspensionItem](ctx, c.http, "/api/bsnSp.json", p.toMap())
}

// RehabilitationItem 은 회생절차 개시신청 (ctrcvsBgrq) 한 건.
type RehabilitationItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Apcnt    string `json:"apcnt"`     // 신청인 (회사와의 관계)
	Cpct     string `json:"cpct"`      // 관할법원
	RqRs     string `json:"rq_rs"`     // 신청사유
	Rqd      string `json:"rqd"`       // 신청일자
	FtCtpSc  string `json:"ft_ctp_sc"` // 향후대책 및 일정
}

// RehabilitationApplications 는 회생절차 개시신청을 조회한다.
func (c *Client) RehabilitationApplications(ctx context.Context, p MaterialParams) ([]RehabilitationItem, error) {
	return httpclient.GetList[RehabilitationItem](ctx, c.http, "/api/ctrcvsBgrq.json", p.toMap())
}

// DissolutionItem 은 해산사유 발생 (dsRsOcr) 한 건.
type DissolutionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DsRs     string `json:"ds_rs"`     // 해산사유
	DsRsd    string `json:"ds_rsd"`    // 해산사유발생일(결정일)
	OdAAtT   string `json:"od_a_at_t"` // 사외이사 참석여부(참석)
	OdAAtB   string `json:"od_a_at_b"` // 사외이사 참석여부(불참)
	AdtAAtn  string `json:"adt_a_atn"` // 감사(감사위원) 참석여부
}

// DissolutionCauses 는 해산사유 발생을 조회한다.
func (c *Client) DissolutionCauses(ctx context.Context, p MaterialParams) ([]DissolutionItem, error) {
	return httpclient.GetList[DissolutionItem](ctx, c.http, "/api/dsRsOcr.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./material/ -v`
Expected: 전체 PASS (Task 1 포함). `go vet ./material/` clean, `gofmt -l material/` no output.

- [ ] **Step 6: Commit**

```bash
git add material/distress.go material/distress_test.go material/testdata/
git commit -m "feat(material): 영업정지·회생절차·해산사유"
```

---

### Task 3: 채권은행 관리절차 개시·중단·소송 (3 엔드포인트)

**Files:**
- Modify: `material/distress.go`, `material/distress_test.go`
- Create: `material/testdata/bnkMngtPcbg.json`, `material/testdata/bnkMngtPcsp.json`, `material/testdata/lwstLg.json`

- [ ] **Step 1: fixture 작성**

`material/testdata/bnkMngtPcbg.json` (실 — 태영건설):
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20250217001661",
            "corp_cls": "Y",
            "corp_code": "00153861",
            "corp_name": "태영건설",
            "mngt_pcbg_dd": "2024년 01월 11일",
            "mngt_int": "(주)태영건설 금융채권자협의회 (주채권은행: 한국산업은행)",
            "mngt_pd": "2024년 01월 11일 ~ 2027년 05월 30일",
            "mngt_rs": "경영정상화",
            "cfd": "2025년 02월 17일"
        }
    ]
}
```

`material/testdata/bnkMngtPcsp.json` (doc/sibling-derived — 실 주요사항보고서 데이터 미발견, docs 필드 스키마 + 실 sibling bnkMngtPcbg 구조 기반의 합성 값):
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20200410000123",
            "corp_cls": "K",
            "corp_code": "00245481",
            "corp_name": "티엔오디앤씨",
            "mngt_pcsp_dd": "2020년 04월 10일",
            "mngt_int": "채권금융기관협의회",
            "sp_rs": "경영정상화 완료",
            "ft_ctp": "자체 경영 정상화 추진",
            "cfd": "2020년 04월 10일"
        }
    ]
}
```

`material/testdata/lwstLg.json` (실 — 올리패스, rq_cn 축약):
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240513000890",
            "corp_cls": "K",
            "corp_code": "01070149",
            "corp_name": "올리패스",
            "icnm": "신주발행금지 등 임시의 지위를 구하는 가처분",
            "ac_ap": "윤용빈",
            "rq_cn": "신주발행무효의 소 본안 판결 확정시까지 신주발행 금지 등을 구하는 가처분 신청",
            "cpct": "수원지방법원",
            "ft_ctp": "법적 절차에 따라 적극적으로 대응할 예정",
            "lgd": "2024년 05월 10일",
            "cfd": "2024년 05월 13일"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `material/distress_test.go` 에 추가:
```go
func TestCreditorBankManagementStart(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bnkMngtPcbg.json": "bnkMngtPcbg.json"})
	items, err := c.CreditorBankManagementStart(context.Background(), MaterialParams{CorpCode: "00153861", BgnDe: "20240101", EndDe: "20251231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "태영건설", items[0].CorpName)
	assert.Equal(t, "경영정상화", items[0].MngtRs)
	assert.Equal(t, "2024년 01월 11일", items[0].MngtPcbgDd)
}

func TestCreditorBankManagementStop(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bnkMngtPcsp.json": "bnkMngtPcsp.json"})
	items, err := c.CreditorBankManagementStop(context.Background(), MaterialParams{CorpCode: "00245481", BgnDe: "20200101", EndDe: "20201231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "경영정상화 완료", items[0].SpRs)
	assert.Equal(t, "2020년 04월 10일", items[0].MngtPcspDd)
}

func TestLawsuits(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/lwstLg.json": "lwstLg.json"})
	items, err := c.Lawsuits(context.Background(), MaterialParams{CorpCode: "01070149", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "신주발행금지 등 임시의 지위를 구하는 가처분", items[0].Icnm)
	assert.Equal(t, "수원지방법원", items[0].Cpct)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./material/ -run 'TestCreditorBank|TestLawsuits' -v`
Expected: FAIL — `undefined: ... CreditorBankManagementStart` 등.

- [ ] **Step 4: 구현** — `material/distress.go` 에 추가:
```go
// CreditorBankMgmtStartItem 은 채권은행 등의 관리절차 개시 (bnkMngtPcbg) 한 건.
type CreditorBankMgmtStartItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcbgDd string `json:"mngt_pcbg_dd"` // 관리절차개시 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	MngtPd     string `json:"mngt_pd"`      // 관리기간
	MngtRs     string `json:"mngt_rs"`      // 관리사유
	Cfd        string `json:"cfd"`          // 확인일자
}

// CreditorBankManagementStart 는 채권은행 등의 관리절차 개시를 조회한다.
func (c *Client) CreditorBankManagementStart(ctx context.Context, p MaterialParams) ([]CreditorBankMgmtStartItem, error) {
	return httpclient.GetList[CreditorBankMgmtStartItem](ctx, c.http, "/api/bnkMngtPcbg.json", p.toMap())
}

// CreditorBankMgmtStopItem 은 채권은행 등의 관리절차 중단 (bnkMngtPcsp) 한 건.
type CreditorBankMgmtStopItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcspDd string `json:"mngt_pcsp_dd"` // 관리절차중단 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	SpRs       string `json:"sp_rs"`        // 중단사유
	FtCtp      string `json:"ft_ctp"`       // 향후대책
	Cfd        string `json:"cfd"`          // 확인일자
}

// CreditorBankManagementStop 은 채권은행 등의 관리절차 중단을 조회한다.
func (c *Client) CreditorBankManagementStop(ctx context.Context, p MaterialParams) ([]CreditorBankMgmtStopItem, error) {
	return httpclient.GetList[CreditorBankMgmtStopItem](ctx, c.http, "/api/bnkMngtPcsp.json", p.toMap())
}

// LawsuitItem 은 소송 등의 제기 (lwstLg) 한 건.
type LawsuitItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Icnm     string `json:"icnm"`      // 사건의 명칭
	AcAp     string `json:"ac_ap"`     // 원고·신청인
	RqCn     string `json:"rq_cn"`     // 청구내용
	Cpct     string `json:"cpct"`      // 관할법원
	FtCtp    string `json:"ft_ctp"`    // 향후대책
	Lgd      string `json:"lgd"`       // 제기일자
	Cfd      string `json:"cfd"`       // 확인일자
}

// Lawsuits 는 소송 등의 제기를 조회한다.
func (c *Client) Lawsuits(ctx context.Context, p MaterialParams) ([]LawsuitItem, error) {
	return httpclient.GetList[LawsuitItem](ctx, c.http, "/api/lwstLg.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./material/ -v`
Expected: 전체 PASS (7개 메서드 + toMap). `go vet ./material/` clean, `gofmt -l material/` no output.

- [ ] **Step 6: Commit**

```bash
git add material/distress.go material/distress_test.go material/testdata/
git commit -m "feat(material): 채권은행 관리절차 개시·중단 + 소송 제기"
```

---

### Task 4: root 와이어링 · README · 통합 테스트 · 최종 검증

**Files:**
- Modify: `client.go`, `client_test.go`, `README.md`, `integration_test.go`

- [ ] **Step 1: root 와이어링 + 테스트** — `client_test.go` 의 `TestNewClient_WiresSubClients` 를 다음으로 교체:
```go
func TestNewClient_WiresSubClients(t *testing.T) {
	c, err := NewClient("KEY", WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)
	assert.NotNil(t, c.Report)
	assert.NotNil(t, c.Ownership)
	assert.NotNil(t, c.Material)
}
```
그리고 `client.go` 수정: import 블록에 `"github.com/kenshin579/opendart/material"` 추가; `Client` struct 의 `Ownership *ownership.Client` 다음 줄에 `Material *material.Client // DS005 주요사항보고서 주요정보` 추가; `NewClient` 의 `c.Ownership = ownership.New(hc)` 다음 줄에 `c.Material = material.New(hc)` 추가.

- [ ] **Step 2: 와이어링 테스트 확인**

Run: `go test . -run TestNewClient -v`
Expected: PASS. `go build ./...` 성공.

- [ ] **Step 3: README 커버리지 갱신** — `README.md` 의 DS004 줄 다음에 추가:
```markdown
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기
```
그리고 `- (예정)` 줄을 다음으로 교체:
```markdown
- (예정) DS005 나머지(증자·감자/사채발행/자기주식/양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```

- [ ] **Step 4: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` 유지; `material` 패키지 타입을 직접 참조하므로 import 블록에 `"github.com/kenshin579/opendart/material"` 추가):
```go
func TestIntegration_DefaultOccurrences(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Material.DefaultOccurrences(context.Background(), material.MaterialParams{
		CorpCode: "00126089", // DH오토넥스 (실제 부도 사례)
		BgnDe:    "20230101",
		EndDe:    "20231231",
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```

- [ ] **Step 5: 통합 빌드 + 최종 검증**

Run:
```bash
go vet -tags integration ./...
go build ./...
go vet ./...
go test ./...
gofmt -l . | grep -v '^scripts/crawl' || echo "clean"
```
Expected: 전체 PASS, integration 미실행, gofmt 신규 파일 차이 없음.

- [ ] **Step 6: Commit**

```bash
git add client.go client_test.go README.md integration_test.go
git commit -m "feat(opendart): wire Material (DS005) + 커버리지/통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** material 패키지+MaterialParams = Task1 · 7개 메서드/struct = Task1(1)+Task2(3)+Task3(3) · root 와이어링 = Task4 · 테스트(toMap+7 fixture) = Task1~3 · README/통합 = Task4. 모두 매핑됨.
- **Type consistency:** `material.{Client,New,MaterialParams(+toMap)}` + 7개 `XItem`/메서드(`DefaultOccurrences`/`BusinessSuspensions`/`RehabilitationApplications`/`DissolutionCauses`/`CreditorBankManagementStart`/`CreditorBankManagementStop`/`Lawsuits`), 시그니처 `(ctx, MaterialParams) ([]XItem, error)` 일관. `httpclient.GetList[T]` 재사용. 필드·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 6개 실 API 사례(DH오토넥스 부도/회생, 태광산업 영업정지, 코리아퍼시픽01호 해산, 태영건설 채권은행개시, 올리패스 소송) + 1개(bnkMngtPcsp 채권은행중단) doc/sibling-derived(실 주요사항보고서 데이터 미발견 — 매우 드문 이벤트; 키/구조는 docs+실 sibling bnkMngtPcbg 와 동일, 값만 합성). 긴 텍스트 필드는 축약.
- **통합 테스트는 실 데이터(DefaultOccurrences/DH오토넥스)** 로 — 합성 fixture(bnkMngtPcsp)에 의존하지 않음.
- **새 추상화 없음:** 기존 httpclient.GetList 재사용. MaterialParams 는 DS005 36개 공통이라 후속 그룹도 재사용.
