# OpenDART DS002 증권 발행·미상환 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS002 정기보고서 주요정보의 증권 발행·미상환 6개 API를 `report` 패키지에 추가한다.

**Architecture:** PR #3에서 확립한 `report` 패키지의 제네릭 `getList[T]` + `ReportParams` 패턴을 그대로 재사용한다. 6개 모두 표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답이라 새 추상화가 없다. 새 파일 `report/securities.go` 에 item struct + 한 줄 메서드를 추가한다. `client.Report` 는 이미 root 에 와이어링되어 있어 root 변경은 없다.

**Tech Stack:** Go 1.25+ (제네릭), 표준 net/http (internal/httpclient 재사용), testify.

**Spec:** `docs/superpowers/specs/2026-05-24-opendart-ds002-securities-design.md`

**검증된 사실 (실 API, 삼성전자 00126380 / 2023 / 11011):** 6개 모두 status 000 + list 반환(발행실적 16행, 미상환 잔액 각 3행). 숫자 콤마 문자열, 빈 값 "-". 아래 fixture 는 실 응답 첫 항목.

**기존 재사용 심볼 (PR #3, report 패키지):** `Client`, `ReportParams{CorpCode,BsnsYear,ReprtCode}`, `ReportCode`/`AnnualReport`, `getList[T](ctx, c.http, path, p)`, 그리고 `report/client_test.go` 의 `newTestClient(t, routes map[string]string) *Client`.

---

## File Structure

```
report/
  securities.go        # 6개 메서드 + item struct (신규)
  securities_test.go   # 6개 fixture 테스트 (신규, newTestClient 재사용)
  testdata/            # 6개 실 응답 JSON fixture 추가
README.md              # (수정) DS002 커버리지에 증권 발행·미상환
integration_test.go    # (수정) DebtSecuritiesIssuance 통합 케이스
```

---

### Task 1: 채무증권 발행실적·회사채·기업어음 미상환 (3 엔드포인트)

**Files:**
- Create: `report/securities.go`, `report/securities_test.go`
- Create: `report/testdata/detScritsIsuAcmslt.json`, `report/testdata/cprndNrdmpBlce.json`, `report/testdata/entrprsBilScritsNrdmpBlce.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/detScritsIsuAcmslt.json`:
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
            "isu_cmpny": "삼성전자㈜",
            "scrits_knd_nm": "회사채",
            "isu_mth_nm": "공모",
            "isu_de": "1997.10.02",
            "facvalu_totamt": "128,940,000,000",
            "intrt": "7.7%",
            "evl_grad_instt": "Aa2 (Moody's), AA- (S&P)",
            "mtd": "2027.10.01",
            "repy_at": "일부상환",
            "mngt_cmpny": "Goldman Sachs 등",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/cprndNrdmpBlce.json`:
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
            "sm": "541,548,000,000",
            "remndr_exprtn1": "미상환잔액",
            "remndr_exprtn2": "공모",
            "yy1_below": "6,447,000,000",
            "yy1_excess_yy2_below": "522,207,000,000",
            "yy2_excess_yy3_below": "6,447,000,000",
            "yy3_excess_yy4_below": "6,447,000,000",
            "yy4_excess_yy5_below": "-",
            "yy5_excess_yy10_below": "-",
            "yy10_excess": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/entrprsBilScritsNrdmpBlce.json`:
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
            "sm": "-",
            "remndr_exprtn1": "미상환잔액",
            "remndr_exprtn2": "공모",
            "de10_below": "-",
            "de10_excess_de30_below": "-",
            "de30_excess_de90_below": "-",
            "de90_excess_de180_below": "-",
            "de180_excess_yy1_below": "-",
            "yy1_excess_yy2_below": "-",
            "yy2_excess_yy3_below": "-",
            "yy3_excess": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `report/securities_test.go` (기존 `newTestClient` 재사용):
```go
package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebtSecuritiesIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/detScritsIsuAcmslt.json": "detScritsIsuAcmslt.json"})
	items, err := c.DebtSecuritiesIssuance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "회사채", items[0].ScritsKndNm)
	assert.Equal(t, "128,940,000,000", items[0].FacvaluTotamt)
	assert.Equal(t, "2027.10.01", items[0].Mtd)
}

func TestCorporateBondBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cprndNrdmpBlce.json": "cprndNrdmpBlce.json"})
	items, err := c.CorporateBondBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "541,548,000,000", items[0].Sm)
	assert.Equal(t, "522,207,000,000", items[0].Yy1ExcessYy2Below)
}

func TestCommercialPaperBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/entrprsBilScritsNrdmpBlce.json": "entrprsBilScritsNrdmpBlce.json"})
	items, err := c.CommercialPaperBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미상환잔액", items[0].RemndrExprtn1)
	assert.Equal(t, "-", items[0].Yy3Excess)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestDebtSecuritiesIssuance|TestCorporateBondBalance|TestCommercialPaperBalance' -v`
Expected: FAIL — `undefined: ... DebtSecuritiesIssuance` 등.

- [ ] **Step 4: 구현** — `report/securities.go`:
```go
package report

import "context"

// DebtSecuritiesIssuanceItem 은 채무증권 발행실적 (detScritsIsuAcmslt) 한 건.
type DebtSecuritiesIssuanceItem struct {
	RceptNo       string `json:"rcept_no"`       // 접수번호
	CorpCls       string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`      // 고유번호
	CorpName      string `json:"corp_name"`      // 회사명
	IsuCmpny      string `json:"isu_cmpny"`      // 발행회사
	ScritsKndNm   string `json:"scrits_knd_nm"`  // 증권종류
	IsuMthNm      string `json:"isu_mth_nm"`     // 발행방법
	IsuDe         string `json:"isu_de"`         // 발행일자
	FacvaluTotamt string `json:"facvalu_totamt"` // 권면(전자등록)총액
	Intrt         string `json:"intrt"`          // 이자율
	EvlGradInstt  string `json:"evl_grad_instt"` // 평가등급(평가기관)
	Mtd           string `json:"mtd"`            // 만기일
	RepyAt        string `json:"repy_at"`        // 상환여부
	MngtCmpny     string `json:"mngt_cmpny"`     // 주관회사
	StlmDt        string `json:"stlm_dt"`        // 결산기준일
}

