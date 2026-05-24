# OpenDART DS005 양수도 Sub-1 «실물 양수도» Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 양수도 그룹 중 실물 양수도 8개 API(영업/유형자산/타법인주식/주권사채권 양수·양도 결정)를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 를 그대로 재사용한다. 신규 파일 `material/transfer.go` 에 8개 item struct + 8개 한 줄 메서드를 추가한다. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-transfer-design.md` 에 8개 struct 전체 필드가 EXACT Go 코드(json 태그 + 한글 코멘트)로 정의돼 있다. 각 Task 의 struct 는 그 파일의 동일 이름 struct 를 **그대로 복사**한다(필드 가감 금지). 아래 각 Task 는 fixture·test·메서드 시그니처를 명시하고 struct 는 spec 위치를 가리킨다.

---

## File Structure

- Create: `material/transfer.go` — 8 item struct + 8 메서드.
- Create: `material/transfer_test.go` — fixture 디코딩 테스트(기존 `newTestClient` 재사용).
- Create: `material/testdata/{bsnInhDecsn,bsnTrfDecsn,tgastInhDecsn,tgastTrfDecsn,otcprStkInvscrInhDecsn,otcprStkInvscrTrfDecsn,stkrtbdInhDecsn,stkrtbdTrfDecsn}.json` — 8 fixture.
- Modify: `integration_test.go` — 통합 케이스 1~2개(`//go:build integration`, ErrNoData skip).
- Modify: `README.md` — DS005 커버리지에 양수도(Sub-1) 추가.

기존 컨벤션(변경 금지):
- `material/treasury.go` 메서드 형태: `func (c *Client) TreasuryStockAcquisition(ctx context.Context, p MaterialParams) ([]TreasuryStockAcquisitionItem, error) { return httpclient.GetList[TreasuryStockAcquisitionItem](ctx, c.http, "/api/tsstkAqDecsn.json", p.toMap()) }`. Client http 필드는 `c.http`.
- `newTestClient(t, routes map[string]string)`: route map 값은 **bare 파일명**(예: `"bsnInhDecsn.json"`) — 내부에서 `filepath.Join("testdata", fixture)` 로 testdata/ 를 붙인다.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

**공통 작업 절차 (모든 Task 1~8 동일):**
1. fixture `material/testdata/<code>.json` 작성(아래 Task별 JSON).
2. `material/transfer_test.go` 에 테스트 함수 추가(아래 Task별). (Task 1 은 파일을 새로 만들고 `package material` + import 블록 포함; Task 2~8 은 함수만 append.)
3. `go test ./material/ -run <TestName>` → FAIL(undefined) 확인.
4. `material/transfer.go` 에 spec 의 해당 struct + 메서드 추가. (Task 1 은 파일을 새로 만들고 `package material` + `import ("context"; "github.com/kenshin579/opendart/internal/httpclient")` 포함; Task 2~8 은 append.)
5. `go test ./material/ -run <TestName>` → PASS.
6. `gofmt -l material/transfer.go material/transfer_test.go` → 출력 없으면 OK(있으면 `gofmt -w`).
7. 실 API 캡처 BEST EFFORT: `$OPENDART_API_KEY` 있으면 `curl -s "https://opendart.fss.or.kr/api/<code>.json?crtfc_key=$OPENDART_API_KEY&corp_code=<corp>&bgn_de=20200101&end_de=20241231"`. 서버 TLS1.2 RSA cipher 만 지원(curl handshake 실패 가능). clean `"status":"000"` + non-empty list 면 fixture 교체 + assert 갱신(+실 날짜 bracket) 후 step 5 재실행. 1회 시도, 실패 시 샘플 유지.
8. 커밋(no push): `git add material/transfer.go material/transfer_test.go material/testdata/<code>.json && git commit -m "feat(material): add DS005 <Method> (<한글>)"`.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 영업양수 결정 (BusinessAcquisition) — bsnInhDecsn, 48필드

