# OpenDART DS004 지분공시 종합정보 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS004 지분공시 종합정보 2개 API** — 신규 `ownership` 패키지 + 공용 list 헬퍼 추출

## 배경 & 목표

DS001(공시정보)·DS002(정기보고서 주요정보 28/30)·DS003(재무정보 7/7)은 main 에 머지됨. 이 spec 은
DS004 지분공시 종합정보다. 2개 API(대량보유 상황보고, 임원·주요주주 소유보고) 모두 `corp_code` 만
받고 JSON list 를 반환한다. 정기보고서가 아니므로 새 sub-package `ownership/` 로 분리한다.

이 시점에 두 번째 카테고리가 동일한 list 패턴을 요구하므로, `report` 의 unexported 제네릭 list
헬퍼를 공용 `internal/httpclient.GetList[T]` 로 승격하고 `report`·`ownership` 둘 다 재사용한다.

## API 표면 (docs 기반 사실)

- **majorstock** (대량보유 상황보고, 5% 룰): `GET /api/majorstock.json`, 파라미터 `corp_code`. JSON `list[]`.
- **elestock** (임원·주요주주 소유보고): `GET /api/elestock.json`, 파라미터 `corp_code`. JSON `list[]`.
- 두 API 모두 `bsns_year`/`reprt_code` 없음 — `corp_code` 단일 파라미터. 숫자/비율은 문자열.

## 아키텍처

### 1) 공용 list 헬퍼 추출 (`internal/httpclient`)

`report` 의 `getListParams[T]`/`listResponse[T]` 로직을 `internal/httpclient` 로 승격한다(신규
`httpclient/list.go`):

```go
// listEnvelope 는 OpenDART list 응답 공통 형태.
type listEnvelope[T any] struct {
	Envelope
	List []T `json:"list"`
}

// GetList 는 JSON list 응답을 디코드해 list 만 반환한다.
// status 검사(013→ErrNoData, 그 외→*APIError)는 GetJSON 이 수행한다.
func GetList[T any](ctx context.Context, c *Client, path string, params map[string]string) ([]T, error) {
	var resp listEnvelope[T]
	if err := c.GetJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
```
`listEnvelope` 가 `Envelope` 를 임베드하므로 `StatusReader` 를 자동 충족한다.

`report/client.go` 는 위임으로 변경(공개 API·동작 불변; 기존 39개 테스트가 회귀 검증):
```go
func getListParams[T any](ctx context.Context, hc *httpclient.Client, path string, params map[string]string) ([]T, error) {
	return httpclient.GetList[T](ctx, hc, path, params)
}
// getList(ReportParams) wrapper 유지. report 의 기존 listResponse[T] 정의는 제거(미사용).
```

### 2) 신규 sub-package `ownership/` (DS004)

`client.Ownership` 로 노출. root `client.go` 에 와이어링 1줄.

```
ownership/
  client.go         # Client + New
  ownership.go      # 2개 메서드 + item struct
  client_test.go    # newTestClient 헬퍼
  ownership_test.go # 2개 fixture 테스트
  testdata/         # majorstock/elestock 실 응답 fixture
internal/httpclient/
  list.go           # GetList[T] + listEnvelope[T] (신규)
report/client.go    # (수정) getListParams 위임, listResponse 제거
client.go           # (수정) Ownership 필드 + 와이어링
README.md           # (수정) DS004 커버리지
integration_test.go # (수정) MajorStockReports 통합 케이스
```

## 2개 메서드 (ownership/ownership.go)

