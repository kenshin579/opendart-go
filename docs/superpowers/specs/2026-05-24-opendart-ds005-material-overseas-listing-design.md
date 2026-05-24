# OpenDART DS005 주요사항보고서 주요정보 — 해외상장 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 해외상장 4개 API** (`material` 패키지 확장 — DS005 마지막 그룹)

## 배경 & 목표

DS005 부실·법적·증자·감자·사채 발행·자기주식·양수도·합병분할(PR #10~#16)이 main 머지됨. 이 spec 은
DS005 **마지막 그룹 — 해외상장 4개**다(해외 증권시장 주권등 상장 결정·상장·상장폐지 결정·상장폐지).
기존 `MaterialParams` + `httpclient.GetList[T]` 재사용, root 변경 없음. 신규 파일
`material/overseas_listing.go`. **이 그룹 완료 시 DS005 36/36 완성.**

## API 표면 (docs 기반 사실)

- 4개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 상장결정(ovLstDecsn 20), 상장(ovLst 10), 상장폐지결정(ovDlstDecsn 14), 상장폐지(ovDlst 10).
- 공통: 머리(rcept_no/corp_cls/corp_code/corp_name) + 상장거래소(lstex_nt). 결정 계열(상장결정·상장폐지결정)은 거버넌스(bddd/od_a_at_t/b/adt_a_atn) 포함.
- 값은 문자열(주식수 콤마, 빈 값 "-").

## 아키텍처

```
material/
  overseas_listing.go       # 4개 메서드 + item struct (신규)
  overseas_listing_test.go  # 4개 fixture 테스트 (신규)
  testdata/                 # 4개 fixture
README.md                   # (수정) DS005 커버리지에 해외상장 (DS005 완료 표기)
integration_test.go         # (수정) 통합 케이스 1~2개 (ErrNoData skip)
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 4개 메서드 (material/overseas_listing.go)

| 메서드 | 한글 | 엔드포인트 | 필드 |
|--------|------|-----------|------|
| `OverseasListingDecision` | 해외 증권시장 주권등 상장 결정 | `/api/ovLstDecsn.json` | 20 |
| `OverseasListing` | 해외 증권시장 주권등 상장 | `/api/ovLst.json` | 10 |
| `OverseasDelistingDecision` | 해외 증권시장 주권등 상장폐지 결정 | `/api/ovDlstDecsn.json` | 14 |
| `OverseasDelisting` | 해외 증권시장 주권등 상장폐지 | `/api/ovDlst.json` | 10 |

```go
// OverseasListingDecisionItem 은 해외 증권시장 주권등 상장 결정 (ovLstDecsn) 한 건.
type OverseasListingDecisionItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	CorpCls         string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`         // 고유번호
	CorpName        string `json:"corp_name"`         // 회사명
	LstprstkOstkCnt string `json:"lstprstk_ostk_cnt"` // 상장예정주식 종류ㆍ수(주)(보통주식)
	LstprstkEstkCnt string `json:"lstprstk_estk_cnt"` // 상장예정주식 종류ㆍ수(주)(기타주식)
	TisstkOstk      string `json:"tisstk_ostk"`       // 발행주식 총수(주)(보통주식)
	TisstkEstk      string `json:"tisstk_estk"`       // 발행주식 총수(주)(기타주식)
	PsmthNstkSl     string `json:"psmth_nstk_sl"`     // 공모방법(신주발행 (주))
	PsmthOstkSl     string `json:"psmth_ostk_sl"`     // 공모방법(구주매출 (주))
	Fdpp            string `json:"fdpp"`              // 자금조달(신주발행) 목적
	LststkOrlst     string `json:"lststk_orlst"`      // 상장증권(원주상장 (주))
	LststkDrlst     string `json:"lststk_drlst"`      // 상장증권(DR상장 (주))
	LstexNt         string `json:"lstex_nt"`          // 상장거래소(소재국가)
	Lstpp           string `json:"lstpp"`             // 해외상장목적
	Lstprd          string `json:"lstprd"`            // 상장예정일자
	Bddd            string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`         // 감사(감사위원) 참석여부
}

// OverseasListingItem 은 해외 증권시장 주권등 상장 (ovLst) 한 건.
type OverseasListingItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	LststkOstkCnt string `json:"lststk_ostk_cnt"` // 상장주식 종류 및 수(보통주식(주))
	LststkEstkCnt string `json:"lststk_estk_cnt"` // 상장주식 종류 및 수(기타주식(주))
	LstexNt       string `json:"lstex_nt"`        // 상장거래소(소재국가)
	StkCd         string `json:"stk_cd"`          // 종목 명 (code)
	Lstd          string `json:"lstd"`            // 상장일자
	Cfd           string `json:"cfd"`             // 확인일자
}

// OverseasDelistingDecisionItem 은 해외 증권시장 주권등 상장폐지 결정 (ovDlstDecsn) 한 건.
type OverseasDelistingDecisionItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	DlststkOstkCnt string `json:"dlststk_ostk_cnt"` // 상장폐지주식 종류ㆍ수(주)(보통주식)
	DlststkEstkCnt string `json:"dlststk_estk_cnt"` // 상장폐지주식 종류ㆍ수(주)(기타주식)
	LstexNt        string `json:"lstex_nt"`         // 상장거래소(소재국가)
	DlstrqPrd      string `json:"dlstrq_prd"`       // 폐지신청예정일자
	DlstPrd        string `json:"dlst_prd"`         // 폐지(예정)일자
	DlstRs         string `json:"dlst_rs"`          // 폐지사유
	Bddd           string `json:"bddd"`             // 이사회결의일(확인일)
	OdAAtT         string `json:"od_a_at_t"`        // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`        // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`        // 감사(감사위원) 참석여부
}

// OverseasDelistingItem 은 해외 증권시장 주권등 상장폐지 (ovDlst) 한 건.
type OverseasDelistingItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	LstexNt        string `json:"lstex_nt"`         // 상장거래소 및 소재국가
	DlststkOstkCnt string `json:"dlststk_ostk_cnt"` // 상장폐지주식의 종류(보통주식(주))
	DlststkEstkCnt string `json:"dlststk_estk_cnt"` // 상장폐지주식의 종류(기타주식(주))
	Tredd          string `json:"tredd"`            // 매매거래종료일
	DlstRs         string `json:"dlst_rs"`          // 폐지사유
	Cfd            string `json:"cfd"`              // 확인일자
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/overseas_listing_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용(route map 값은 bare 파일명).
- 4개 메서드 각각 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).
- `integration_test.go` 에 통합 케이스 1~2개(`//go:build integration`, ErrNoData skip 허용).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어. (상장결정 bddd=결정일, 상장폐지결정 bddd=확인일 — docs 그대로.)
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "해외 증권시장 주권등 상장 결정·상장·상장폐지 결정·상장폐지" 추가, 예정 줄에서 DS005 제거(DS005 36/36 완료).

## 비범위 (후속 plan)

- DS006 증권신고서(지분증권/채무증권/증권예탁증권/합병/분할/주식의포괄적교환·이전 6) — 신규 패키지(예: registration).
- DS002 개인별 보수 Ver 2.0 2종(데이터 확보 시).
