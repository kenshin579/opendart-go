# OpenDART DS005 주요사항보고서 주요정보 — 사채 발행 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 사채 발행 4개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 부실·법적(PR #10)·증자·감자(PR #11)는 main 에 머지됨, `material` 패키지 + 공통
`MaterialParams` 확립. 이 spec 은 DS005 세 번째 그룹 — **사채 발행 4개**다(전환사채/신주인수권부사채/
교환사채/상각형 조건부자본증권 발행결정). 기존 `MaterialParams` + `httpclient.GetList[T]` 재사용,
새 추상화·root 변경 없음(`client.Material` 기존). 신규 파일 `material/bond.go`.

## API 표면 (docs 기반 사실)

- 4개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 전환사채(cvbdIsDecsn 46필드), 신주인수권부사채(bdwtIsDecsn 49), 교환사채(exbdIsDecsn 42),
  상각형 조건부자본증권(wdCocobdIsDecsn 34).
- 공통 백본: 머리(rcept_no/corp_cls/corp_code/corp_name) + 사채 종류/총액(bd_tm/bd_knd/bd_fta) +
  해외발행(ovis_*) + 자금조달목적(fdpp_*) + 이율(bd_intr_ex/sf) + 만기/발행방법(bd_mtd/bdis_mthn) +
  청약/납입/주관/보증(sbd/pymd/rpmcmp/grint) + 이사회/사외이사/감사(bddd/od_a_at_*/adt_a_atn) +
  증권신고서(rs_sm_atn/ex_sm_r) + 해외대차(ovis_ltdtl) + 공정위(ftc_stt_atn). 타입별 고유 필드는 각 struct.
- 값은 문자열(금액/주식수 콤마, 비율 %, 빈 값 "-").

## 아키텍처

