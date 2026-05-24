# OpenDART DS005 해외상장 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS005 해외상장 4개 API(해외 증권시장 주권등 상장 결정·상장·상장폐지 결정·상장폐지)를 `material` 패키지에 추가한다. (DS005 마지막 그룹 — 완료 시 36/36.)

**Architecture:** 기존 `material.Client` + 공통 `MaterialParams{CorpCode,BgnDe,EndDe}` + `httpclient.GetList[T]` 재사용. 신규 파일 `material/overseas_listing.go` 에 4개 item struct + 4개 한 줄 메서드. root `opendart` 패키지 변경 없음(`client.Material` 기존 와이어링 유지).

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-24-opendart-ds005-material-overseas-listing-design.md` 에 4개 struct 전체 필드가 EXACT Go 코드(json 태그 + 한글 코멘트)로 정의돼 있다. 각 Task 의 struct 는 그 파일의 동일 이름 struct 를 **그대로 복사**한다(필드 가감 금지).

---

## File Structure

- Create: `material/overseas_listing.go` — 4 item struct(`OverseasListingDecisionItem` 20, `OverseasListingItem` 10, `OverseasDelistingDecisionItem` 14, `OverseasDelistingItem` 10) + 4 메서드.
- Create: `material/overseas_listing_test.go` — fixture 디코딩 테스트(기존 `newTestClient` 재사용).
- Create: `material/testdata/{ovLstDecsn,ovLst,ovDlstDecsn,ovDlst}.json` — 4 fixture.
- Modify: `integration_test.go` — 통합 케이스 2개(`//go:build integration`, ErrNoData skip).
- Modify: `README.md` — DS005 커버리지에 해외상장 추가, 예정 줄에서 DS005 제거(36/36 완료).

기존 컨벤션(변경 금지):
- `material/merger.go` 메서드 형태: `func (c *Client) CompanyMerger(ctx context.Context, p MaterialParams) ([]CompanyMergerItem, error) { return httpclient.GetList[CompanyMergerItem](ctx, c.http, "/api/cmpMgDecsn.json", p.toMap()) }`. Client http 필드는 `c.http`.
- `newTestClient(t, routes map[string]string)`: route map 값은 **bare 파일명** — 내부에서 `filepath.Join("testdata", fixture)` 로 testdata/ 를 붙인다.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 해외 증권시장 주권등 상장 결정 (OverseasListingDecision) — ovLstDecsn, 20필드

**Files:** Create `material/overseas_listing.go`, `material/overseas_listing_test.go`, `material/testdata/ovLstDecsn.json`.

Struct: spec 의 `OverseasListingDecisionItem` (20필드) 그대로. Method:
```go
// OverseasListingDecision 은 해외 증권시장 주권등 상장 결정(주요사항보고서)을 조회한다.
func (c *Client) OverseasListingDecision(ctx context.Context, p MaterialParams) ([]OverseasListingDecisionItem, error) {
	return httpclient.GetList[OverseasListingDecisionItem](ctx, c.http, "/api/ovLstDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/ovLstDecsn.json`** (20개 json 키 전체 포함)

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20230410000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트해외상장결정",
      "lstprstk_ostk_cnt": "1,000,000", "lstprstk_estk_cnt": "-",
      "tisstk_ostk": "100,000,000", "tisstk_estk": "-",
      "psmth_nstk_sl": "1,000,000", "psmth_ostk_sl": "-",
      "fdpp": "해외 생산시설 투자",
      "lststk_orlst": "-", "lststk_drlst": "1,000,000",
      "lstex_nt": "미국 (NASDAQ)", "lstpp": "글로벌 자금조달", "lstprd": "2023년 09월 30일",
      "bddd": "2023년 04월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1"
    }
  ]
}
```

- [ ] **Step 2: create `material/overseas_listing_test.go`** (`package material` + import `context`,`testing`,`github.com/stretchr/testify/assert`,`github.com/stretchr/testify/require`)

```go
func TestOverseasListingDecision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLstDecsn.json": "ovLstDecsn.json"})
	items, err := c.OverseasListingDecision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230410000111", got.RceptNo)
	assert.Equal(t, "1,000,000", got.LstprstkOstkCnt)
	assert.Equal(t, "미국 (NASDAQ)", got.LstexNt)
	assert.Equal(t, "2023년 09월 30일", got.Lstprd)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestOverseasListingDecision`** — confirm FAIL (undefined).
- [ ] **Step 4: create `material/overseas_listing.go`** (`package material` + `import ("context"; "github.com/kenshin579/opendart/internal/httpclient")`), then paste `OverseasListingDecisionItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: `gofmt -l material/overseas_listing.go material/overseas_listing_test.go`** — no output (else `gofmt -w`).
- [ ] **Step 7: verify** `grep -c 'json:"' material/overseas_listing.go` should be 20.
- [ ] **Step 8: real API BEST EFFORT** — `$OPENDART_API_KEY` set → `curl -s "https://opendart.fss.or.kr/api/ovLstDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. 서버 TLS1.2 RSA cipher 만(curl 실패 가능). clean `"status":"000"` + non-empty list 면 fixture 교체 + assert 갱신(실 날짜 bracket) 후 step 5 재실행. 1회 시도, 실패 시 샘플 유지.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/overseas_listing.go material/overseas_listing_test.go material/testdata/ovLstDecsn.json
git commit -m "feat(material): add DS005 OverseasListingDecision (해외 증권시장 주권등 상장 결정)"
```

