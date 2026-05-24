# OpenDART DS005 주요사항보고서 주요정보 — 증자·감자 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 증자·감자 4개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 부실·법적 이벤트 7개는 main 에 머지됨(PR #10), `material` 패키지 + 공통 `MaterialParams`
확립. 이 spec 은 DS005 두 번째 그룹 — **증자·감자 4개**다(유상/무상/유무상 증자 결정, 감자 결정).
기존 `MaterialParams` + `httpclient.GetList[T]` 재사용, 새 추상화 없음. root 변경 없음(`client.Material`
기존). 신규 파일 `material/capital.go`.

## API 표면 (docs 기반 사실)

- 4개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 유상증자(piicDecsn), 무상증자(fricDecsn), 유무상증자(pifricDecsn), 감자(crDecsn).
- 공통 머리: `rcept_no`/`corp_cls`/`corp_code`/`corp_name`. 값은 문자열(금액/주식수 콤마, 비율 %, 빈 값 "-").

## 아키텍처

```
material/
  capital.go        # 4개 메서드 + item struct (신규)
  capital_test.go   # 4개 fixture 테스트 (신규)
  testdata/         # 4개 실 응답 fixture
README.md           # (수정) DS005 커버리지에 증자·감자
integration_test.go # (수정) PaidInCapitalIncrease 통합 케이스
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 4개 메서드 (material/capital.go)

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `PaidInCapitalIncrease` | 유상증자 결정 | `/api/piicDecsn.json` |
| `FreeCapitalIncrease` | 무상증자 결정 | `/api/fricDecsn.json` |
| `PaidFreeCapitalIncrease` | 유무상증자 결정 | `/api/pifricDecsn.json` |
| `CapitalReduction` | 감자 결정 | `/api/crDecsn.json` |

```go
// PaidInCapitalIncreaseItem 은 유상증자 결정 (piicDecsn) 한 건.
type PaidInCapitalIncreaseItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 회사명
	NstkOstkCnt    string `json:"nstk_ostk_cnt"`    // 신주의 종류와 수(보통주식)
	NstkEstkCnt    string `json:"nstk_estk_cnt"`    // 신주의 종류와 수(기타주식)
	FvPs           string `json:"fv_ps"`            // 1주당 액면가액 (원)
	BficTisstkOstk string `json:"bfic_tisstk_ostk"` // 증자전 발행주식총수(보통주식)
	BficTisstkEstk string `json:"bfic_tisstk_estk"` // 증자전 발행주식총수(기타주식)
	FdppFclt       string `json:"fdpp_fclt"`        // 자금조달목적(시설자금)
	FdppBsninh     string `json:"fdpp_bsninh"`      // 자금조달목적(영업양수자금)
	FdppOp         string `json:"fdpp_op"`          // 자금조달목적(운영자금)
	FdppDtrp       string `json:"fdpp_dtrp"`        // 자금조달목적(채무상환자금)
	FdppOcsa       string `json:"fdpp_ocsa"`        // 자금조달목적(타법인 증권 취득자금)
	FdppEtc        string `json:"fdpp_etc"`         // 자금조달목적(기타자금)
	IcMthn         string `json:"ic_mthn"`          // 증자방식
	SslAt          string `json:"ssl_at"`           // 공매도 해당여부
	SslBgd         string `json:"ssl_bgd"`          // 공매도 시작일
	SslEdd         string `json:"ssl_edd"`          // 공매도 종료일
}

// FreeCapitalIncreaseItem 은 무상증자 결정 (fricDecsn) 한 건.
type FreeCapitalIncreaseItem struct {
	RceptNo         string `json:"rcept_no"`            // 접수번호
	CorpCls         string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`           // 고유번호
	CorpName        string `json:"corp_name"`           // 회사명
	NstkOstkCnt     string `json:"nstk_ostk_cnt"`       // 신주의 종류와 수(보통주식)
	NstkEstkCnt     string `json:"nstk_estk_cnt"`       // 신주의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`               // 1주당 액면가액 (원)
	BficTisstkOstk  string `json:"bfic_tisstk_ostk"`    // 증자전 발행주식총수(보통주식)
	BficTisstkEstk  string `json:"bfic_tisstk_estk"`    // 증자전 발행주식총수(기타주식)
	NstkAsstd       string `json:"nstk_asstd"`          // 신주배정기준일
	NstkAscntPsOstk string `json:"nstk_ascnt_ps_ostk"`  // 1주당 신주배정 주식수(보통주식)
	NstkAscntPsEstk string `json:"nstk_ascnt_ps_estk"`  // 1주당 신주배정 주식수(기타주식)
	NstkDividrk     string `json:"nstk_dividrk"`        // 신주의 배당기산일
	NstkDlprd       string `json:"nstk_dlprd"`          // 신주권교부예정일
	NstkLstprd      string `json:"nstk_lstprd"`         // 신주의 상장 예정일
	Bddd            string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`           // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`           // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`           // 감사(감사위원) 참석여부
}