```
material/
  bond.go        # 4개 메서드 + item struct (신규)
  bond_test.go   # 4개 fixture 테스트 (신규)
  testdata/      # 4개 실 응답 fixture
README.md        # (수정) DS005 커버리지에 사채발행
integration_test.go # (수정) ConvertibleBondIssuance 통합 케이스
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 4개 메서드 (material/bond.go)

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `ConvertibleBondIssuance` | 전환사채권 발행결정 | `/api/cvbdIsDecsn.json` |
| `BondWithWarrantIssuance` | 신주인수권부사채권 발행결정 | `/api/bdwtIsDecsn.json` |
| `ExchangeableBondIssuance` | 교환사채권 발행결정 | `/api/exbdIsDecsn.json` |
| `ContingentConvertibleBondIssuance` | 상각형 조건부자본증권 발행결정 | `/api/wdCocobdIsDecsn.json` |

```go
// ConvertibleBondItem 은 전환사채권 발행결정 (cvbdIsDecsn) 한 건.
type ConvertibleBondItem struct {
	RceptNo                   string `json:"rcept_no"`                       // 접수번호
	CorpCls                   string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                      // 고유번호
	CorpName                  string `json:"corp_name"`                      // 회사명
	BdTm                      string `json:"bd_tm"`                          // 사채의 종류(회차)
	BdKnd                     string `json:"bd_knd"`                         // 사채의 종류(종류)
	BdFta                     string `json:"bd_fta"`                         // 사채의 권면(전자등록)총액 (원)
	AtcscRmislmt              string `json:"atcsc_rmislmt"`                  // 정관상 잔여 발행한도 (원)
	OvisFta                   string `json:"ovis_fta"`                       // 해외발행(권면(전자등록)총액)
	OvisFtaCrn                string `json:"ovis_fta_crn"`                   // 해외발행(권면총액 통화단위)
	OvisSter                  string `json:"ovis_ster"`                      // 해외발행(기준환율등)
	OvisIsar                  string `json:"ovis_isar"`                      // 해외발행(발행지역)
	OvisMktnm                 string `json:"ovis_mktnm"`                     // 해외발행(해외상장시 시장 명칭)
	FdppFclt                  string `json:"fdpp_fclt"`                      // 자금조달목적(시설자금)
	FdppBsninh                string `json:"fdpp_bsninh"`                    // 자금조달목적(영업양수자금)
	FdppOp                    string `json:"fdpp_op"`                        // 자금조달목적(운영자금)
	FdppDtrp                  string `json:"fdpp_dtrp"`                      // 자금조달목적(채무상환자금)
	FdppOcsa                  string `json:"fdpp_ocsa"`                      // 자금조달목적(타법인 증권 취득자금)
	FdppEtc                   string `json:"fdpp_etc"`                       // 자금조달목적(기타자금)
	BdIntrEx                  string `json:"bd_intr_ex"`                     // 사채의 이율(표면이자율 %)
	BdIntrSf                  string `json:"bd_intr_sf"`                     // 사채의 이율(만기이자율 %)
	BdMtd                     string `json:"bd_mtd"`                         // 사채만기일
	BdisMthn                  string `json:"bdis_mthn"`                      // 사채발행방법
	CvRt                      string `json:"cv_rt"`                          // 전환비율 (%)
	CvPrc                     string `json:"cv_prc"`                         // 전환가액 (원/주)
	CvisstkKnd                string `json:"cvisstk_knd"`                    // 전환에 따라 발행할 주식(종류)
	CvisstkCnt                string `json:"cvisstk_cnt"`                    // 전환에 따라 발행할 주식(주식수)
	CvisstkTisstkVs           string `json:"cvisstk_tisstk_vs"`              // 전환에 따라 발행할 주식(주식총수 대비 %)
	CvrqpdBgd                 string `json:"cvrqpd_bgd"`                     // 전환청구기간(시작일)
	CvrqpdEdd                 string `json:"cvrqpd_edd"`                     // 전환청구기간(종료일)
	ActMktprcflCvprcLwtrsprc  string `json:"act_mktprcfl_cvprc_lwtrsprc"`    // 시가하락 전환가액 조정(최저 조정가액 원)
	ActMktprcflCvprcLwtrsprcBs string `json:"act_mktprcfl_cvprc_lwtrsprc_bs"` // 시가하락 전환가액 조정(최저 조정가액 근거)
	RmislmtLt70p              string `json:"rmislmt_lt70p"`                  // 시가하락 조정(전환가 70% 미만 조정가능 잔여한도 원)
	Abmg                      string `json:"abmg"`                           // 합병 관련 사항
	Sbd                       string `json:"sbd"`                            // 청약일
	Pymd                      string `json:"pymd"`                           // 납입일
	Rpmcmp                    string `json:"rpmcmp"`                         // 대표주관회사
	Grint                     string `json:"grint"`                          // 보증기관
	Bddd                      string `json:"bddd"`                           // 이사회결의일(결정일)
	OdAAtT                    string `json:"od_a_at_t"`                      // 사외이사 참석여부(참석)
	OdAAtB                    string `json:"od_a_at_b"`                      // 사외이사 참석여부(불참)
	AdtAAtn                   string `json:"adt_a_atn"`                      // 감사(감사위원) 참석여부
	RsSmAtn                   string `json:"rs_sm_atn"`                      // 증권신고서 제출대상 여부
	ExSmR                     string `json:"ex_sm_r"`                        // 제출 면제 사유
	OvisLtdtl                 string `json:"ovis_ltdtl"`                     // 해외발행 연계 대차거래 내역
	FtcSttAtn                 string `json:"ftc_stt_atn"`                    // 공정거래위원회 신고대상 여부
}

