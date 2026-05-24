# OpenDART DS005 합병·분할 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 합병·분할 3개 API(회사합병/회사분할/회사분할합병 결정)를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 재사용. 신규 파일 `material/merger.go` 에 3개 item struct + 3개 한 줄 메서드. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-merger-design.md` 에 3개 struct 전체 필드가 EXACT Go 코드(json 태그 + 한글 코멘트)로 정의돼 있다. 각 Task 의 struct 는 그 파일의 동일 이름 struct 를 **그대로 복사**한다(필드 가감 금지).

---

## File Structure

- Create: `material/merger.go` — 3 item struct(`CompanyMergerItem` 69, `CompanyDivisionItem` 49, `CompanyDivisionMergerItem` 90) + 3 메서드.
- Create: `material/merger_test.go` — fixture 디코딩 테스트(기존 `newTestClient` 재사용).
- Create: `material/testdata/{cmpMgDecsn,cmpDvDecsn,cmpDvmgDecsn}.json` — 3 fixture.
- Modify: `integration_test.go` — 통합 케이스 2개(`//go:build integration`, ErrNoData skip).
- Modify: `README.md` — DS005 커버리지에 합병·분할 추가.

기존 컨벤션(변경 금지):
- `material/merger.go` 메서드 형태: `func (c *Client) CompanyMerger(ctx context.Context, p MaterialParams) ([]CompanyMergerItem, error) { return httpclient.GetList[CompanyMergerItem](ctx, c.http, "/api/cmpMgDecsn.json", p.toMap()) }`. Client http 필드는 `c.http`.
- `newTestClient(t, routes map[string]string)`: route map 값은 **bare 파일명** — 내부에서 `filepath.Join("testdata", fixture)` 로 testdata/ 를 붙인다.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 회사합병 결정 (CompanyMerger) — cmpMgDecsn, 69필드

**Files:** Create `material/merger.go`, `material/merger_test.go`, `material/testdata/cmpMgDecsn.json`.

Struct: spec 의 `CompanyMergerItem` (69필드) 그대로. Method:
```go
// CompanyMerger 는 회사합병 결정(주요사항보고서)을 조회한다.
func (c *Client) CompanyMerger(ctx context.Context, p MaterialParams) ([]CompanyMergerItem, error) {
	return httpclient.GetList[CompanyMergerItem](ctx, c.http, "/api/cmpMgDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/cmpMgDecsn.json`**

