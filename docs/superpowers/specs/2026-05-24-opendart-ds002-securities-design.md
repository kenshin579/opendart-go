# OpenDART DS002 정기보고서 주요정보 — 증권 발행·미상환 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS002 증권 발행·미상환 6개 API** (`report` 패키지 확장)

## 배경 & 목표

DS002 공통 추상화(`report` 패키지: `ReportParams`/`ReportCode`/제네릭 `getList`)와 지분·주식·배당
7개는 main 에 머지됨(PR #3). 이 spec 은 DS002 두 번째 그룹인 **증권 발행·미상환 6개**다. 6개 모두
표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답을 가지므로 **새 추상화 없이** 기존
`getList[T]` 패턴을 그대로 재사용한다 — 각 엔드포인트 = item struct + 한 줄 메서드.

## API 표면 (docs 기반 사실)

- 요청: 전 API 공통 `crtfc_key`(자동 주입) + `corp_code` + `bsns_year` + `reprt_code` (추가 파라미터 없음).
- 응답: 공통 envelope(`status`/`message`) + `list[]`. 숫자는 콤마 문자열, 빈 값 "-".
- 미상환 잔액 5개는 `remndr_exprtn1/2`(잔여만기 구분) + 만기 버킷(증권 종류별로 상이) + `sm`(합계).

## 아키텍처

`report` 패키지에 새 파일 `securities.go` 추가 (기존 `equity.go` 는 그대로). `client.Report` 는 PR #3
에서 이미 root 에 와이어링됨 → **root 변경 불필요**. 메서드는 기존 `getList[T]`/`ReportParams` 재사용.

```
report/
  securities.go       # 6개 메서드 + item struct
  securities_test.go  # 6개 fixture 테스트 (newTestClient 재사용)
  testdata/           # 6개 실 응답 JSON fixture 추가
README.md             # (수정) DS002 커버리지에 증권 발행·미상환 추가
integration_test.go   # (수정) DebtSecuritiesIssuance 통합 케이스 추가
```

## 6개 엔드포인트 (report/securities.go)

각 메서드: `func (c *Client) X(ctx, p ReportParams) ([]XItem, error) { return getList[XItem](ctx, c.http, "<path>", p) }`.

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `DebtSecuritiesIssuance` | 채무증권 발행실적 | `/api/detScritsIsuAcmslt.json` |
| `CorporateBondBalance` | 회사채 미상환 잔액 | `/api/cprndNrdmpBlce.json` |
| `CommercialPaperBalance` | 기업어음증권 미상환 잔액 | `/api/entrprsBilScritsNrdmpBlce.json` |
| `ShortTermBondBalance` | 단기사채 미상환 잔액 | `/api/srtpdPsndbtNrdmpBlce.json` |
| `HybridSecuritiesBalance` | 신종자본증권 미상환 잔액 | `/api/newCaplScritsNrdmpBlce.json` |
| `ContingentCapitalBalance` | 조건부 자본증권 미상환 잔액 | `/api/cndlCaplScritsNrdmpBlce.json` |

```go
// DebtSecuritiesIssuanceItem 은 채무증권 발행실적 (detScritsIsuAcmslt) 한 건.
type DebtSecuritiesIssuanceItem struct {
	RceptNo       string `json:"rcept_no"`       // 접수번호
	CorpCls       string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`      // 고유번호
	CorpName      string `json:"corp_name"`      // 회사명
	IsuCmpny      string `json:"isu_cmpny"`      // 발행회사
	ScritsKndNm   string `json:"scrits_knd_nm"`  // 증권종류
	IsuMthNm      string `json:"isu_mth_nm"`     // 발행방법
	IsuDe         string `json:"isu_de"`         // 발행일자
	FacvaluTotamt string `json:"facvalu_totamt"` // 권면(전자등록)총액
	Intrt         string `json:"intrt"`          // 이자율
	EvlGradInstt  string `json:"evl_grad_instt"` // 평가등급(평가기관)
	Mtd           string `json:"mtd"`            // 만기일
	RepyAt        string `json:"repy_at"`        // 상환여부
	MngtCmpny     string `json:"mngt_cmpny"`     // 주관회사
	StlmDt        string `json:"stlm_dt"`        // 결산기준일
}

// CorporateBondBalanceItem 은 회사채 미상환 잔액 (cprndNrdmpBlce) 한 건.
type CorporateBondBalanceItem struct {
	RceptNo            string `json:"rcept_no"`              // 접수번호
	CorpCls            string `json:"corp_cls"`              // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`             // 고유번호
	CorpName           string `json:"corp_name"`             // 회사명
	RemndrExprtn1      string `json:"remndr_exprtn1"`        // 잔여만기 (구분1)
	RemndrExprtn2      string `json:"remndr_exprtn2"`        // 잔여만기 (구분2)
	Yy1Below           string `json:"yy1_below"`             // 1년 이하
	Yy1ExcessYy2Below  string `json:"yy1_excess_yy2_below"`  // 1년초과 2년이하
	Yy2ExcessYy3Below  string `json:"yy2_excess_yy3_below"`  // 2년초과 3년이하
	Yy3ExcessYy4Below  string `json:"yy3_excess_yy4_below"`  // 3년초과 4년이하
	Yy4ExcessYy5Below  string `json:"yy4_excess_yy5_below"`  // 4년초과 5년이하
	Yy5ExcessYy10Below string `json:"yy5_excess_yy10_below"` // 5년초과 10년이하
	Yy10Excess         string `json:"yy10_excess"`           // 10년초과
	Sm                 string `json:"sm"`                    // 합계
	StlmDt             string `json:"stlm_dt"`               // 결산기준일
}

