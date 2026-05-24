# OpenDART DS005 양수도 Sub-2 «기타» Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 양수도 그룹 중 기타 2개 API(자산양수도(기타)·풋백옵션 + 주식교환·이전 결정)를 `material` 패키지에 추가한다.

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 를 그대로 재사용한다. 신규 파일 `material/transfer_etc.go` 에 2개 item struct + 2개 한 줄 메서드를 추가한다. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-transfer-etc-design.md` 에 2개 struct 전체 필드가 EXACT Go 코드(json 태그 + 한글 코멘트)로 정의돼 있다. 각 Task 의 struct 는 그 파일의 동일 이름 struct 를 **그대로 복사**한다(필드 가감 금지).

---

## File Structure

- Create: `material/transfer_etc.go` — 2 item struct + 2 메서드.
- Create: `material/transfer_etc_test.go` — fixture 디코딩 테스트(기존 `newTestClient` 재사용).
- Create: `material/testdata/{astInhtrfEtcPtbkOpt,stkExtrDecsn}.json` — 2 fixture.
- Modify: `integration_test.go` — 통합 케이스 2개(`//go:build integration`, ErrNoData skip).
- Modify: `README.md` — DS005 커버리지에 양수도 Sub-2 추가.

기존 컨벤션(변경 금지):
- `material/transfer.go` 메서드 형태: `func (c *Client) BusinessAcquisition(ctx context.Context, p MaterialParams) ([]BusinessAcquisitionItem, error) { return httpclient.GetList[BusinessAcquisitionItem](ctx, c.http, "/api/bsnInhDecsn.json", p.toMap()) }`. Client http 필드는 `c.http`.
- `newTestClient(t, routes map[string]string)`: route map 값은 **bare 파일명** — 내부에서 `filepath.Join("testdata", fixture)` 로 testdata/ 를 붙인다.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 자산양수도(기타), 풋백옵션 (OtherAssetTransferPutbackOption) — astInhtrfEtcPtbkOpt, 6필드

**Files:** Create `material/transfer_etc.go`, `material/transfer_etc_test.go`, `material/testdata/astInhtrfEtcPtbkOpt.json`.

