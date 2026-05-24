# OpenDART DS002 임원·직원·보수 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS002 정기보고서 주요정보의 임원·직원·보수 원본 9개 API를 `report` 패키지에 추가한다 (개인별 보수 Ver 2.0 2종은 보류).

**Architecture:** PR #3/#4/#5에서 확립한 `report` 패키지의 제네릭 `getList[T]` + `ReportParams` 패턴을 그대로 재사용한다. 9개 모두 표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답이라 새 추상화가 없다. 새 파일 `report/personnel.go` 에 item struct + 한 줄 메서드를 추가한다. `client.Report` 는 이미 root 에 와이어링되어 있어 root 변경은 없다.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds002-personnel-design.md`

**검증된 사실 (실 API, 삼성전자 00126380 / 2023 / 11011):** 9개 모두 status 000 + list 반환(임원 11, 직원 6, 보수 승인 4, 유형별 4, 개인별 각 5행 등). 숫자 콤마 문자열, 빈 값 "-". 아래 fixture 는 실 응답 첫 항목(임원 현황의 다중행 텍스트는 가독성 위해 단일행으로 정리 — 키/구조 동일). **Ver 2.0 2종(`indvdlByPayV2`/`hmvAuditIndvdlBySttusV2`)은 데이터 미존재로 비범위.**

