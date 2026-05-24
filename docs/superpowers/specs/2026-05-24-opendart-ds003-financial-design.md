# OpenDART DS003 정기보고서 재무정보 — 재무 핵심(JSON) 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS003 재무정보 중 JSON list 5개 API** (`report` 패키지 확장). XBRL 바이너리·택사노미 2개는 후속.

## 배경 & 목표

DS001(공시정보)·DS002(정기보고서 주요정보 28/30)는 main 에 머지됨. 이 spec 은 DS003 정기보고서
재무정보의 첫 그룹 — **재무 핵심 JSON 5개**다. DS003 은 DS002 와 달리 추가 파라미터(`fs_div`,
`idx_cl_code`)와 다중회사(콤마 corp_code)가 있어, 기존 `getList` 를 raw-map 코어로 한 번
리팩토링하고 엔드포인트별 파라미터 struct 를 도입한다. XBRL 원본파일(바이너리 ZIP)과 택사노미는
형태가 근본적으로 달라 별도 그룹으로 보류한다.

## API 표면 (docs 기반 사실)

- 5개 모두 JSON `list[]` 응답. 금액은 콤마 문자열.
- 요청 파라미터:
  - `fnlttSinglAcnt`(단일 주요계정), `fnlttMultiAcnt`(다중 주요계정): `corp_code`+`bsns_year`+`reprt_code` (= 기존 ReportParams). 다중은 `corp_code` 를 콤마로 여러 개.
  - `fnlttSinglAcntAll`(단일 전체 재무제표): + `fs_div` (OFS 개별/CFS 연결).
  - `fnlttSinglIndx`(단일 주요 재무지표), `fnlttCmpnyIndx`(다중 주요 재무지표): + `idx_cl_code` (M210000 수익성/M220000 안정성/M230000 성장성/M240000 활동성). 다중은 `corp_code` 콤마.
- **응답 필드 중복:** 단일 주요계정 = 다중 주요계정(동일 필드) → `AccountItem` 공유. 단일 재무지표 = 다중 재무지표(동일) → `FinancialIndexItem` 공유. 전체 재무제표만 별도(`FullStatementItem`). 즉 5개 메서드 = 3개 struct.

## 아키텍처

`report` 패키지에 새 파일 `financial.go` 추가. `client.Report` 는 이미 root 에 와이어링됨 → root 변경
불필요. 기존 `report/client.go` 의 `getList` 를 raw-map 코어로 리팩토링한다(기존 DS002 호출 영향 없음).

```
report/
  client.go        # (수정) getList → getListParams[T] raw-map 코어 + thin wrapper
  financial.go     # 5개 메서드 + 3 item struct + 2 파라미터 struct + 타입 상수 (신규)
  financial_test.go
  testdata/        # 5개 실 응답 JSON fixture 추가
README.md          # (수정) DS003 커버리지
integration_test.go  # (수정) SingleAccount 통합 케이스
```

### getList 리팩토링 (report/client.go)

```go
// getListParams 는 raw 파라미터 맵으로 list 를 조회하는 코어 헬퍼.
func getListParams[T any](ctx context.Context, hc *httpclient.Client, path string, params map[string]string) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}

// getList 는 ReportParams 기반 thin wrapper (기존 DS002 메서드가 사용).
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	return getListParams[T](ctx, hc, path, p.toMap())
}
```

### 파라미터 struct + 타입 상수 (financial.go)

```go
// FsDiv 는 개별/연결 구분.
type FsDiv string

const (
	FsDivSeparate     FsDiv = "OFS" // 재무제표(개별)
	FsDivConsolidated FsDiv = "CFS" // 연결재무제표
)

// IndexClass 는 재무지표 분류 코드.
type IndexClass string

const (
	IndexProfitability IndexClass = "M210000" // 수익성지표
	IndexStability     IndexClass = "M220000" // 안정성지표
	IndexGrowth        IndexClass = "M230000" // 성장성지표
	IndexActivity      IndexClass = "M240000" // 활동성지표
)

// FinancialStatementParams 는 전체 재무제표 요청 인자.
type FinancialStatementParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
	FsDiv     FsDiv      // 개별(OFS)/연결(CFS)
}

func (p FinancialStatementParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
		"fs_div":     string(p.FsDiv),
	}
}

// FinancialIndexParams 는 주요 재무지표 요청 인자. 다중회사는 CorpCode 를 콤마로 구분.
type FinancialIndexParams struct {
	CorpCode  string     // 고유번호 (8자리; 다중은 콤마 구분 "code1,code2")
	BsnsYear  string     // 사업연도 (4자리)
	ReprtCode ReportCode // 보고서 코드
	IdxClCode IndexClass // 지표분류코드
}

func (p FinancialIndexParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":   p.CorpCode,
		"bsns_year":   p.BsnsYear,
		"reprt_code":  string(p.ReprtCode),
		"idx_cl_code": string(p.IdxClCode),
	}
}
```