// BondWithWarrantItem 은 신주인수권부사채권 발행결정 (bdwtIsDecsn) 한 건.
type BondWithWarrantItem struct {
	RceptNo                    string `json:"rcept_no"`                       // 접수번호
	CorpCls                    string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                   string `json:"corp_code"`                      // 고유번호
	CorpName                   string `json:"corp_name"`                      // 회사명
	BdTm                       string `json:"bd_tm"`                          // 사채의 종류(회차)
	BdKnd                      string `json:"bd_knd"`                         // 사채의 종류(종류)
	BdFta                      string `json:"bd_fta"`                         // 사채의 권면(전자등록)총액 (원)
	AtcscRmislmt               string `json:"atcsc_rmislmt"`                  // 정관상 잔여 발행한도 (원)
	OvisFta                    string `json:"ovis_fta"`                       // 해외발행(권면총액)
	OvisFtaCrn                 string `json:"ovis_fta_crn"`                   // 해외발행(권면총액 통화단위)
	OvisSter                   string `json:"ovis_ster"`                      // 해외발행(기준환율등)
	OvisIsar                   string `json:"ovis_isar"`                      // 해외발행(발행지역)
	OvisMktnm                  string `json:"ovis_mktnm"`                     // 해외발행(해외상장시 시장 명칭)
	FdppFclt                   string `json:"fdpp_fclt"`                      // 자금조달목적(시설자금)
	FdppBsninh                 string `json:"fdpp_bsninh"`                    // 자금조달목적(영업양수자금)
	FdppOp                     string `json:"fdpp_op"`                        // 자금조달목적(운영자금)
	FdppDtrp                   string `json:"fdpp_dtrp"`                      // 자금조달목적(채무상환자금)
	FdppOcsa                   string `json:"fdpp_ocsa"`                      // 자금조달목적(타법인 증권 취득자금)
	FdppEtc                    string `json:"fdpp_etc"`                       // 자금조달목적(기타자금)
	BdIntrEx                   string `json:"bd_intr_ex"`                     // 사채의 이율(표면이자율 %)
	BdIntrSf                   string `json:"bd_intr_sf"`                     // 사채의 이율(만기이자율 %)
	BdMtd                      string `json:"bd_mtd"`                         // 사채만기일
	BdisMthn                   string `json:"bdis_mthn"`                      // 사채발행방법
	ExRt                       string `json:"ex_rt"`                          // 신주인수권 행사비율 (%)
	ExPrc                      string `json:"ex_prc"`                         // 신주인수권 행사가액 (원/주)
	ExPrcDmth                  string `json:"ex_prc_dmth"`                    // 신주인수권 행사가액 결정방법
	BdwtDivAtn                 string `json:"bdwt_div_atn"`                   // 사채와 인수권의 분리여부
	NstkPymMth                 string `json:"nstk_pym_mth"`                   // 신주대금 납입방법
	NstkIsstkKnd               string `json:"nstk_isstk_knd"`                 // 행사에 따라 발행할 주식(종류)
	NstkIsstkCnt               string `json:"nstk_isstk_cnt"`                 // 행사에 따라 발행할 주식(주식수)
	NstkIsstkTisstkVs          string `json:"nstk_isstk_tisstk_vs"`           // 행사에 따라 발행할 주식(주식총수 대비 %)
	ExpdBgd                    string `json:"expd_bgd"`                       // 권리행사기간(시작일)
	ExpdEdd                    string `json:"expd_edd"`                       // 권리행사기간(종료일)
	ActMktprcflCvprcLwtrsprc   string `json:"act_mktprcfl_cvprc_lwtrsprc"`    // 시가하락 행사가액 조정(최저 조정가액 원)
	ActMktprcflCvprcLwtrsprcBs string `json:"act_mktprcfl_cvprc_lwtrsprc_bs"` // 시가하락 행사가액 조정(최저 조정가액 근거)
	RmislmtLt70p               string `json:"rmislmt_lt70p"`                  // 시가하락 조정(행사가 70% 미만 조정가능 잔여한도 원)
	Abmg                       string `json:"abmg"`                           // 합병 관련 사항
	Sbd                        string `json:"sbd"`                            // 청약일
	Pymd                       string `json:"pymd"`                           // 납입일
	Rpmcmp                     string `json:"rpmcmp"`                         // 대표주관회사
	Grint                      string `json:"grint"`                          // 보증기관
	Bddd                       string `json:"bddd"`                           // 이사회결의일(결정일)
	OdAAtT                     string `json:"od_a_at_t"`                      // 사외이사 참석여부(참석)
	OdAAtB                     string `json:"od_a_at_b"`                      // 사외이사 참석여부(불참)
	AdtAAtn                    string `json:"adt_a_atn"`                      // 감사(감사위원) 참석여부
	RsSmAtn                    string `json:"rs_sm_atn"`                      // 증권신고서 제출대상 여부
	ExSmR                      string `json:"ex_sm_r"`                        // 제출 면제 사유
	OvisLtdtl                  string `json:"ovis_ltdtl"`                     // 해외발행 연계 대차거래 내역
	FtcSttAtn                  string `json:"ftc_stt_atn"`                    // 공정거래위원회 신고대상 여부
}

