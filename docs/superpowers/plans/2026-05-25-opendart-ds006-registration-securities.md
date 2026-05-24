# OpenDART DS006 증권신고서 Sub-1 «증권» Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DS006 증권신고서 Sub-1 «증권» 3개 API(지분증권/채무증권/증권예탁증권)를 신규 `registration` 패키지에 추가한다. 그룹형 응답 디코더 `httpclient.GetGroups` 도 신설한다.

**Architecture:** DS006 응답은 그룹형(`{status, message, group:[{title, list:[...]}]}`). 신규 `internal/httpclient.GetGroups` 가 status 엔벨로프(013→ErrNoData) 처리 후 `[]Group{Title, List json.RawMessage}` 반환. 신규 `registration` 패키지: `Client`/`New`/`Params`. 각 엔드포인트 메서드는 GetGroups 호출 → title 별 switch 로 `json.Unmarshal` → 그룹 슬라이스 wrapper 반환. root `client.go` 에 `Registration *registration.Client` 와이어링.

**Tech Stack:** Go, 표준 net/http(`internal/httpclient`), `encoding/json`, testify, httptest. fixture 는 실 API 캡처 권장(남양유업 estkRs 확인됨; 불가 시 docs 스키마 일치 샘플, 그룹 배열 형태 유지).

**SINGLE SOURCE OF TRUTH for struct definitions:** `docs/superpowers/specs/2026-05-25-opendart-ds006-registration-securities-design.md` 에 인프라 코드 + 10개 item struct + 3개 wrapper + 3개 메서드가 EXACT Go 코드로 정의돼 있다. 각 Task 는 그 파일의 동일 이름 정의를 **그대로 복사**한다(필드 가감 금지).

---

## File Structure

- Create: `internal/httpclient/groups.go` — `Group`, `groupEnvelope`, `GetGroups`.
- Create: `registration/client.go` — `Client`, `New`, `Params`+`toMap`.
- Create: `registration/securities.go` — 10 item struct + 3 wrapper + 3 메서드.
- Create: `registration/client_test.go` — `newTestClient` (httptest 단일 fixture 서빙).
- Create: `registration/securities_test.go` — 3개 fixture 테스트.
- Create: `registration/testdata/{estkRs,bdRs,stkdpRs}.json` — 3 fixture.
- Modify: `client.go` (root) — `Registration *registration.Client` 필드 + `c.Registration = registration.New(hc)`.
- Modify: `integration_test.go` — 통합 케이스 2개(`//go:build integration`).
- Modify: `README.md` — DS006 커버리지 줄 신설.

기존 컨벤션(변경 금지):
- root `client.go` 에 이미 `Disclosure`/`Report`/`Ownership`/`Material` 필드가 있고 `New` 류에서 `material.New(hc)` 식으로 와이어링됨. 같은 패턴으로 `Registration` 추가.
- `internal/httpclient` 의 `Client.GetJSON(ctx, path, params, out)` + `Envelope`(APIStatus) + `GetList[T]` 가 이미 존재. `GetGroups` 는 그 옆에 추가.
- 모든 응답 필드 string + 한글 코멘트. UTF-8. testify.

`go test` 의 작업 디렉터리는 항상 `cd /Users/user/src/workspace_moneyflow/opendart`.

---

## Task 1: 인프라 — GetGroups 헬퍼 + registration 패키지 골격 + 지분증권(EquitySecurities)

이 Task 는 신규 인프라를 세우고 첫 엔드포인트까지 동작시킨다.

**Files:** Create `internal/httpclient/groups.go`, `registration/client.go`, `registration/securities.go`, `registration/client_test.go`, `registration/securities_test.go`, `registration/testdata/estkRs.json`; Modify root `client.go`.

