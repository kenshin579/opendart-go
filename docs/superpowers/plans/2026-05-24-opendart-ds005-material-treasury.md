# OpenDART DS005 자기주식 그룹 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 주요사항보고서 주요정보의 자기주식 4개 API(취득/처분/신탁계약 체결/신탁계약 해지 결정)를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 를 그대로 재사용한다. 신규 파일 `material/treasury.go` 에 4개 item struct + 4개 한 줄 메서드를 추가한다. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 후 임베드(불가 시 docs 스키마 일치 샘플).

---

## File Structure

- Create: `material/treasury.go` — 4 item struct(`TreasuryStockAcquisitionItem` 29필드, `TreasuryStockDisposalItem` 32, `TreasuryStockTrustContractItem` 23, `TreasuryStockTrustCancellationItem` 24) + 4 메서드.
- Create: `material/treasury_test.go` — fixture 디코딩 테스트(기존 `material/client_test.go` 의 `newTestClient` 재사용).
- Create: `material/testdata/tsstkAqDecsn.json`, `tsstkDpDecsn.json`, `tsstkAqTrctrCnsDecsn.json`, `tsstkAqTrctrCcDecsn.json` — 실 응답 fixture.
- Modify: `integration_test.go` — `TreasuryStockAcquisition` 통합 케이스(`//go:build integration`, ErrNoData skip).
- Modify: `README.md` — DS005 커버리지에 자기주식 추가.

기존 컨벤션(참고용, 변경 금지):
- `material/bond.go` 메서드 형태: `func (c *Client) ConvertibleBondIssuance(ctx context.Context, p MaterialParams) ([]ConvertibleBondItem, error) { return httpclient.GetList[ConvertibleBondItem](ctx, c.http, "/api/cvbdIsDecsn.json", p.toMap()) }` — Client 의 http 필드는 `c.http`.
- `material/client.go` 의 `MaterialParams.toMap()`: corp_code 항상 포함, 빈 bgn_de/end_de omit.
- `material/client_test.go` 의 `newTestClient(t, routes map[string]string) *Client`: route map 값은 **bare 파일명**(예: `"tsstkAqDecsn.json"`) — 내부에서 `filepath.Join("testdata", fixture)` 로 testdata/ 를 붙인다. 절대 `"testdata/..."` 로 쓰지 말 것.

---

## Task 1: 자기주식 취득 결정 (TreasuryStockAcquisition)

**Files:**
- Create: `material/treasury.go`
- Create: `material/testdata/tsstkAqDecsn.json`
- Create: `material/treasury_test.go`