// CommercialPaperBalanceItem 은 기업어음증권 미상환 잔액 (entrprsBilScritsNrdmpBlce) 한 건.
type CommercialPaperBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Yy1ExcessYy2Below    string `json:"yy1_excess_yy2_below"`    // 1년초과 2년이하
	Yy2ExcessYy3Below    string `json:"yy2_excess_yy3_below"`    // 2년초과 3년이하
	Yy3Excess            string `json:"yy3_excess"`              // 3년 초과
	Sm                   string `json:"sm"`                      // 합계
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// ShortTermBondBalanceItem 은 단기사채 미상환 잔액 (srtpdPsndbtNrdmpBlce) 한 건.
type ShortTermBondBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Sm                   string `json:"sm"`                      // 합계
	IsuLmt               string `json:"isu_lmt"`                 // 발행 한도
	RemndrLmt            string `json:"remndr_lmt"`              // 잔여 한도
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// HybridSecuritiesBalanceItem 은 신종자본증권 미상환 잔액 (newCaplScritsNrdmpBlce) 한 건.
type HybridSecuritiesBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy5Below   string `json:"yy1_excess_yy5_below"`   // 1년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy15Below string `json:"yy10_excess_yy15_below"` // 10년초과 15년이하
	Yy15ExcessYy20Below string `json:"yy15_excess_yy20_below"` // 15년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// ContingentCapitalBalanceItem 은 조건부 자본증권 미상환 잔액 (cndlCaplScritsNrdmpBlce) 한 건.
type ContingentCapitalBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy2Below   string `json:"yy1_excess_yy2_below"`   // 1년초과 2년이하
	Yy2ExcessYy3Below   string `json:"yy2_excess_yy3_below"`   // 2년초과 3년이하
	Yy3ExcessYy4Below   string `json:"yy3_excess_yy4_below"`   // 3년초과 4년이하
	Yy4ExcessYy5Below   string `json:"yy4_excess_yy5_below"`   // 4년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy20Below string `json:"yy10_excess_yy20_below"` // 10년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}
```

각 메서드 (securities.go):
```go
func (c *Client) DebtSecuritiesIssuance(ctx context.Context, p ReportParams) ([]DebtSecuritiesIssuanceItem, error) {
	return getList[DebtSecuritiesIssuanceItem](ctx, c.http, "/api/detScritsIsuAcmslt.json", p)
}
func (c *Client) CorporateBondBalance(ctx context.Context, p ReportParams) ([]CorporateBondBalanceItem, error) {
	return getList[CorporateBondBalanceItem](ctx, c.http, "/api/cprndNrdmpBlce.json", p)
}
func (c *Client) CommercialPaperBalance(ctx context.Context, p ReportParams) ([]CommercialPaperBalanceItem, error) {
	return getList[CommercialPaperBalanceItem](ctx, c.http, "/api/entrprsBilScritsNrdmpBlce.json", p)
}
func (c *Client) ShortTermBondBalance(ctx context.Context, p ReportParams) ([]ShortTermBondBalanceItem, error) {
	return getList[ShortTermBondBalanceItem](ctx, c.http, "/api/srtpdPsndbtNrdmpBlce.json", p)
}
func (c *Client) HybridSecuritiesBalance(ctx context.Context, p ReportParams) ([]HybridSecuritiesBalanceItem, error) {
	return getList[HybridSecuritiesBalanceItem](ctx, c.http, "/api/newCaplScritsNrdmpBlce.json", p)
}
func (c *Client) ContingentCapitalBalance(ctx context.Context, p ReportParams) ([]ContingentCapitalBalanceItem, error) {
	return getList[ContingentCapitalBalanceItem](ctx, c.http, "/api/cndlCaplScritsNrdmpBlce.json", p)
}
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `report/securities_test.go`: 기존 `report/client_test.go` 의 `newTestClient` 재사용.
- 6개 메서드 각각: 실 응답 JSON fixture 디코딩 → 대표 필드 매핑 검증.
- fixture 는 실 API 로 캡처해 임베드(계획 작성 단계). 데이터가 없는 종목은 데이터 있는 회사/연도로 캡처.
- `integration_test.go` 에 `DebtSecuritiesIssuance` 통합 케이스 추가(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(콤마 string 유지), UTF-8.
- README "커버리지" DS002 줄에 "증권 발행·미상환" 추가.

## 비범위 (후속 plan)

- DS002 나머지 2개 그룹: 임원·보수(11), 감사·자금·출자(6).
- DS003~DS006 카테고리.
- 신규 예제(기존 `examples/report` 로 충분).