- [ ] **Step 1: `internal/httpclient/groups.go`** — spec 의 "internal/httpclient.GetGroups (신규)" 코드 그대로. 파일:
```go
package httpclient

import (
	"context"
	"encoding/json"
)

// Group 은 증권신고서(DS006) 응답의 그룹 하나(title + raw list).
type Group struct {
	Title string          `json:"title"` // 그룹명칭(예: 일반사항)
	List  json.RawMessage `json:"list"`  // 그룹 항목 배열(타입은 호출측에서 결정)
}

type groupEnvelope struct {
	Envelope
	Group []Group `json:"group"`
}

// GetGroups 는 그룹형(DS006) 응답을 디코딩해 그룹 목록을 반환한다.
func GetGroups(ctx context.Context, c *Client, path string, params map[string]string) ([]Group, error) {
	var env groupEnvelope
	if err := c.GetJSON(ctx, path, params, &env); err != nil {
		return nil, err
	}
	return env.Group, nil
}
```
주의: `Envelope` 가 `APIStatus()`(StatusReader)를 제공하는지 기존 `internal/httpclient/client.go`(또는 list.go)에서 확인. `GetList` 의 `listEnvelope` 와 동일하게 `Envelope` 임베드로 status 검사가 동작해야 함. 만약 `GetJSON` 이 `StatusReader` 를 요구하면 `groupEnvelope` 가 이를 만족하는지(임베드된 `Envelope.APIStatus()`) 확인하고, 아니라면 `listEnvelope` 와 동일한 방식을 따른다.

- [ ] **Step 2: `internal/httpclient/groups.go` 컴파일 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./internal/httpclient/` → 성공. (status 검사 인터페이스 미충족 시 `listEnvelope` 패턴에 맞춰 수정.)

- [ ] **Step 3: `registration/client.go`** — spec 의 "registration 패키지" 코드 그대로:
```go
package registration

import (
	"encoding/json"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Client 는 DS006 증권신고서 주요정보 API 클라이언트.
type Client struct {
	http *httpclient.Client
}

// New 는 registration.Client 를 만든다.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}

// Params 는 DS006 증권신고서 공통 요청 파라미터(corp_code + 기간).
type Params struct {
	CorpCode string // 고유번호(8자리)
	BgnDe    string // 검색시작 접수일자(YYYYMMDD)
	EndDe    string // 검색종료 접수일자(YYYYMMDD)
}

func (p Params) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" {
		m["bgn_de"] = p.BgnDe
	}
	if p.EndDe != "" {
		m["end_de"] = p.EndDe
	}
	return m
}

var _ = json.Unmarshal // securities.go 에서 사용
```
(마지막 `var _` 줄은 client.go 에 json import 가 떠서 unused 가 되면 제거; securities.go 가 같은 패키지에서 json 을 쓰면 client.go 의 json import 는 불필요하니 client.go 에서 `encoding/json` import 를 빼는 게 맞다. **권장**: client.go 에서 json import 제거하고 `var _` 줄도 제거. json 은 securities.go 에서 import.)

정리된 client.go (권장):
```go
package registration

import "github.com/kenshin579/opendart/internal/httpclient"

type Client struct{ http *httpclient.Client }

func New(http *httpclient.Client) *Client { return &Client{http: http} }

type Params struct {
	CorpCode string // 고유번호(8자리)
	BgnDe    string // 검색시작 접수일자(YYYYMMDD)
	EndDe    string // 검색종료 접수일자(YYYYMMDD)
}