- Test: `TestBusinessAcquisition`. Struct: spec 의 `BusinessAcquisitionItem` (48필드) 그대로. Method:
```go
// BusinessAcquisition 은 영업양수 결정(주요사항보고서)을 조회한다.
func (c *Client) BusinessAcquisition(ctx context.Context, p MaterialParams) ([]BusinessAcquisitionItem, error) {
	return httpclient.GetList[BusinessAcquisitionItem](ctx, c.http, "/api/bsnInhDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/bsnInhDecsn.json`**

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트영업양수",
      "inh_bsn": "반도체 사업부문", "inh_bsn_mc": "메모리 반도체 제조", "inh_prc": "500,000,000,000", "absn_inh_atn": "아니오",
      "ast_inh_bsn": "300,000,000,000", "ast_cmp_all": "3,000,000,000,000", "ast_rt": "10.0",
      "sl_inh_bsn": "200,000,000,000", "sl_cmp_all": "2,000,000,000,000", "sl_rt": "10.0",
      "dbt_inh_bsn": "100,000,000,000", "dbt_cmp_all": "1,000,000,000,000", "dbt_rt": "10.0",
      "inh_pp": "사업 경쟁력 강화", "inh_af": "매출 증대 기대",
      "inh_prd_ctr_cnsd": "2023년 04월 10일", "inh_prd_inh_std": "2023년 06월 30일",
      "dlptn_cmpnm": "상대회사", "dlptn_cpt": "50,000,000,000", "dlptn_mbsn": "반도체", "dlptn_hoadd": "경기도 수원시", "dlptn_rl_cmpn": "특수관계 없음",
      "inh_pym": "현금 지급",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법 제165조의4", "exevl_intn": "삼일회계법인", "exevl_pd": "2023년 03월", "exevl_op": "적정",
      "gmtsck_spd_atn": "예", "gmtsck_prd": "2023년 05월 30일",
      "aprskh_plnprc": "70,000", "aprskh_pym_plpd_mth": "2023년 06월, 현금", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bddd": "2023년 04월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "bdlst_atn": "아니오", "n6m_tpai_plann": "없음", "otcpr_bdlst_sf_atn": "해당없음",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (transfer_test.go 신규 — `package material` + import `context`,`testing`,`github.com/stretchr/testify/assert`,`github.com/stretchr/testify/require`)

