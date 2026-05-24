# OpenDART DS005 증자·감자 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 주요사항보고서 주요정보의 증자·감자 4개 API를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material` 패키지의 `MaterialParams`(+toMap) + 공유 `httpclient.GetList[T]` 재사용. 신규 파일 `material/capital.go` 에 4개 item struct + 한 줄 메서드. `client.Material` 는 이미 root 에 와이어링됨(DS005 부실 그룹에서) → **root 변경 없음**.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-capital-design.md`

**검증된 사실 (실 API, 공시검색으로 사례 탐색):** 4개 모두 동일 요청(`corp_code`+`bgn_de`+`end_de`), JSON list. 실 사례 fixture: 유상증자(남양유업 00107598/2023), 무상증자(TYM 00117230/2023), 유무상증자(에코캡 00870481/2022), 감자(코다코 00295857/2024). 금액/주식수 콤마, 비율 %, 빈 값 "-".

**기존 재사용 심볼:** `material.{Client,New,MaterialParams(+toMap)}`, `material/client_test.go` 의 `newTestClient`, `httpclient.GetList[T]`.

---

## File Structure

```
material/
  capital.go        # 4개 메서드 + item struct (신규)
  capital_test.go   # 4개 fixture 테스트 (신규)
  testdata/         # piicDecsn/fricDecsn/pifricDecsn/crDecsn fixture (신규)
README.md           # (수정) DS005 커버리지에 증자·감자
integration_test.go # (수정) PaidInCapitalIncrease 통합 케이스
```

---

### Task 1: 유상증자·무상증자 (PaidInCapitalIncrease / FreeCapitalIncrease)

**Files:**
- Create: `material/capital.go`, `material/capital_test.go`
- Create: `material/testdata/piicDecsn.json`, `material/testdata/fricDecsn.json`

- [ ] **Step 1: fixture 작성**

`material/testdata/piicDecsn.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20230530000410",
            "corp_cls": "Y",
            "corp_code": "00107598",
            "corp_name": "남양유업",
            "nstk_ostk_cnt": "-",
            "nstk_estk_cnt": "33,338",
            "fv_ps": "5,000",
            "bfic_tisstk_ostk": "720,000",
            "bfic_tisstk_estk": "166,662",
            "fdpp_fclt": "-",
            "fdpp_bsninh": "-",
            "fdpp_op": "7,184,339,000",
            "fdpp_dtrp": "-",
            "fdpp_ocsa": "-",
            "fdpp_etc": "-",
            "ic_mthn": "주주우선공모증자",
            "ssl_at": "Y",
            "ssl_bgd": "20230403",
            "ssl_edd": "20230526"
        }
    ]
}
```

`material/testdata/fricDecsn.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20230517000192",
            "corp_cls": "Y",
            "corp_code": "00117230",
            "corp_name": "TYM",
            "nstk_ostk_cnt": "14,580,207",
            "nstk_estk_cnt": "-",
            "fv_ps": "2,500",
            "bfic_tisstk_ostk": "30,470,749",
            "bfic_tisstk_estk": "-",
            "nstk_asstd": "2023년 06월 01일",
            "nstk_ascnt_ps_ostk": "0.5",
            "nstk_ascnt_ps_estk": "-",
            "nstk_dividrk": "2023년 01월 01일",
            "nstk_dlprd": "-",
            "nstk_lstprd": "2023년 06월 16일",
            "bddd": "2023년 03월 30일",
            "od_a_at_t": "3",
            "od_a_at_b": "0",
            "adt_a_atn": "참석"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `material/capital_test.go` (기존 `newTestClient` 재사용):
```go
package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaidInCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/piicDecsn.json": "piicDecsn.json"})
	items, err := c.PaidInCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "남양유업", items[0].CorpName)
	assert.Equal(t, "주주우선공모증자", items[0].IcMthn)
	assert.Equal(t, "7,184,339,000", items[0].FdppOp)
}

func TestFreeCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fricDecsn.json": "fricDecsn.json"})
	items, err := c.FreeCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00117230", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "14,580,207", items[0].NstkOstkCnt)
	assert.Equal(t, "0.5", items[0].NstkAscntPsOstk)
	assert.Equal(t, "2,500", items[0].FvPs)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./material/ -run 'TestPaidInCapitalIncrease|TestFreeCapitalIncrease' -v`
Expected: FAIL — `undefined: ... PaidInCapitalIncrease` 등.

- [ ] **Step 4: 구현** — `material/capital.go`:
```go
package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// PaidInCapitalIncreaseItem 은 유상증자 결정 (piicDecsn) 한 건.
type PaidInCapitalIncreaseItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	NstkOstkCnt    string `json:"nstk_ostk_cnt"`    // 신주의 종류와 수(보통주식)
	NstkEstkCnt    string `json:"nstk_estk_cnt"`    // 신주의 종류와 수(기타주식)
	FvPs           string `json:"fv_ps"`            // 1주당 액면가액 (원)
	BficTisstkOstk string `json:"bfic_tisstk_ostk"` // 증자전 발행주식총수(보통주식)
	BficTisstkEstk string `json:"bfic_tisstk_estk"` // 증자전 발행주식총수(기타주식)
	FdppFclt       string `json:"fdpp_fclt"`        // 자금조달목적(시설자금)
	FdppBsninh     string `json:"fdpp_bsninh"`      // 자금조달목적(영업양수자금)
	FdppOp         string `json:"fdpp_op"`          // 자금조달목적(운영자금)
	FdppDtrp       string `json:"fdpp_dtrp"`        // 자금조달목적(채무상환자금)
	FdppOcsa       string `json:"fdpp_ocsa"`        // 자금조달목적(타법인 증권 취득자금)
	FdppEtc        string `json:"fdpp_etc"`         // 자금조달목적(기타자금)
	IcMthn         string `json:"ic_mthn"`          // 증자방식
	SslAt          string `json:"ssl_at"`           // 공매도 해당여부
	SslBgd         string `json:"ssl_bgd"`          // 공매도 시작일
	SslEdd         string `json:"ssl_edd"`          // 공매도 종료일
}

// PaidInCapitalIncrease 는 유상증자 결정을 조회한다.
func (c *Client) PaidInCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidInCapitalIncreaseItem, error) {
	return httpclient.GetList[PaidInCapitalIncreaseItem](ctx, c.http, "/api/piicDecsn.json", p.toMap())
}

// FreeCapitalIncreaseItem 은 무상증자 결정 (fricDecsn) 한 건.
type FreeCapitalIncreaseItem struct {
	RceptNo         string `json:"rcept_no"`           // 접수번호
	CorpCls         string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`          // 고유번호
	CorpName        string `json:"corp_name"`          // 회사명
	NstkOstkCnt     string `json:"nstk_ostk_cnt"`      // 신주의 종류와 수(보통주식)
	NstkEstkCnt     string `json:"nstk_estk_cnt"`      // 신주의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`              // 1주당 액면가액 (원)
	BficTisstkOstk  string `json:"bfic_tisstk_ostk"`   // 증자전 발행주식총수(보통주식)
	BficTisstkEstk  string `json:"bfic_tisstk_estk"`   // 증자전 발행주식총수(기타주식)
	NstkAsstd       string `json:"nstk_asstd"`         // 신주배정기준일
	NstkAscntPsOstk string `json:"nstk_ascnt_ps_ostk"` // 1주당 신주배정 주식수(보통주식)
	NstkAscntPsEstk string `json:"nstk_ascnt_ps_estk"` // 1주당 신주배정 주식수(기타주식)
	NstkDividrk     string `json:"nstk_dividrk"`       // 신주의 배당기산일
	NstkDlprd       string `json:"nstk_dlprd"`         // 신주권교부예정일
	NstkLstprd      string `json:"nstk_lstprd"`        // 신주의 상장 예정일
	Bddd            string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`          // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`          // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`          // 감사(감사위원) 참석여부
}