func (p Params) toMap() map[string]string {
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

- [ ] **Step 4: root `client.go` 와이어링** — 먼저 `client.go` 를 읽어 기존 필드/생성 패턴 확인. `Material *material.Client` 옆에 `Registration *registration.Client` 필드를 추가하고, `c.Material = material.New(hc)` 옆에 `c.Registration = registration.New(hc)` 를 추가. import 에 `"github.com/kenshin579/opendart/registration"` 추가. (hc 는 기존에 쓰이는 `*httpclient.Client` 인스턴스명 — 실제 코드의 이름을 그대로 사용.)

- [ ] **Step 5: fixture `registration/testdata/estkRs.json`** — 그룹형. 남양유업 실 데이터 기반 샘플(6개 그룹 모두 포함, 각 그룹 list 1건):
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "sbd": "2023년 06월 01일 ~ 2023년 06월 02일", "pymd": "2023년 06월 12일", "sband": "-", "asand": "2023년 05월 17일", "asstd": "2023년 05월 08일", "exstk": "-", "exprc": "-", "expd": "-", "rpt_rcpn": "20230331003707"}
    ]},
    {"title": "증권의종류", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "stksen": "기명식 보통주", "stkcnt": "1,000,000", "fv": "5,000", "slprc": "50,000", "slta": "50,000,000,000", "slmthn": "주주배정후 실권주 일반공모"}
    ]},
    {"title": "인수인정보", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "actsen": "대표주관회사", "actnmn": "한국투자증권", "stksen": "기명식 보통주", "udtcnt": "1,000,000", "udtamt": "50,000,000,000", "udtprc": "잔액인수", "udtmth": "총액인수"}
    ]},
    {"title": "자금의사용목적", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "se": "운영자금", "amt": "50,000,000,000"}
    ]},
    {"title": "매출인에관한사항", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "hdr": "-", "rl_cmp": "-", "bfsl_hdstk": "-", "slstk": "-", "atsl_hdstk": "-"}
    ]},
    {"title": "일반청약자환매청구권", "list": [
      {"rcept_no": "20230515002454", "corp_cls": "Y", "corp_code": "00107598", "corp_name": "남양유업", "grtrs": "-", "exavivr": "-", "grtcnt": "-", "expd": "-", "exprc": "-"}
    ]}
  ]
}
```

- [ ] **Step 6: `registration/client_test.go`** — httptest 단일 fixture 서빙 헬퍼:
```go
package registration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 지정한 testdata fixture 를 모든 요청에 서빙하는 registration.Client 를 만든다.
func newTestClient(t *testing.T, fixture string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(filepath.Join("testdata", fixture))
		require.NoError(t, err)
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	return New(hc)
}
```
주의: `httpclient.New(httpclient.Config{...})` 의 정확한 필드명은 기존 `material/client_test.go` 의 `newTestClient` 를 참고해 동일하게 맞춘다(APIKey/BaseURL/HTTPClient).

- [ ] **Step 7: securities_test.go 에 지분증권 테스트** — `registration/securities_test.go` 생성:
```go
package registration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquitySecurities(t *testing.T) {
	c := newTestClient(t, "estkRs.json")
	res, err := c.EquitySecurities(context.Background(), Params{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.SecurityTypes, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 1)
	require.Len(t, res.Sellers, 1)
	require.Len(t, res.RetailPutbackOption, 1)
	assert.Equal(t, "20230515002454", res.General[0].RceptNo)
	assert.Equal(t, "남양유업", res.General[0].CorpName)
	assert.Equal(t, "기명식 보통주", res.SecurityTypes[0].Stksen)
	assert.Equal(t, "한국투자증권", res.Underwriters[0].Actnmn)
	assert.Equal(t, "운영자금", res.FundUsage[0].Se)
}
```

- [ ] **Step 8: `go test ./registration/ -run TestEquitySecurities`** — confirm FAIL (EquitySecurities/types undefined).

- [ ] **Step 9: `registration/securities.go`** — spec 에서 **공유 item 6종**(RsGeneralItem, RsSecurityTypeItem, RsUnderwriterItem, RsFundUsageItem, RsSellerItem, EquityRetailPutbackOptionItem) + **EquitySecuritiesRegistration wrapper + EquitySecurities 메서드** 를 VERBATIM 복사. 파일 헤더:
```go
package registration

import (
	"context"
	"encoding/json"

	"github.com/kenshin579/opendart/internal/httpclient"
)
```
(이 Task 에서는 채무/예탁 타입·메서드는 아직 추가하지 않음 — Task 2/3.)

- [ ] **Step 10: `go test ./registration/ -run TestEquitySecurities`** — PASS.
- [ ] **Step 11: gofmt** — `gofmt -l internal/httpclient/groups.go registration/` 출력 없으면 OK(있으면 `gofmt -w`).
- [ ] **Step 12: 전체 빌드** — `go build ./...` 성공(root 와이어링 포함).
- [ ] **Step 13: 실 API fixture 캡처(권장)** — `$OPENDART_API_KEY` 있으면 `curl -s "https://opendart.fss.or.kr/api/estkRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00107598&bgn_de=20180101&end_de=20241231"` (남양유업, 데이터 있음 확인됨). 서버 TLS1.2 RSA. clean `"status":"000"` + group 있으면 그 응답을 fixture 로 저장하고 test assert 를 실제 값에 맞춰 갱신 후 step 10 재실행. 1회 시도.
- [ ] **Step 14: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add internal/httpclient/groups.go registration/ client.go
git commit -m "feat(registration): add DS006 GetGroups infra + EquitySecurities (지분증권)"
```

---

## Task 2: 채무증권 (DebtSecurities) — bdRs

**Files:** Modify `registration/securities.go`, `registration/securities_test.go`; Create `registration/testdata/bdRs.json`.

Struct: spec 의 **DebtGeneralItem(40)/DebtUnderwriterItem(12)/DebtFundUsageItem(7)/DebtSellerItem(10)** + **DebtSecuritiesRegistration wrapper + DebtSecurities 메서드** VERBATIM.

- [ ] **Step 1: fixture `registration/testdata/bdRs.json`** — 그룹형, 4개 그룹(일반사항/인수인정보/자금의사용목적/매출인에관한사항) 각 list 1건. 일반사항 40필드 전체 키 포함(대표값 일부, 나머지 "-"):
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230601000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트채무증권", "tm": "10", "bdnmn": "제10회 무보증사채", "slmth": "공모", "fta": "100,000,000,000", "slta": "100,000,000,000", "isprc": "10,000", "intr": "4.5", "isrr": "4.5", "rpd": "2026년 06월 01일", "print_pymint": "국민은행", "mngt_cmp": "한국투자증권", "cdrt_int": "AA(한국신용평가)", "sbd": "2023년 05월 30일", "pymd": "2023년 06월 01일", "sband": "-", "asand": "-", "asstd": "-", "dpcrn": "원화", "dpcr_amt": "-", "usarn": "국내", "usntn": "대한민국", "wnexpl_at": "-", "udtintnm": "한국투자증권", "grt_int": "-", "grt_amt": "-", "icmg_mgknd": "-", "icmg_mgamt": "-", "estk_exstk": "-", "estk_exrt": "-", "estk_exprc": "-", "estk_expd": "-", "rpt_rcpn": "-", "drcb_at": "미해당", "drcb_uast": "-", "drcb_optknd": "-", "drcb_mtd": "-"}
    ]},
    {"title": "인수인정보", "list": [
      {"rcept_no": "20230601000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트채무증권", "tm": "10", "actsen": "대표주관회사", "actnmn": "한국투자증권", "stksen": "무보증사채", "udtcnt": "50,000,000,000", "udtamt": "50,000,000,000", "udtprc": "총액인수", "udtmth": "총액인수"}
    ]},
    {"title": "자금의사용목적", "list": [
      {"rcept_no": "20230601000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트채무증권", "tm": "10", "se": "채무상환자금", "amt": "100,000,000,000"}
    ]},
    {"title": "매출인에관한사항", "list": [
      {"rcept_no": "20230601000111", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트채무증권", "tm": "10", "hdr": "-", "rl_cmp": "-", "bfsl_hdstk": "-", "slstk": "-", "atsl_hdstk": "-"}
    ]}
  ]
}
```

- [ ] **Step 2: append test to `registration/securities_test.go`**
```go
func TestDebtSecurities(t *testing.T) {
	c := newTestClient(t, "bdRs.json")
	res, err := c.DebtSecurities(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 1)
	require.Len(t, res.Sellers, 1)
	assert.Equal(t, "20230601000111", res.General[0].RceptNo)
	assert.Equal(t, "10", res.General[0].Tm)
	assert.Equal(t, "제10회 무보증사채", res.General[0].Bdnmn)
	assert.Equal(t, "4.5", res.General[0].Intr)
	assert.Equal(t, "채무상환자금", res.FundUsage[0].Se)
}
```

- [ ] **Step 3: run `go test ./registration/ -run TestDebtSecurities`** — FAIL (undefined).
- [ ] **Step 4: append to `registration/securities.go`**: Debt* 4 item struct + DebtSecuritiesRegistration wrapper + DebtSecurities 메서드 VERBATIM from spec.
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l registration/` 출력 없으면 OK.
- [ ] **Step 7: verify** DebtGeneralItem has 40 json tags.
- [ ] **Step 8: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/bdRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20180101&end_de=20241231"` (삼성전자 회사채 발행이력 있을 수 있음). clean 000+group → fixture 교체 + assert 갱신, 재실행. 안되면 다른 회사채 발행사 1~2개 시도(예: 00164779). 모두 013 이면 샘플 유지. 최대 2회.
- [ ] **Step 9: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add registration/securities.go registration/securities_test.go registration/testdata/bdRs.json
git commit -m "feat(registration): add DS006 DebtSecurities (채무증권)"
```

---

## Task 3: 증권예탁증권 (DepositaryReceipts) — stkdpRs

**Files:** Modify `registration/securities.go`, `registration/securities_test.go`; Create `registration/testdata/stkdpRs.json`.

Struct: spec 의 **DepositaryReceiptsRegistration wrapper + DepositaryReceipts 메서드** VERBATIM. (item 타입은 공유 Rs* 재사용 — 신규 타입 없음.)

- [ ] **Step 1: fixture `registration/testdata/stkdpRs.json`** — 그룹형 5개 그룹(일반사항/증권의종류/인수인정보/자금의사용목적/매출인에관한사항), 각 list 1건. 스키마는 Rs* 공유 타입과 동일:
```json
{
  "status": "000",
  "message": "정상",
  "group": [
    {"title": "일반사항", "list": [
      {"rcept_no": "20230701000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트예탁증권", "sbd": "2023년 07월 10일", "pymd": "2023년 07월 20일", "sband": "-", "asand": "-", "asstd": "-", "exstk": "-", "exprc": "-", "expd": "-", "rpt_rcpn": "-"}
    ]},
    {"title": "증권의종류", "list": [
      {"rcept_no": "20230701000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트예탁증권", "stksen": "주식예탁증서(DR)", "stkcnt": "2,000,000", "fv": "5,000", "slprc": "60,000", "slta": "120,000,000,000", "slmthn": "해외모집"}
    ]},
    {"title": "인수인정보", "list": [
      {"rcept_no": "20230701000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트예탁증권", "actsen": "대표주관회사", "actnmn": "외국계증권사", "stksen": "DR", "udtcnt": "2,000,000", "udtamt": "120,000,000,000", "udtprc": "총액인수", "udtmth": "총액인수"}
    ]},
    {"title": "자금의사용목적", "list": [
      {"rcept_no": "20230701000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트예탁증권", "se": "시설자금", "amt": "120,000,000,000"}
    ]},
    {"title": "매출인에관한사항", "list": [
      {"rcept_no": "20230701000222", "corp_cls": "Y", "corp_code": "00126380", "corp_name": "테스트예탁증권", "hdr": "-", "rl_cmp": "-", "bfsl_hdstk": "-", "slstk": "-", "atsl_hdstk": "-"}
    ]}
  ]
}
```

- [ ] **Step 2: append test to `registration/securities_test.go`**
```go
func TestDepositaryReceipts(t *testing.T) {
	c := newTestClient(t, "stkdpRs.json")
	res, err := c.DepositaryReceipts(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.SecurityTypes, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 1)
	require.Len(t, res.Sellers, 1)
	assert.Equal(t, "20230701000222", res.General[0].RceptNo)
	assert.Equal(t, "주식예탁증서(DR)", res.SecurityTypes[0].Stksen)
	assert.Equal(t, "시설자금", res.FundUsage[0].Se)
}
```

- [ ] **Step 3: run `go test ./registration/ -run TestDepositaryReceipts`** — FAIL (undefined).
- [ ] **Step 4: append to `registration/securities.go`**: DepositaryReceiptsRegistration wrapper + DepositaryReceipts 메서드 VERBATIM from spec (공유 Rs* 타입 재사용).
- [ ] **Step 5: run test** — PASS.
- [ ] **Step 6: gofmt** — `gofmt -l registration/` 출력 없으면 OK.
- [ ] **Step 7: real API BEST EFFORT** — `curl -s "https://opendart.fss.or.kr/api/stkdpRs.json?crtfc_key=$OPENDART_API_KEY&corp_code=00126380&bgn_de=20180101&end_de=20241231"`. 대부분 013(국내사 DR 발행 드뭄). 안되면 샘플 유지. 1회.
- [ ] **Step 8: commit (no push)**
```bash
cd /Users/user/src/workspace_moneyflow/opendart
git add registration/securities.go registration/securities_test.go registration/testdata/stkdpRs.json
git commit -m "feat(registration): add DS006 DepositaryReceipts (증권예탁증권)"
```

---

## Task 4: 통합 테스트 + README

**Files:** Modify `integration_test.go`, `README.md`.

- [ ] **Step 1: 통합 테스트 추가**

`integration_test.go` 를 먼저 읽어 패턴 확인(`//go:build integration`, `package opendart`, `NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))`, `ErrNoData` 직접 참조). import 에 `"github.com/kenshin579/opendart/registration"` 추가. 파일 끝에 2개 추가:
```go
func TestIntegration_EquitySecurities(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	// 남양유업: 2023 유상증자 지분증권 신고서 존재
	res, err := c.Registration.EquitySecurities(context.Background(), registration.Params{CorpCode: "00107598", BgnDe: "20180101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 지분증권 신고서 데이터 없음")
	}
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, it := range res.General {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_DebtSecurities(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	res, err := c.Registration.DebtSecurities(context.Background(), registration.Params{CorpCode: corp, BgnDe: "20180101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 채무증권 신고서 데이터 없음")
	}
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, it := range res.General {
		require.NotEmpty(t, it.RceptNo)
	}
}
```