// DebtSecuritiesIssuance 는 채무증권 발행실적을 조회한다.
func (c *Client) DebtSecuritiesIssuance(ctx context.Context, p ReportParams) ([]DebtSecuritiesIssuanceItem, error) {
	return getList[DebtSecuritiesIssuanceItem](ctx, c.http, "/api/detScritsIsuAcmslt.json", p)
}

// CorporateBondBalanceItem 은 회사채 미상환 잔액 (cprndNrdmpBlce) 한 건.
type CorporateBondBalanceItem struct {
	RceptNo            string `json:"rcept_no"`              // 접수번호
	CorpCls            string `json:"corp_cls"`              // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`             // 고유번호
	CorpName           string `json:"corp_name"`             // 회사명
	RemndrExprtn1      string `json:"remndr_exprtn1"`        // 잔여만기 (구분1)
	RemndrExprtn2      string `json:"remndr_exprtn2"`        // 잔여만기 (구분2)
	Yy1Below           string `json:"yy1_below"`             // 1년 이하
	Yy1ExcessYy2Below  string `json:"yy1_excess_yy2_below"`  // 1년초과 2년이하
	Yy2ExcessYy3Below  string `json:"yy2_excess_yy3_below"`  // 2년초과 3년이하
	Yy3ExcessYy4Below  string `json:"yy3_excess_yy4_below"`  // 3년초과 4년이하
	Yy4ExcessYy5Below  string `json:"yy4_excess_yy5_below"`  // 4년초과 5년이하
	Yy5ExcessYy10Below string `json:"yy5_excess_yy10_below"` // 5년초과 10년이하
	Yy10Excess         string `json:"yy10_excess"`           // 10년초과
	Sm                 string `json:"sm"`                    // 합계
	StlmDt             string `json:"stlm_dt"`               // 결산기준일
}

// CorporateBondBalance 는 회사채 미상환 잔액을 조회한다.
func (c *Client) CorporateBondBalance(ctx context.Context, p ReportParams) ([]CorporateBondBalanceItem, error) {
	return getList[CorporateBondBalanceItem](ctx, c.http, "/api/cprndNrdmpBlce.json", p)
}