---

## Task 2: 해외 증권시장 주권등 상장 (OverseasListing) — ovLst, 10필드

**Files:** Modify `material/overseas_listing.go`, `material/overseas_listing_test.go`; Create `material/testdata/ovLst.json`.

Struct: spec 의 `OverseasListingItem` (10필드) 그대로. Method:
```go
// OverseasListing 은 해외 증권시장 주권등 상장(주요사항보고서)을 조회한다.
func (c *Client) OverseasListing(ctx context.Context, p MaterialParams) ([]OverseasListingItem, error) {
	return httpclient.GetList[OverseasListingItem](ctx, c.http, "/api/ovLst.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/ovLst.json`** (10개 json 키 전체 포함)

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20231010000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트해외상장",
      "lststk_ostk_cnt": "1,000,000", "lststk_estk_cnt": "-",
      "lstex_nt": "미국 (NASDAQ)", "stk_cd": "TEST", "lstd": "2023년 10월 02일", "cfd": "2023년 10월 03일"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/overseas_listing_test.go`**

```go
func TestOverseasListing(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLst.json": "ovLst.json"})
	items, err := c.OverseasListing(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231010000222", got.RceptNo)
	assert.Equal(t, "1,000,000", got.LststkOstkCnt)
	assert.Equal(t, "TEST", got.StkCd)
	assert.Equal(t, "2023년 10월 02일", got.Lstd)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestOverseasListing`** — FAIL.
- [ ] **Step 4: append to `material/overseas_listing.go`**: `OverseasListingItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l material/overseas_listing.go material/overseas_listing_test.go` 출력 없으면 OK.
- [ ] **Step 7: verify** `OverseasListingItem` has 10 json tags (file total 20+10=30).
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/ovLst.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/overseas_listing.go material/overseas_listing_test.go material/testdata/ovLst.json
git commit -m "feat(material): add DS005 OverseasListing (해외 증권시장 주권등 상장)"
```

---

## Task 3: 해외 증권시장 주권등 상장폐지 결정 (OverseasDelistingDecision) — ovDlstDecsn, 14필드

**Files:** Modify `material/overseas_listing.go`, `material/overseas_listing_test.go`; Create `material/testdata/ovDlstDecsn.json`.

Struct: spec 의 `OverseasDelistingDecisionItem` (14필드) 그대로. **주의**: bddd 코멘트는 "이사회결의일(확인일)". Method:
```go
// OverseasDelistingDecision 은 해외 증권시장 주권등 상장폐지 결정(주요사항보고서)을 조회한다.
func (c *Client) OverseasDelistingDecision(ctx context.Context, p MaterialParams) ([]OverseasDelistingDecisionItem, error) {
	return httpclient.GetList[OverseasDelistingDecisionItem](ctx, c.http, "/api/ovDlstDecsn.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/ovDlstDecsn.json`** (14개 json 키 전체 포함)

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20231110000333", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트상장폐지결정",
      "dlststk_ostk_cnt": "1,000,000", "dlststk_estk_cnt": "-",
      "lstex_nt": "미국 (NASDAQ)", "dlstrq_prd": "2023년 12월 01일", "dlst_prd": "2023년 12월 15일", "dlst_rs": "거래량 부족",
      "bddd": "2023년 11월 10일", "od_a_at_t": "3", "od_a_at_b": "0", "adt_a_atn": "1"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/overseas_listing_test.go`**

```go
func TestOverseasDelistingDecision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovDlstDecsn.json": "ovDlstDecsn.json"})
	items, err := c.OverseasDelistingDecision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231110000333", got.RceptNo)
	assert.Equal(t, "1,000,000", got.DlststkOstkCnt)
	assert.Equal(t, "거래량 부족", got.DlstRs)
	assert.Equal(t, "2023년 11월 10일", got.Bddd)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestOverseasDelistingDecision`** — FAIL.
- [ ] **Step 4: append to `material/overseas_listing.go`**: `OverseasDelistingDecisionItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l material/overseas_listing.go material/overseas_listing_test.go` 출력 없으면 OK.
- [ ] **Step 7: verify** `OverseasDelistingDecisionItem` has 14 json tags (file total 30+14=44).
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/ovDlstDecsn.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/overseas_listing.go material/overseas_listing_test.go material/testdata/ovDlstDecsn.json
git commit -m "feat(material): add DS005 OverseasDelistingDecision (해외 증권시장 주권등 상장폐지 결정)"
```

---

## Task 4: 해외 증권시장 주권등 상장폐지 (OverseasDelisting) — ovDlst, 10필드

**Files:** Modify `material/overseas_listing.go`, `material/overseas_listing_test.go`; Create `material/testdata/ovDlst.json`.

Struct: spec 의 `OverseasDelistingItem` (10필드) 그대로. Method:
```go
// OverseasDelisting 은 해외 증권시장 주권등 상장폐지(주요사항보고서)을 조회한다.
func (c *Client) OverseasDelisting(ctx context.Context, p MaterialParams) ([]OverseasDelistingItem, error) {
	return httpclient.GetList[OverseasDelistingItem](ctx, c.http, "/api/ovDlst.json", p.toMap())
}
```

- [ ] **Step 1: fixture `material/testdata/ovDlst.json`** (10개 json 키 전체 포함)

```json
{
  "status": "000",
  "message": "정상",
  "list": [
    {
      "rcept_no": "20231215000444", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트상장폐지",
      "lstex_nt": "미국 (NASDAQ)",
      "dlststk_ostk_cnt": "1,000,000", "dlststk_estk_cnt": "-",
      "tredd": "2023년 12월 14일", "dlst_rs": "거래량 부족", "cfd": "2023년 12월 15일"
    }
  ]
}
```

- [ ] **Step 2: append test to `material/overseas_listing_test.go`**

```go
func TestOverseasDelisting(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovDlst.json": "ovDlst.json"})
	items, err := c.OverseasDelisting(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231215000444", got.RceptNo)
	assert.Equal(t, "미국 (NASDAQ)", got.LstexNt)
	assert.Equal(t, "2023년 12월 14일", got.Tredd)
	assert.Equal(t, "거래량 부족", got.DlstRs)
}
```

- [ ] **Step 3: run `go test ./material/ -run TestOverseasDelisting`** — FAIL.
- [ ] **Step 4: append to `material/overseas_listing.go`**: `OverseasDelistingItem` struct VERBATIM from spec, then the method above.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l material/overseas_listing.go material/overseas_listing_test.go` 출력 없으면 OK.
- [ ] **Step 7: verify** `OverseasDelistingItem` has 10 json tags (file total 44+10=54).
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/ovDlst.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20200101&end_de=20241231"`. TLS1.2 RSA may fail. clean 000+non-empty → replace+re-assert, else keep sample. 1회.
- [ ] **Step 9: commit (no push)**

```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add material/overseas_listing.go material/overseas_listing_test.go material/testdata/ovDlst.json
git commit -m "feat(material): add DS005 OverseasDelisting (해외 증권시장 주권등 상장폐지)"
```

---

## Task 5: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 는 같은 패키지라 직접 참조, `errors`/`material` import 기존 존재). 파일 끝에 2개 추가:

```go
func TestIntegration_OverseasListingDecision(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OverseasListingDecision(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 해외상장 결정 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_OverseasDelisting(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OverseasDelisting(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 해외상장폐지 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_OverseasListingDecision|TestIntegration_OverseasDelisting" ./...` → PASS 또는 SKIP.

- [ ] **Step 4: README 커버리지 갱신**

현재 DS005 줄 끝은 `... · 회사합병·분할·분할합병 결정` 이고, 예정 줄은 `(예정) DS005 나머지(해외상장) · DS006 · DS002 개인별 보수 Ver2.0`. DS005 줄 끝에 ` · 해외 증권시장 주권등 상장 결정·상장·상장폐지 결정·상장폐지` 를 추가하고, 예정 줄을 `(예정) DS006 · DS002 개인별 보수 Ver2.0` 로 변경(DS005 나머지(해외상장) 제거 — DS005 36/36 완료). 그 두 줄만 변경.

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l material/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(material): add DS005 해외상장 통합 테스트 + README 커버리지 (DS005 완료)"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 4개 메서드(OverseasListingDecision/OverseasListing/OverseasDelistingDecision/OverseasDelisting) → Task 1~4 매핑(메서드 시그니처·엔드포인트 verbatim, struct 는 spec 동일 이름 참조). 통합 테스트·README = Task 5. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 메서드·fixture·test 완전. struct body 는 committed spec(EXACT Go 코드, 단일 출처) 참조 — 레포 내 존재, 필드 수 명시(20/10/14/10)로 검증.

**3. Type consistency:** 메서드명·struct명·필드명이 spec 표와 1:1. route map 값은 bare 파일명. 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조. 상장폐지결정 bddd 코멘트(확인일) vs 상장결정(결정일) 구분 유지. fixture json 키는 struct 태그와 1:1.
