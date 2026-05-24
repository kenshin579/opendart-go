# OpenDART DS005 사채 발행 그룹 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 주요사항보고서 주요정보의 사채 발행 4개 API(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정)를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 를 그대로 재사용한다. 신규 파일 `material/bond.go` 에 4개 item struct + 4개 한 줄 메서드를 추가한다. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 후 임베드.

---

## File Structure

- Create: `material/bond.go` — 4 item struct(`ConvertibleBondItem` 46필드, `BondWithWarrantItem` 49, `ExchangeableBondItem` 42, `ContingentConvertibleBondItem` 34) + 4 메서드(`ConvertibleBondIssuance`/`BondWithWarrantIssuance`/`ExchangeableBondIssuance`/`ContingentConvertibleBondIssuance`).
- Create: `material/bond_test.go` — fixture 디코딩 테스트(기존 `material/client_test.go` 의 `newTestClient` 재사용).
- Create: `material/testdata/cvbdIsDecsn.json`, `bdwtIsDecsn.json`, `exbdIsDecsn.json`, `wdCocobdIsDecsn.json` — 실 응답 fixture.
- Modify: `integration_test.go` — `ConvertibleBondIssuance` 통합 케이스 추가(`//go:build integration`).
- Modify: `README.md:33-34` — DS005 커버리지에 "사채 발행" 추가.

기존 컨벤션(참고용, 변경 금지):
- `material/capital.go` 의 메서드 형태: `func (c *Client) PaidInCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidInCapitalIncreaseItem, error) { return httpclient.GetList[PaidInCapitalIncreaseItem](ctx, c.http, "/api/piicDecsn.json", p.toMap()) }`
- `material/client.go` 의 `MaterialParams.toMap()` 는 corp_code 항상 포함, 빈 bgn_de/end_de 는 omit.
- `material/client_test.go` 의 `newTestClient(t, routes map[string]string) *Client` 는 routes 의 path→testdata 파일을 httptest 로 서빙한다.

---

## Task 1: 전환사채권 발행결정 (ConvertibleBondIssuance)

**Files:**
- Create: `material/bond.go`
- Create: `material/testdata/cvbdIsDecsn.json`
- Create: `material/bond_test.go`

- [ ] **Step 1: fixture 파일 작성**

