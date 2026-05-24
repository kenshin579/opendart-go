# OpenDART DS006 증권신고서 주요정보 — Sub-2 «신고» 설계

- 작성일: 2026-05-25
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS006 증권신고서 Sub-2 3개 API** (`registration` 패키지 확장 — DS006 마지막)

## 배경 & 목표

DS006 Sub-1 «증권» 3개(지분/채무/예탁)는 PR #18 main 머지됨 — `registration` 패키지 + `httpclient.GetGroups`
그룹형 디코더 확립. 이 spec 은 **Sub-2 «신고» 3개**(합병/분할/주식의포괄적교환·이전). 세 엔드포인트의
응답 그룹(일반사항/발행증권/당사회사에관한사항) 스키마가 **완전히 동일** → 공유 item 타입 3종으로 DRY.
기존 인프라 그대로, 신규 파일 `registration/restructuring.go`, root 변경 없음. **완료 시 DS006 6/6.**

## API 표면 (docs 기반 사실)

- 3개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `registration.Params`), 그룹형 응답(`group:[{title,list}]`).
- 합병(mgRs)/분할(dvRs)/주식의포괄적교환·이전(extrRs) — 모두 3그룹: 일반사항(17)·발행증권(9)·당사회사에관한사항(10), 스키마 동일.

## 아키텍처

```
registration/
  restructuring.go       # 3 공유 item + 3 wrapper + 3 메서드 (신규)
  restructuring_test.go  # 3개 fixture 테스트 (신규)
  testdata/              # 3개 fixture
README.md                # (수정) DS006 줄에 합병·분할·교환이전 추가 (DS006 완료)
integration_test.go      # (수정) 통합 케이스 1~2개 (ErrNoData skip)
```

각 메서드: `GetGroups` 호출 → title switch → `json.Unmarshal` → wrapper 반환 (Sub-1 패턴 동일).

## API 표면 (registration/restructuring.go)

| 메서드 | 한글 | 엔드포인트 | 반환 |
|--------|------|-----------|------|
| `Merger` | 합병 | `/api/mgRs.json` | `*MergerRegistration` |
| `Division` | 분할 | `/api/dvRs.json` | `*DivisionRegistration` |
| `StockExchangeTransfer` | 주식의포괄적교환·이전 | `/api/extrRs.json` | `*StockExchangeTransferRegistration` |

**공유 item 타입 3종**(세 엔드포인트 동일 그룹):
```go
// RestructuringGeneralItem 은 증권신고서(합병/분할/교환이전) 일반사항 그룹 항목.
type RestructuringGeneralItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	Stn           string `json:"stn"`             // 형태
	Bddd          string `json:"bddd"`            // 이사회 결의일
	Ctrd          string `json:"ctrd"`            // 계약일
	GmtsckShddstd string `json:"gmtsck_shddstd"`  // 주주총회를 위한 주주확정일
	ApGmtsck      string `json:"ap_gmtsck"`       // 승인을 위한 주주총회일
	AprskhPdBgd   string `json:"aprskh_pd_bgd"`   // 주식매수청구권 행사 기간 및 가격(시작일)
	AprskhPdEdd   string `json:"aprskh_pd_edd"`   // 주식매수청구권 행사 기간 및 가격(종료일)
	AprskhPrc     string `json:"aprskh_prc"`      // 주식매수청구권 행사 기간 및 가격(주식매수청구가격-회사제시)
	MgdtEtc       string `json:"mgdt_etc"`        // 합병기일등
	RtVl          string `json:"rt_vl"`           // 비율 또는 가액
	ExevlInt      string `json:"exevl_int"`       // 외부평가기관
	GrtmnEtc      string `json:"grtmn_etc"`       // 지급 교부금 등
	RptRcpn       string `json:"rpt_rcpn"`        // 주요사항보고서(접수번호)
}

// RestructuringIssuedSecurityItem 은 증권신고서(합병/분할/교환이전) 발행증권 그룹 항목.
type RestructuringIssuedSecurityItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Kndn     string `json:"kndn"`      // 종류
	Cnt      string `json:"cnt"`       // 수량
	Fv       string `json:"fv"`        // 액면가액
	Slprc    string `json:"slprc"`     // 모집(매출)가액
	Slta     string `json:"slta"`      // 모집(매출)총액
}

// RestructuringPartyCompanyItem 은 증권신고서(합병/분할/교환이전) 당사회사에관한사항 그룹 항목.
type RestructuringPartyCompanyItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Cmpnm     string `json:"cmpnm"`      // 회사명
	Sen       string `json:"sen"`        // 구분
	Tast      string `json:"tast"`       // 총자산
	Cpt       string `json:"cpt"`        // 자본금
	IsstkKnd  string `json:"isstk_knd"`  // 발행주식수(주식의종류)
	IsstkCnt  string `json:"isstk_cnt"`  // 발행주식수(주식수)
}
```