// FreeCapitalIncrease 는 무상증자 결정을 조회한다.
func (c *Client) FreeCapitalIncrease(ctx context.Context, p MaterialParams) ([]FreeCapitalIncreaseItem, error) {
	return httpclient.GetList[FreeCapitalIncreaseItem](ctx, c.http, "/api/fricDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./material/ -v`
Expected: 전체 PASS (기존 8 + 신규 2, 회귀 없음). `go vet ./material/` clean, `gofmt -l material/` no output.

- [ ] **Step 6: Commit**

```bash
git add material/capital.go material/capital_test.go material/testdata/
git commit -m "feat(material): 유상증자·무상증자 결정"
```

---

### Task 2: 유무상증자·감자 (PaidFreeCapitalIncrease / CapitalReduction)

**Files:**
- Modify: `material/capital.go`, `material/capital_test.go`
- Create: `material/testdata/pifricDecsn.json`, `material/testdata/crDecsn.json`

- [ ] **Step 1: fixture 작성**

`material/testdata/pifricDecsn.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20220527000497",
            "corp_cls": "K",
            "corp_code": "00870481",
            "corp_name": "에코캡",
            "piic_nstk_ostk_cnt": "5,600,000",
            "piic_nstk_estk_cnt": "-",
            "piic_fv_ps": "100",
            "piic_bfic_tisstk_ostk": "14,558,800",
            "piic_bfic_tisstk_estk": "-",
            "piic_fdpp_fclt": "8,779,904,500",
            "piic_fdpp_bsninh": "-",
            "piic_fdpp_op": "7,424,095,500",
            "piic_fdpp_dtrp": "23,500,000,000",
            "piic_fdpp_ocsa": "-",
            "piic_fdpp_etc": "-",
            "piic_ic_mthn": "주주배정후 실권주 일반공모",
            "fric_nstk_ostk_cnt": "5,874,630",
            "fric_nstk_estk_cnt": "-",
            "fric_fv_ps": "100",
            "fric_bfic_tisstk_ostk": "20,158,800",
            "fric_bfic_tisstk_estk": "-",
            "fric_nstk_asstd": "2022년 06월 14일",
            "fric_nstk_ascnt_ps_ostk": "0.3",
            "fric_nstk_ascnt_ps_estk": "-",
            "fric_nstk_dividrk": "2022년 01월 01일",
            "fric_nstk_dlprd": "-",
            "fric_nstk_lstprd": "2022년 07월 01일",
            "fric_bddd": "2022년 03월 03일",
            "fric_od_a_at_t": "3",
            "fric_od_a_at_b": "-",
            "fric_adt_a_atn": "참석",
            "ssl_at": "Y",
            "ssl_bgd": "20220304",
            "ssl_edd": "20220526"
        }
    ]
}
```

`material/testdata/crDecsn.json`:
```json
{
    "status": "000",
    "message": "정상",
    "list": [
        {
            "rcept_no": "20240531000831",
            "corp_cls": "K",
            "corp_code": "00295857",
            "corp_name": "코다코",
            "crstk_ostk_cnt": "281,338,637",
            "crstk_estk_cnt": "-",
            "fv_ps": "500",
            "bfcr_cpt": "140,669,318,500",
            "atcr_cpt": "7,033,374,500",
            "bfcr_tisstk_ostk": "281,338,637",
            "atcr_tisstk_ostk": "14,066,749",
            "bfcr_tisstk_estk": "-",
            "atcr_tisstk_estk": "-",
            "cr_rt_ostk": "5.00",
            "cr_rt_estk": "-",
            "cr_std": "2024년 05월 23일",
            "cr_mth": "액면가 500원의 보통주식 20주를 1주로 주식병합",
            "cr_rs": "재무구조 개선",
            "crsc_gmtsck_prd": "-",
            "crsc_trnmsppd": "-",
            "crsc_osprpd": "-",
            "crsc_trspprpd": "-",
            "crsc_osprpd_bgd": "-",
            "crsc_osprpd_edd": "-",
            "crsc_trspprpd_bgd": "-",
            "crsc_trspprpd_edd": "-",
            "crsc_nstkdlprd": "-",
            "crsc_nstklstprd": "-",
            "cdobprpd_bgd": "-",
            "cdobprpd_edd": "-",
            "ospr_nstkdl_pl": "국민은행 증권대행부",
            "bddd": "2024년 05월 21일",
            "od_a_at_t": "-",
            "od_a_at_b": "-",
            "adt_a_atn": "-",
            "ftc_stt_atn": "미해당"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `material/capital_test.go` 에 추가:
```go
func TestPaidFreeCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/pifricDecsn.json": "pifricDecsn.json"})
	items, err := c.PaidFreeCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00870481", BgnDe: "20220101", EndDe: "20221231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "주주배정후 실권주 일반공모", items[0].PiicIcMthn)
	assert.Equal(t, "23,500,000,000", items[0].PiicFdppDtrp)
	assert.Equal(t, "0.3", items[0].FricNstkAscntPsOstk)
}