- [ ] **Step 2: 통합 빌드 확인** — `cd /Users/user/src/workspace_moneyflow/opendart && go vet -tags integration ./...` → 출력 없음.
- [ ] **Step 3: 통합 테스트 실행(키 있으면)** — `go test -tags integration -run "TestIntegration_EquitySecurities|TestIntegration_DebtSecurities" ./...` → PASS 또는 SKIP. (남양유업 EquitySecurities 는 PASS 기대.)

- [ ] **Step 4: README 커버리지 갱신**

`README.md` 의 DS004/DS005 줄 아래(또는 DS005 다음)에 DS006 줄 신설:
```
- DS006 증권신고서 주요정보: 지분증권 · 채무증권 · 증권예탁증권
```
그리고 `(예정)` 줄을 `(예정) DS006 나머지(합병/분할/주식의포괄적교환·이전) · DS002 개인별 보수 Ver2.0` 로 갱신(기존 "(예정) DS006 · DS002 ..." 에서 DS006 일부 구현 반영). 실제 파일의 현재 문구를 확인해 동등하게 편집.

- [ ] **Step 5: 전체 게이트** — `cd /Users/user/src/workspace_moneyflow/opendart && go build ./... && go test ./... && gofmt -l internal/httpclient/ registration/ integration_test.go` → 빌드 OK, 전체 PASS, gofmt 출력 없음.
- [ ] **Step 6: README UTF-8** — `file -I README.md` → `charset=utf-8`.
- [ ] **Step 7: 커밋** — `git add integration_test.go README.md && git commit -m "test(registration): add DS006 증권 통합 테스트 + README 커버리지"`.

---

## Self-Review (작성자 점검 결과)

**1. Spec coverage:** spec 의 인프라(GetGroups/Params/패키지) → Task 1. 3개 메서드(EquitySecurities/DebtSecurities/DepositaryReceipts) → Task 1/2/3. 공유 Rs* 6종은 Task 1, Debt* 4종은 Task 2, 예탁은 공유 재사용(Task 3). 통합·README = Task 4. 누락 없음.

**2. Placeholder scan:** TBD/TODO 없음. 인프라·메서드·fixture·test 완전. struct body 는 committed spec(EXACT Go 코드) 참조 — 레포 내 존재, 필드 수 명시. client.go json import 정리 주의 명시. httpclient status 인터페이스 충족 확인 단계 포함.

**3. Type consistency:** 메서드/wrapper/item 이름 spec 표와 1:1. 지분·예탁 공유 Rs* 타입, 채무 Debt* 별도. fixture 의 group title 한글이 메서드 switch case 와 정확히 일치(일반사항/증권의종류/인수인정보/자금의사용목적/매출인에관한사항/일반청약자환매청구권). 통합 테스트는 같은 패키지라 `ErrNoData` 직접 참조, registration import 추가.