// ExchangeableBondItem 은 교환사채권 발행결정 (exbdIsDecsn) 한 건.
type ExchangeableBondItem struct {
	RceptNo      string `json:"rcept_no"`      // 접수번호
	CorpCls      string `json:"corp_cls"`      // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`     // 고유번호
	CorpName     string `json:"corp_name"`     // 회사명
	BdTm         string `json:"bd_tm"`         // 사채의 종류(회차)
	BdKnd        string `json:"bd_knd"`        // 사채의 종류(종류)
	BdFta        string `json:"bd_fta"`        // 사채의 권면(전자등록)총액 (원)
	OvisFta      string `json:"ovis_fta"`      // 해외발행(권면총액)
	OvisFtaCrn   string `json:"ovis_fta_crn"`  // 해외발행(권면총액 통화단위)
	OvisSter     string `json:"ovis_ster"`     // 해외발행(기준환율등)
	OvisIsar     string `json:"ovis_isar"`     // 해외발행(발행지역)
	OvisMktnm    string `json:"ovis_mktnm"`    // 해외발행(해외상장시 시장 명칭)
	FdppFclt     string `json:"fdpp_fclt"`     // 자금조달목적(시설자금)
	FdppBsninh   string `json:"fdpp_bsninh"`   // 자금조달목적(영업양수자금)
	FdppOp       string `json:"fdpp_op"`       // 자금조달목적(운영자금)
	FdppDtrp     string `json:"fdpp_dtrp"`     // 자금조달목적(채무상환자금)
	FdppOcsa     string `json:"fdpp_ocsa"`     // 자금조달목적(타법인 증권 취득자금)
	FdppEtc      string `json:"fdpp_etc"`      // 자금조달목적(기타자금)
	BdIntrEx     string `json:"bd_intr_ex"`    // 사채의 이율(표면이자율 %)
	BdIntrSf     string `json:"bd_intr_sf"`    // 사채의 이율(만기이자율 %)
	BdMtd        string `json:"bd_mtd"`        // 사채만기일
	BdisMthn     string `json:"bdis_mthn"`     // 사채발행방법
	ExRt         string `json:"ex_rt"`         // 교환비율 (%)
	ExPrc        string `json:"ex_prc"`        // 교환가액 (원/주)
	ExPrcDmth    string `json:"ex_prc_dmth"`   // 교환가액 결정방법
	Extg         string `json:"extg"`          // 교환대상(종류)
	ExtgStkcnt   string `json:"extg_stkcnt"`   // 교환대상(주식수)
	ExtgTisstkVs string `json:"extg_tisstk_vs"`// 교환대상(주식총수 대비 %)
	ExrqpdBgd    string `json:"exrqpd_bgd"`    // 교환청구기간(시작일)
	ExrqpdEdd    string `json:"exrqpd_edd"`    // 교환청구기간(종료일)
	Sbd          string `json:"sbd"`           // 청약일
	Pymd         string `json:"pymd"`          // 납입일
	Rpmcmp       string `json:"rpmcmp"`        // 대표주관회사
	Grint        string `json:"grint"`         // 보증기관
	Bddd         string `json:"bddd"`          // 이사회결의일(결정일)
	OdAAtT       string `json:"od_a_at_t"`     // 사외이사 참석여부(참석)
	OdAAtB       string `json:"od_a_at_b"`     // 사외이사 참석여부(불참)
	AdtAAtn      string `json:"adt_a_atn"`     // 감사(감사위원) 참석여부
	RsSmAtn      string `json:"rs_sm_atn"`     // 증권신고서 제출대상 여부
	ExSmR        string `json:"ex_sm_r"`       // 제출 면제 사유
	OvisLtdtl    string `json:"ovis_ltdtl"`    // 해외발행 연계 대차거래 내역
	FtcSttAtn    string `json:"ftc_stt_atn"`   // 공정거래위원회 신고대상 여부
}