`material/testdata/cvbdIsDecsn.json` 생성(실 응답 캡처 전 단계에서는 아래 스키마 일치 샘플; Task 1 Step 7에서 실 API로 교체):

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230815000123",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트전환사채",
      "bd_tm": "10",
      "bd_knd": "기명식 무보증 사모 전환사채",
      "bd_fta": "30,000,000,000",
      "atcsc_rmislmt": "100,000,000,000",
      "ovis_fta": "-",
      "ovis_fta_crn": "-",
      "ovis_ster": "-",
      "ovis_isar": "-",
      "ovis_mktnm": "-",
      "fdpp_fclt": "-",
      "fdpp_bsninh": "-",
      "fdpp_op": "30,000,000,000",
      "fdpp_dtrp": "-",
      "fdpp_ocsa": "-",
      "fdpp_etc": "-",
      "bd_intr_ex": "0.0",
      "bd_intr_sf": "2.0",
      "bd_mtd": "2028년 08월 16일",
      "bdis_mthn": "사모",
      "cv_rt": "100",
      "cv_prc": "15,000",
      "cvisstk_knd": "기명식 보통주",
      "cvisstk_cnt": "2,000,000",
      "cvisstk_tisstk_vs": "3.5",
      "cvrqpd_bgd": "2024년 08월 16일",
      "cvrqpd_edd": "2028년 07월 16일",
      "act_mktprcfl_cvprc_lwtrsprc": "10,500",
      "act_mktprcfl_cvprc_lwtrsprc_bs": "전환가액의 70%",
      "rmislmt_lt70p": "-",
      "abmg": "-",
      "sbd": "2023년 08월 16일",
      "pymd": "2023년 08월 16일",
      "rpmcmp": "-",
      "grint": "-",
      "bddd": "2023년 08월 15일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "rs_sm_atn": "미제출",
      "ex_sm_r": "사모발행",
      "ovis_ltdtl": "-",
      "ftc_stt_atn": "미해당"
    }
  ]
}
```

- [ ] **Step 2: 테스트 작성 (실패 확인용)**

`material/bond_test.go` 생성:

```go
package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertibleBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/cvbdIsDecsn.json": "testdata/cvbdIsDecsn.json",
	})

	items, err := c.ConvertibleBondIssuance(context.Background(), MaterialParams{CorpCode: "00126380"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230815000123", got.RceptNo)
	assert.Equal(t, "30,000,000,000", got.BdFta)
	assert.Equal(t, "15,000", got.CvPrc)
	assert.Equal(t, "2024년 08월 16일", got.CvrqpdBgd)
	assert.Equal(t, "미제출", got.RsSmAtn)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestConvertibleBondIssuance`
Expected: FAIL — `c.ConvertibleBondIssuance undefined` / `ConvertibleBondItem` 미정의 컴파일 에러.

- [ ] **Step 4: bond.go 작성 (struct + 메서드)**

`material/bond.go` 생성:

```go
package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// ConvertibleBondItem 은 전환사채권 발행결정 (cvbdIsDecsn) 한 건.
type ConvertibleBondItem struct {
	RceptNo                    string `json:"rcept_no"`                       // 접수번호
	CorpCls                    string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                   string `json:"corp_code"`                      // 고유번호
	CorpName                   string `json:"corp_name"`                      // 회사명
	BdTm                       string `json:"bd_tm"`                          // 사채의 종류(회차)
	BdKnd                      string `json:"bd_knd"`                         // 사채의 종류(종류)
	BdFta                      string `json:"bd_fta"`                         // 사채의 권면(전자등록)총액 (원)
	AtcscRmislmt               string `json:"atcsc_rmislmt"`                  // 정관상 잔여 발행한도 (원)
	OvisFta                    string `json:"ovis_fta"`                       // 해외발행(권면(전자등록)총액)
	OvisFtaCrn                 string `json:"ovis_fta_crn"`                   // 해외발행(권면총액 통화단위)
	OvisSter                   string `json:"ovis_ster"`                      // 해외발행(기준환율등)
	OvisIsar                   string `json:"ovis_isar"`                      // 해외발행(발행지역)
	OvisMktnm                  string `json:"ovis_mktnm"`                     // 해외발행(해외상장시 시장 명칭)
	FdppFclt                   string `json:"fdpp_fclt"`                      // 자금조달목적(시설자금)
	FdppBsninh                 string `json:"fdpp_bsninh"`                    // 자금조달목적(영업양수자금)
	FdppOp                     string `json:"fdpp_op"`                        // 자금조달목적(운영자금)
	FdppDtrp                   string `json:"fdpp_dtrp"`                      // 자금조달목적(채무상환자금)
	FdppOcsa                   string `json:"fdpp_ocsa"`                      // 자금조달목적(타법인 증권 취득자금)
	FdppEtc                    string `json:"fdpp_etc"`                       // 자금조달목적(기타자금)
	BdIntrEx                   string `json:"bd_intr_ex"`                     // 사채의 이율(표면이자율 %)
	BdIntrSf                   string `json:"bd_intr_sf"`                     // 사채의 이율(만기이자율 %)
	BdMtd                      string `json:"bd_mtd"`                         // 사채만기일
	BdisMthn                   string `json:"bdis_mthn"`                      // 사채발행방법
	CvRt                       string `json:"cv_rt"`                          // 전환비율 (%)
	CvPrc                      string `json:"cv_prc"`                         // 전환가액 (원/주)
	CvisstkKnd                 string `json:"cvisstk_knd"`                    // 전환에 따라 발행할 주식(종류)
	CvisstkCnt                 string `json:"cvisstk_cnt"`                    // 전환에 따라 발행할 주식(주식수)
	CvisstkTisstkVs            string `json:"cvisstk_tisstk_vs"`              // 전환에 따라 발행할 주식(주식총수 대비 %)
	CvrqpdBgd                  string `json:"cvrqpd_bgd"`                     // 전환청구기간(시작일)
	CvrqpdEdd                  string `json:"cvrqpd_edd"`                     // 전환청구기간(종료일)
	ActMktprcflCvprcLwtrsprc   string `json:"act_mktprcfl_cvprc_lwtrsprc"`    // 시가하락 전환가액 조정(최저 조정가액 원)
	ActMktprcflCvprcLwtrsprcBs string `json:"act_mktprcfl_cvprc_lwtrsprc_bs"` // 시가하락 전환가액 조정(최저 조정가액 근거)
	RmislmtLt70p               string `json:"rmislmt_lt70p"`                  // 시가하락 조정(전환가 70% 미만 조정가능 잔여한도 원)
	Abmg                       string `json:"abmg"`                           // 합병 관련 사항
	Sbd                        string `json:"sbd"`                            // 청약일
	Pymd                       string `json:"pymd"`                           // 납입일
	Rpmcmp                     string `json:"rpmcmp"`                         // 대표주관회사
	Grint                      string `json:"grint"`                          // 보증기관
	Bddd                       string `json:"bddd"`                           // 이사회결의일(결정일)
	OdAAtT                     string `json:"od_a_at_t"`                      // 사외이사 참석여부(참석)
	OdAAtB                     string `json:"od_a_at_b"`                      // 사외이사 참석여부(불참)
	AdtAAtn                    string `json:"adt_a_atn"`                      // 감사(감사위원) 참석여부
	RsSmAtn                    string `json:"rs_sm_atn"`                      // 증권신고서 제출대상 여부
	ExSmR                      string `json:"ex_sm_r"`                        // 제출 면제 사유
	OvisLtdtl                  string `json:"ovis_ltdtl"`                     // 해외발행 연계 대차거래 내역
	FtcSttAtn                  string `json:"ftc_stt_atn"`                    // 공정거래위원회 신고대상 여부
}

// ConvertibleBondIssuance 는 전환사채권 발행결정(주요사항보고서)을 조회한다.
func (c *Client) ConvertibleBondIssuance(ctx context.Context, p MaterialParams) ([]ConvertibleBondItem, error) {
	return httpclient.GetList[ConvertibleBondItem](ctx, c.http, "/api/cvbdIsDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestConvertibleBondIssuance`
Expected: PASS.

- [ ] **Step 6: gofmt 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && gofmt -l material/bond.go`
Expected: 출력 없음(정렬됨). 출력 있으면 `gofmt -w material/bond.go`.

- [ ] **Step 7: 실 API fixture 캡처(권장)**

OPENDART_API_KEY 가 있으면 실 응답으로 fixture 교체. 전환사채 발행 사례는 공시검색으로 corp_code 탐색:
```bash
# 예: 공시검색에서 "전환사채권발행결정" 보고서가 있는 종목 찾기
curl -s "https://opendart.fss.or.kr/api/cvbdIsDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231" | python3 -m json.tool
```
status "000" 이고 list 가 비어있지 않으면 그 응답을 `testdata/cvbdIsDecsn.json` 에 저장하고, Step 2 의 assert 값을 실제 값에 맞춰 갱신 후 Step 5 재실행. 실 데이터 확보 불가하면 Step 1 샘플 유지(스키마는 docs 와 일치).

- [ ] **Step 8: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/bond.go material/bond_test.go material/testdata/cvbdIsDecsn.json
git commit -m "feat(material): add DS005 ConvertibleBondIssuance (전환사채 발행결정)"
```

---

## Task 2: 신주인수권부사채권 발행결정 (BondWithWarrantIssuance)

**Files:**
- Modify: `material/bond.go`
- Create: `material/testdata/bdwtIsDecsn.json`
- Modify: `material/bond_test.go`

- [ ] **Step 1: fixture 파일 작성**

`material/testdata/bdwtIsDecsn.json` (스키마는 BondWithWarrantItem 와 1:1):

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230910000222",
      "corp_cls": "K",
      "corp_code": "00164779",
      "corp_name": "테스트신주인수권부사채",
      "bd_tm": "5",
      "bd_knd": "기명식 무보증 분리형 신주인수권부사채",
      "bd_fta": "20,000,000,000",
      "atcsc_rmislmt": "80,000,000,000",
      "ovis_fta": "-",
      "ovis_fta_crn": "-",
      "ovis_ster": "-",
      "ovis_isar": "-",
      "ovis_mktnm": "-",
      "fdpp_fclt": "10,000,000,000",
      "fdpp_bsninh": "-",
      "fdpp_op": "10,000,000,000",
      "fdpp_dtrp": "-",
      "fdpp_ocsa": "-",
      "fdpp_etc": "-",
      "bd_intr_ex": "1.0",
      "bd_intr_sf": "3.0",
      "bd_mtd": "2026년 09월 12일",
      "bdis_mthn": "사모",
      "ex_rt": "100",
      "ex_prc": "8,000",
      "ex_prc_dmth": "이사회결의일 전 기산일 가중산술평균주가",
      "bdwt_div_atn": "분리형",
      "nstk_pym_mth": "현금납입",
      "nstk_isstk_knd": "기명식 보통주",
      "nstk_isstk_cnt": "2,500,000",
      "nstk_isstk_tisstk_vs": "4.0",
      "expd_bgd": "2024년 09월 12일",
      "expd_edd": "2026년 08월 12일",
      "act_mktprcfl_cvprc_lwtrsprc": "5,600",
      "act_mktprcfl_cvprc_lwtrsprc_bs": "행사가액의 70%",
      "rmislmt_lt70p": "-",
      "abmg": "-",
      "sbd": "2023년 09월 12일",
      "pymd": "2023년 09월 12일",
      "rpmcmp": "-",
      "grint": "-",
      "bddd": "2023년 09월 10일",
      "od_a_at_t": "2",
      "od_a_at_b": "1",
      "adt_a_atn": "1",
      "rs_sm_atn": "미제출",
      "ex_sm_r": "사모발행",
      "ovis_ltdtl": "-",
      "ftc_stt_atn": "미해당"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가 (실패 확인용)**

`material/bond_test.go` 에 함수 추가:

```go
func TestBondWithWarrantIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/bdwtIsDecsn.json": "testdata/bdwtIsDecsn.json",
	})

	items, err := c.BondWithWarrantIssuance(context.Background(), MaterialParams{CorpCode: "00164779"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230910000222", got.RceptNo)
	assert.Equal(t, "8,000", got.ExPrc)
	assert.Equal(t, "분리형", got.BdwtDivAtn)
	assert.Equal(t, "2024년 09월 12일", got.ExpdBgd)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestBondWithWarrantIssuance`
Expected: FAIL — `c.BondWithWarrantIssuance undefined`.

- [ ] **Step 4: bond.go 에 struct + 메서드 추가**

`material/bond.go` 끝에 추가:

```go
// BondWithWarrantItem 은 신주인수권부사채권 발행결정 (bdwtIsDecsn) 한 건.
type BondWithWarrantItem struct {
	RceptNo                    string `json:"rcept_no"`                       // 접수번호
	CorpCls                    string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                   string `json:"corp_code"`                      // 고유번호
	CorpName                   string `json:"corp_name"`                      // 회사명
	BdTm                       string `json:"bd_tm"`                          // 사채의 종류(회차)
	BdKnd                      string `json:"bd_knd"`                         // 사채의 종류(종류)
	BdFta                      string `json:"bd_fta"`                         // 사채의 권면(전자등록)총액 (원)
	AtcscRmislmt               string `json:"atcsc_rmislmt"`                  // 정관상 잔여 발행한도 (원)
	OvisFta                    string `json:"ovis_fta"`                       // 해외발행(권면총액)
	OvisFtaCrn                 string `json:"ovis_fta_crn"`                   // 해외발행(권면총액 통화단위)
	OvisSter                   string `json:"ovis_ster"`                      // 해외발행(기준환율등)
	OvisIsar                   string `json:"ovis_isar"`                      // 해외발행(발행지역)
	OvisMktnm                  string `json:"ovis_mktnm"`                     // 해외발행(해외상장시 시장 명칭)
	FdppFclt                   string `json:"fdpp_fclt"`                      // 자금조달목적(시설자금)
	FdppBsninh                 string `json:"fdpp_bsninh"`                    // 자금조달목적(영업양수자금)
	FdppOp                     string `json:"fdpp_op"`                        // 자금조달목적(운영자금)
	FdppDtrp                   string `json:"fdpp_dtrp"`                      // 자금조달목적(채무상환자금)
	FdppOcsa                   string `json:"fdpp_ocsa"`                      // 자금조달목적(타법인 증권 취득자금)
	FdppEtc                    string `json:"fdpp_etc"`                       // 자금조달목적(기타자금)
	BdIntrEx                   string `json:"bd_intr_ex"`                     // 사채의 이율(표면이자율 %)
	BdIntrSf                   string `json:"bd_intr_sf"`                     // 사채의 이율(만기이자율 %)
	BdMtd                      string `json:"bd_mtd"`                         // 사채만기일
	BdisMthn                   string `json:"bdis_mthn"`                      // 사채발행방법
	ExRt                       string `json:"ex_rt"`                          // 신주인수권 행사비율 (%)
	ExPrc                      string `json:"ex_prc"`                         // 신주인수권 행사가액 (원/주)
	ExPrcDmth                  string `json:"ex_prc_dmth"`                    // 신주인수권 행사가액 결정방법
	BdwtDivAtn                 string `json:"bdwt_div_atn"`                   // 사채와 인수권의 분리여부
	NstkPymMth                 string `json:"nstk_pym_mth"`                   // 신주대금 납입방법
	NstkIsstkKnd               string `json:"nstk_isstk_knd"`                 // 행사에 따라 발행할 주식(종류)
	NstkIsstkCnt               string `json:"nstk_isstk_cnt"`                 // 행사에 따라 발행할 주식(주식수)
	NstkIsstkTisstkVs          string `json:"nstk_isstk_tisstk_vs"`           // 행사에 따라 발행할 주식(주식총수 대비 %)
	ExpdBgd                    string `json:"expd_bgd"`                       // 권리행사기간(시작일)
	ExpdEdd                    string `json:"expd_edd"`                       // 권리행사기간(종료일)
	ActMktprcflCvprcLwtrsprc   string `json:"act_mktprcfl_cvprc_lwtrsprc"`    // 시가하락 행사가액 조정(최저 조정가액 원)
	ActMktprcflCvprcLwtrsprcBs string `json:"act_mktprcfl_cvprc_lwtrsprc_bs"` // 시가하락 행사가액 조정(최저 조정가액 근거)
	RmislmtLt70p               string `json:"rmislmt_lt70p"`                  // 시가하락 조정(행사가 70% 미만 조정가능 잔여한도 원)
	Abmg                       string `json:"abmg"`                           // 합병 관련 사항
	Sbd                        string `json:"sbd"`                            // 청약일
	Pymd                       string `json:"pymd"`                           // 납입일
	Rpmcmp                     string `json:"rpmcmp"`                         // 대표주관회사
	Grint                      string `json:"grint"`                          // 보증기관
	Bddd                       string `json:"bddd"`                           // 이사회결의일(결정일)
	OdAAtT                     string `json:"od_a_at_t"`                      // 사외이사 참석여부(참석)
	OdAAtB                     string `json:"od_a_at_b"`                      // 사외이사 참석여부(불참)
	AdtAAtn                    string `json:"adt_a_atn"`                      // 감사(감사위원) 참석여부
	RsSmAtn                    string `json:"rs_sm_atn"`                      // 증권신고서 제출대상 여부
	ExSmR                      string `json:"ex_sm_r"`                        // 제출 면제 사유
	OvisLtdtl                  string `json:"ovis_ltdtl"`                     // 해외발행 연계 대차거래 내역
	FtcSttAtn                  string `json:"ftc_stt_atn"`                    // 공정거래위원회 신고대상 여부
}

// BondWithWarrantIssuance 는 신주인수권부사채권 발행결정(주요사항보고서)을 조회한다.
func (c *Client) BondWithWarrantIssuance(ctx context.Context, p MaterialParams) ([]BondWithWarrantItem, error) {
	return httpclient.GetList[BondWithWarrantItem](ctx, c.http, "/api/bdwtIsDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestBondWithWarrantIssuance`
Expected: PASS.

- [ ] **Step 6: 실 API fixture 캡처(권장)**

Task 1 Step 7 과 동일 방식, 엔드포인트 `/api/bdwtIsDecsn.json`. 실 데이터 확보 시 fixture·assert 갱신 후 Step 5 재실행.

- [ ] **Step 7: gofmt + 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
gofmt -w material/bond.go
git add material/bond.go material/bond_test.go material/testdata/bdwtIsDecsn.json
git commit -m "feat(material): add DS005 BondWithWarrantIssuance (신주인수권부사채 발행결정)"
```

---

## Task 3: 교환사채권 발행결정 (ExchangeableBondIssuance)

**Files:**
- Modify: `material/bond.go`
- Create: `material/testdata/exbdIsDecsn.json`
- Modify: `material/bond_test.go`

- [ ] **Step 1: fixture 파일 작성**

`material/testdata/exbdIsDecsn.json` (스키마는 ExchangeableBondItem 와 1:1, atcsc_rmislmt 없음):

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230705000333",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트교환사채",
      "bd_tm": "3",
      "bd_knd": "기명식 무보증 사모 교환사채",
      "bd_fta": "50,000,000,000",
      "ovis_fta": "-",
      "ovis_fta_crn": "-",
      "ovis_ster": "-",
      "ovis_isar": "-",
      "ovis_mktnm": "-",
      "fdpp_fclt": "-",
      "fdpp_bsninh": "-",
      "fdpp_op": "-",
      "fdpp_dtrp": "50,000,000,000",
      "fdpp_ocsa": "-",
      "fdpp_etc": "-",
      "bd_intr_ex": "0.0",
      "bd_intr_sf": "1.5",
      "bd_mtd": "2027년 07월 06일",
      "bdis_mthn": "사모",
      "ex_rt": "100",
      "ex_prc": "25,000",
      "ex_prc_dmth": "이사회결의일 전 기산일 가중산술평균주가",
      "extg": "자기주식(기명식 보통주)",
      "extg_stkcnt": "2,000,000",
      "extg_tisstk_vs": "2.8",
      "exrqpd_bgd": "2024년 07월 06일",
      "exrqpd_edd": "2027년 06월 06일",
      "sbd": "2023년 07월 06일",
      "pymd": "2023년 07월 06일",
      "rpmcmp": "-",
      "grint": "-",
      "bddd": "2023년 07월 05일",
      "od_a_at_t": "3",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "rs_sm_atn": "미제출",
      "ex_sm_r": "사모발행",
      "ovis_ltdtl": "-",
      "ftc_stt_atn": "미해당"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가 (실패 확인용)**

`material/bond_test.go` 에 함수 추가:

```go
func TestExchangeableBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/exbdIsDecsn.json": "testdata/exbdIsDecsn.json",
	})

	items, err := c.ExchangeableBondIssuance(context.Background(), MaterialParams{CorpCode: "00126380"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230705000333", got.RceptNo)
	assert.Equal(t, "25,000", got.ExPrc)
	assert.Equal(t, "자기주식(기명식 보통주)", got.Extg)
	assert.Equal(t, "2024년 07월 06일", got.ExrqpdBgd)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestExchangeableBondIssuance`
Expected: FAIL — `c.ExchangeableBondIssuance undefined`.

- [ ] **Step 4: bond.go 에 struct + 메서드 추가**

`material/bond.go` 끝에 추가:

```go
// ExchangeableBondItem 은 교환사채권 발행결정 (exbdIsDecsn) 한 건.
type ExchangeableBondItem struct {
	RceptNo      string `json:"rcept_no"`       // 접수번호
	CorpCls      string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`      // 고유번호
	CorpName     string `json:"corp_name"`      // 회사명
	BdTm         string `json:"bd_tm"`          // 사채의 종류(회차)
	BdKnd        string `json:"bd_knd"`         // 사채의 종류(종류)
	BdFta        string `json:"bd_fta"`         // 사채의 권면(전자등록)총액 (원)
	OvisFta      string `json:"ovis_fta"`       // 해외발행(권면총액)
	OvisFtaCrn   string `json:"ovis_fta_crn"`   // 해외발행(권면총액 통화단위)
	OvisSter     string `json:"ovis_ster"`      // 해외발행(기준환율등)
	OvisIsar     string `json:"ovis_isar"`      // 해외발행(발행지역)
	OvisMktnm    string `json:"ovis_mktnm"`     // 해외발행(해외상장시 시장 명칭)
	FdppFclt     string `json:"fdpp_fclt"`      // 자금조달목적(시설자금)
	FdppBsninh   string `json:"fdpp_bsninh"`    // 자금조달목적(영업양수자금)
	FdppOp       string `json:"fdpp_op"`        // 자금조달목적(운영자금)
	FdppDtrp     string `json:"fdpp_dtrp"`      // 자금조달목적(채무상환자금)
	FdppOcsa     string `json:"fdpp_ocsa"`      // 자금조달목적(타법인 증권 취득자금)
	FdppEtc      string `json:"fdpp_etc"`       // 자금조달목적(기타자금)
	BdIntrEx     string `json:"bd_intr_ex"`     // 사채의 이율(표면이자율 %)
	BdIntrSf     string `json:"bd_intr_sf"`     // 사채의 이율(만기이자율 %)
	BdMtd        string `json:"bd_mtd"`         // 사채만기일
	BdisMthn     string `json:"bdis_mthn"`      // 사채발행방법
	ExRt         string `json:"ex_rt"`          // 교환비율 (%)
	ExPrc        string `json:"ex_prc"`         // 교환가액 (원/주)
	ExPrcDmth    string `json:"ex_prc_dmth"`    // 교환가액 결정방법
	Extg         string `json:"extg"`           // 교환대상(종류)
	ExtgStkcnt   string `json:"extg_stkcnt"`    // 교환대상(주식수)
	ExtgTisstkVs string `json:"extg_tisstk_vs"` // 교환대상(주식총수 대비 %)
	ExrqpdBgd    string `json:"exrqpd_bgd"`     // 교환청구기간(시작일)
	ExrqpdEdd    string `json:"exrqpd_edd"`     // 교환청구기간(종료일)
	Sbd          string `json:"sbd"`            // 청약일
	Pymd         string `json:"pymd"`           // 납입일
	Rpmcmp       string `json:"rpmcmp"`         // 대표주관회사
	Grint        string `json:"grint"`          // 보증기관
	Bddd         string `json:"bddd"`           // 이사회결의일(결정일)
	OdAAtT       string `json:"od_a_at_t"`      // 사외이사 참석여부(참석)
	OdAAtB       string `json:"od_a_at_b"`      // 사외이사 참석여부(불참)
	AdtAAtn      string `json:"adt_a_atn"`      // 감사(감사위원) 참석여부
	RsSmAtn      string `json:"rs_sm_atn"`      // 증권신고서 제출대상 여부
	ExSmR        string `json:"ex_sm_r"`        // 제출 면제 사유
	OvisLtdtl    string `json:"ovis_ltdtl"`     // 해외발행 연계 대차거래 내역
	FtcSttAtn    string `json:"ftc_stt_atn"`    // 공정거래위원회 신고대상 여부
}

// ExchangeableBondIssuance 는 교환사채권 발행결정(주요사항보고서)을 조회한다.
func (c *Client) ExchangeableBondIssuance(ctx context.Context, p MaterialParams) ([]ExchangeableBondItem, error) {
	return httpclient.GetList[ExchangeableBondItem](ctx, c.http, "/api/exbdIsDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestExchangeableBondIssuance`
Expected: PASS.

- [ ] **Step 6: 실 API fixture 캡처(권장)**

Task 1 Step 7 과 동일, 엔드포인트 `/api/exbdIsDecsn.json`.

- [ ] **Step 7: gofmt + 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
gofmt -w material/bond.go
git add material/bond.go material/bond_test.go material/testdata/exbdIsDecsn.json
git commit -m "feat(material): add DS005 ExchangeableBondIssuance (교환사채 발행결정)"
```

---

## Task 4: 상각형 조건부자본증권 발행결정 (ContingentConvertibleBondIssuance)

**Files:**
- Modify: `material/bond.go`
- Create: `material/testdata/wdCocobdIsDecsn.json`
- Modify: `material/bond_test.go`

- [ ] **Step 1: fixture 파일 작성**

`material/testdata/wdCocobdIsDecsn.json` (스키마는 ContingentConvertibleBondItem 와 1:1; bd_intr_sf=표면이자율, bd_intr_ex=만기이자율 — docs 표기 그대로):

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230601000444",
      "corp_cls": "Y",
      "corp_code": "00164779",
      "corp_name": "테스트조건부자본증권",
      "bd_tm": "1",
      "bd_knd": "기명식 무보증 상각형 조건부자본증권",
      "bd_fta": "100,000,000,000",
      "ovis_fta": "-",
      "ovis_fta_crn": "-",
      "ovis_ster": "-",
      "ovis_isar": "-",
      "ovis_mktnm": "-",
      "fdpp_fclt": "-",
      "fdpp_bsninh": "-",
      "fdpp_op": "100,000,000,000",
      "fdpp_dtrp": "-",
      "fdpp_ocsa": "-",
      "fdpp_etc": "-",
      "bd_intr_sf": "4.5",
      "bd_intr_ex": "4.5",
      "bd_mtd": "2053년 06월 02일",
      "dbtrs_sc": "발행회사가 부실금융기관으로 지정되는 경우 사채의 상환 및 이자지급의무가 전액 영구적으로 감면됨",
      "sbd": "2023년 06월 02일",
      "pymd": "2023년 06월 02일",
      "rpmcmp": "-",
      "grint": "-",
      "bddd": "2023년 06월 01일",
      "od_a_at_t": "4",
      "od_a_at_b": "0",
      "adt_a_atn": "1",
      "rs_sm_atn": "제출",
      "ex_sm_r": "-",
      "ovis_ltdtl": "-",
      "ftc_stt_atn": "미해당"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가 (실패 확인용)**

`material/bond_test.go` 에 함수 추가:

```go
func TestContingentConvertibleBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/wdCocobdIsDecsn.json": "testdata/wdCocobdIsDecsn.json",
	})

	items, err := c.ContingentConvertibleBondIssuance(context.Background(), MaterialParams{CorpCode: "00164779"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230601000444", got.RceptNo)
	assert.Equal(t, "100,000,000,000", got.BdFta)
	assert.Equal(t, "4.5", got.BdIntrSf)
	assert.Contains(t, got.DbtrsSc, "부실금융기관")
	assert.Equal(t, "제출", got.RsSmAtn)
}
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestContingentConvertibleBondIssuance`
Expected: FAIL — `c.ContingentConvertibleBondIssuance undefined`.

- [ ] **Step 4: bond.go 에 struct + 메서드 추가**

`material/bond.go` 끝에 추가:

```go
// ContingentConvertibleBondItem 은 상각형 조건부자본증권 발행결정 (wdCocobdIsDecsn) 한 건.
type ContingentConvertibleBondItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	BdTm       string `json:"bd_tm"`        // 사채의 종류(회차)
	BdKnd      string `json:"bd_knd"`       // 사채의 종류(종류)
	BdFta      string `json:"bd_fta"`       // 사채의 권면(전자등록)총액 (원)
	OvisFta    string `json:"ovis_fta"`     // 해외발행(권면총액)
	OvisFtaCrn string `json:"ovis_fta_crn"` // 해외발행(권면총액 통화단위)
	OvisSter   string `json:"ovis_ster"`    // 해외발행(기준환율등)
	OvisIsar   string `json:"ovis_isar"`    // 해외발행(발행지역)
	OvisMktnm  string `json:"ovis_mktnm"`   // 해외발행(해외상장시 시장 명칭)
	FdppFclt   string `json:"fdpp_fclt"`    // 자금조달목적(시설자금)
	FdppBsninh string `json:"fdpp_bsninh"`  // 자금조달목적(영업양수자금)
	FdppOp     string `json:"fdpp_op"`      // 자금조달목적(운영자금)
	FdppDtrp   string `json:"fdpp_dtrp"`    // 자금조달목적(채무상환자금)
	FdppOcsa   string `json:"fdpp_ocsa"`    // 자금조달목적(타법인 증권 취득자금)
	FdppEtc    string `json:"fdpp_etc"`     // 자금조달목적(기타자금)
	BdIntrSf   string `json:"bd_intr_sf"`   // 사채의 이율(표면이자율 %)
	BdIntrEx   string `json:"bd_intr_ex"`   // 사채의 이율(만기이자율 %)
	BdMtd      string `json:"bd_mtd"`       // 사채만기일
	DbtrsSc    string `json:"dbtrs_sc"`     // 채무재조정의 범위
	Sbd        string `json:"sbd"`          // 청약일
	Pymd       string `json:"pymd"`         // 납입일
	Rpmcmp     string `json:"rpmcmp"`       // 대표주관회사
	Grint      string `json:"grint"`        // 보증기관
	Bddd       string `json:"bddd"`         // 이사회결의일(결정일)
	OdAAtT     string `json:"od_a_at_t"`    // 사외이사 참석여부(참석)
	OdAAtB     string `json:"od_a_at_b"`    // 사외이사 참석여부(불참)
	AdtAAtn    string `json:"adt_a_atn"`    // 감사(감사위원) 참석여부
	RsSmAtn    string `json:"rs_sm_atn"`    // 증권신고서 제출대상 여부
	ExSmR      string `json:"ex_sm_r"`      // 제출 면제 사유
	OvisLtdtl  string `json:"ovis_ltdtl"`   // 해외발행 연계 대차거래 내역
	FtcSttAtn  string `json:"ftc_stt_atn"`  // 공정거래위원회 신고대상 여부
}

// ContingentConvertibleBondIssuance 는 상각형 조건부자본증권 발행결정(주요사항보고서)을 조회한다.
func (c *Client) ContingentConvertibleBondIssuance(ctx context.Context, p MaterialParams) ([]ContingentConvertibleBondItem, error) {
	return httpclient.GetList[ContingentConvertibleBondItem](ctx, c.http, "/api/wdCocobdIsDecsn.json", p.toMap())
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test ./material/ -run TestContingentConvertibleBondIssuance`
Expected: PASS.

- [ ] **Step 6: 실 API fixture 캡처(권장)**

Task 1 Step 7 과 동일, 엔드포인트 `/api/wdCocobdIsDecsn.json`. 상각형 조건부자본증권은 주로 금융지주·은행 발행(예: corp_code 탐색 시 금융사 위주).

- [ ] **Step 7: gofmt + 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
gofmt -w material/bond.go
git add material/bond.go material/bond_test.go material/testdata/wdCocobdIsDecsn.json
git commit -m "feat(material): add DS005 ContingentConvertibleBondIssuance (상각형 조건부자본증권 발행결정)"
```

---

## Task 5: 통합 테스트 + README

**Files:**
- Modify: `integration_test.go`
- Modify: `README.md:33-34`

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 에 함수 추가(기존 통합 테스트의 클라이언트 생성 패턴 — `newIntegrationClient(t)` 또는 동일 헬퍼 — 을 따른다. 파일을 먼저 읽어 기존 헬퍼/패턴을 확인하고 그대로 사용할 것). 예시(기존 `TestIntegration_DefaultOccurrences` 형태에 맞춤):

```go
func TestIntegration_ConvertibleBondIssuance(t *testing.T) {
	c := newIntegrationClient(t)
	// 전환사채 발행 사례가 있는 종목/기간으로 조회 (없으면 ErrNoData 허용)
	items, err := c.Material.ConvertibleBondIssuance(context.Background(), material.MaterialParams{
		CorpCode: "00126380",
		BgnDe:    "20200101",
		EndDe:    "20241231",
	})
	if errors.Is(err, opendart.ErrNoData) {
		t.Skip("해당 기간 전환사채 발행 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		assert.NotEmpty(t, it.RceptNo)
	}
}
```

주의: import 와 헬퍼 이름은 기존 `integration_test.go` 와 정확히 일치시킬 것. 기존 파일에 `newIntegrationClient` 가 없으면 그 파일에서 실제 사용하는 생성 방식을 그대로 복제한다.

- [ ] **Step 2: 통합 테스트 컴파일 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...`
Expected: 출력 없음(컴파일·vet 통과).

- [ ] **Step 3: 통합 테스트 실행(키 있으면)**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go test -tags integration -run TestIntegration_ConvertibleBondIssuance ./...`
Expected: PASS 또는 SKIP(데이터 없음). 실패하면 fixture/필드 점검.

- [ ] **Step 4: README 커버리지 갱신**

`README.md` 의 DS005 줄(`README.md:33`)에 사채 발행 추가, 예정 줄(`README.md:34`)에서 사채발행 제거:

```markdown
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정)
- (예정) DS005 나머지(자기주식/양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```

- [ ] **Step 5: 전체 테스트 + 빌드 확인**

Run: `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/`
Expected: 빌드 성공, 전체 테스트 PASS, gofmt 출력 없음.

- [ ] **Step 6: 커밋**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add integration_test.go README.md
git commit -m "test(material): add DS005 사채 발행 통합 테스트 + README 커버리지"
```

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 4개 메서드(ConvertibleBondIssuance/BondWithWarrantIssuance/ExchangeableBondIssuance/ContingentConvertibleBondIssuance) → Task 1~4 각각 매핑. struct 4종 전 필드 = spec 의 4 struct 와 1:1. MaterialParams+GetList 재사용·root 무변경 = Architecture 일치. 테스트 전략(fixture+통합 1개)·README 갱신 = Task 5. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 모든 코드 step 에 완전한 코드 포함. fixture 는 "실 API 캡처로 교체(권장)" 명시 — DS005 부실/증자 plan 과 동일 정책(드문 이벤트는 스키마 일치 샘플 허용).

**3. Type consistency:** 메서드명·struct명·필드명이 Task 간 일관(전부 bond.go 단일 파일에 누적). 상각형의 bd_intr_sf=표면/bd_intr_ex=만기 순서는 docs 표기 그대로(다른 3종과 의미 반대 — 의도된 것, spec 에 명시됨).