**기존 재사용 심볼 (PR #3/#4/#5, report 패키지):** `Client`, `ReportParams{CorpCode,BsnsYear,ReprtCode}`, `ReportCode`/`AnnualReport`, `getList[T](ctx, c.http, path, p)`, `report/client_test.go` 의 `newTestClient(t, routes map[string]string) *Client`.

---

## File Structure

```
report/
  personnel.go        # 9개 메서드 + item struct (신규)
  personnel_test.go   # 9개 fixture 테스트 (신규, newTestClient 재사용)
  testdata/           # 9개 실 응답 JSON fixture 추가
README.md             # (수정) DS002 커버리지 + V2 보류 명시
integration_test.go   # (수정) Executives 통합 케이스
```

---

### Task 1: 임원·직원·미등기임원보수·사외이사 변동 (4 엔드포인트)

**Files:**
- Create: `report/personnel.go`, `report/personnel_test.go`
- Create: `report/testdata/exctvSttus.json`, `report/testdata/empSttus.json`, `report/testdata/unrstExctvMendngSttus.json`, `report/testdata/outcmpnyDrctrNdChangeSttus.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/exctvSttus.json`:
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
            "nm": "한종희",
            "sexdstn": "남",
            "birth_ym": "1962년 03월",
            "ofcps": "부회장",
            "rgist_exctv_at": "사내이사",
            "fte_at": "상근",
            "chrg_job": "대표이사 (DX 부문 경영전반 총괄)",
            "main_career": "삼성전자 DX부문장",
            "mxmm_shrholdr_relate": "계열회사 임원",
            "hffc_pd": "46개월",
            "tenure_end_on": "2026년 03월 17일",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/empSttus.json`:
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
            "sexdstn": "남",
            "fo_bbm": "DX",
            "reform_bfe_emp_co_rgllbr": "-",
            "reform_bfe_emp_co_cnttk": "-",
            "reform_bfe_emp_co_etc": "-",
            "rgllbr_co": "37,962",
            "rgllbr_abacpt_labrr_co": "-",
            "cnttk_co": "324",
            "cnttk_abacpt_labrr_co": "-",
            "sm": "38,286",
            "avrg_cnwk_sdytrn": "16.5",
            "fyer_salary_totamt": "-",
            "jan_salary_am": "-",
            "rm": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/unrstExctvMendngSttus.json`:
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
            "se": "미등기임원",
            "fyer_salary_totamt": "705,410,000,000",
            "jan_salary_am": "726,000,000",
            "nmpr": "1,015",
            "rm": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/outcmpnyDrctrNdChangeSttus.json`:
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
            "drctr_co": "11",
            "otcmp_drctr_co": "6",
            "apnt": "-",
            "rlsofc": "-",
            "mdstrm_resig": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/personnel_test.go` (기존 `newTestClient` 재사용):
```go
package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutives(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/exctvSttus.json": "exctvSttus.json"})
	items, err := c.Executives(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "한종희", items[0].Nm)
	assert.Equal(t, "부회장", items[0].Ofcps)
	assert.Equal(t, "사내이사", items[0].RgistExctvAt)
}

func TestEmployees(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/empSttus.json": "empSttus.json"})
	items, err := c.Employees(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DX", items[0].FoBbm)
	assert.Equal(t, "37,962", items[0].RgllbrCo)
	assert.Equal(t, "38,286", items[0].Sm)
}

func TestUnregisteredExecutiveCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/unrstExctvMendngSttus.json": "unrstExctvMendngSttus.json"})
	items, err := c.UnregisteredExecutiveCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미등기임원", items[0].Se)
	assert.Equal(t, "1,015", items[0].Nmpr)
}

func TestOutsideDirectorChanges(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/outcmpnyDrctrNdChangeSttus.json": "outcmpnyDrctrNdChangeSttus.json"})
	items, err := c.OutsideDirectorChanges(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "11", items[0].DrctrCo)
	assert.Equal(t, "6", items[0].OtcmpDrctrCo)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestExecutives|TestEmployees|TestUnregisteredExecutiveCompensation|TestOutsideDirectorChanges' -v`
Expected: FAIL — `undefined: ... Executives` 등.

- [ ] **Step 4: 구현** — `report/personnel.go`:
```go
package report

import "context"

// ExecutiveItem 은 임원 현황 (exctvSttus) 한 건.
type ExecutiveItem struct {
	RceptNo            string `json:"rcept_no"`             // 접수번호
	CorpCls            string `json:"corp_cls"`             // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`            // 고유번호
	CorpName           string `json:"corp_name"`            // 법인명
	Nm                 string `json:"nm"`                   // 성명
	Sexdstn            string `json:"sexdstn"`              // 성별
	BirthYm            string `json:"birth_ym"`             // 출생 년월
	Ofcps              string `json:"ofcps"`                // 직위
	RgistExctvAt       string `json:"rgist_exctv_at"`       // 등기 임원 여부
	FteAt              string `json:"fte_at"`               // 상근 여부
	ChrgJob            string `json:"chrg_job"`             // 담당 업무
	MainCareer         string `json:"main_career"`          // 주요 경력
	MxmmShrholdrRelate string `json:"mxmm_shrholdr_relate"`  // 최대 주주 관계
	HffcPd             string `json:"hffc_pd"`              // 재직 기간
	TenureEndOn        string `json:"tenure_end_on"`        // 임기 만료 일
	StlmDt             string `json:"stlm_dt"`              // 결산기준일
}

// Executives 는 임원 현황을 조회한다.
func (c *Client) Executives(ctx context.Context, p ReportParams) ([]ExecutiveItem, error) {
	return getList[ExecutiveItem](ctx, c.http, "/api/exctvSttus.json", p)
}