// CommercialPaperBalanceItem 은 기업어음증권 미상환 잔액 (entrprsBilScritsNrdmpBlce) 한 건.
type CommercialPaperBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Yy1ExcessYy2Below    string `json:"yy1_excess_yy2_below"`    // 1년초과 2년이하
	Yy2ExcessYy3Below    string `json:"yy2_excess_yy3_below"`    // 2년초과 3년이하
	Yy3Excess            string `json:"yy3_excess"`              // 3년 초과
	Sm                   string `json:"sm"`                      // 합계
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// CommercialPaperBalance 는 기업어음증권 미상환 잔액을 조회한다.
func (c *Client) CommercialPaperBalance(ctx context.Context, p ReportParams) ([]CommercialPaperBalanceItem, error) {
	return getList[CommercialPaperBalanceItem](ctx, c.http, "/api/entrprsBilScritsNrdmpBlce.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (PR #3 기존 9개 + 신규 3개, 회귀 없음). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/securities.go report/securities_test.go report/testdata/
git commit -m "feat(report): 채무증권 발행실적·회사채·기업어음 미상환"
```

---

### Task 2: 단기사채·신종자본증권·조건부자본증권 미상환 (3 엔드포인트)

**Files:**
- Modify: `report/securities.go` (3개 메서드+struct 추가)
- Modify: `report/securities_test.go` (3개 테스트 추가)
- Create: `report/testdata/srtpdPsndbtNrdmpBlce.json`, `report/testdata/newCaplScritsNrdmpBlce.json`, `report/testdata/cndlCaplScritsNrdmpBlce.json`

- [ ] **Step 1: fixture 작성**

`report/testdata/srtpdPsndbtNrdmpBlce.json`:
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
            "sm": "-",
            "remndr_exprtn1": "미상환잔액",
            "remndr_exprtn2": "공모",
            "de10_below": "-",
            "de10_excess_de30_below": "-",
            "de30_excess_de90_below": "-",
            "de90_excess_de180_below": "-",
            "de180_excess_yy1_below": "-",
            "isu_lmt": "-",
            "remndr_lmt": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/newCaplScritsNrdmpBlce.json`:
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
            "sm": "-",
            "remndr_exprtn1": "미상환잔액",
            "remndr_exprtn2": "공모",
            "yy1_below": "-",
            "yy1_excess_yy5_below": "-",
            "yy5_excess_yy10_below": "-",
            "yy10_excess_yy15_below": "-",
            "yy15_excess_yy20_below": "-",
            "yy20_excess_yy30_below": "-",
            "yy30_excess": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

`report/testdata/cndlCaplScritsNrdmpBlce.json`:
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
            "sm": "-",
            "remndr_exprtn1": "미상환잔액",
            "remndr_exprtn2": "공모",
            "yy1_below": "-",
            "yy1_excess_yy2_below": "-",
            "yy2_excess_yy3_below": "-",
            "yy3_excess_yy4_below": "-",
            "yy4_excess_yy5_below": "-",
            "yy5_excess_yy10_below": "-",
            "yy10_excess_yy20_below": "-",
            "yy20_excess_yy30_below": "-",
            "yy30_excess": "-",
            "stlm_dt": "2023-12-31"
        }
    ]
}
```

- [ ] **Step 2: 실패하는 테스트 추가** — `report/securities_test.go` 에 추가:
```go
func TestShortTermBondBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/srtpdPsndbtNrdmpBlce.json": "srtpdPsndbtNrdmpBlce.json"})
	items, err := c.ShortTermBondBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미상환잔액", items[0].RemndrExprtn1)
	assert.Equal(t, "-", items[0].IsuLmt)
	assert.Equal(t, "-", items[0].RemndrLmt)
}

func TestHybridSecuritiesBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/newCaplScritsNrdmpBlce.json": "newCaplScritsNrdmpBlce.json"})
	items, err := c.HybridSecuritiesBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "공모", items[0].RemndrExprtn2)
	assert.Equal(t, "-", items[0].Yy30Excess)
}

func TestContingentCapitalBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cndlCaplScritsNrdmpBlce.json": "cndlCaplScritsNrdmpBlce.json"})
	items, err := c.ContingentCapitalBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미상환잔액", items[0].RemndrExprtn1)
	assert.Equal(t, "-", items[0].Yy10ExcessYy20Below)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./report/ -run 'TestShortTermBondBalance|TestHybridSecuritiesBalance|TestContingentCapitalBalance' -v`
Expected: FAIL — `undefined: ... ShortTermBondBalance` 등.

- [ ] **Step 4: 구현** — `report/securities.go` 에 추가:
```go
// ShortTermBondBalanceItem 은 단기사채 미상환 잔액 (srtpdPsndbtNrdmpBlce) 한 건.
type ShortTermBondBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Sm                   string `json:"sm"`                      // 합계
	IsuLmt               string `json:"isu_lmt"`                 // 발행 한도
	RemndrLmt            string `json:"remndr_lmt"`              // 잔여 한도
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// ShortTermBondBalance 는 단기사채 미상환 잔액을 조회한다.
func (c *Client) ShortTermBondBalance(ctx context.Context, p ReportParams) ([]ShortTermBondBalanceItem, error) {
	return getList[ShortTermBondBalanceItem](ctx, c.http, "/api/srtpdPsndbtNrdmpBlce.json", p)
}