```go
// MajorStockItem 은 대량보유 상황보고 (majorstock) 한 건.
type MajorStockItem struct {
	RceptNo    string `json:"rcept_no"`    // 접수번호
	RceptDt    string `json:"rcept_dt"`    // 접수일자
	CorpCode   string `json:"corp_code"`   // 고유번호
	CorpName   string `json:"corp_name"`   // 회사명
	ReportTp   string `json:"report_tp"`   // 보고구분
	Repror     string `json:"repror"`      // 대표보고자
	Stkqy      string `json:"stkqy"`       // 보유주식등의 수
	StkqyIrds  string `json:"stkqy_irds"`  // 보유주식등의 증감
	Stkrt      string `json:"stkrt"`       // 보유비율
	StkrtIrds  string `json:"stkrt_irds"`  // 보유비율 증감
	CtrStkqy   string `json:"ctr_stkqy"`   // 주요체결 주식등의 수
	CtrStkrt   string `json:"ctr_stkrt"`   // 주요체결 보유비율
	ReportResn string `json:"report_resn"` // 보고사유
}

// MajorStockReports 는 대량보유 상황보고(5% 룰)를 조회한다.
func (c *Client) MajorStockReports(ctx context.Context, corpCode string) ([]MajorStockItem, error) {
	return httpclient.GetList[MajorStockItem](ctx, c.http, "/api/majorstock.json", map[string]string{"corp_code": corpCode})
}

// ExecutiveStockItem 은 임원·주요주주 소유보고 (elestock) 한 건.
type ExecutiveStockItem struct {
	RceptNo            string `json:"rcept_no"`              // 접수번호
	RceptDt            string `json:"rcept_dt"`              // 접수일자
	CorpCode           string `json:"corp_code"`             // 고유번호
	CorpName           string `json:"corp_name"`             // 회사명
	Repror             string `json:"repror"`                // 보고자
	IsuExctvRgistAt    string `json:"isu_exctv_rgist_at"`    // 발행회사 관계 임원(등기여부)
	IsuExctvOfcps      string `json:"isu_exctv_ofcps"`       // 발행회사 관계 임원 직위
	IsuMainShrholdr    string `json:"isu_main_shrholdr"`     // 발행회사 관계 주요주주
	SpStockLmpCnt      string `json:"sp_stock_lmp_cnt"`      // 특정증권등 소유 수
	SpStockLmpIrdsCnt  string `json:"sp_stock_lmp_irds_cnt"` // 특정증권등 소유 증감 수
	SpStockLmpRate     string `json:"sp_stock_lmp_rate"`     // 특정증권등 소유 비율
	SpStockLmpIrdsRate string `json:"sp_stock_lmp_irds_rate"`// 특정증권등 소유 증감 비율
}

// ExecutiveStockReports 는 임원·주요주주 소유보고를 조회한다.
func (c *Client) ExecutiveStockReports(ctx context.Context, corpCode string) ([]ExecutiveStockItem, error) {
	return httpclient.GetList[ExecutiveStockItem](ctx, c.http, "/api/elestock.json", map[string]string{"corp_code": corpCode})
}
```

`ownership/client.go`:
```go
// Package ownership 는 OpenDART DS004 지분공시 종합정보 API sub-client 다.
package ownership

import "github.com/kenshin579/opendart/internal/httpclient"

type Client struct {
	http *httpclient.Client
}

func New(http *httpclient.Client) *Client { return &Client{http: http} }
```

## root 와이어링 (client.go)

```go
type Client struct {
	http *httpclient.Client
	corp *corpcode.Cache

	Disclosure *disclosure.Client // DS001 공시정보
	Report     *report.Client     // DS002 정기보고서 주요정보 + DS003 재무정보
	Ownership  *ownership.Client  // DS004 지분공시 종합정보
}
// NewClient 내부: c.Ownership = ownership.New(hc)
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError` (httpclient.GetList 가 GetJSON 통해 처리).

## 테스트 전략

- `internal/httpclient` 의 `GetList` 는 기존 GetJSON 테스트 + report/ownership 사용처가 커버 (별도 단위
  테스트는 선택). report 위임 후 기존 report 테스트 39개로 회귀 검증.
- `ownership/client_test.go`: disclosure/report 와 동일한 fixture 서빙 `newTestClient` 헬퍼.
- `ownership/ownership_test.go`: 2개 메서드 각각 실 응답 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 로 캡처해 임베드(데이터 있는 종목/접수).
- `integration_test.go` 에 `MajorStockReports` 통합 케이스(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(string 유지), UTF-8.
- README "커버리지" 에 DS004 지분공시(대량보유·임원/주요주주 소유) 추가.

## 비범위 (후속 plan)

- DS005 주요사항보고서 주요정보, DS006 증권신고서 주요정보.
- DS002 개인별 보수 Ver 2.0 2종.
