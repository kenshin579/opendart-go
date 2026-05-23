# OpenDART DS002 정기보고서 주요정보 — 지분·주식·배당 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS002 정기보고서 주요정보**의 공통 추상화(`report` 패키지) + **지분·주식·배당 7개 API**

## 배경 & 목표

라이브러리 기반 + DS001 공시정보는 main 에 머지됨(PR #2). 이 spec 은 다음 카테고리인
**DS002 정기보고서 주요정보**(총 30개 API)의 첫 수직 슬라이스다. DS002 의 30개 API 는 모두
동일한 요청 패턴(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답을 가지므로, 공통
추상화(`report` 패키지의 `ReportParams`/`ReportCode`/제네릭 `getList`)를 세우고 **지분·주식·배당
7개**를 끝까지 구현해 패턴을 확정한다. 나머지 23개는 동일 패턴 후속 plan 으로 양산한다.

## API 표면 (docs 기반 사실)

- 전 API 공통 요청: `crtfc_key`(자동 주입) + `corp_code`(8자리) + `bsns_year`(4자리, 2015 이후) + `reprt_code`(5자리).
- `reprt_code`: `11013`(1분기) / `11012`(반기) / `11014`(3분기) / `11011`(사업).
- 응답: 공통 envelope(`status`/`message`) + `list[]` 배열. 페이지네이션 없음.
- 숫자 필드는 콤마 포함 문자열(예: `9,999,999,999`) → `string` 유지(숫자 coercion 없음).
- 조회 데이터 없음 → status `013` → `ErrNoData`.

## 아키텍처

신규 sub-package `report/` (DS002). root `Client` 에 `Report *report.Client` 추가, `client.Report` 로 노출.
DS001 `disclosure` 와 동일한 sub-client 패턴. list 응답이 균일하므로 Go 제네릭으로 공통 헬퍼를 둔다.

```
report/
  client.go      # Client, New, ReportCode 상수, ReportParams, listResponse[T], getList[T]
  equity.go      # 7개 메서드 + item struct
  equity_test.go
  testdata/      # 7개 실 응답 JSON fixture
client.go        # (수정) Report 필드 + 와이어링
examples/report/main.go
README.md        # (수정) 커버리지에 DS002 7개 추가
integration_test.go  # (수정) Dividend 통합 케이스 추가
```

### 공통 헬퍼 (report/client.go)

```go
// Package report 는 OpenDART DS002 정기보고서 주요정보 API sub-client 다.
package report

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Client 는 정기보고서 주요정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// ReportCode 는 정기보고서 종류 코드.
type ReportCode string

const (
	Q1Report     ReportCode = "11013" // 1분기보고서
	HalfReport   ReportCode = "11012" // 반기보고서
	Q3Report     ReportCode = "11014" // 3분기보고서
	AnnualReport ReportCode = "11011" // 사업보고서
)

// ReportParams 는 DS002 공통 요청 인자 (전부 필수).
type ReportParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
}

func (p ReportParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
	}
}

// listResponse 는 DS002 공통 list 응답 envelope.
type listResponse[T any] struct {
	httpclient.Envelope
	List []T `json:"list"`
}

// getList 는 공통 list 조회 헬퍼. status 검사 후 list 만 반환한다.
// 조회 데이터 없음(013)은 httpclient 가 ErrNoData 로 변환한다.
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, p.toMap(), &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
```

> 메서드는 Go 메서드에 타입 파라미터를 둘 수 없으므로, 패키지 레벨 제네릭 함수 `getList[T]` 를
> 각 메서드가 `c.http` 를 넘겨 호출한다.

## 7개 엔드포인트 (report/equity.go)

각 메서드는 `func (c *Client) X(ctx, p ReportParams) ([]XItem, error) { return getList[XItem](ctx, c.http, "<path>", p) }`.
모든 item 은 공통 머리(`rcept_no`/`corp_cls`/`corp_code`/`corp_name`)와 꼬리(`stlm_dt`)를 가진다.

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `CapitalChange` | 증자(감자) 현황 | `/api/irdsSttus.json` |
| `Dividend` | 배당에 관한 사항 | `/api/alotMatter.json` |
| `TreasuryStock` | 자기주식 취득·처분 현황 | `/api/tesstkAcqsDspsSttus.json` |
| `TotalStock` | 주식의 총수 현황 | `/api/stockTotqySttus.json` |
| `MajorShareholders` | 최대주주 현황 | `/api/hyslrSttus.json` |
| `MajorShareholderChanges` | 최대주주 변동현황 | `/api/hyslrChgSttus.json` |
| `MinorityShareholders` | 소액주주 현황 | `/api/mrhlSttus.json` |

```go
// CapitalChangeItem 은 증자(감자) 현황 (irdsSttus) 한 건.
type CapitalChangeItem struct {
	RceptNo                 string `json:"rcept_no"`                    // 접수번호 (14자리)
	CorpCls                 string `json:"corp_cls"`                    // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                   // 고유번호 (8자리)
	CorpName                string `json:"corp_name"`                   // 법인명
	IsuDcrsDe               string `json:"isu_dcrs_de"`                 // 주식발행 감소일자
	IsuDcrsStle             string `json:"isu_dcrs_stle"`               // 발행 감소 형태
	IsuDcrsStockKnd         string `json:"isu_dcrs_stock_knd"`          // 발행 감소 주식 종류
	IsuDcrsQy               string `json:"isu_dcrs_qy"`                 // 발행 감소 수량
	IsuDcrsMstvdvFvalAmount string `json:"isu_dcrs_mstvdv_fval_amount"` // 발행 감소 주당 액면 가액
	IsuDcrsMstvdvAmount     string `json:"isu_dcrs_mstvdv_amount"`      // 발행 감소 주당 가액
	StlmDt                  string `json:"stlm_dt"`                     // 결산기준일
}

// DividendItem 은 배당에 관한 사항 (alotMatter) 한 건.
type DividendItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호 (14자리)
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호 (8자리)
	CorpName string `json:"corp_name"` // 법인명
	Se       string `json:"se"`        // 구분 (유상증자(주주배정), 전환권행사 등)
	StockKnd string `json:"stock_knd"` // 주식 종류 (보통주 등)
	Thstrm   string `json:"thstrm"`    // 당기
	Frmtrm   string `json:"frmtrm"`    // 전기
	Lwfr     string `json:"lwfr"`      // 전전기
	StlmDt   string `json:"stlm_dt"`   // 결산기준일 (YYYY-MM-DD)
}

// TreasuryStockItem 은 자기주식 취득 및 처분 현황 (tesstkAcqsDspsSttus) 한 건.
type TreasuryStockItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	AcqsMth1      string `json:"acqs_mth1"`       // 취득방법 대분류
	AcqsMth2      string `json:"acqs_mth2"`       // 취득방법 중분류
	AcqsMth3      string `json:"acqs_mth3"`       // 취득방법 소분류
	StockKnd      string `json:"stock_knd"`       // 주식 종류
	BsisQy        string `json:"bsis_qy"`         // 기초 수량
	ChangeQyAcqs  string `json:"change_qy_acqs"`  // 변동 수량 취득
	ChangeQyDsps  string `json:"change_qy_dsps"`  // 변동 수량 처분
	ChangeQyIncnr string `json:"change_qy_incnr"` // 변동 수량 소각
	TrmendQy      string `json:"trmend_qy"`       // 기말 수량
	Rm            string `json:"rm"`              // 비고
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// TotalStockItem 은 주식의 총수 현황 (stockTotqySttus) 한 건.
type TotalStockItem struct {
	RceptNo            string `json:"rcept_no"`               // 접수번호
	CorpCls            string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`              // 고유번호
	CorpName           string `json:"corp_name"`              // 회사명
	Se                 string `json:"se"`                     // 구분
	IsuStockTotqy      string `json:"isu_stock_totqy"`        // 발행할 주식의 총수
	NowToIsuStockTotqy string `json:"now_to_isu_stock_totqy"` // 현재까지 발행한 주식의 총수
	NowToDcrsStockTotqy string `json:"now_to_dcrs_stock_totqy"` // 현재까지 감소한 주식의 총수
	Redc               string `json:"redc"`                   // 감자
	ProfitIncnr        string `json:"profit_incnr"`           // 이익소각
	RdmstkRepy         string `json:"rdmstk_repy"`            // 상환주식의 상환
	Etc                string `json:"etc"`                    // 기타
	IstcTotqy          string `json:"istc_totqy"`             // 발행주식의 총수
	TesstkCo           string `json:"tesstk_co"`              // 자기주식수
	DistbStockCo       string `json:"distb_stock_co"`         // 유통주식수
	StlmDt             string `json:"stlm_dt"`                // 결산기준일
}

// MajorShareholderItem 은 최대주주 현황 (hyslrSttus) 한 건.
type MajorShareholderItem struct {
	RceptNo                 string `json:"rcept_no"`                   // 접수번호
	CorpCls                 string `json:"corp_cls"`                   // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                  // 고유번호
	CorpName                string `json:"corp_name"`                  // 법인명
	Nm                      string `json:"nm"`                         // 성명
	Relate                  string `json:"relate"`                     // 관계
	StockKnd                string `json:"stock_knd"`                  // 주식 종류
	BsisPosesnStockCo       string `json:"bsis_posesn_stock_co"`       // 기초 소유 주식 수
	BsisPosesnStockQotaRt   string `json:"bsis_posesn_stock_qota_rt"`  // 기초 소유 주식 지분율
	TrmendPosesnStockCo     string `json:"trmend_posesn_stock_co"`     // 기말 소유 주식 수
	TrmendPosesnStockQotaRt string `json:"trmend_posesn_stock_qota_rt"`// 기말 소유 주식 지분율
	Rm                      string `json:"rm"`                         // 비고
	StlmDt                  string `json:"stlm_dt"`                    // 결산기준일
}

// MajorShareholderChangeItem 은 최대주주 변동현황 (hyslrChgSttus) 한 건.
type MajorShareholderChangeItem struct {
	RceptNo       string `json:"rcept_no"`         // 접수번호
	CorpCls       string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`        // 고유번호
	CorpName      string `json:"corp_name"`        // 법인명
	ChangeOn      string `json:"change_on"`        // 변동 일
	MxmmShrholdrNm string `json:"mxmm_shrholdr_nm"` // 최대 주주 명
	PosesnStockCo string `json:"posesn_stock_co"`  // 소유 주식 수
	QotaRt        string `json:"qota_rt"`          // 지분율
	ChangeCause   string `json:"change_cause"`     // 변동 원인
	Rm            string `json:"rm"`               // 비고
	StlmDt        string `json:"stlm_dt"`          // 결산기준일
}

// MinorityShareholderItem 은 소액주주 현황 (mrhlSttus) 한 건.
type MinorityShareholderItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	Se            string `json:"se"`              // 구분
	ShrholdrCo    string `json:"shrholdr_co"`     // 주주수
	ShrholdrTotCo string `json:"shrholdr_tot_co"` // 전체 주주수
	ShrholdrRate  string `json:"shrholdr_rate"`   // 주주 비율
	HoldStockCo   string `json:"hold_stock_co"`   // 보유 주식수
	StockTotCo    string `json:"stock_tot_co"`    // 총발행 주식수
	HoldStockRate string `json:"hold_stock_rate"` // 보유 주식 비율
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}
```

각 메서드 (equity.go):
```go
func (c *Client) CapitalChange(ctx context.Context, p ReportParams) ([]CapitalChangeItem, error) {
	return getList[CapitalChangeItem](ctx, c.http, "/api/irdsSttus.json", p)
}
func (c *Client) Dividend(ctx context.Context, p ReportParams) ([]DividendItem, error) {
	return getList[DividendItem](ctx, c.http, "/api/alotMatter.json", p)
}
func (c *Client) TreasuryStock(ctx context.Context, p ReportParams) ([]TreasuryStockItem, error) {
	return getList[TreasuryStockItem](ctx, c.http, "/api/tesstkAcqsDspsSttus.json", p)
}
func (c *Client) TotalStock(ctx context.Context, p ReportParams) ([]TotalStockItem, error) {
	return getList[TotalStockItem](ctx, c.http, "/api/stockTotqySttus.json", p)
}
func (c *Client) MajorShareholders(ctx context.Context, p ReportParams) ([]MajorShareholderItem, error) {
	return getList[MajorShareholderItem](ctx, c.http, "/api/hyslrSttus.json", p)
}
func (c *Client) MajorShareholderChanges(ctx context.Context, p ReportParams) ([]MajorShareholderChangeItem, error) {
	return getList[MajorShareholderChangeItem](ctx, c.http, "/api/hyslrChgSttus.json", p)
}
func (c *Client) MinorityShareholders(ctx context.Context, p ReportParams) ([]MinorityShareholderItem, error) {
	return getList[MinorityShareholderItem](ctx, c.http, "/api/mrhlSttus.json", p)
}
```

## root 와이어링 (client.go 수정)

```go
type Client struct {
	http *httpclient.Client
	corp *corpcode.Cache

	Disclosure *disclosure.Client // DS001 공시정보
	Report     *report.Client     // DS002 정기보고서 주요정보
}
// NewClient 내부:
c.Report = report.New(hc)
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`(`errors.Is`), 그 외 status → `*opendart.APIError`.
`report` 패키지는 별도 에러 타입을 두지 않는다.

## 테스트 전략

- `report/equity_test.go`: `disclosure` 의 `newTestClient` 와 같은 방식의 path별 fixture 서빙 httptest 헬퍼 + 7개 실 응답 JSON fixture(`report/testdata/`).
- 7개 메서드 각각: fixture 디코딩 → item 필드 매핑 검증(대표 필드).
- 공통 헬퍼: `ReportParams.toMap`(3키 정확), `getList` 의 `013→ErrNoData` 경로(`{"status":"013",...}` 응답을 서빙해 `errors.Is(err, opendart.ErrNoData)` 검증).
- fixture 는 실 API 로 캡처해 임베드(연도/보고서 지정; 계획 작성 단계에서 캡처).
- `integration_test.go` 에 `Dividend` 대표 통합 케이스 추가(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(콤마 string 유지).
- README "커버리지" 에 DS002 7개 추가, `examples/report/` 예제 1개, 모든 파일 UTF-8.

## 비범위 (후속 plan)

- DS002 나머지 23개: 임원·직원·보수 그룹, 증권발행·미상환 그룹, 감사·자금·출자 그룹.
- DS003~DS006 카테고리.
- 숫자/날짜 coercion 헬퍼.