// HybridSecuritiesBalanceItem 은 신종자본증권 미상환 잔액 (newCaplScritsNrdmpBlce) 한 건.
type HybridSecuritiesBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy5Below   string `json:"yy1_excess_yy5_below"`   // 1년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy15Below string `json:"yy10_excess_yy15_below"` // 10년초과 15년이하
	Yy15ExcessYy20Below string `json:"yy15_excess_yy20_below"` // 15년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// HybridSecuritiesBalance 는 신종자본증권 미상환 잔액을 조회한다.
func (c *Client) HybridSecuritiesBalance(ctx context.Context, p ReportParams) ([]HybridSecuritiesBalanceItem, error) {
	return getList[HybridSecuritiesBalanceItem](ctx, c.http, "/api/newCaplScritsNrdmpBlce.json", p)
}

// ContingentCapitalBalanceItem 은 조건부 자본증권 미상환 잔액 (cndlCaplScritsNrdmpBlce) 한 건.
type ContingentCapitalBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy2Below   string `json:"yy1_excess_yy2_below"`   // 1년초과 2년이하
	Yy2ExcessYy3Below   string `json:"yy2_excess_yy3_below"`   // 2년초과 3년이하
	Yy3ExcessYy4Below   string `json:"yy3_excess_yy4_below"`   // 3년초과 4년이하
	Yy4ExcessYy5Below   string `json:"yy4_excess_yy5_below"`   // 4년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy20Below string `json:"yy10_excess_yy20_below"` // 10년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// ContingentCapitalBalance 는 조건부 자본증권 미상환 잔액을 조회한다.
func (c *Client) ContingentCapitalBalance(ctx context.Context, p ReportParams) ([]ContingentCapitalBalanceItem, error) {
	return getList[ContingentCapitalBalanceItem](ctx, c.http, "/api/cndlCaplScritsNrdmpBlce.json", p)
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./report/ -v`
Expected: 전체 PASS (Task 1 포함 신규 6개 + PR #3 기존). `go vet ./report/` clean.

- [ ] **Step 6: Commit**

```bash
git add report/securities.go report/securities_test.go report/testdata/
git commit -m "feat(report): 단기사채·신종자본증권·조건부자본증권 미상환"
```

---

### Task 3: README 커버리지 · 통합 테스트 · 최종 검증

**Files:**
- Modify: `README.md`
- Modify: `integration_test.go`

- [ ] **Step 1: README 커버리지 갱신** — `README.md` 의 DS002 줄을 다음으로 교체:
```markdown
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황 · 증권 발행실적 · 미상환 잔액(회사채/기업어음/단기사채/신종자본증권/조건부자본증권)
```
(바로 다음 줄 `- (예정) DS002 나머지 · DS003~DS006` 은 그대로 둔다.)

- [ ] **Step 2: 통합 테스트 추가** — `integration_test.go` 에 함수 추가 (기존 `//go:build integration` · `report` import 유지):
```go
func TestIntegration_DebtSecuritiesIssuance(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.DebtSecuritiesIssuance(context.Background(), report.ReportParams{
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
git commit -m "docs(report): DS002 증권 발행·미상환 커버리지 + 통합 테스트"
```

---

## Self-Review Notes

- **Spec coverage:** 6개 메서드+struct = Task1(3)+Task2(3) · 테스트(fixture) = Task1·2 · README 커버리지 = Task3 · 통합 테스트 = Task3. root 와이어링은 PR #3 에서 이미 완료(변경 없음). 모두 매핑됨.
- **Type consistency:** 6개 `XItem`/메서드(`DebtSecuritiesIssuance`/`CorporateBondBalance`/`CommercialPaperBalance`/`ShortTermBondBalance`/`HybridSecuritiesBalance`/`ContingentCapitalBalance`), 시그니처 `(ctx, ReportParams) ([]XItem, error)` 일관. `getList[T]`/`ReportParams`/`AnnualReport`/`newTestClient` 는 PR #3 기존 심볼 재사용. 필드명·json 태그는 캡처한 실 응답과 1:1.
- **검증된 fixture:** 6개 모두 실 API(삼성전자/2023/사업보고서) 응답 첫 항목. 숫자 콤마 string, 빈 값 "-".
- **새 추상화 없음:** 기존 제네릭 getList 재사용만. root 변경 없음(client.Report 기존).