머리 4 + 69필드 전체를 담되 대표 필드만 의미 있는 값, 나머지는 "-" 로. 최소 다음 키는 의미값:
```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트회사합병",
      "mg_mth": "흡수합병", "mg_stn": "소규모합병", "mg_pp": "경영효율화", "mg_rt": "1:0.5", "mg_rt_bs": "본질가치 평가",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "삼일회계법인", "exevl_pd": "2023년 03월", "exevl_op": "적정",
      "mgnstk_ostk_cnt": "1,000,000", "mgnstk_cstk_cnt": "-",
      "mgptncmp_cmpnm": "합병상대회사", "mgptncmp_mbsn": "제조", "mgptncmp_rl_cmpn": "특수관계 없음",
      "rbsnfdtl_tast": "500,000,000,000", "rbsnfdtl_tdbt": "200,000,000,000", "rbsnfdtl_teqt": "300,000,000,000", "rbsnfdtl_cpt": "50,000,000,000", "rbsnfdtl_sl": "400,000,000,000", "rbsnfdtl_nic": "30,000,000,000",
      "eadtat_intn": "한영회계법인", "eadtat_op": "적정",
      "nmgcmp_cmpnm": "-", "ffdtl_tast": "-", "ffdtl_tdbt": "-", "ffdtl_teqt": "-", "ffdtl_cpt": "-", "ffdtl_std": "-", "nmgcmp_nbsn_rsl": "-", "nmgcmp_mbsn": "-", "nmgcmp_rlst_atn": "-",
      "mgsc_mgctrd": "2023년 04월 10일", "mgsc_shddstd": "2023년 05월 10일", "mgsc_shclspd_bgd": "2023년 05월 11일", "mgsc_shclspd_edd": "2023년 05월 15일", "mgsc_mgop_rcpd_bgd": "2023년 05월 16일", "mgsc_mgop_rcpd_edd": "2023년 06월 05일", "mgsc_gmtsck_prd": "2023년 06월 10일", "mgsc_aprskh_expd_bgd": "2023년 06월 10일", "mgsc_aprskh_expd_edd": "2023년 06월 30일", "mgsc_osprpd_bgd": "2023년 07월 01일", "mgsc_osprpd_edd": "2023년 07월 14일", "mgsc_trspprpd_bgd": "2023년 07월 13일", "mgsc_trspprpd_edd": "2023년 07월 28일", "mgsc_cdobprpd_bgd": "2023년 06월 11일", "mgsc_cdobprpd_edd": "2023년 07월 11일", "mgsc_mgdt": "2023년 07월 15일", "mgsc_ergmd": "2023년 07월 20일", "mgsc_mgrgsprd": "2023년 07월 25일", "mgsc_nstkdlprd": "2023년 07월 28일", "mgsc_nstklstprd": "2023년 07월 29일",
      "bdlst_atn": "아니오", "otcpr_bdlst_sf_atn": "해당없음",
      "aprskh_plnprc": "65,000", "aprskh_pym_plpd_mth": "2023년 07월, 현금", "aprskh_ctref": "-",
      "bddd": "2023년 04월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "popt_ctr_atn": "아니오", "popt_ctr_cn": "-", "rs_sm_atn": "제출", "ex_sm_r": "-"
    }
  ]
}
```
(모든 69개 json 키가 빠짐없이 포함되어야 한다 — 누락 키 없이 작성. struct 와 1:1.)

- [ ] **Step 2: create `material/merger_test.go`** (`package material` + import `context`,`testing`,`github.com/stretchr/testify/assert`,`github.com/stretchr/testify/require`)

```go
func TestCompanyMerger(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cmpMgDecsn.json": "cmpMgDecsn.json"})
	items, err := c.CompanyMerger(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230410000111", got.RceptNo)
	assert.Equal(t, "흡수합병", got.MgMth)
	assert.Equal(t, "1:0.5", got.MgRt)
	assert.Equal(t, "합병상대회사", got.MgptncmpCmpnm)
	assert.Equal(t, "2023년 06월 10일", got.MgscGmtsckPrd)
	assert.Equal(t, "제출", got.RsSmAtn)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestCompanyMerger`** — confirm FAIL (undefined).