// ContingentConvertibleBondItem 은 상각형 조건부자본증권 발행결정 (wdCocobdIsDecsn) 한 건.
type ContingentConvertibleBondItem struct {
	RceptNo    string `json:"rcept_no"`    // 접수번호
	CorpCls    string `json:"corp_cls"`    // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`   // 고유번호
	CorpName   string `json:"corp_name"`   // 회사명
	BdTm       string `json:"bd_tm"`       // 사채의 종류(회차)
	BdKnd      string `json:"bd_knd"`      // 사채의 종류(종류)
	BdFta      string `json:"bd_fta"`      // 사채의 권면(전자등록)총액 (원)
	OvisFta    string `json:"ovis_fta"`    // 해외발행(권면총액)
	OvisFtaCrn string `json:"ovis_fta_crn"`// 해외발행(권면총액 통화단위)
	OvisSter   string `json:"ovis_ster"`   // 해외발행(기준환율등)
	OvisIsar   string `json:"ovis_isar"`   // 해외발행(발행지역)
	OvisMktnm  string `json:"ovis_mktnm"`  // 해외발행(해외상장시 시장 명칭)
	FdppFclt   string `json:"fdpp_fclt"`   // 자금조달목적(시설자금)
	FdppBsninh string `json:"fdpp_bsninh"` // 자금조달목적(영업양수자금)
	FdppOp     string `json:"fdpp_op"`     // 자금조달목적(운영자금)
	FdppDtrp   string `json:"fdpp_dtrp"`   // 자금조달목적(채무상환자금)
	FdppOcsa   string `json:"fdpp_ocsa"`   // 자금조달목적(타법인 증권 취득자금)
	FdppEtc    string `json:"fdpp_etc"`    // 자금조달목적(기타자금)
	BdIntrSf   string `json:"bd_intr_sf"`  // 사채의 이율(표면이자율 %)
	BdIntrEx   string `json:"bd_intr_ex"`  // 사채의 이율(만기이자율 %)
	BdMtd      string `json:"bd_mtd"`      // 사채만기일
	DbtrsSc    string `json:"dbtrs_sc"`    // 채무재조정의 범위
	Sbd        string `json:"sbd"`         // 청약일
	Pymd       string `json:"pymd"`        // 납입일
	Rpmcmp     string `json:"rpmcmp"`      // 대표주관회사
	Grint      string `json:"grint"`       // 보증기관
	Bddd       string `json:"bddd"`        // 이사회결의일(결정일)
	OdAAtT     string `json:"od_a_at_t"`   // 사외이사 참석여부(참석)
	OdAAtB     string `json:"od_a_at_b"`   // 사외이사 참석여부(불참)
	AdtAAtn    string `json:"adt_a_atn"`   // 감사(감사위원) 참석여부
	RsSmAtn    string `json:"rs_sm_atn"`   // 증권신고서 제출대상 여부
	ExSmR      string `json:"ex_sm_r"`     // 제출 면제 사유
	OvisLtdtl  string `json:"ovis_ltdtl"`  // 해외발행 연계 대차거래 내역
	FtcSttAtn  string `json:"ftc_stt_atn"` // 공정거래위원회 신고대상 여부
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/bond_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용.
- 4개 메서드 각각 실 응답 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 로 캡처해 임베드(사채 발행 사례 있는 종목/기간; 공시검색으로 corp_code 탐색).
- `integration_test.go` 에 `ConvertibleBondIssuance` 통합 케이스(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "사채 발행(전환/신주인수권부/교환/상각형 조건부자본증권)" 추가.

## 비범위 (후속 plan)

- DS005 나머지 그룹: 자기주식 / 영업·자산 양수도 / 합병·분할 / 해외상장 (~17개).
- DS006 증권신고서 주요정보. DS002 개인별 보수 Ver 2.0 2종.
