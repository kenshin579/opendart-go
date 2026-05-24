# OpenDART DS005 주요사항보고서 주요정보 — 자기주식 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 자기주식 4개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 부실·법적(PR #10)·증자·감자(PR #11)·사채 발행(PR #12)이 main 에 머지됨, `material`
패키지 + 공통 `MaterialParams` 확립. 이 spec 은 DS005 네 번째 그룹 — **자기주식 4개**다
(취득/처분/신탁계약 체결/신탁계약 해지 결정). 기존 `MaterialParams` + `httpclient.GetList[T]`
재사용, 새 추상화·root 변경 없음(`client.Material` 기존). 신규 파일 `material/treasury.go`.

## API 표면 (docs 기반 사실)

- 4개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 취득(tsstkAqDecsn 29필드), 처분(tsstkDpDecsn 32), 신탁체결(tsstkAqTrctrCnsDecsn 23), 신탁해지(tsstkAqTrctrCcDecsn 24).
- 공통 블록: 머리(rcept_no/corp_cls/corp_code/corp_name) + 직전 자기주식 보유현황 8(aq_wtn_div_* 배당가능범위내, eaq_* 기타취득; 각 보통/기타 주식수+비율) + 거버넌스(od_a_at_t/b, adt_a_atn). 타입별 고유 필드는 각 struct.
- 값은 문자열(주식수/금액 콤마, 비율 %, 빈 값 "-").

## 아키텍처

```
material/
  treasury.go       # 4개 메서드 + item struct (신규)
  treasury_test.go  # 4개 fixture 테스트 (신규)
  testdata/         # 4개 실 응답 fixture
README.md           # (수정) DS005 커버리지에 자기주식
integration_test.go # (수정) TreasuryStockAcquisition 통합 케이스
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 4개 메서드 (material/treasury.go)

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `TreasuryStockAcquisition` | 자기주식 취득 결정 | `/api/tsstkAqDecsn.json` |
| `TreasuryStockDisposal` | 자기주식 처분 결정 | `/api/tsstkDpDecsn.json` |
| `TreasuryStockTrustContract` | 자기주식취득 신탁계약 체결 결정 | `/api/tsstkAqTrctrCnsDecsn.json` |
| `TreasuryStockTrustCancellation` | 자기주식취득 신탁계약 해지 결정 | `/api/tsstkAqTrctrCcDecsn.json` |

```go
// TreasuryStockAcquisitionItem 은 자기주식 취득 결정 (tsstkAqDecsn) 한 건.
type TreasuryStockAcquisitionItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	AqplnStkOstk   string `json:"aqpln_stk_ostk"`   // 취득예정주식(주)(보통주식)
	AqplnStkEstk   string `json:"aqpln_stk_estk"`   // 취득예정주식(주)(기타주식)
	AqplnPrcOstk   string `json:"aqpln_prc_ostk"`   // 취득예정금액(원)(보통주식)
	AqplnPrcEstk   string `json:"aqpln_prc_estk"`   // 취득예정금액(원)(기타주식)
	AqexpdBgd      string `json:"aqexpd_bgd"`       // 취득예상기간(시작일)
	AqexpdEdd      string `json:"aqexpd_edd"`       // 취득예상기간(종료일)
	HdexpdBgd      string `json:"hdexpd_bgd"`       // 보유예상기간(시작일)
	HdexpdEdd      string `json:"hdexpd_edd"`       // 보유예상기간(종료일)
	AqPp           string `json:"aq_pp"`            // 취득목적
	AqMth          string `json:"aq_mth"`           // 취득방법
	CsIvBk         string `json:"cs_iv_bk"`         // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`  // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`  // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`         // 취득 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`      // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`         // 취득 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`      // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	AqDd           string `json:"aq_dd"`            // 취득결정일
	OdAAtT         string `json:"od_a_at_t"`        // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`        // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`        // 감사(사외이사가 아닌 감사위원) 참석여부
	D1ProdlmOstk   string `json:"d1_prodlm_ostk"`   // 1일 매수 주문수량 한도(보통주식)
	D1ProdlmEstk   string `json:"d1_prodlm_estk"`   // 1일 매수 주문수량 한도(기타주식)
}