```go
func TestBusinessAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnInhDecsn.json": "bsnInhDecsn.json"})
	items, err := c.BusinessAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230410000111", got.RceptNo)
	assert.Equal(t, "반도체 사업부문", got.InhBsn)
	assert.Equal(t, "500,000,000,000", got.InhPrc)
	assert.Equal(t, "삼일회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차(위). corp 캡처는 `00126380`. struct 는 spec `BusinessAcquisitionItem` 복사.

---

## Task 2: 영업양도 결정 (BusinessTransfer) — bsnTrfDecsn, 41필드

- Test: `TestBusinessTransfer`. Struct: spec `BusinessTransferItem` (41필드; 부채액 dbt_* 없음). Method:
```go
// BusinessTransfer 는 영업양도 결정(주요사항보고서)을 조회한다.
func (c *Client) BusinessTransfer(ctx context.Context, p MaterialParams) ([]BusinessTransferItem, error) {
	return httpclient.GetList[BusinessTransferItem](ctx, c.http, "/api/bsnTrfDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/bsnTrfDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20230510000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트영업양도",
      "trf_bsn": "디스플레이 사업부문", "trf_bsn_mc": "LCD 패널 제조", "trf_prc": "300,000,000,000",
      "ast_trf_bsn": "200,000,000,000", "ast_cmp_all": "3,000,000,000,000", "ast_rt": "6.7",
      "sl_trf_bsn": "150,000,000,000", "sl_cmp_all": "2,000,000,000,000", "sl_rt": "7.5",
      "trf_pp": "사업 구조조정", "trf_af": "비핵심사업 정리",
      "trf_prd_ctr_cnsd": "2023년 05월 10일", "trf_prd_trf_std": "2023년 07월 31일",
      "dlptn_cmpnm": "양수회사", "dlptn_cpt": "30,000,000,000", "dlptn_mbsn": "디스플레이", "dlptn_hoadd": "충청남도 아산시", "dlptn_rl_cmpn": "특수관계 없음",
      "trf_pym": "현금 수령",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법 제165조의4", "exevl_intn": "안진회계법인", "exevl_pd": "2023년 04월", "exevl_op": "적정",
      "gmtsck_spd_atn": "예", "gmtsck_prd": "2023년 06월 30일",
      "aprskh_plnprc": "68,000", "aprskh_pym_plpd_mth": "2023년 07월, 현금", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bddd": "2023년 05월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestBusinessTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnTrfDecsn.json": "bsnTrfDecsn.json"})
	items, err := c.BusinessTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230510000222", got.RceptNo)
	assert.Equal(t, "디스플레이 사업부문", got.TrfBsn)
	assert.Equal(t, "300,000,000,000", got.TrfPrc)
	assert.Equal(t, "안진회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `BusinessTransferItem` 복사.

---

## Task 3: 유형자산 양수 결정 (TangibleAssetAcquisition) — tgastInhDecsn, 40필드

- Test: `TestTangibleAssetAcquisition`. Struct: spec `TangibleAssetAcquisitionItem` (40필드). Method:
```go
// TangibleAssetAcquisition 은 유형자산 양수 결정(주요사항보고서)을 조회한다.
func (c *Client) TangibleAssetAcquisition(ctx context.Context, p MaterialParams) ([]TangibleAssetAcquisitionItem, error) {
	return httpclient.GetList[TangibleAssetAcquisitionItem](ctx, c.http, "/api/tgastInhDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/tgastInhDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20230610000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트유형자산양수",
      "ast_sen": "토지 및 건물", "ast_nm": "평택 공장부지",
      "inhdtl_inhprc": "150,000,000,000", "inhdtl_tast": "3,000,000,000,000", "inhdtl_tast_vs": "5.0",
      "inh_pp": "생산능력 확충", "inh_af": "생산설비 증대",
      "inh_prd_ctr_cnsd": "2023년 06월 10일", "inh_prd_inh_std": "2023년 07월 31일", "inh_prd_rgs_prd": "2023년 08월 15일",
      "dlptn_cmpnm": "매도법인", "dlptn_cpt": "20,000,000,000", "dlptn_mbsn": "부동산", "dlptn_hoadd": "경기도 평택시", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 지급",
      "exevl_atn": "예", "exevl_bs_rs": "감정평가", "exevl_intn": "한국감정원", "exevl_pd": "2023년 05월", "exevl_op": "적정",
      "gmtsck_spd_atn": "아니오", "gmtsck_prd": "-",
      "aprskh_exrq": "-", "aprskh_plnprc": "-", "aprskh_ex_pc_mth_pd_pl": "-", "aprskh_pym_plpd_mth": "-", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bddd": "2023년 06월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestTangibleAssetAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tgastInhDecsn.json": "tgastInhDecsn.json"})
	items, err := c.TangibleAssetAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230610000333", got.RceptNo)
	assert.Equal(t, "토지 및 건물", got.AstSen)
	assert.Equal(t, "150,000,000,000", got.InhdtlInhprc)
	assert.Equal(t, "한국감정원", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `TangibleAssetAcquisitionItem` 복사.

---

## Task 4: 유형자산 양도 결정 (TangibleAssetTransfer) — tgastTrfDecsn, 40필드

- Test: `TestTangibleAssetTransfer`. Struct: spec `TangibleAssetTransferItem` (40필드). Method:
```go
// TangibleAssetTransfer 는 유형자산 양도 결정(주요사항보고서)을 조회한다.
func (c *Client) TangibleAssetTransfer(ctx context.Context, p MaterialParams) ([]TangibleAssetTransferItem, error) {
	return httpclient.GetList[TangibleAssetTransferItem](ctx, c.http, "/api/tgastTrfDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/tgastTrfDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20230710000444", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트유형자산양도",
      "ast_sen": "토지", "ast_nm": "구미 공장부지",
      "trfdtl_trfprc": "80,000,000,000", "trfdtl_tast": "3,000,000,000,000", "trfdtl_tast_vs": "2.7",
      "trf_pp": "유휴자산 처분", "trf_af": "자금 확보",
      "trf_prd_ctr_cnsd": "2023년 07월 10일", "trf_prd_trf_std": "2023년 08월 31일", "trf_prd_rgs_prd": "2023년 09월 15일",
      "dlptn_cmpnm": "매수법인", "dlptn_cpt": "10,000,000,000", "dlptn_mbsn": "부동산개발", "dlptn_hoadd": "경상북도 구미시", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 수령",
      "exevl_atn": "예", "exevl_bs_rs": "감정평가", "exevl_intn": "한국감정원", "exevl_pd": "2023년 06월", "exevl_op": "적정",
      "gmtsck_spd_atn": "아니오", "gmtsck_prd": "-",
      "aprskh_exrq": "-", "aprskh_plnprc": "-", "aprskh_ex_pc_mth_pd_pl": "-", "aprskh_pym_plpd_mth": "-", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bddd": "2023년 07월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestTangibleAssetTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tgastTrfDecsn.json": "tgastTrfDecsn.json"})
	items, err := c.TangibleAssetTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230710000444", got.RceptNo)
	assert.Equal(t, "토지", got.AstSen)
	assert.Equal(t, "80,000,000,000", got.TrfdtlTrfprc)
	assert.Equal(t, "한국감정원", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `TangibleAssetTransferItem` 복사.

---

## Task 5: 타법인 주식·출자증권 양수결정 (OtherCorpStockAcquisition) — otcprStkInvscrInhDecsn, 43필드

- Test: `TestOtherCorpStockAcquisition`. Struct: spec `OtherCorpStockAcquisitionItem` (43필드). Method:
```go
// OtherCorpStockAcquisition 은 타법인 주식 및 출자증권 양수결정(주요사항보고서)을 조회한다.
func (c *Client) OtherCorpStockAcquisition(ctx context.Context, p MaterialParams) ([]OtherCorpStockAcquisitionItem, error) {
	return httpclient.GetList[OtherCorpStockAcquisitionItem](ctx, c.http, "/api/otcprStkInvscrInhDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/otcprStkInvscrInhDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20230810000555", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트타법인주식양수",
      "iscmp_cmpnm": "대상회사", "iscmp_nt": "대한민국", "iscmp_rp": "홍길동", "iscmp_cpt": "10,000,000,000", "iscmp_rl_cmpn": "특수관계 없음", "iscmp_tisstk": "2,000,000", "iscmp_mbsn": "소프트웨어",
      "l6m_tpa_nstkaq_atn": "아니오",
      "inhdtl_stkcnt": "1,000,000", "inhdtl_inhprc": "50,000,000,000", "inhdtl_tast": "3,000,000,000,000", "inhdtl_tast_vs": "1.7", "inhdtl_ecpt": "2,000,000,000,000", "inhdtl_ecpt_vs": "2.5",
      "atinh_owstkcnt": "1,000,000", "atinh_eqrt": "50.0",
      "inh_pp": "사업 시너지", "inh_prd": "2023년 09월 30일",
      "dlptn_cmpnm": "매도주주", "dlptn_cpt": "-", "dlptn_mbsn": "-", "dlptn_hoadd": "서울시 강남구", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 지급",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "삼정회계법인", "exevl_pd": "2023년 07월", "exevl_op": "적정",
      "bddd": "2023년 08월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "bdlst_atn": "아니오", "n6m_tpai_plann": "없음", "iscmp_bdlst_sf_atn": "해당없음",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestOtherCorpStockAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/otcprStkInvscrInhDecsn.json": "otcprStkInvscrInhDecsn.json"})
	items, err := c.OtherCorpStockAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230810000555", got.RceptNo)
	assert.Equal(t, "대상회사", got.IscmpCmpnm)
	assert.Equal(t, "1,000,000", got.InhdtlStkcnt)
	assert.Equal(t, "50.0", got.AtinhEqrt)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `OtherCorpStockAcquisitionItem` 복사.

---

## Task 6: 타법인 주식·출자증권 양도결정 (OtherCorpStockTransfer) — otcprStkInvscrTrfDecsn, 39필드

- Test: `TestOtherCorpStockTransfer`. Struct: spec `OtherCorpStockTransferItem` (39필드). Method:
```go
// OtherCorpStockTransfer 는 타법인 주식 및 출자증권 양도결정(주요사항보고서)을 조회한다.
func (c *Client) OtherCorpStockTransfer(ctx context.Context, p MaterialParams) ([]OtherCorpStockTransferItem, error) {
	return httpclient.GetList[OtherCorpStockTransferItem](ctx, c.http, "/api/otcprStkInvscrTrfDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/otcprStkInvscrTrfDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20230910000666", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트타법인주식양도",
      "iscmp_cmpnm": "처분대상회사", "iscmp_nt": "대한민국", "iscmp_rp": "김철수", "iscmp_cpt": "5,000,000,000", "iscmp_rl_cmpn": "특수관계 없음", "iscmp_tisstk": "1,000,000", "iscmp_mbsn": "유통",
      "trfdtl_stkcnt": "500,000", "trfdtl_trfprc": "30,000,000,000", "trfdtl_tast": "3,000,000,000,000", "trfdtl_tast_vs": "1.0", "trfdtl_ecpt": "2,000,000,000,000", "trfdtl_ecpt_vs": "1.5",
      "attrf_owstkcnt": "0", "attrf_eqrt": "0.0",
      "trf_pp": "투자자산 회수", "trf_prd": "2023년 10월 31일",
      "dlptn_cmpnm": "매수자", "dlptn_cpt": "-", "dlptn_mbsn": "-", "dlptn_hoadd": "서울시 서초구", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 수령",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "한영회계법인", "exevl_pd": "2023년 08월", "exevl_op": "적정",
      "bddd": "2023년 09월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestOtherCorpStockTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/otcprStkInvscrTrfDecsn.json": "otcprStkInvscrTrfDecsn.json"})
	items, err := c.OtherCorpStockTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230910000666", got.RceptNo)
	assert.Equal(t, "처분대상회사", got.IscmpCmpnm)
	assert.Equal(t, "500,000", got.TrfdtlStkcnt)
	assert.Equal(t, "한영회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `OtherCorpStockTransferItem` 복사.

---

## Task 7: 주권 관련 사채권 양수 결정 (StockRelatedBondAcquisition) — stkrtbdInhDecsn, 41필드

- Test: `TestStockRelatedBondAcquisition`. Struct: spec `StockRelatedBondAcquisitionItem` (41필드; l6m_tpa 있음, aqd 없음). Method:
```go
// StockRelatedBondAcquisition 은 주권 관련 사채권 양수 결정(주요사항보고서)을 조회한다.
func (c *Client) StockRelatedBondAcquisition(ctx context.Context, p MaterialParams) ([]StockRelatedBondAcquisitionItem, error) {
	return httpclient.GetList[StockRelatedBondAcquisitionItem](ctx, c.http, "/api/stkrtbdInhDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/stkrtbdInhDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20231010000777", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트사채권양수",
      "stkrtbd_kndn": "전환사채", "tm": "5", "knd": "기명식 무보증 전환사채",
      "bdiscmp_cmpnm": "발행회사", "bdiscmp_nt": "대한민국", "bdiscmp_rp": "이영희", "bdiscmp_cpt": "8,000,000,000", "bdiscmp_rl_cmpn": "특수관계 없음", "bdiscmp_tisstk": "1,500,000", "bdiscmp_mbsn": "바이오",
      "l6m_tpa_nstkaq_atn": "아니오",
      "inhdtl_bd_fta": "20,000,000,000", "inhdtl_inhprc": "20,000,000,000", "inhdtl_tast": "3,000,000,000,000", "inhdtl_tast_vs": "0.7", "inhdtl_ecpt": "2,000,000,000,000", "inhdtl_ecpt_vs": "1.0",
      "inh_pp": "전략적 투자", "inh_prd": "2023년 11월 30일",
      "dlptn_cmpnm": "매도자", "dlptn_cpt": "-", "dlptn_mbsn": "-", "dlptn_hoadd": "대전시 유성구", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 지급",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "삼일회계법인", "exevl_pd": "2023년 09월", "exevl_op": "적정",
      "bddd": "2023년 10월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestStockRelatedBondAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stkrtbdInhDecsn.json": "stkrtbdInhDecsn.json"})
	items, err := c.StockRelatedBondAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231010000777", got.RceptNo)
	assert.Equal(t, "전환사채", got.StkrtbdKndn)
	assert.Equal(t, "20,000,000,000", got.InhdtlBdFta)
	assert.Equal(t, "삼일회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `StockRelatedBondAcquisitionItem` 복사.

---

## Task 8: 주권 관련 사채권 양도 결정 (StockRelatedBondTransfer) — stkrtbdTrfDecsn, 41필드

- Test: `TestStockRelatedBondTransfer`. Struct: spec `StockRelatedBondTransferItem` (41필드; aqd 있음, l6m_tpa 없음). Method:
```go
// StockRelatedBondTransfer 는 주권 관련 사채권 양도 결정(주요사항보고서)을 조회한다.
func (c *Client) StockRelatedBondTransfer(ctx context.Context, p MaterialParams) ([]StockRelatedBondTransferItem, error) {
	return httpclient.GetList[StockRelatedBondTransferItem](ctx, c.http, "/api/stkrtbdTrfDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/stkrtbdTrfDecsn.json`**

```json
{
  "status": "000", "message": "정상",
  "list": [
    {
      "rcept_no": "20231110000888", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트사채권양도",
      "stkrtbd_kndn": "신주인수권부사채", "tm": "3", "knd": "기명식 무보증 신주인수권부사채", "aqd": "2022년 05월 01일",
      "bdiscmp_cmpnm": "발행사", "bdiscmp_nt": "대한민국", "bdiscmp_rp": "박민수", "bdiscmp_cpt": "6,000,000,000", "bdiscmp_rl_cmpn": "특수관계 없음", "bdiscmp_tisstk": "1,200,000", "bdiscmp_mbsn": "제약",
      "trfdtl_bd_fta": "15,000,000,000", "trfdtl_trfprc": "16,000,000,000", "trfdtl_tast": "3,000,000,000,000", "trfdtl_tast_vs": "0.5", "trfdtl_ecpt": "2,000,000,000,000", "trfdtl_ecpt_vs": "0.8",
      "trf_pp": "투자 회수", "trf_prd": "2023년 12월 31일",
      "dlptn_cmpnm": "매수자", "dlptn_cpt": "-", "dlptn_mbsn": "-", "dlptn_hoadd": "서울시 종로구", "dlptn_rl_cmpn": "특수관계 없음",
      "dl_pym": "현금 수령",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "안진회계법인", "exevl_pd": "2023년 10월", "exevl_op": "적정",
      "bddd": "2023년 11월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "ftc_stt_atn": "미해당", "popt_ctr_atn": "아니오", "popt_ctr_cn": "-"
    }
  ]
}
```

- [ ] **Step 2: test 함수** (append)

```go
func TestStockRelatedBondTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stkrtbdTrfDecsn.json": "stkrtbdTrfDecsn.json"})
	items, err := c.StockRelatedBondTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231110000888", got.RceptNo)
	assert.Equal(t, "신주인수권부사채", got.StkrtbdKndn)
	assert.Equal(t, "2022년 05월 01일", got.Aqd)
	assert.Equal(t, "16,000,000,000", got.TrfdtlTrfprc)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
```

- [ ] **Step 3~8**: 공통 절차. struct 는 spec `StockRelatedBondTransferItem` 복사.

---

## Task 9: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 는 같은 패키지라 `ErrNoData` 직접 참조, `errors`/`material` import 기존 존재). 양수도는 흔치 않아 ErrNoData skip. 파일 끝에 2개 추가:

```go
func TestIntegration_BusinessAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.BusinessAcquisition(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 영업양수 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_OtherCorpStockAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OtherCorpStockAcquisition(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 타법인주식 양수 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_BusinessAcquisition|TestIntegration_OtherCorpStockAcquisition" ./...` → PASS 또는 SKIP.

- [ ] **Step 4: README 커버리지 갱신**

현재(자기주식 PR 머지 후) 두 줄:
```
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정) · 자기주식(취득/처분/신탁계약 체결·해지 결정)
- (예정) DS005 나머지(양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
다음으로 교체:
```
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정) · 자기주식(취득/처분/신탁계약 체결·해지 결정) · 양수도(영업/유형자산/타법인주식/주권사채권 양수·양도 결정)
- (예정) DS005 나머지(자산양수도·풋백옵션/주식교환·이전/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
(그 두 줄만 변경.)

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(material): add DS005 양수도(Sub-1) 통합 테스트 + README 커버리지"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 8개 메서드 → Task 1~8 각각 매핑(메서드 시그니처·엔드포인트 verbatim, struct 는 spec 동일 이름 참조). 통합 테스트·README = Task 9. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 메서드·fixture·test 는 완전한 코드. struct body 는 의도적으로 committed spec(EXACT Go 코드, 단일 출처)을 가리킴 — "show the code" 는 spec 파일이 충족(레포 내 존재). 각 Task 에 필드 수 명시로 검증 가능.

**3. Type consistency:** 8개 struct·메서드명이 spec 표와 1:1. 엔드포인트 차이 반영: 영업양수만 dbt_*, 영업양도엔 없음 / 유형자산은 aprskh_exrq·aprskh_ex_pc_mth_pd_pl 추가 / 타법인주식·주권사채권 양수는 l6m_tpa_nstkaq_atn, 양도는 미포함(주권사채권 양도는 aqd 포함) / 양수 계열만 bdlst 관련(영업양수·타법인주식양수). 모든 struct 거버넌스 od_a_at_*/adt_a_atn 동일. route map 값은 bare 파일명. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조.