// EmployeeItem 은 직원 현황 (empSttus) 한 건.
type EmployeeItem struct {
	RceptNo              string `json:"rcept_no"`                 // 접수번호
	CorpCls              string `json:"corp_cls"`                 // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`                // 고유번호
	CorpName             string `json:"corp_name"`                // 법인명
	FoBbm                string `json:"fo_bbm"`                   // 사업부문
	Sexdstn              string `json:"sexdstn"`                  // 성별
	ReformBfeEmpCoRgllbr string `json:"reform_bfe_emp_co_rgllbr"` // 개정 전 직원 수 정규직
	ReformBfeEmpCoCnttk  string `json:"reform_bfe_emp_co_cnttk"`  // 개정 전 직원 수 계약직
	ReformBfeEmpCoEtc    string `json:"reform_bfe_emp_co_etc"`    // 개정 전 직원 수 기타
	RgllbrCo             string `json:"rgllbr_co"`                // 정규직 수
	RgllbrAbacptLabrrCo  string `json:"rgllbr_abacpt_labrr_co"`   // 정규직 단시간 근로자 수
	CnttkCo              string `json:"cnttk_co"`                 // 계약직 수
	CnttkAbacptLabrrCo   string `json:"cnttk_abacpt_labrr_co"`    // 계약직 단시간 근로자 수
	Sm                   string `json:"sm"`                       // 합계
	AvrgCnwkSdytrn       string `json:"avrg_cnwk_sdytrn"`         // 평균 근속 연수
	FyerSalaryTotamt     string `json:"fyer_salary_totamt"`       // 연간 급여 총액
	JanSalaryAm          string `json:"jan_salary_am"`            // 1인평균 급여 액
	Rm                   string `json:"rm"`                       // 비고
	StlmDt               string `json:"stlm_dt"`                  // 결산기준일
}

// Employees 는 직원 현황을 조회한다.
func (c *Client) Employees(ctx context.Context, p ReportParams) ([]EmployeeItem, error) {
	return getList[EmployeeItem](ctx, c.http, "/api/empSttus.json", p)
}

// UnregisteredExecutiveCompensationItem 은 미등기임원 보수현황 (unrstExctvMendngSttus) 한 건.
type UnregisteredExecutiveCompensationItem struct {
	RceptNo          string `json:"rcept_no"`           // 접수번호
	CorpCls          string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`          // 고유번호
	CorpName         string `json:"corp_name"`          // 회사명
	Se               string `json:"se"`                 // 구분
	FyerSalaryTotamt string `json:"fyer_salary_totamt"`  // 연간급여 총액
	JanSalaryAm      string `json:"jan_salary_am"`       // 1인평균 급여액
	Nmpr             string `json:"nmpr"`               // 인원수
	Rm               string `json:"rm"`                 // 비고
	StlmDt           string `json:"stlm_dt"`            // 결산기준일
}

// UnregisteredExecutiveCompensation 은 미등기임원 보수현황을 조회한다.
func (c *Client) UnregisteredExecutiveCompensation(ctx context.Context, p ReportParams) ([]UnregisteredExecutiveCompensationItem, error) {
	return getList[UnregisteredExecutiveCompensationItem](ctx, c.http, "/api/unrstExctvMendngSttus.json", p)
}

// OutsideDirectorChangeItem 은 사외이사 및 그 변동현황 (outcmpnyDrctrNdChangeSttus) 한 건.
type OutsideDirectorChangeItem struct {
	RceptNo      string `json:"rcept_no"`       // 접수번호
	CorpCls      string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`      // 고유번호
	CorpName     string `json:"corp_name"`      // 회사명
	DrctrCo      string `json:"drctr_co"`       // 이사의 수
	OtcmpDrctrCo string `json:"otcmp_drctr_co"` // 사외이사 수
	Apnt         string `json:"apnt"`           // 사외이사 변동현황(선임)
	Rlsofc       string `json:"rlsofc"`         // 사외이사 변동현황(해임)
	MdstrmResig  string `json:"mdstrm_resig"`   // 사외이사 변동현황(중도퇴임)
	StlmDt       string `json:"stlm_dt"`        // 결산기준일
}