// PaidFreeCapitalIncreaseItem 은 유무상증자 결정 (pifricDecsn) 한 건. piic_* 유상 / fric_* 무상.
type PaidFreeCapitalIncreaseItem struct {
	RceptNo              string `json:"rcept_no"`                 // 접수번호
	CorpCls              string `json:"corp_cls"`                 // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`                // 고유번호
	CorpName             string `json:"corp_name"`                // 회사명
	PiicNstkOstkCnt      string `json:"piic_nstk_ostk_cnt"`       // 유상증자 신주수(보통주식)
	PiicNstkEstkCnt      string `json:"piic_nstk_estk_cnt"`       // 유상증자 신주수(기타주식)
	PiicFvPs             string `json:"piic_fv_ps"`               // 유상증자 1주당 액면가액
	PiicBficTisstkOstk   string `json:"piic_bfic_tisstk_ostk"`    // 유상증자 증자전 발행총수(보통주식)
	PiicBficTisstkEstk   string `json:"piic_bfic_tisstk_estk"`    // 유상증자 증자전 발행총수(기타주식)
	PiicFdppFclt         string `json:"piic_fdpp_fclt"`           // 유상증자 자금조달목적(시설자금)
	PiicFdppBsninh       string `json:"piic_fdpp_bsninh"`         // 유상증자 자금조달목적(영업양수자금)
	PiicFdppOp           string `json:"piic_fdpp_op"`             // 유상증자 자금조달목적(운영자금)
	PiicFdppDtrp         string `json:"piic_fdpp_dtrp"`           // 유상증자 자금조달목적(채무상환자금)
	PiicFdppOcsa         string `json:"piic_fdpp_ocsa"`           // 유상증자 자금조달목적(타법인 증권 취득자금)
	PiicFdppEtc          string `json:"piic_fdpp_etc"`            // 유상증자 자금조달목적(기타자금)
	PiicIcMthn           string `json:"piic_ic_mthn"`             // 유상증자 증자방식
	FricNstkOstkCnt      string `json:"fric_nstk_ostk_cnt"`       // 무상증자 신주수(보통주식)
	FricNstkEstkCnt      string `json:"fric_nstk_estk_cnt"`       // 무상증자 신주수(기타주식)
	FricFvPs             string `json:"fric_fv_ps"`               // 무상증자 1주당 액면가액
	FricBficTisstkOstk   string `json:"fric_bfic_tisstk_ostk"`    // 무상증자 증자전 발행총수(보통주식)
	FricBficTisstkEstk   string `json:"fric_bfic_tisstk_estk"`    // 무상증자 증자전 발행총수(기타주식)
	FricNstkAsstd        string `json:"fric_nstk_asstd"`          // 무상증자 신주배정기준일
	FricNstkAscntPsOstk  string `json:"fric_nstk_ascnt_ps_ostk"`  // 무상증자 1주당 신주배정수(보통주식)
	FricNstkAscntPsEstk  string `json:"fric_nstk_ascnt_ps_estk"`  // 무상증자 1주당 신주배정수(기타주식)
	FricNstkDividrk      string `json:"fric_nstk_dividrk"`        // 무상증자 신주 배당기산일
	FricNstkDlprd        string `json:"fric_nstk_dlprd"`          // 무상증자 신주권교부예정일
	FricNstkLstprd       string `json:"fric_nstk_lstprd"`         // 무상증자 신주 상장예정일
	FricBddd             string `json:"fric_bddd"`                // 무상증자 이사회결의일(결정일)
	FricOdAAtT           string `json:"fric_od_a_at_t"`           // 무상증자 사외이사 참석(참석)
	FricOdAAtB           string `json:"fric_od_a_at_b"`           // 무상증자 사외이사 참석(불참)
	FricAdtAAtn          string `json:"fric_adt_a_atn"`           // 무상증자 감사 참석여부
	SslAt                string `json:"ssl_at"`                   // 공매도 해당여부
	SslBgd               string `json:"ssl_bgd"`                  // 공매도 시작일
	SslEdd               string `json:"ssl_edd"`                  // 공매도 종료일
}