// TreasuryStockDisposalItem 은 자기주식 처분 결정 (tsstkDpDecsn) 한 건.
type TreasuryStockDisposalItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	DpplnStkOstk   string `json:"dppln_stk_ostk"`   // 처분예정주식(주)(보통주식)
	DpplnStkEstk   string `json:"dppln_stk_estk"`   // 처분예정주식(주)(기타주식)
	DpstkPrcOstk   string `json:"dpstk_prc_ostk"`   // 처분 대상 주식가격(원)(보통주식)
	DpstkPrcEstk   string `json:"dpstk_prc_estk"`   // 처분 대상 주식가격(원)(기타주식)
	DpplnPrcOstk   string `json:"dppln_prc_ostk"`   // 처분예정금액(원)(보통주식)
	DpplnPrcEstk   string `json:"dppln_prc_estk"`   // 처분예정금액(원)(기타주식)
	DpprpdBgd      string `json:"dpprpd_bgd"`       // 처분예정기간(시작일)
	DpprpdEdd      string `json:"dpprpd_edd"`       // 처분예정기간(종료일)
	DpPp           string `json:"dp_pp"`            // 처분목적
	DpMMkt         string `json:"dp_m_mkt"`         // 처분방법(시장을 통한 매도(주))
	DpMOvtm        string `json:"dp_m_ovtm"`        // 처분방법(시간외대량매매(주))
	DpMOtc         string `json:"dp_m_otc"`         // 처분방법(장외처분(주))
	DpMEtc         string `json:"dp_m_etc"`         // 처분방법(기타(주))
	CsIvBk         string `json:"cs_iv_bk"`         // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`  // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`  // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`         // 처분 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`      // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`         // 처분 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`      // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	DpDd           string `json:"dp_dd"`            // 처분결정일
	OdAAtT         string `json:"od_a_at_t"`        // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`        // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`        // 감사(사외이사가 아닌 감사위원) 참석여부
	D1SlodlmOstk   string `json:"d1_slodlm_ostk"`   // 1일 매도 주문수량 한도(보통주식)
	D1SlodlmEstk   string `json:"d1_slodlm_estk"`   // 1일 매도 주문수량 한도(기타주식)
}

// TreasuryStockTrustContractItem 은 자기주식취득 신탁계약 체결 결정 (tsstkAqTrctrCnsDecsn) 한 건.
type TreasuryStockTrustContractItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	CtrPrc         string `json:"ctr_prc"`          // 계약금액(원)
	CtrPdBgd       string `json:"ctr_pd_bgd"`       // 계약기간(시작일)
	CtrPdEdd       string `json:"ctr_pd_edd"`       // 계약기간(종료일)
	CtrPp          string `json:"ctr_pp"`           // 계약목적
	CtrCnsInt      string `json:"ctr_cns_int"`      // 계약체결기관
	CtrCnsPrd      string `json:"ctr_cns_prd"`      // 계약체결 예정일자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`  // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`  // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`         // 계약 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`      // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`         // 계약 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`      // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`             // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`        // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`        // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`        // 감사(사외이사가 아닌 감사위원) 참석여부
	CsIvBk         string `json:"cs_iv_bk"`         // 위탁투자중개업자
}

// TreasuryStockTrustCancellationItem 은 자기주식취득 신탁계약 해지 결정 (tsstkAqTrctrCcDecsn) 한 건.
type TreasuryStockTrustCancellationItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	CtrPrcBfcc     string `json:"ctr_prc_bfcc"`     // 계약금액(원)(해지 전)
	CtrPrcAtcc     string `json:"ctr_prc_atcc"`     // 계약금액(원)(해지 후)
	CtrPdBfccBgd   string `json:"ctr_pd_bfcc_bgd"`  // 해지 전 계약기간(시작일)
	CtrPdBfccEdd   string `json:"ctr_pd_bfcc_edd"`  // 해지 전 계약기간(종료일)
	CcPp           string `json:"cc_pp"`            // 해지목적
	CcInt          string `json:"cc_int"`           // 해지기관
	CcPrd          string `json:"cc_prd"`           // 해지예정일자
	TpRmAtcc       string `json:"tp_rm_atcc"`       // 해지후 신탁재산의 반환방법
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`  // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`  // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`         // 해지 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`      // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`         // 해지 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`      // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`             // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`        // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`        // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`        // 감사(사외이사가 아닌 감사위원) 참석여부
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/treasury_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용
  (route map 값은 bare 파일명 — `newTestClient` 가 testdata/ 를 붙임).
- 4개 메서드 각각 실 응답 fixture 디코딩 → 대표 필드 검증(머리 + 타입 고유 필드 + 보유현황 1 + 거버넌스).
- fixture 는 실 API 로 캡처해 임베드(자기주식 취득/처분은 흔하므로 실 데이터 확보 가능성 높음; 안되면 docs 스키마 일치 샘플).
- `integration_test.go` 에 `TreasuryStockAcquisition` 통합 케이스(`//go:build integration`, ErrNoData skip 허용).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "자기주식(취득/처분/신탁계약 체결·해지)" 추가.

## 비범위 (후속 plan)

- DS005 나머지 그룹: 영업·자산 양수도 / 합병·분할 / 해외상장 (~13개).
- DS006 증권신고서 주요정보. DS002 개인별 보수 Ver 2.0 2종.