// OutsideDirectorChanges 는 사외이사 및 그 변동현황을 조회한다.
func (c *Client) OutsideDirectorChanges(ctx context.Context, p ReportParams) ([]OutsideDirectorChangeItem, error) {
	return getList[OutsideDirectorChangeItem](ctx, c.http, "/api/outcmpnyDrctrNdChangeSttus.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (기존 21개 + 신규 4개, 회귀 없음). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/personnel.go report/personnel_test.go report/testdata/
git commit -m "feat(report): 임원·직원·미등기임원보수·사외이사 변동현황"
```

---

### Task 2: 이사·감사 전체 보수 3종 + 개인별 보수 2종 (5 엔드포인트)

**Files:**
- Modify: `report/personnel.go` (5개 메서드+struct 추가)
- Modify: `report/personnel_test.go` (5개 테스트 추가)
- Create: `report/testdata/drctrAdtAllMendngSttusGmtsckConfmAmount.json`, `report/testdata/hmvAuditAllSttus.json`, `report/testdata/drctrAdtAllMendngSttusMendngPymntamtTyCl.json`, `report/testdata/hmvAuditIndvdlBySttus.json`, `report/testdata/indvdlByPay.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/drctrAdtAllMendngSttusGmtsckConfmAmount.json`:
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
            "se": "등기이사",
            "nmpr": "5",
            "gmtsck_confm_amount": "-",
            "rm": "-",
            "stlm_dt": "2023-12-31",
            "fscl_year": "-"
        }
    ]
}
```

`report/testdata/hmvAuditAllSttus.json`:
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
            "nmpr": "11",
            "mendng_totamt": "23,227,000,000",
            "jan_avrg_mendng_am": "2,112,000,000",
            "rm": "-",
            "stlm_dt": "2023-12-31",
            "fscl_year": "-",
            "stk_bsd_pd_mendng_totamt": "-",
            "stk_opt_exrcsbl_qty": "-",
            "stk_opt_unexrcsbl_qty": "-",
            "stk_opt_rmn_blce": "-",
            "othr_stk_bsd_cmpn_unpyd_qty": "-",
            "othr_stk_bsd_cmpn_mkt_vl": "-"
        }
    ]
}
```

`report/testdata/drctrAdtAllMendngSttusMendngPymntamtTyCl.json`:
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
            "se": "등기이사(사외이사, 감사위원회 위원 제외)",
            "nmpr": "5",
            "pymnt_totamt": "22,009,000,000",
            "psn1_avrg_pymntamt": "4,402,000,000",
            "rm": "-",
            "stlm_dt": "2023-12-31",
            "fscl_year": "-",
            "stk_bsd_pd_mendng_totamt": "-",
            "stk_opt_exrcsbl_qty": "-",
            "stk_opt_unexrcsbl_qty": "-",
            "stk_opt_rmn_blce": "-",
            "othr_stk_bsd_cmpn_unpyd_qty": "-",
            "othr_stk_bsd_cmpn_mkt_vl": "-"
        }
    ]
}
```

`report/testdata/hmvAuditIndvdlBySttus.json`:
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
            "nm": "경계현",
            "ofcps": "대표이사",
            "mendng_totamt": "2,403,000,000",
            "mendng_totamt_ct_incls_mendng": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/indvdlByPay.json`:
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
            "nm": "김기남",
            "ofcps": "고문",
            "mendng_totamt": "17,265,000,000",
            "mendng_totamt_ct_incls_mendng": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/personnel_test.go` 에 추가:
```go
func TestDirectorAuditorApprovedCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json": "drctrAdtAllMendngSttusGmtsckConfmAmount.json"})
	items, err := c.DirectorAuditorApprovedCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "등기이사", items[0].Se)
	assert.Equal(t, "5", items[0].Nmpr)
}

func TestDirectorAuditorTotalCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hmvAuditAllSttus.json": "hmvAuditAllSttus.json"})
	items, err := c.DirectorAuditorTotalCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "11", items[0].Nmpr)
	assert.Equal(t, "23,227,000,000", items[0].MendngTotamt)
}

func TestDirectorAuditorCompensationByType(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json": "drctrAdtAllMendngSttusMendngPymntamtTyCl.json"})
	items, err := c.DirectorAuditorCompensationByType(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "22,009,000,000", items[0].PymntTotamt)
	assert.Equal(t, "4,402,000,000", items[0].Psn1AvrgPymntamt)
}

func TestIndividualDirectorAuditorCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hmvAuditIndvdlBySttus.json": "hmvAuditIndvdlBySttus.json"})
	items, err := c.IndividualDirectorAuditorCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "경계현", items[0].Nm)
	assert.Equal(t, "2,403,000,000", items[0].MendngTotamt)
}

func TestIndividualTop5Compensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/indvdlByPay.json": "indvdlByPay.json"})
	items, err := c.IndividualTop5Compensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "김기남", items[0].Nm)
	assert.Equal(t, "고문", items[0].Ofcps)
	assert.Equal(t, "17,265,000,000", items[0].MendngTotamt)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestDirectorAuditor|TestIndividual' -v`