// CapitalReductionItem 은 감자 결정 (crDecsn) 한 건.
type CapitalReductionItem struct {
	RceptNo         string `json:"rcept_no"`           // 접수번호
	CorpCls         string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`          // 고유번호
	CorpName        string `json:"corp_name"`          // 회사명
	CrstkOstkCnt    string `json:"crstk_ostk_cnt"`     // 감자주식의 종류와 수(보통주식)
	CrstkEstkCnt    string `json:"crstk_estk_cnt"`     // 감자주식의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`              // 1주당 액면가액 (원)
	BfcrCpt         string `json:"bfcr_cpt"`           // 감자전 자본금 (원)
	AtcrCpt         string `json:"atcr_cpt"`           // 감자후 자본금 (원)
	BfcrTisstkOstk  string `json:"bfcr_tisstk_ostk"`   // 감자전 발행주식수(보통주식)
	AtcrTisstkOstk  string `json:"atcr_tisstk_ostk"`   // 감자후 발행주식수(보통주식)
	BfcrTisstkEstk  string `json:"bfcr_tisstk_estk"`   // 감자전 발행주식수(기타주식)
	AtcrTisstkEstk  string `json:"atcr_tisstk_estk"`   // 감자후 발행주식수(기타주식)
	CrRtOstk        string `json:"cr_rt_ostk"`         // 감자비율(보통주식 %)
	CrRtEstk        string `json:"cr_rt_estk"`         // 감자비율(기타주식 %)
	CrStd           string `json:"cr_std"`             // 감자기준일
	CrMth           string `json:"cr_mth"`             // 감자방법
	CrRs            string `json:"cr_rs"`              // 감자사유
	CrscGmtsckPrd   string `json:"crsc_gmtsck_prd"`    // 감자일정(주주총회 예정일)
	CrscTrnmsppd    string `json:"crsc_trnmsppd"`      // 감자일정(명의개서정지기간)
	CrscOsprpd      string `json:"crsc_osprpd"`        // 감자일정(구주권 제출기간)
	CrscTrspprpd    string `json:"crsc_trspprpd"`      // 감자일정(매매거래 정지예정기간)
	CrscOsprpdBgd   string `json:"crsc_osprpd_bgd"`    // 감자일정(구주권 제출기간 시작일)
	CrscOsprpdEdd   string `json:"crsc_osprpd_edd"`    // 감자일정(구주권 제출기간 종료일)
	CrscTrspprpdBgd string `json:"crsc_trspprpd_bgd"`  // 감자일정(매매거래 정지예정기간 시작일)
	CrscTrspprpdEdd string `json:"crsc_trspprpd_edd"`  // 감자일정(매매거래 정지예정기간 종료일)
	CrscNstkdlprd   string `json:"crsc_nstkdlprd"`     // 감자일정(신주권교부예정일)
	CrscNstklstprd  string `json:"crsc_nstklstprd"`    // 감자일정(신주상장예정일)
	CdobprpdBgd     string `json:"cdobprpd_bgd"`       // 채권자 이의제출기간(시작일)
	CdobprpdEdd     string `json:"cdobprpd_edd"`       // 채권자 이의제출기간(종료일)
	OsprNstkdlPl    string `json:"ospr_nstkdl_pl"`     // 구주권제출 및 신주권교부장소
	Bddd            string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`          // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`          // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`          // 감사(감사위원) 참석여부
	FtcSttAtn       string `json:"ftc_stt_atn"`        // 공정거래위원회 신고대상 여부
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/capital_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용.
- 4개 메서드 각각 실 응답 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 로 캡처해 임베드(증자/감자 사례 있는 종목/기간; 공시검색으로 corp_code 탐색).
- `integration_test.go` 에 `PaidInCapitalIncrease` 통합 케이스(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "증자(유상/무상/유무상)·감자 결정" 추가.

## 비범위 (후속 plan)

- DS005 나머지 그룹: 사채 발행 / 자기주식 / 영업·자산 양수도 / 합병·분할 / 해외상장 (~25개).
- DS006 증권신고서 주요정보. DS002 개인별 보수 Ver 2.0 2종.