Struct: spec 의 `OtherAssetTransferPutbackOptionItem` (6필드) 그대로. Method:
```go
// OtherAssetTransferPutbackOption 은 자산양수도(기타), 풋백옵션(주요사항보고서)을 조회한다.
func (c *Client) OtherAssetTransferPutbackOption(ctx context.Context, p MaterialParams) ([]OtherAssetTransferPutbackOptionItem, error) {
	return httpclient.GetList[OtherAssetTransferPutbackOptionItem](ctx, c.http, "/api/astInhtrfEtcPtbkOpt.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/astInhtrfEtcPtbkOpt.json`**

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230315000111",
      "corp_cls": "Y",
      "corp_code": "00126380",
      "corp_name": "테스트자산양수도",
      "rp_rsn": "자산양수도(기타)에 해당",
      "ast_inhtrf_prc": "120,000,000,000"
    }
  ]
}
```

- [ ] **Step 2: create `material/transfer_etc_test.go`** (`package material` + import `context`,`testing`,`github.com/stretchr/testify/assert`,`github.com/stretchr/testify/require`)

```go
func TestOtherAssetTransferPutbackOption(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/astInhtrfEtcPtbkOpt.json": "astInhtrfEtcPtbkOpt.json"})
	items, err := c.OtherAssetTransferPutbackOption(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230315000111", got.RceptNo)
	assert.Equal(t, "자산양수도(기타)에 해당", got.RpRsn)
	assert.Equal(t, "120,000,000,000", got.AstInhtrfPrc)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestOtherAssetTransferPutbackOption`** — confirm FAIL (undefined).

- [ ] **Step 4: create `material/transfer_etc.go`** (`package material` + `import ("context"; "github.com/kenshin579/opendart/internal/httpclient")`), then paste `OtherAssetTransferPutbackOptionItem` struct VERBATIM from spec, then the method above.

- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: `gofmt -l material/transfer_etc.go material/transfer_etc_test.go`** — no output (else `gofmt -w`).
- [ ] **Step 7: verify** `OtherAssetTransferPutbackOptionItem` has 6 json tags.
- [ ] **Step 8: real API BEST EFFORT** — `$OPENDART_API_KEY` set → `curl -s "https://opendart.fss.or.kr/api/astInhtrfEtcPtbkOpt.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. 서버 TLS1.2 RSA cipher 만 지원(curl handshake 실패 가능). clean `"status":"000"` + non-empty list 면 fixture 교체 + assert 갱신(실 날짜 bracket) 후 step 5 재실행. 1회 시도, 실패 시 샘플 유지.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/transfer_etc.go material/transfer_etc_test.go material/testdata/astInhtrfEtcPtbkOpt.json
git commit -m "feat(material): add DS005 OtherAssetTransferPutbackOption (자산양수도(기타), 풋백옵션)"
```

---

## Task 2: 주식교환·이전 결정 (StockExchangeTransfer) — stkExtrDecsn, 56필드

**Files:** Modify `material/transfer_etc.go`, `material/transfer_etc_test.go`; Create `material/testdata/stkExtrDecsn.json`.

Struct: spec 의 `StockExchangeTransferItem` (56필드) 그대로. Method:
```go
// StockExchangeTransfer 는 주식교환·이전 결정(주요사항보고서)을 조회한다.
func (c *Client) StockExchangeTransfer(ctx context.Context, p MaterialParams) ([]StockExchangeTransferItem, error) {
	return httpclient.GetList[StockExchangeTransferItem](ctx, c.http, "/api/stkExtrDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/stkExtrDecsn.json`**

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230420000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트주식교환이전",
      "extr_sen": "주식의 포괄적 교환", "extr_stn": "완전자회사",
      "extr_tgcmp_cmpnm": "대상법인", "extr_tgcmp_rp": "홍길동", "extr_tgcmp_mbsn": "소프트웨어", "extr_tgcmp_rl_cmpn": "특수관계 없음",
      "extr_tgcmp_tisstk_ostk": "1,000,000", "extr_tgcmp_tisstk_cstk": "-",
      "rbsnfdtl_tast": "500,000,000,000", "rbsnfdtl_tdbt": "200,000,000,000", "rbsnfdtl_teqt": "300,000,000,000", "rbsnfdtl_cpt": "50,000,000,000",
      "extr_rt": "1:0.5", "extr_rt_bs": "본질가치 평가",
      "exevl_atn": "예", "exevl_bs_rs": "자본시장법 제165조의4", "exevl_intn": "삼일회계법인", "exevl_pd": "2023년 03월", "exevl_op": "적정",
      "extr_pp": "경영 효율화",
      "extrsc_extrctrd": "2023년 04월 20일", "extrsc_shddstd": "2023년 05월 10일",
      "extrsc_shclspd_bgd": "2023년 05월 11일", "extrsc_shclspd_edd": "2023년 05월 15일",
      "extrsc_extrop_rcpd_bgd": "2023년 05월 16일", "extrsc_extrop_rcpd_edd": "2023년 06월 05일",
      "extrsc_gmtsck_prd": "2023년 06월 10일",
      "extrsc_aprskh_expd_bgd": "2023년 06월 10일", "extrsc_aprskh_expd_edd": "2023년 06월 30일",
      "extrsc_osprpd_bgd": "2023년 07월 01일", "extrsc_osprpd_edd": "2023년 07월 14일",
      "extrsc_trspprpd": "2023년 07월 13일 ~ 2023년 07월 28일", "extrsc_trspprpd_bgd": "2023년 07월 13일", "extrsc_trspprpd_edd": "2023년 07월 28일",
      "extrsc_extrdt": "2023년 07월 15일", "extrsc_nstkdlprd": "2023년 07월 28일", "extrsc_nstklstprd": "2023년 07월 29일",
      "atextr_cpcmpnm": "테스트주식교환이전",
      "aprskh_plnprc": "65,000", "aprskh_pym_plpd_mth": "2023년 07월, 현금", "aprskh_lmt": "-", "aprskh_ctref": "-",
      "bdlst_atn": "아니오", "otcpr_bdlst_sf_atn": "해당없음",
      "bddd": "2023년 04월 20일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1",
      "popt_ctr_atn": "아니오", "popt_ctr_cn": "-",
      "rs_sm_atn": "제출", "ex_sm_r": "-"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/transfer_etc_test.go`**

```go
func TestStockExchangeTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stkExtrDecsn.json": "stkExtrDecsn.json"})
	items, err := c.StockExchangeTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230420000222", got.RceptNo)
	assert.Equal(t, "주식의 포괄적 교환", got.ExtrSen)
	assert.Equal(t, "대상법인", got.ExtrTgcmpCmpnm)
	assert.Equal(t, "1:0.5", got.ExtrRt)
	assert.Equal(t, "2023년 06월 10일", got.ExtrscGmtsckPrd)
	assert.Equal(t, "제출", got.RsSmAtn)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestStockExchangeTransfer`** — confirm FAIL (undefined).
- [ ] **Step 4: append to `material/transfer_etc.go`**: `StockExchangeTransferItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: `gofmt -l material/transfer_etc.go material/transfer_etc_test.go`** — no output (else `gofmt -w`).
- [ ] **Step 7: verify** `StockExchangeTransferItem` has 56 json tags.
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/stkExtrDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회 시도.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/transfer_etc.go material/transfer_etc_test.go material/testdata/stkExtrDecsn.json
git commit -m "feat(material): add DS005 StockExchangeTransfer (주식교환·이전 결정)"
```

---

## Task 3: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 는 같은 패키지라 직접 참조, `errors`/`material` import 기존 존재). 파일 끝에 2개 추가:

```go
func TestIntegration_OtherAssetTransferPutbackOption(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OtherAssetTransferPutbackOption(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 자산양수도(기타) 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_StockExchangeTransfer(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.StockExchangeTransfer(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 주식교환·이전 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_OtherAssetTransferPutbackOption|TestIntegration_StockExchangeTransfer" ./...` → PASS 또는 SKIP.

- [ ] **Step 4: README 커버리지 갱신**

현재(양수도 Sub-1 PR 머지 후) 두 줄:
```
- DS005 주요사항보고서 주요정보: ... · 자기주식(취득/처분/신탁계약 체결·해지 결정) · 양수도(영업/유형자산/타법인주식/주권사채권 양수·양도 결정)
- (예정) DS005 나머지(자산양수도·풋백옵션/주식교환·이전/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
다음으로 교체(DS005 줄 끝에 추가, 예정 줄에서 자산양수도·풋백옵션/주식교환·이전 제거):
```
- DS005 주요사항보고서 주요정보: ... · 자기주식(취득/처분/신탁계약 체결·해지 결정) · 양수도(영업/유형자산/타법인주식/주권사채권 양수·양도 결정) · 자산양수도(기타)·풋백옵션 · 주식교환·이전 결정
- (예정) DS005 나머지(합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0
```
(그 두 줄만 변경. DS005 줄의 앞부분 "..." 은 실제 파일 내용 그대로 두고 끝에만 추가.)

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(material): add DS005 양수도 Sub-2 통합 테스트 + README 커버리지"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 2개 메서드(OtherAssetTransferPutbackOption/StockExchangeTransfer) → Task 1~2 매핑(메서드 시그니처·엔드포인트 verbatim, struct 는 spec 동일 이름 참조). 통합 테스트·README = Task 3. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 메서드·fixture·test 완전. struct body 는 committed spec(EXACT Go 코드, 단일 출처) 참조 — 레포 내 존재, 필드 수 명시(6/56)로 검증.

**3. Type consistency:** 메서드명·struct명·필드명이 spec 표와 1:1. route map 값은 bare 파일명. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조. fixture json 키는 struct 태그의 부분집합.