Expected: FAIL — `undefined: ... DirectorAuditorApprovedCompensation` 등.

- [ ] **Step 4: 구현** — `report/personnel.go` 에 추가:
```go
// DirectorAuditorApprovedCompensationItem 은 이사·감사 전체의 보수현황(주주총회 승인금액)
// (drctrAdtAllMendngSttusGmtsckConfmAmount) 한 건.
type DirectorAuditorApprovedCompensationItem struct {
	RceptNo           string `json:"rcept_no"`            // 접수번호
	CorpCls           string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode          string `json:"corp_code"`           // 고유번호
	CorpName          string `json:"corp_name"`           // 회사명
	Se                string `json:"se"`                  // 구분
	Nmpr              string `json:"nmpr"`                // 인원수
	GmtsckConfmAmount string `json:"gmtsck_confm_amount"`  // 주주총회 승인금액
	Rm                string `json:"rm"`                  // 비고
	StlmDt            string `json:"stlm_dt"`             // 결산기준일
	FsclYear          string `json:"fscl_year"`           // 사업연도
}

// DirectorAuditorApprovedCompensation 은 이사·감사 전체의 보수현황(주주총회 승인금액)을 조회한다.
func (c *Client) DirectorAuditorApprovedCompensation(ctx context.Context, p ReportParams) ([]DirectorAuditorApprovedCompensationItem, error) {
	return getList[DirectorAuditorApprovedCompensationItem](ctx, c.http, "/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json", p)
}

// DirectorAuditorTotalCompensationItem 은 이사·감사 전체의 보수현황(보수지급금액 - 이사·감사 전체)
// (hmvAuditAllSttus) 한 건.
type DirectorAuditorTotalCompensationItem struct {
	RceptNo                string `json:"rcept_no"`                   // 접수번호
	CorpCls                string `json:"corp_cls"`                   // 법인구분 (Y/K/N/E)
	CorpCode               string `json:"corp_code"`                  // 고유번호
	CorpName               string `json:"corp_name"`                  // 법인명
	Nmpr                   string `json:"nmpr"`                       // 인원수
	MendngTotamt           string `json:"mendng_totamt"`              // 보수 총액
	JanAvrgMendngAm        string `json:"jan_avrg_mendng_am"`         // 1인 평균 보수 액
	Rm                     string `json:"rm"`                         // 비고
	StlmDt                 string `json:"stlm_dt"`                    // 결산기준일
	FsclYear               string `json:"fscl_year"`                  // 사업연도
	StkBsdPdMendngTotamt   string `json:"stk_bsd_pd_mendng_totamt"`   // 보수총액 중 주식기준보상 지급액
	StkOptExrcsblQty       string `json:"stk_opt_exrcsbl_qty"`        // 주식매수선택권 행사가능수량
	StkOptUnexrcsblQty     string `json:"stk_opt_unexrcsbl_qty"`      // 주식매수선택권 행사불가수량
	StkOptRmnBlce          string `json:"stk_opt_rmn_blce"`           // 주식매수선택권 잔여금액
	OthrStkBsdCmpnUnpydQty string `json:"othr_stk_bsd_cmpn_unpyd_qty"` // 그 외 주식기준 보상 미지급수량
	OthrStkBsdCmpnMktVl    string `json:"othr_stk_bsd_cmpn_mkt_vl"`   // 그 외 주식기준 보상 시장가치
}

// DirectorAuditorTotalCompensation 은 이사·감사 전체의 보수현황(보수지급금액 - 이사·감사 전체)을 조회한다.
func (c *Client) DirectorAuditorTotalCompensation(ctx context.Context, p ReportParams) ([]DirectorAuditorTotalCompensationItem, error) {
	return getList[DirectorAuditorTotalCompensationItem](ctx, c.http, "/api/hmvAuditAllSttus.json", p)
}

// DirectorAuditorCompensationByTypeItem 은 이사·감사 전체의 보수현황(보수지급금액 - 유형별)
// (drctrAdtAllMendngSttusMendngPymntamtTyCl) 한 건.
type DirectorAuditorCompensationByTypeItem struct {
	RceptNo                string `json:"rcept_no"`                   // 접수번호
	CorpCls                string `json:"corp_cls"`                   // 법인구분 (Y/K/N/E)
	CorpCode               string `json:"corp_code"`                  // 고유번호
	CorpName               string `json:"corp_name"`                  // 회사명
	Se                     string `json:"se"`                         // 구분 (등기이사/사외이사/감사위원회 위원 등)
	Nmpr                   string `json:"nmpr"`                       // 인원수
	PymntTotamt            string `json:"pymnt_totamt"`               // 보수총액
	Psn1AvrgPymntamt       string `json:"psn1_avrg_pymntamt"`         // 1인당 평균보수액
	Rm                     string `json:"rm"`                         // 비고
	StlmDt                 string `json:"stlm_dt"`                    // 결산기준일
	FsclYear               string `json:"fscl_year"`                  // 사업연도
	StkBsdPdMendngTotamt   string `json:"stk_bsd_pd_mendng_totamt"`   // 보수총액 중 주식기준보상 지급액
	StkOptExrcsblQty       string `json:"stk_opt_exrcsbl_qty"`        // 주식매수선택권 행사가능수량
	StkOptUnexrcsblQty     string `json:"stk_opt_unexrcsbl_qty"`      // 주식매수선택권 행사불가수량
	StkOptRmnBlce          string `json:"stk_opt_rmn_blce"`           // 주식매수선택권 잔여금액
	OthrStkBsdCmpnUnpydQty string `json:"othr_stk_bsd_cmpn_unpyd_qty"` // 그 외 주식기준 보상 미지급수량
	OthrStkBsdCmpnMktVl    string `json:"othr_stk_bsd_cmpn_mkt_vl"`   // 그 외 주식기준 보상 시장가치
}

// DirectorAuditorCompensationByType 은 이사·감사 전체의 보수현황(보수지급금액 - 유형별)을 조회한다.
func (c *Client) DirectorAuditorCompensationByType(ctx context.Context, p ReportParams) ([]DirectorAuditorCompensationByTypeItem, error) {
	return getList[DirectorAuditorCompensationByTypeItem](ctx, c.http, "/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json", p)
}

// IndividualDirectorAuditorCompensationItem 은 이사·감사의 개인별 보수현황(5억원 이상)
// (hmvAuditIndvdlBySttus) 한 건.
type IndividualDirectorAuditorCompensationItem struct {
	RceptNo                   string `json:"rcept_no"`                      // 접수번호
	CorpCls                   string `json:"corp_cls"`                      // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                     // 고유번호
	CorpName                  string `json:"corp_name"`                     // 법인명
	Nm                        string `json:"nm"`                            // 이름
	Ofcps                     string `json:"ofcps"`                         // 직위
	MendngTotamt              string `json:"mendng_totamt"`                 // 보수 총액
	MendngTotamtCtInclsMendng string `json:"mendng_totamt_ct_incls_mendng"` // 보수 총액 비 포함 보수
	StlmDt                    string `json:"stlm_dt"`                       // 결산기준일
}

// IndividualDirectorAuditorCompensation 은 이사·감사의 개인별 보수현황(5억원 이상)을 조회한다.
func (c *Client) IndividualDirectorAuditorCompensation(ctx context.Context, p ReportParams) ([]IndividualDirectorAuditorCompensationItem, error) {
	return getList[IndividualDirectorAuditorCompensationItem](ctx, c.http, "/api/hmvAuditIndvdlBySttus.json", p)
}

// IndividualTop5CompensationItem 은 개인별 보수지급 금액(5억이상 상위5인) (indvdlByPay) 한 건.
type IndividualTop5CompensationItem struct {
	RceptNo                   string `json:"rcept_no"`                      // 접수번호
	CorpCls                   string `json:"corp_cls"`                      // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                     // 고유번호
	CorpName                  string `json:"corp_name"`                     // 법인명
	Nm                        string `json:"nm"`                            // 이름
	Ofcps                     string `json:"ofcps"`                         // 직위
	MendngTotamt              string `json:"mendng_totamt"`                 // 보수 총액
	MendngTotamtCtInclsMendng string `json:"mendng_totamt_ct_incls_mendng"` // 보수 총액 비 포함 보수
	StlmDt                    string `json:"stlm_dt"`                       // 결산기준일
}

// IndividualTop5Compensation 은 개인별 보수지급 금액(5억이상 상위5인)을 조회한다.
func (c *Client) IndividualTop5Compensation(ctx context.Context, p ReportParams) ([]IndividualTop5CompensationItem, error) {
	return getList[IndividualTop5CompensationItem](ctx, c.http, "/api/indvdlByPay.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 1 포함 신규 9개 + 기존). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/personnel.go report/personnel_test.go report/testdata/
git commit -m "feat(report): 이사·감사 보수현황 3종 + 개인별 보수 2종"
```