- [ ] **Step 4: create `material/merger.go`** (`package material` + `import ("context"; "github.com/kenshin579/opendart/internal/httpclient")`), then paste `CompanyMergerItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: `gofmt -l material/merger.go material/merger_test.go`** — no output (else `gofmt -w`).
- [ ] **Step 7: verify** `grep -c 'json:"' material/merger.go` should be 69.
- [ ] **Step 8: real API BEST EFFORT** — `$OPENDART_API_KEY` set → `curl -s "https://opendart.fss.or.kr/api/cmpMgDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. 서버 TLS1.2 RSA cipher 만(curl 실패 가능). clean `"status":"000"` + non-empty list 면 fixture 교체 + assert 갱신(실 날짜 bracket) 후 step 5 재실행. 1회 시도, 실패 시 샘플 유지.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/merger.go material/merger_test.go material/testdata/cmpMgDecsn.json
git commit -m "feat(material): add DS005 CompanyMerger (회사합병 결정)"
```

---

## Task 2: 회사분할 결정 (CompanyDivision) — cmpDvDecsn, 49필드

**Files:** Modify `material/merger.go`, `material/merger_test.go`; Create `material/testdata/cmpDvDecsn.json`.

Struct: spec 의 `CompanyDivisionItem` (49필드) 그대로. Method:
```go
// CompanyDivision 은 회사분할 결정(주요사항보고서)을 조회한다.
func (c *Client) CompanyDivision(ctx context.Context, p MaterialParams) ([]CompanyDivisionItem, error) {
	return httpclient.GetList[CompanyDivisionItem](ctx, c.http, "/api/cmpDvDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/cmpDvDecsn.json`** — 49개 json 키 전체 포함(누락 없이). 대표값:

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230510000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트회사분할",
      "dv_mth": "인적분할", "dv_impef": "사업 전문화", "dv_rt": "0.7:0.3", "dv_trfbsnprt_cn": "배터리 사업부문",
      "atdv_excmp_cmpnm": "존속회사", "atdvfdtl_tast": "2,000,000,000,000", "atdvfdtl_tdbt": "800,000,000,000", "atdvfdtl_teqt": "1,200,000,000,000", "atdvfdtl_cpt": "100,000,000,000", "atdvfdtl_std": "2022년 12월 31일", "atdv_excmp_exbsn_rsl": "1,500,000,000,000", "atdv_excmp_mbsn": "전자", "atdv_excmp_atdv_lstmn_atn": "예",
      "dvfcmp_cmpnm": "신설회사", "ffdtl_tast": "1,000,000,000,000", "ffdtl_tdbt": "300,000,000,000", "ffdtl_teqt": "700,000,000,000", "ffdtl_cpt": "50,000,000,000", "ffdtl_std": "2023년 07월 01일", "dvfcmp_nbsn_rsl": "500,000,000,000", "dvfcmp_mbsn": "배터리", "dvfcmp_rlst_atn": "예",
      "abcr_crrt": "30.0", "abcr_osprpd_bgd": "2023년 06월 01일", "abcr_osprpd_edd": "2023년 06월 14일", "abcr_trspprpd_bgd": "2023년 06월 28일", "abcr_trspprpd_edd": "2023년 07월 03일", "abcr_nstkascnd": "분할비율에 따라 배정", "abcr_shstkcnt_rt_at_rs": "비례 배정", "abcr_nstkasstd": "2023년 06월 30일", "abcr_nstkdlprd": "2023년 07월 14일", "abcr_nstklstprd": "2023년 07월 17일",
      "gmtsck_prd": "2023년 05월 30일", "cdobprpd_bgd": "2023년 05월 31일", "cdobprpd_edd": "2023년 06월 30일", "dvdt": "2023년 07월 01일", "dvrgsprd": "2023년 07월 03일",
      "bddd": "2023년 05월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "popt_ctr_atn": "아니오", "popt_ctr_cn": "-", "rs_sm_atn": "제출", "ex_sm_r": "-"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/merger_test.go`**

```go
func TestCompanyDivision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cmpDvDecsn.json": "cmpDvDecsn.json"})
	items, err := c.CompanyDivision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230510000222", got.RceptNo)
	assert.Equal(t, "인적분할", got.DvMth)
	assert.Equal(t, "배터리 사업부문", got.DvTrfbsnprtCn)
	assert.Equal(t, "30.0", got.AbcrCrrt)
	assert.Equal(t, "2023년 07월 01일", got.Dvdt)
	assert.Equal(t, "제출", got.RsSmAtn)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestCompanyDivision`** — FAIL.
- [ ] **Step 4: append to `material/merger.go`**: `CompanyDivisionItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l material/merger.go material/merger_test.go` 출력 없으면 OK.
- [ ] **Step 7: verify** `CompanyDivisionItem` has 49 json tags (file total now 69+49=118).
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/cmpDvDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/merger.go material/merger_test.go material/testdata/cmpDvDecsn.json
git commit -m "feat(material): add DS005 CompanyDivision (회사분할 결정)"
```

---

## Task 3: 회사분할합병 결정 (CompanyDivisionMerger) — cmpDvmgDecsn, 90필드

**Files:** Modify `material/merger.go`, `material/merger_test.go`; Create `material/testdata/cmpDvmgDecsn.json`.

Struct: spec 의 `CompanyDivisionMergerItem` (90필드) 그대로. **주의**: `dvfcmp_atdv_lstmn_at`(분할설립, `_at`로 끝)와 `atdv_excmp_atdv_lstmn_atn`(존속, `_atn`)은 다른 키이므로 spec 그대로 복사할 것. Method:
```go
// CompanyDivisionMerger 는 회사분할합병 결정(주요사항보고서)을 조회한다.
func (c *Client) CompanyDivisionMerger(ctx context.Context, p MaterialParams) ([]CompanyDivisionMergerItem, error) {
	return httpclient.GetList[CompanyDivisionMergerItem](ctx, c.http, "/api/cmpDvmgDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/cmpDvmgDecsn.json`** — 90개 json 키 전체 포함(누락 없이; 대표 키만 의미값, 나머지 "-"). 최소 다음 키는 의미값:

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230610000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트분할합병",
      "dvmg_mth": "분할합병", "dvmg_impef": "사업 재편",
      "dv_trfbsnprt_cn": "물류 사업부문",
      "atdv_excmp_cmpnm": "존속회사", "atdvfdtl_tast": "2,000,000,000,000", "atdvfdtl_tdbt": "-", "atdvfdtl_teqt": "-", "atdvfdtl_cpt": "-", "atdvfdtl_std": "-", "atdv_excmp_exbsn_rsl": "-", "atdv_excmp_mbsn": "전자", "atdv_excmp_atdv_lstmn_atn": "예",
      "dvfcmp_cmpnm": "분할설립회사", "ffdtl_tast": "-", "ffdtl_tdbt": "-", "ffdtl_teqt": "-", "ffdtl_cpt": "-", "ffdtl_std": "-", "dvfcmp_nbsn_rsl": "-", "dvfcmp_mbsn": "물류", "dvfcmp_atdv_lstmn_at": "아니오",
      "abcr_crrt": "20.0", "abcr_osprpd_bgd": "-", "abcr_osprpd_edd": "-", "abcr_trspprpd_bgd": "-", "abcr_trspprpd_edd": "-", "abcr_nstkascnd": "-", "abcr_shstkcnt_rt_at_rs": "-", "abcr_nstkasstd": "-", "abcr_nstkdlprd": "-", "abcr_nstklstprd": "-",
      "mg_stn": "흡수합병", "mgptncmp_cmpnm": "합병상대회사", "mgptncmp_mbsn": "물류", "mgptncmp_rl_cmpn": "특수관계 없음",
      "rbsnfdtl_tast": "800,000,000,000", "rbsnfdtl_tdbt": "-", "rbsnfdtl_teqt": "-", "rbsnfdtl_cpt": "-", "rbsnfdtl_sl": "-", "rbsnfdtl_nic": "-", "eadtat_intn": "삼정회계법인", "eadtat_op": "적정",
      "dvmgnstk_ostk_cnt": "2,000,000", "dvmgnstk_cstk_cnt": "-",
      "nmgcmp_cmpnm": "-", "nmgcmp_cpt": "-", "nmgcmp_mbsn": "-", "nmgcmp_rlst_atn": "-",
      "dvmg_rt": "1:0.8", "dvmg_rt_bs": "본질가치 평가",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법", "exevl_intn": "삼일회계법인", "exevl_pd": "2023년 05월", "exevl_op": "적정",
      "dvmgsc_dvmgctrd": "2023년 06월 10일", "dvmgsc_shddstd": "2023년 06월 30일", "dvmgsc_shclspd_bgd": "-", "dvmgsc_shclspd_edd": "-", "dvmgsc_dvmgop_rcpd_bgd": "-", "dvmgsc_dvmgop_rcpd_edd": "-", "dvmgsc_gmtsck_prd": "2023년 07월 20일", "dvmgsc_aprskh_expd_bgd": "-", "dvmgsc_aprskh_expd_edd": "-", "dvmgsc_cdobprpd_bgd": "-", "dvmgsc_cdobprpd_edd": "-", "dvmgsc_dvmgdt": "2023년 08월 01일", "dvmgsc_ergmd": "-", "dvmgsc_dvmgrgsprd": "2023년 08월 05일",
      "bdlst_atn": "아니오", "otcpr_bdlst_sf_atn": "해당없음",
      "aprskh_exrq": "-", "aprskh_plnprc": "70,000", "aprskh_ex_pc_mth_pd_pl": "-", "aprskh_pym_plpd_mth": "2023년 08월, 현금", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bddd": "2023년 06월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "popt_ctr_atn": "아니오", "popt_ctr_cn": "-", "rs_sm_atn": "제출", "ex_sm_r": "-"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/merger_test.go`**

```go
func TestCompanyDivisionMerger(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cmpDvmgDecsn.json": "cmpDvmgDecsn.json"})
	items, err := c.CompanyDivisionMerger(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230610000333", got.RceptNo)
	assert.Equal(t, "분할합병", got.DvmgMth)
	assert.Equal(t, "물류 사업부문", got.DvTrfbsnprtCn)
	assert.Equal(t, "1:0.8", got.DvmgRt)
	assert.Equal(t, "아니오", got.DvfcmpAtdvLstmnAt)
	assert.Equal(t, "2023년 08월 01일", got.DvmgscDvmgdt)
	assert.Equal(t, "제출", got.RsSmAtn)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestCompanyDivisionMerger`** — FAIL.
- [ ] **Step 4: append to `material/merger.go`**: `CompanyDivisionMergerItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l material/merger.go material/merger_test.go` 출력 없으면 OK.
- [ ] **Step 7: verify** `CompanyDivisionMergerItem` has 90 json tags (file total now 69+49+90=208). Confirm `dvfcmp_atdv_lstmn_at` present and `atdv_excmp_atdv_lstmn_atn` present (both, distinct).
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/cmpDvmgDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/merger.go material/merger_test.go material/testdata/cmpDvmgDecsn.json
git commit -m "feat(material): add DS005 CompanyDivisionMerger (회사분할합병 결정)"
```

---

## Task 4: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 는 같은 패키지라 직접 참조, `errors`/`material` import 기존 존재). 파일 끝에 2개 추가:

```go
func TestIntegration_CompanyMerger(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.CompanyMerger(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 회사합병 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_CompanyDivision(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.CompanyDivision(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 회사분할 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_CompanyMerger|TestIntegration_CompanyDivision" ./...` → PASS 또는 SKIP.

- [ ] **Step 4: README 커버리지 갱신**

현재(양수도 Sub-2 PR 머지 후) DS005 줄 끝은 `... · 자산양수도(기타)·풋백옵션 · 주식교환·이전 결정` 이고, 예정 줄은 `(예정) DS005 나머지(합병·분할/해외상장) · DS006 · ...`. DS005 줄 끝에 ` · 회사합병·분할·분할합병 결정` 을 추가하고, 예정 줄을 `(예정) DS005 나머지(해외상장) · DS006 · ...` 로 변경(합병·분할/ 제거). 그 두 줄만 변경.

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(material): add DS005 합병·분할 통합 테스트 + README 커버리지"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 3개 메서드(CompanyMerger/CompanyDivision/CompanyDivisionMerger) → Task 1~3 매핑(메서드 시그니처·엔드포인트 verbatim, struct 는 spec 동일 이름 참조). 통합 테스트·README = Task 4. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 메서드·fixture·test 완전(fixture 는 전 json 키 포함 명시). struct body 는 committed spec(EXACT Go 코드, 단일 출처) 참조 — 레포 내 존재, 필드 수 명시(69/49/90)로 검증.

**3. Type consistency:** 메서드명·struct명·필드명이 spec 표와 1:1. route map 값은 bare 파일명. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조. 분할합병의 `dvfcmp_atdv_lstmn_at`(_at) vs 존속 `atdv_excmp_atdv_lstmn_atn`(_atn) 구분을 Task 3 Step 7 에서 명시 검증. fixture json 키는 struct 태그와 1:1.