**Wrapper 3종 + 메서드 3개**(구조 동일, 엔드포인트별 분리 — Sub-1 패턴):
```go
// MergerRegistration 은 합병 증권신고서(mgRs)의 그룹별 항목.
type MergerRegistration struct {
	General          []RestructuringGeneralItem         // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem  // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem    // 당사회사에관한사항
}

// Merger 는 합병 증권신고서(DS006)를 조회한다.
func (c *Client) Merger(ctx context.Context, p Params) (*MergerRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/mgRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &MergerRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// DivisionRegistration 은 분할 증권신고서(dvRs)의 그룹별 항목.
type DivisionRegistration struct {
	General          []RestructuringGeneralItem
	IssuedSecurities []RestructuringIssuedSecurityItem
	PartyCompanies   []RestructuringPartyCompanyItem
}

// Division 은 분할 증권신고서(DS006)를 조회한다.
func (c *Client) Division(ctx context.Context, p Params) (*DivisionRegistration, error) {
	// (Merger 와 동일 패턴, 엔드포인트 "/api/dvRs.json", out *DivisionRegistration)
	...
}

// StockExchangeTransferRegistration 은 주식의포괄적교환·이전 증권신고서(extrRs)의 그룹별 항목.
type StockExchangeTransferRegistration struct {
	General          []RestructuringGeneralItem
	IssuedSecurities []RestructuringIssuedSecurityItem
	PartyCompanies   []RestructuringPartyCompanyItem
}

// StockExchangeTransfer 는 주식의포괄적교환·이전 증권신고서(DS006)를 조회한다.
func (c *Client) StockExchangeTransfer(ctx context.Context, p Params) (*StockExchangeTransferRegistration, error) {
	// (Merger 와 동일 패턴, 엔드포인트 "/api/extrRs.json", out *StockExchangeTransferRegistration)
	...
}
```
(plan 에서 Division/StockExchangeTransfer 메서드 전체 코드를 Merger 와 동일 패턴으로 명시한다. 위 `...` 는 spec 가독성 위한 축약이며 plan/구현은 완전한 코드.)

## 에러 처리

`GetGroups` 가 status 처리: 013→`opendart.ErrNoData`, 그 외→`*opendart.APIError`. 빈/누락 그룹은 해당 슬라이스 nil/empty.

## 테스트 전략

- `registration/restructuring_test.go`: 기존 `registration/client_test.go` 의 `newTestClient(t, fixture)` 재사용.
- 3개 메서드 각각 fixture 디코딩 → 각 그룹 슬라이스 len + 대표 필드 검증.
- fixture 는 실 API 캡처 권장(합병/분할 사례 종목 탐색; 불가 시 docs 스키마 일치 샘플, 그룹 배열 형태 유지).
- `integration_test.go` 에 `Merger`·`Division` 통합 케이스(`//go:build integration`, ErrNoData skip).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS006 줄에 "합병 · 분할 · 주식의포괄적교환·이전" 추가, 예정 줄에서 DS006 제거(DS006 6/6 완료).

## 비범위 (후속)

- DS002 개인별 보수 Ver 2.0 2종(`indvdlByPayV2`/`hmvAuditIndvdlBySttusV2`, 데이터 확보 시). → 이후 OpenDART 전체 API 커버 완료.