---

### Task 3: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS002 줄과 그 다음 "(예정)" 줄을 다음 2줄로 교체:
```markdown
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황 · 증권 발행실적 · 미상환 잔액(회사채/기업어음/단기사채/신종자본증권/조건부자본증권) · 감사의견 · 감사/비감사용역 · 타법인 출자 · 공모/사모자금 사용내역 · 임원/직원 현황 · 임원·이사·감사 보수현황
- (예정) DS002 개인별 보수 Ver2.0 2종 · DS003~DS006
```

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `report` import 유지):
```go
func TestIntegration_Executives(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.Executives(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
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
git commit -m "docs(report): DS002 임원·직원·보수 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** 9개 메서드+struct = Task1(4)+Task2(5) · 테스트(fixture) = Task1·2 · README 커버리지 + V2 보류 명시 = Task3 · 통합 테스트 = Task3. root 와이어링은 PR #3 에서 완료(변경 없음). V2 2종은 비범위(spec 명시). 모두 매핑됨.
- **Type consistency:** 9개 `XItem`/메서드(`Executives`/`Employees`/`UnregisteredExecutiveCompensation`/`OutsideDirectorChanges`/`DirectorAuditorApprovedCompensation`/`DirectorAuditorTotalCompensation`/`DirectorAuditorCompensationByType`/`IndividualDirectorAuditorCompensation`/`IndividualTop5Compensation`), 시그니처 `(ctx, ReportParams) ([]XItem, error)` 일관. `getList[T]`/`ReportParams`/`AnnualReport`/`newTestClient` 는 기존 심볼 재사용. 필드명·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 9개 모두 실 API(삼성전자/2023/사업보고서) 응답 첫 항목. 임원 현황 다중행 텍스트는 단일행으로 정리(키/구조 동일). 숫자 콤마 string, 빈 값 "-".
- **새 추상화 없음:** 기존 제네릭 getList 재사용만. root 변경 없음(client.Report 기존).
- **IndividualDirectorAuditorCompensationItem 과 IndividualTop5CompensationItem 은 필드 동일하나 의미·엔드포인트가 달라 별도 타입 유지**(spec 명시).