## 5개 메서드 (3 struct 공유)

| 메서드 | 한글 | 엔드포인트 | 파라미터 | item |
|--------|------|-----------|----------|------|
| `SingleAccount` | 단일회사 주요계정 | `/api/fnlttSinglAcnt.json` | `ReportParams` | `AccountItem` |
| `MultiAccount` | 다중회사 주요계정 | `/api/fnlttMultiAcnt.json` | `ReportParams` (corp_code 콤마) | `AccountItem` |
| `SingleFullStatement` | 단일회사 전체 재무제표 | `/api/fnlttSinglAcntAll.json` | `FinancialStatementParams` | `FullStatementItem` |
| `SingleIndex` | 단일회사 주요 재무지표 | `/api/fnlttSinglIndx.json` | `FinancialIndexParams` | `FinancialIndexItem` |
| `MultiIndex` | 다중회사 주요 재무지표 | `/api/fnlttCmpnyIndx.json` | `FinancialIndexParams` (corp_code 콤마) | `FinancialIndexItem` |

```go
// AccountItem 은 단일/다중회사 주요계정 (fnlttSinglAcnt / fnlttMultiAcnt) 한 건.
type AccountItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	StockCode       string `json:"stock_code"`        // 종목 코드
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	AccountNm       string `json:"account_nm"`        // 계정명
	FsDiv           string `json:"fs_div"`            // 개별/연결구분
	FsNm            string `json:"fs_nm"`             // 개별/연결명
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmDt        string `json:"thstrm_dt"`         // 당기일자
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmDt        string `json:"frmtrm_dt"`         // 전기일자
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmDt     string `json:"bfefrmtrm_dt"`      // 전전기일자
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}

// FinancialIndexItem 은 단일/다중회사 주요 재무지표 (fnlttSinglIndx / fnlttCmpnyIndx) 한 건.
type FinancialIndexItem struct {
	ReprtCode string `json:"reprt_code"`  // 보고서 코드
	BsnsYear  string `json:"bsns_year"`   // 사업 연도
	CorpCode  string `json:"corp_code"`   // 고유번호
	StockCode string `json:"stock_code"`  // 종목 코드
	StlmDt    string `json:"stlm_dt"`     // 결산기준일
	IdxClCode string `json:"idx_cl_code"` // 지표분류코드
	IdxClNm   string `json:"idx_cl_nm"`   // 지표분류명
	IdxCode   string `json:"idx_code"`    // 지표코드
	IdxNm     string `json:"idx_nm"`      // 지표명
	IdxVal    string `json:"idx_val"`     // 지표값
}

// FullStatementItem 은 단일회사 전체 재무제표 (fnlttSinglAcntAll) 한 건.
type FullStatementItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	CorpCode        string `json:"corp_code"`         // 고유번호
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	AccountId       string `json:"account_id"`        // 계정ID
	AccountNm       string `json:"account_nm"`        // 계정명
	AccountDetail   string `json:"account_detail"`    // 계정상세
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmQNm       string `json:"frmtrm_q_nm"`       // 전기명(분/반기)
	FrmtrmQAmount   string `json:"frmtrm_q_amount"`   // 전기금액(분/반기)
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}
```

메서드 (financial.go):
```go
func (c *Client) SingleAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttSinglAcnt.json", p)
}
func (c *Client) MultiAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttMultiAcnt.json", p)
}
func (c *Client) SingleFullStatement(ctx context.Context, p FinancialStatementParams) ([]FullStatementItem, error) {
	return getListParams[FullStatementItem](ctx, c.http, "/api/fnlttSinglAcntAll.json", p.toMap())
}
func (c *Client) SingleIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttSinglIndx.json", p.toMap())
}
func (c *Client) MultiIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttCmpnyIndx.json", p.toMap())
}
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `report/financial_test.go`: 기존 `report/client_test.go` 의 `newTestClient` 재사용.
- 5개 메서드 각각: 실 응답 JSON fixture 디코딩 → 대표 필드 매핑 검증. 파라미터 struct 의 `toMap`
  (fs_div/idx_cl_code 포함) 도 별도 테스트.
- fixture 는 실 API 로 캡처해 임베드(계획 작성 단계, 삼성전자 2023 사업보고서; 다중회사는 삼성+SK하이닉스 등 콤마).
- `integration_test.go` 에 `SingleAccount` 통합 케이스 추가(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item/파라미터 struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(콤마 string 유지), UTF-8.
- README "커버리지" 에 DS003 재무정보(주요계정·전체재무제표·재무지표) 추가.

## 비범위 (후속 plan)

- DS003 XBRL 그룹: 재무제표 원본파일(`fnlttXbrl`, 바이너리 ZIP, params rcept_no+reprt_code),
  XBRL택사노미 재무제표양식(`xbrlTaxonomy`, param sj_div).
- DS004~DS006 카테고리.
- 숫자/통화 coercion 헬퍼.