- [ ] **Step 1: fixture 파일 작성** `material/testdata/tsstkAqDecsn.json`

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230510000111",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트취득",
      "aqpln_stk_ostk": "1,000,000",
      "aqpln_stk_estk": "-",
      "aqpln_prc_ostk": "70,000,000,000",
      "aqpln_prc_estk": "-",
      "aqexpd_bgd": "2023년 05월 11일",
      "aqexpd_edd": "2023년 08월 10일",
      "hdexpd_bgd": "-",
      "hdexpd_edd": "-",
      "aq_pp": "주주가치 제고",
      "aq_mth": "유가증권시장을 통한 장내매수",
      "cs_iv_bk": "한국투자증권",
      "aq_wtn_div_ostk": "5,000,000",
      "aq_wtn_div_ostk_rt": "2.5",
      "aq_wtn_div_estk": "-",
      "aq_wtn_div_estk_rt": "-",
      "eaq_ostk": "-",
      "eaq_ostk_rt": "-",
      "eaq_estk": "-",
      "eaq_estk_rt": "-",
      "aq_dd": "2023년 05월 10일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "d1_prodlm_ostk": "250,000",
      "d1_prodlm_estk": "-"
    }
  ]
}
```

- [ ] **Step 2: 테스트 작성** `material/treasury_test.go`

```go
package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreasuryStockAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqDecsn.json": "tsstkAqDecsn.json",
	})

	items, err := c.TreasuryStockAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230510000111", got.RceptNo)
	assert.Equal(t, "1,000,000", got.AqplnStkOstk)
	assert.Equal(t, "유가증권시장을 통한 장내매수", got.AqMth)
	assert.Equal(t, "2.5", got.AqWtnDivOstkRt)
	assert.Equal(t, "250,000", got.D1ProdlmOstk)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestTreasuryStockAcquisition`
Expected: FAIL — `c.TreasuryStockAcquisition undefined` / `TreasuryStockAcquisitionItem` 미정의.

- [ ] **Step 4: treasury.go 작성** `material/treasury.go`

```go
package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// TreasuryStockAcquisitionItem 은 자기주식 취득 결정 (tsstkAqDecsn) 한 건.
type TreasuryStockAcquisitionItem struct {
	RceptNo        string `json:"rcept_no"`          // 접수번호
	CorpCls        string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`         // 고유번호
	CorpName       string `json:"corp_name"`         // 회사명
	AqplnStkOstk   string `json:"aqpln_stk_ostk"`    // 취득예정주식(주)(보통주식)
	AqplnStkEstk   string `json:"aqpln_stk_estk"`    // 취득예정주식(주)(기타주식)
	AqplnPrcOstk   string `json:"aqpln_prc_ostk"`    // 취득예정금액(원)(보통주식)
	AqplnPrcEstk   string `json:"aqpln_prc_estk"`    // 취득예정금액(원)(기타주식)
	AqexpdBgd      string `json:"aqexpd_bgd"`        // 취득예상기간(시작일)
	AqexpdEdd      string `json:"aqexpd_edd"`        // 취득예상기간(종료일)
	HdexpdBgd      string `json:"hdexpd_bgd"`        // 보유예상기간(시작일)
	HdexpdEdd      string `json:"hdexpd_edd"`        // 보유예상기간(종료일)
	AqPp           string `json:"aq_pp"`             // 취득목적
	AqMth          string `json:"aq_mth"`            // 취득방법
	CsIvBk         string `json:"cs_iv_bk"`          // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`   // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`   // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`          // 취득 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`       // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`          // 취득 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`       // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	AqDd           string `json:"aq_dd"`             // 취득결정일
	OdAAtT         string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`         // 감사(사외이사가 아닌 감사위원) 참석여부
	D1ProdlmOstk   string `json:"d1_prodlm_ostk"`    // 1일 매수 주문수량 한도(보통주식)
	D1ProdlmEstk   string `json:"d1_prodlm_estk"`    // 1일 매수 주문수량 한도(기타주식)
}

// TreasuryStockAcquisition 은 자기주식 취득 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockAcquisition(ctx context.Context, p MaterialParams) ([]TreasuryStockAcquisitionItem, error) {
	return httpclient.GetList[TreasuryStockAcquisitionItem](ctx, c.http, "/api/tsstkAqDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestTreasuryStockAcquisition`
Expected: PASS.

- [ ] **Step 6: gofmt** `gofmt -l material/treasury.go material/treasury_test.go` — 출력 없으면 OK(있으면 `gofmt -w`).

- [ ] **Step 7: 실 API fixture 캡처(권장)**

`$OPENDART_API_KEY` 있으면 실 응답으로 교체. 자기주식 취득은 흔함:
```bash
curl -s "https://opendart.fss.or.kr/api/tsstkAqDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231" | python3 -m json.tool
```
(서버는 TLS1.2 RSA cipher 만 지원 → curl handshake 실패할 수 있음. 실패하면 샘플 유지.) status "000" + list 비어있지 않으면 그 응답을 fixture 로 저장하고 Step 2 assert 값을 실제 값에 맞춰 갱신 후 Step 5 재실행. 1회 시도.

- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/treasury.go material/treasury_test.go material/testdata/tsstkAqDecsn.json
git commit -m "feat(material): add DS005 TreasuryStockAcquisition (자기주식 취득 결정)"
```

---

## Task 2: 자기주식 처분 결정 (TreasuryStockDisposal)

**Files:**
- Modify: `material/treasury.go`
- Create: `material/testdata/tsstkDpDecsn.json`
- Modify: `material/treasury_test.go`

- [ ] **Step 1: fixture** `material/testdata/tsstkDpDecsn.json`

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230610000222",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트처분",
      "dppln_stk_ostk": "500,000",
      "dppln_stk_estk": "-",
      "dpstk_prc_ostk": "75,000",
      "dpstk_prc_estk": "-",
      "dppln_prc_ostk": "37,500,000,000",
      "dppln_prc_estk": "-",
      "dpprpd_bgd": "2023년 06월 11일",
      "dpprpd_edd": "2023년 09월 10일",
      "dp_pp": "임직원 상여 지급",
      "dp_m_mkt": "500,000",
      "dp_m_ovtm": "-",
      "dp_m_otc": "-",
      "dp_m_etc": "-",
      "cs_iv_bk": "한국투자증권",
      "aq_wtn_div_ostk": "5,000,000",
      "aq_wtn_div_ostk_rt": "2.5",
      "aq_wtn_div_estk": "-",
      "aq_wtn_div_estk_rt": "-",
      "eaq_ostk": "-",
      "eaq_ostk_rt": "-",
      "eaq_estk": "-",
      "eaq_estk_rt": "-",
      "dp_dd": "2023년 06월 10일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "d1_slodlm_ostk": "125,000",
      "d1_slodlm_estk": "-"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가** `material/treasury_test.go` 에 함수 추가

```go
func TestTreasuryStockDisposal(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkDpDecsn.json": "tsstkDpDecsn.json",
	})

	items, err := c.TreasuryStockDisposal(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230610000222", got.RceptNo)
	assert.Equal(t, "500,000", got.DpplnStkOstk)
	assert.Equal(t, "37,500,000,000", got.DpplnPrcOstk)
	assert.Equal(t, "임직원 상여 지급", got.DpPp)
	assert.Equal(t, "500,000", got.DpMMkt)
	assert.Equal(t, "125,000", got.D1SlodlmOstk)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestTreasuryStockDisposal`
Expected: FAIL — `c.TreasuryStockDisposal undefined`.

- [ ] **Step 4: treasury.go 에 struct + 메서드 추가** (파일 끝에 append)

```go
// TreasuryStockDisposalItem 은 자기주식 처분 결정 (tsstkDpDecsn) 한 건.
type TreasuryStockDisposalItem struct {
	RceptNo        string `json:"rcept_no"`          // 접수번호
	CorpCls        string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`         // 고유번호
	CorpName       string `json:"corp_name"`         // 회사명
	DpplnStkOstk   string `json:"dppln_stk_ostk"`    // 처분예정주식(주)(보통주식)
	DpplnStkEstk   string `json:"dppln_stk_estk"`    // 처분예정주식(주)(기타주식)
	DpstkPrcOstk   string `json:"dpstk_prc_ostk"`    // 처분 대상 주식가격(원)(보통주식)
	DpstkPrcEstk   string `json:"dpstk_prc_estk"`    // 처분 대상 주식가격(원)(기타주식)
	DpplnPrcOstk   string `json:"dppln_prc_ostk"`    // 처분예정금액(원)(보통주식)
	DpplnPrcEstk   string `json:"dppln_prc_estk"`    // 처분예정금액(원)(기타주식)
	DpprpdBgd      string `json:"dpprpd_bgd"`        // 처분예정기간(시작일)
	DpprpdEdd      string `json:"dpprpd_edd"`        // 처분예정기간(종료일)
	DpPp           string `json:"dp_pp"`             // 처분목적
	DpMMkt         string `json:"dp_m_mkt"`          // 처분방법(시장을 통한 매도(주))
	DpMOvtm        string `json:"dp_m_ovtm"`         // 처분방법(시간외대량매매(주))
	DpMOtc         string `json:"dp_m_otc"`          // 처분방법(장외처분(주))
	DpMEtc         string `json:"dp_m_etc"`          // 처분방법(기타(주))
	CsIvBk         string `json:"cs_iv_bk"`          // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`   // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`   // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`          // 처분 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`       // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`          // 처분 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`       // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	DpDd           string `json:"dp_dd"`             // 처분결정일
	OdAAtT         string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`         // 감사(사외이사가 아닌 감사위원) 참석여부
	D1SlodlmOstk   string `json:"d1_slodlm_ostk"`    // 1일 매도 주문수량 한도(보통주식)
	D1SlodlmEstk   string `json:"d1_slodlm_estk"`    // 1일 매도 주문수량 한도(기타주식)
}

// TreasuryStockDisposal 은 자기주식 처분 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockDisposal(ctx context.Context, p MaterialParams) ([]TreasuryStockDisposalItem, error) {
	return httpclient.GetList[TreasuryStockDisposalItem](ctx, c.http, "/api/tsstkDpDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인** — `go test ./material/ -run TestTreasuryStockDisposal` → PASS.
- [ ] **Step 6: gofmt** `gofmt -l material/treasury.go material/treasury_test.go` — 출력 없으면 OK.
- [ ] **Step 7: 실 API 캡처(권장)** — Task 1 Step 7 과 동일, 엔드포인트 `/api/tsstkDpDecsn.json`. 1회 시도, 실패 시 샘플 유지.
- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/treasury.go material/treasury_test.go material/testdata/tsstkDpDecsn.json
git commit -m "feat(material): add DS005 TreasuryStockDisposal (자기주식 처분 결정)"
```

---

## Task 3: 자기주식취득 신탁계약 체결 결정 (TreasuryStockTrustContract)

**Files:**
- Modify: `material/treasury.go`
- Create: `material/testdata/tsstkAqTrctrCnsDecsn.json`
- Modify: `material/treasury_test.go`

- [ ] **Step 1: fixture** `material/testdata/tsstkAqTrctrCnsDecsn.json`

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230310000333",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트신탁체결",
      "ctr_prc": "50,000,000,000",
      "ctr_pd_bgd": "2023년 03월 11일",
      "ctr_pd_edd": "2023년 09월 10일",
      "ctr_pp": "주가 안정 및 주주가치 제고",
      "ctr_cns_int": "한국투자증권",
      "ctr_cns_prd": "2023년 03월 11일",
      "aq_wtn_div_ostk": "5,000,000",
      "aq_wtn_div_ostk_rt": "2.5",
      "aq_wtn_div_estk": "-",
      "aq_wtn_div_estk_rt": "-",
      "eaq_ostk": "-",
      "eaq_ostk_rt": "-",
      "eaq_estk": "-",
      "eaq_estk_rt": "-",
      "bddd": "2023년 03월 10일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "cs_iv_bk": "한국투자증권"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가** `material/treasury_test.go` 에 함수 추가

```go
func TestTreasuryStockTrustContract(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqTrctrCnsDecsn.json": "tsstkAqTrctrCnsDecsn.json",
	})

	items, err := c.TreasuryStockTrustContract(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230310000333", got.RceptNo)
	assert.Equal(t, "50,000,000,000", got.CtrPrc)
	assert.Equal(t, "한국투자증권", got.CtrCnsInt)
	assert.Equal(t, "2023년 03월 10일", got.Bddd)
	assert.Equal(t, "주가 안정 및 주주가치 제고", got.CtrPp)
}
```

- [ ] **Step 3: 테스트 실패 확인** — `go test ./material/ -run TestTreasuryStockTrustContract` → FAIL (`undefined`).

- [ ] **Step 4: treasury.go 에 struct + 메서드 추가** (파일 끝에 append)

```go
// TreasuryStockTrustContractItem 은 자기주식취득 신탁계약 체결 결정 (tsstkAqTrctrCnsDecsn) 한 건.
type TreasuryStockTrustContractItem struct {
	RceptNo        string `json:"rcept_no"`          // 접수번호
	CorpCls        string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`         // 고유번호
	CorpName       string `json:"corp_name"`         // 회사명
	CtrPrc         string `json:"ctr_prc"`           // 계약금액(원)
	CtrPdBgd       string `json:"ctr_pd_bgd"`        // 계약기간(시작일)
	CtrPdEdd       string `json:"ctr_pd_edd"`        // 계약기간(종료일)
	CtrPp          string `json:"ctr_pp"`            // 계약목적
	CtrCnsInt      string `json:"ctr_cns_int"`       // 계약체결기관
	CtrCnsPrd      string `json:"ctr_cns_prd"`       // 계약체결 예정일자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`   // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`   // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`          // 계약 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`       // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`          // 계약 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`       // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`         // 감사(사외이사가 아닌 감사위원) 참석여부
	CsIvBk         string `json:"cs_iv_bk"`          // 위탁투자중개업자
}

// TreasuryStockTrustContract 는 자기주식취득 신탁계약 체결 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockTrustContract(ctx context.Context, p MaterialParams) ([]TreasuryStockTrustContractItem, error) {
	return httpclient.GetList[TreasuryStockTrustContractItem](ctx, c.http, "/api/tsstkAqTrctrCnsDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인** — `go test ./material/ -run TestTreasuryStockTrustContract` → PASS.
- [ ] **Step 6: gofmt** `gofmt -l material/treasury.go material/treasury_test.go` — 출력 없으면 OK.
- [ ] **Step 7: 실 API 캡처(권장)** — Task 1 Step 7 과 동일, 엔드포인트 `/api/tsstkAqTrctrCnsDecsn.json`. 1회 시도.
- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/treasury.go material/treasury_test.go material/testdata/tsstkAqTrctrCnsDecsn.json
git commit -m "feat(material): add DS005 TreasuryStockTrustContract (자기주식취득 신탁계약 체결 결정)"
```

---

## Task 4: 자기주식취득 신탁계약 해지 결정 (TreasuryStockTrustCancellation)

**Files:**
- Modify: `material/treasury.go`
- Create: `material/testdata/tsstkAqTrctrCcDecsn.json`
- Modify: `material/treasury_test.go`

- [ ] **Step 1: fixture** `material/testdata/tsstkAqTrctrCcDecsn.json`

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230910000444",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트신탁해지",
      "ctr_prc_bfcc": "50,000,000,000",
      "ctr_prc_atcc": "50,000,000,000",
      "ctr_pd_bfcc_bgd": "2023년 03월 11일",
      "ctr_pd_bfcc_edd": "2023년 09월 10일",
      "cc_pp": "신탁계약 기간 만료",
      "cc_int": "한국투자증권",
      "cc_prd": "2023년 09월 11일",
      "tp_rm_atcc": "현금으로 반환",
      "aq_wtn_div_ostk": "5,500,000",
      "aq_wtn_div_ostk_rt": "2.7",
      "aq_wtn_div_estk": "-",
      "aq_wtn_div_estk_rt": "-",
      "eaq_ostk": "-",
      "eaq_ostk_rt": "-",
      "eaq_estk": "-",
      "eaq_estk_rt": "-",
      "bddd": "2023년 09월 10일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가** `material/treasury_test.go` 에 함수 추가

```go
func TestTreasuryStockTrustCancellation(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqTrctrCcDecsn.json": "tsstkAqTrctrCcDecsn.json",
	})

	items, err := c.TreasuryStockTrustCancellation(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230910000444", got.RceptNo)
	assert.Equal(t, "50,000,000,000", got.CtrPrcBfcc)
	assert.Equal(t, "신탁계약 기간 만료", got.CcPp)
	assert.Equal(t, "현금으로 반환", got.TpRmAtcc)
	assert.Equal(t, "2023년 09월 10일", got.Bddd)
}
```

- [ ] **Step 3: 테스트 실패 확인** — `go test ./material/ -run TestTreasuryStockTrustCancellation` → FAIL (`undefined`).

- [ ] **Step 4: treasury.go 에 struct + 메서드 추가** (파일 끝에 append)

```go
// TreasuryStockTrustCancellationItem 은 자기주식취득 신탁계약 해지 결정 (tsstkAqTrctrCcDecsn) 한 건.
type TreasuryStockTrustCancellationItem struct {
	RceptNo        string `json:"rcept_no"`          // 접수번호
	CorpCls        string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`         // 고유번호
	CorpName       string `json:"corp_name"`         // 회사명
	CtrPrcBfcc     string `json:"ctr_prc_bfcc"`      // 계약금액(원)(해지 전)
	CtrPrcAtcc     string `json:"ctr_prc_atcc"`      // 계약금액(원)(해지 후)
	CtrPdBfccBgd   string `json:"ctr_pd_bfcc_bgd"`   // 해지 전 계약기간(시작일)
	CtrPdBfccEdd   string `json:"ctr_pd_bfcc_edd"`   // 해지 전 계약기간(종료일)
	CcPp           string `json:"cc_pp"`             // 해지목적
	CcInt          string `json:"cc_int"`            // 해지기관
	CcPrd          string `json:"cc_prd"`            // 해지예정일자
	TpRmAtcc       string `json:"tp_rm_atcc"`        // 해지후 신탁재산의 반환방법
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`   // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`   // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`          // 해지 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`       // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`          // 해지 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`       // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`         // 감사(사외이사가 아닌 감사위원) 참석여부
}

// TreasuryStockTrustCancellation 은 자기주식취득 신탁계약 해지 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockTrustCancellation(ctx context.Context, p MaterialParams) ([]TreasuryStockTrustCancellationItem, error) {
	return httpclient.GetList[TreasuryStockTrustCancellationItem](ctx, c.http, "/api/tsstkAqTrctrCcDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인** — `go test ./material/ -run TestTreasuryStockTrustCancellation` → PASS.
- [ ] **Step 6: gofmt** `gofmt -l material/treasury.go material/treasury_test.go` — 출력 없으면 OK.
- [ ] **Step 7: 실 API 캡처(권장)** — Task 1 Step 7 과 동일, 엔드포인트 `/api/tsstkAqTrctrCcDecsn.json`. 1회 시도.
- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/treasury.go material/treasury_test.go material/testdata/tsstkAqTrctrCcDecsn.json
git commit -m "feat(material): add DS005 TreasuryStockTrustCancellation (자기주식취득 신탁계약 해지 결정)"
```

---

## Task 5: 통합 테스트 + README

**Files:**
- Modify: `integration_test.go`
- Modify: `README.md`

- [ ] **Step 1: 통합 테스트 추가**

먼저 `integration_test.go` 를 읽어 기존 패턴을 확인할 것. 이 파일은 `//go:build integration` + `package opendart` 이고, 클라이언트는 `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))` 로 직접 생성한다(헬퍼 함수 없음). `ErrNoData` 는 **같은 패키지**라 `ErrNoData` 로 참조(NOT `opendart.ErrNoData`). 이미 `errors` import 가 있다(사채 발행 PR 에서 추가됨). 파일 끝에 함수 추가:

```go
func TestIntegration_TreasuryStockAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Material.TreasuryStockAcquisition(context.Background(), material.MaterialParams{
		CorpCode: corp,
		BgnDe:    "20200101",
		EndDe:    "20241231",
	})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 자기주식 취득 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

만약 `errors` import 가 실제로 없으면 import 블록에 추가한다. 기존 파일의 import 와 헬퍼 사용법을 정확히 일치시킬 것.

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.

- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run TestIntegration_TreasuryStockAcquisition ./...` → PASS 또는 SKIP. 데이터 외 이유로 실패하면 점검.

- [ ] **Step 4: README 커버리지 갱신**

`README.md` 의 DS005 줄에 자기주식 추가, 예정 줄에서 자기주식 제거. 현재(사채 발행 PR 머지 후) 두 줄:
```
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정)
- (예정) DS005 나머지(자기주식/양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
다음으로 교체:
```
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정) · 자기주식(취득/처분/신탁계약 체결·해지 결정)
- (예정) DS005 나머지(양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
(그 두 줄만 변경, README 나머지는 그대로.)

- [ ] **Step 5: 전체 빌드 + 테스트 + 포맷 게이트**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/ integration_test.go`
Expected: 빌드 성공, 전체 테스트 PASS, gofmt 출력 없음.

- [ ] **Step 6: README UTF-8 확인** — `file -I README.md` → `charset=utf-8`.

- [ ] **Step 7: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add integration_test.go README.md
git commit -m "test(material): add DS005 자기주식 통합 테스트 + README 커버리지"
```

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 4개 메서드(TreasuryStockAcquisition/Disposal/TrustContract/TrustCancellation) → Task 1~4 각각 매핑. struct 4종 전 필드(29/32/23/24) = spec 의 4 struct 와 1:1. MaterialParams+GetList 재사용·root 무변경 = Architecture 일치. 테스트 전략(fixture+통합 1개)·README 갱신 = Task 5. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 모든 코드 step 에 완전한 코드 포함. fixture 는 "실 API 캡처 권장, 안되면 샘플 유지" 명시(기존 그룹과 동일 정책).

**3. Type consistency:** 메서드명·struct명·필드명이 Task 간 일관. 보유현황 json 키는 4종 모두 `aq_wtn_div_*`/`eaq_*` (처분/해지에서도 docs 원문대로 동일 키, 코멘트만 "처분 전/해지 전"). route map 값은 전부 bare 파일명. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조.