func TestCapitalReduction(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/crDecsn.json": "crDecsn.json"})
	items, err := c.CapitalReduction(context.Background(), MaterialParams{CorpCode: "00295857", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "5.00", items[0].CrRtOstk)
	assert.Equal(t, "140,669,318,500", items[0].BfcrCpt)
	assert.Equal(t, "7,033,374,500", items[0].AtcrCpt)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./material/ -run 'TestPaidFreeCapitalIncrease|TestCapitalReduction' -v`
Expected: FAIL — `undefined: ... PaidFreeCapitalIncrease` 등.

- [ ] **Step 4: 구현** — `material/capital.go` 에 추가:
```go
// PaidFreeCapitalIncreaseItem 은 유무상증자 결정 (pifricDecsn) 한 건. piic_* 유상 / fric_* 무상.
type PaidFreeCapitalIncreaseItem struct {
	RceptNo             string `json:"rcept_no"`                // 접수번호
	CorpCls             string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`               // 고유번호
	CorpName            string `json:"corp_name"`               // 회사명
	PiicNstkOstkCnt     string `json:"piic_nstk_ostk_cnt"`      // 유상증자 신주수(보통주식)
	PiicNstkEstkCnt     string `json:"piic_nstk_estk_cnt"`      // 유상증자 신주수(기타주식)
	PiicFvPs            string `json:"piic_fv_ps"`              // 유상증자 1주당 액면가액
	PiicBficTisstkOstk  string `json:"piic_bfic_tisstk_ostk"`   // 유상증자 증자전 발행총수(보통주식)
	PiicBficTisstkEstk  string `json:"piic_bfic_tisstk_estk"`   // 유상증자 증자전 발행총수(기타주식)
	PiicFdppFclt        string `json:"piic_fdpp_fclt"`          // 유상증자 자금조달목적(시설자금)
	PiicFdppBsninh      string `json:"piic_fdpp_bsninh"`        // 유상증자 자금조달목적(영업양수자금)
	PiicFdppOp          string `json:"piic_fdpp_op"`            // 유상증자 자금조달목적(운영자금)
	PiicFdppDtrp        string `json:"piic_fdpp_dtrp"`          // 유상증자 자금조달목적(채무상환자금)
	PiicFdppOcsa        string `json:"piic_fdpp_ocsa"`          // 유상증자 자금조달목적(타법인 증권 취득자금)
	PiicFdppEtc         string `json:"piic_fdpp_etc"`           // 유상증자 자금조달목적(기타자금)
	PiicIcMthn          string `json:"piic_ic_mthn"`            // 유상증자 증자방식
	FricNstkOstkCnt     string `json:"fric_nstk_ostk_cnt"`      // 무상증자 신주수(보통주식)
	FricNstkEstkCnt     string `json:"fric_nstk_estk_cnt"`      // 무상증자 신주수(기타주식)
	FricFvPs            string `json:"fric_fv_ps"`              // 무상증자 1주당 액면가액
	FricBficTisstkOstk  string `json:"fric_bfic_tisstk_ostk"`   // 무상증자 증자전 발행총수(보통주식)
	FricBficTisstkEstk  string `json:"fric_bfic_tisstk_estk"`   // 무상증자 증자전 발행총수(기타주식)
	FricNstkAsstd       string `json:"fric_nstk_asstd"`         // 무상증자 신주배정기준일
	FricNstkAscntPsOstk string `json:"fric_nstk_ascnt_ps_ostk"` // 무상증자 1주당 신주배정수(보통주식)
	FricNstkAscntPsEstk string `json:"fric_nstk_ascnt_ps_estk"` // 무상증자 1주당 신주배정수(기타주식)
	FricNstkDividrk     string `json:"fric_nstk_dividrk"`       // 무상증자 신주 배당기산일
	FricNstkDlprd       string `json:"fric_nstk_dlprd"`         // 무상증자 신주권교부예정일
	FricNstkLstprd      string `json:"fric_nstk_lstprd"`        // 무상증자 신주 상장예정일
	FricBddd            string `json:"fric_bddd"`               // 무상증자 이사회결의일(결정일)
	FricOdAAtT          string `json:"fric_od_a_at_t"`          // 무상증자 사외이사 참석(참석)
	FricOdAAtB          string `json:"fric_od_a_at_b"`          // 무상증자 사외이사 참석(불참)
	FricAdtAAtn         string `json:"fric_adt_a_atn"`          // 무상증자 감사 참석여부
	SslAt               string `json:"ssl_at"`                  // 공매도 해당여부
	SslBgd              string `json:"ssl_bgd"`                 // 공매도 시작일
	SslEdd              string `json:"ssl_edd"`                 // 공매도 종료일
}

// PaidFreeCapitalIncrease 는 유무상증자 결정을 조회한다.
func (c *Client) PaidFreeCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidFreeCapitalIncreaseItem, error) {
	return httpclient.GetList[PaidFreeCapitalIncreaseItem](ctx, c.http, "/api/pifricDecsn.json", p.toMap())
}

// CapitalReductionItem 은 감자 결정 (crDecsn) 한 건.
type CapitalReductionItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	CorpCls         string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`         // 고유번호
	CorpName        string `json:"corp_name"`         // 회사명
	CrstkOstkCnt    string `json:"crstk_ostk_cnt"`    // 감자주식의 종류와 수(보통주식)
	CrstkEstkCnt    string `json:"crstk_estk_cnt"`    // 감자주식의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`             // 1주당 액면가액 (원)
	BfcrCpt         string `json:"bfcr_cpt"`          // 감자전 자본금 (원)
	AtcrCpt         string `json:"atcr_cpt"`          // 감자후 자본금 (원)
	BfcrTisstkOstk  string `json:"bfcr_tisstk_ostk"`  // 감자전 발행주식수(보통주식)
	AtcrTisstkOstk  string `json:"atcr_tisstk_ostk"`  // 감자후 발행주식수(보통주식)
	BfcrTisstkEstk  string `json:"bfcr_tisstk_estk"`  // 감자전 발행주식수(기타주식)
	AtcrTisstkEstk  string `json:"atcr_tisstk_estk"`  // 감자후 발행주식수(기타주식)
	CrRtOstk        string `json:"cr_rt_ostk"`        // 감자비율(보통주식 %)
	CrRtEstk        string `json:"cr_rt_estk"`        // 감자비율(기타주식 %)
	CrStd           string `json:"cr_std"`            // 감자기준일
	CrMth           string `json:"cr_mth"`            // 감자방법
	CrRs            string `json:"cr_rs"`             // 감자사유
	CrscGmtsckPrd   string `json:"crsc_gmtsck_prd"`   // 감자일정(주주총회 예정일)
	CrscTrnmsppd    string `json:"crsc_trnmsppd"`     // 감자일정(명의개서정지기간)
	CrscOsprpd      string `json:"crsc_osprpd"`       // 감자일정(구주권 제출기간)
	CrscTrspprpd    string `json:"crsc_trspprpd"`     // 감자일정(매매거래 정지예정기간)
	CrscOsprpdBgd   string `json:"crsc_osprpd_bgd"`   // 감자일정(구주권 제출기간 시작일)
	CrscOsprpdEdd   string `json:"crsc_osprpd_edd"`   // 감자일정(구주권 제출기간 종료일)
	CrscTrspprpdBgd string `json:"crsc_trspprpd_bgd"` // 감자일정(매매거래 정지예정기간 시작일)
	CrscTrspprpdEdd string `json:"crsc_trspprpd_edd"` // 감자일정(매매거래 정지예정기간 종료일)
	CrscNstkdlprd   string `json:"crsc_nstkdlprd"`    // 감자일정(신주권교부예정일)
	CrscNstklstprd  string `json:"crsc_nstklstprd"`   // 감자일정(신주상장예정일)
	CdobprpdBgd     string `json:"cdobprpd_bgd"`      // 채권자 이의제출기간(시작일)
	CdobprpdEdd     string `json:"cdobprpd_edd"`      // 채권자 이의제출기간(종료일)
	OsprNstkdlPl    string `json:"ospr_nstkdl_pl"`    // 구주권제출 및 신주권교부장소
	Bddd            string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`         // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`         // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`         // 감사(감사위원) 참석여부
	FtcSttAtn       string `json:"ftc_stt_atn"`       // 공정거래위원회 신고대상 여부
}

// CapitalReduction 은 감자 결정을 조회한다.
func (c *Client) CapitalReduction(ctx context.Context, p MaterialParams) ([]CapitalReductionItem, error) {
	return httpclient.GetList[CapitalReductionItem](ctx, c.http, "/api/crDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./material/ -v`
Expected: 전체 PASS (Task 1 포함 신규 4개 + 기존 8). `go vet ./material/` clean, `gofmt -l material/` no output.

- [ ] **Step 6: Commit**

```bash
git add material/capital.go material/capital_test.go material/testdata/
git commit -m "feat(material): 유무상증자·감자 결정"
```

---

### Task 3: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS005 줄을 다음으로 교체:
```markdown
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정
```
(바로 다음 `- (예정)` 줄은 그대로 둔다.)

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `material` import 유지):
```go
func TestIntegration_PaidInCapitalIncrease(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Material.PaidInCapitalIncrease(context.Background(), material.MaterialParams{
		CorpCode: "00107598", // 남양유업 (실제 유상증자 사례)
		BgnDe:    "20230101",
		EndDe:    "20231231",
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}
```

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
git commit -m "docs(material): DS005 증자·감자 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** 4개 메서드/struct = Task1(2)+Task2(2) · 테스트(fixture) = Task1·2 · README/통합 = Task3. root 와이어링은 DS005 부실 그룹(PR #10)에서 완료(변경 없음). 모두 매핑됨.
- **Type consistency:** 4개 `XItem`/메서드(`PaidInCapitalIncrease`/`FreeCapitalIncrease`/`PaidFreeCapitalIncrease`/`CapitalReduction`), 시그니처 `(ctx, MaterialParams) ([]XItem, error)` 일관. `MaterialParams`/`httpclient.GetList[T]`/`newTestClient` 는 기존 심볼 재사용. 필드·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 4개 모두 실 API 사례(남양유업 유상증자, TYM 무상증자, 에코캡 유무상증자, 코다코 감자) — 공시검색으로 corp_code 탐색. cr_rs 는 가독성 위해 축약("재무구조 개선").
- **새 추상화 없음:** 기존 MaterialParams + httpclient.GetList 재사용. root 변경 없음(client.Material 기존).
