# OpenDART DS005 주요사항보고서 주요정보 — 합병·분할 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 합병·분할 3개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 부실·법적·증자·감자·사채 발행·자기주식·양수도(PR #10~#15)가 main 머지됨. 이 spec 은 DS005
**합병·분할 3개**다(회사합병/회사분할/회사분할합병 결정). 기존 `MaterialParams` + `httpclient.GetList[T]`
재사용, root 변경 없음. 신규 파일 `material/merger.go`. 필드 수가 큼(69/49/90).

## API 표면 (docs 기반 사실)

- 3개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 회사합병(cmpMgDecsn 69), 회사분할(cmpDvDecsn 49), 회사분할합병(cmpDvmgDecsn 90).
- 공통 블록: 머리 + 외부평가(exevl_*) + 거버넌스(bddd/od_a_at_t/b/adt_a_atn) + 우회상장(bdlst_atn/otcpr_bdlst_sf_atn) + 주식매수청구권(aprskh_*) + 풋옵션(popt_ctr_*) + 증권신고서(rs_sm_atn/ex_sm_r). 타입별 고유 필드(합병상대/신설/존속/분할설립 회사·일정·감자)는 각 struct.
- 값은 문자열(금액/주식수 콤마, 비율 %, 빈 값 "-").

## 아키텍처

```
material/
  merger.go       # 3개 메서드 + item struct (신규)
  merger_test.go  # 3개 fixture 테스트 (신규)
  testdata/       # 3개 fixture
README.md         # (수정) DS005 커버리지에 합병·분할
integration_test.go # (수정) 통합 케이스 1~2개 (ErrNoData skip)
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 3개 메서드 (material/merger.go)

| 메서드 | 한글 | 엔드포인트 | 필드 |
|--------|------|-----------|------|
| `CompanyMerger` | 회사합병 결정 | `/api/cmpMgDecsn.json` | 69 |
| `CompanyDivision` | 회사분할 결정 | `/api/cmpDvDecsn.json` | 49 |
| `CompanyDivisionMerger` | 회사분할합병 결정 | `/api/cmpDvmgDecsn.json` | 90 |

```go
// CompanyMergerItem 은 회사합병 결정 (cmpMgDecsn) 한 건.
type CompanyMergerItem struct {
	RceptNo          string `json:"rcept_no"`            // 접수번호
	CorpCls          string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`           // 고유번호
	CorpName         string `json:"corp_name"`           // 회사명
	MgMth            string `json:"mg_mth"`              // 합병방법
	MgStn            string `json:"mg_stn"`              // 합병형태
	MgPp             string `json:"mg_pp"`               // 합병목적
	MgRt             string `json:"mg_rt"`               // 합병비율
	MgRtBs           string `json:"mg_rt_bs"`            // 합병비율 산출근거
	ExevlAtn         string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs        string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn        string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd          string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp          string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	MgnstkOstkCnt    string `json:"mgnstk_ostk_cnt"`     // 합병신주의 종류와 수(주)(보통주식)
	MgnstkCstkCnt    string `json:"mgnstk_cstk_cnt"`     // 합병신주의 종류와 수(주)(종류주식)
	MgptncmpCmpnm    string `json:"mgptncmp_cmpnm"`      // 합병상대회사(회사명)
	MgptncmpMbsn     string `json:"mgptncmp_mbsn"`       // 합병상대회사(주요사업)
	MgptncmpRlCmpn   string `json:"mgptncmp_rl_cmpn"`    // 합병상대회사(회사와의 관계)
	RbsnfdtlTast     string `json:"rbsnfdtl_tast"`       // 합병상대회사 최근 사업연도 재무(원)(자산총계)
	RbsnfdtlTdbt     string `json:"rbsnfdtl_tdbt"`       // 합병상대회사 최근 사업연도 재무(원)(부채총계)
	RbsnfdtlTeqt     string `json:"rbsnfdtl_teqt"`       // 합병상대회사 최근 사업연도 재무(원)(자본총계)
	RbsnfdtlCpt      string `json:"rbsnfdtl_cpt"`        // 합병상대회사 최근 사업연도 재무(원)(자본금)
	RbsnfdtlSl       string `json:"rbsnfdtl_sl"`         // 합병상대회사 최근 사업연도 재무(원)(매출액)
	RbsnfdtlNic      string `json:"rbsnfdtl_nic"`        // 합병상대회사 최근 사업연도 재무(원)(당기순이익)
	EadtatIntn       string `json:"eadtat_intn"`         // 합병상대회사 외부감사 여부(기관명)
	EadtatOp         string `json:"eadtat_op"`           // 합병상대회사 외부감사 여부(감사의견)
	NmgcmpCmpnm      string `json:"nmgcmp_cmpnm"`        // 신설합병회사(회사명)
	FfdtlTast        string `json:"ffdtl_tast"`          // 신설합병회사 설립시 재무(원)(자산총계)
	FfdtlTdbt        string `json:"ffdtl_tdbt"`          // 신설합병회사 설립시 재무(원)(부채총계)
	FfdtlTeqt        string `json:"ffdtl_teqt"`          // 신설합병회사 설립시 재무(원)(자본총계)
	FfdtlCpt         string `json:"ffdtl_cpt"`           // 신설합병회사 설립시 재무(원)(자본금)
	FfdtlStd         string `json:"ffdtl_std"`           // 신설합병회사 설립시 재무(원)(현재기준)
	NmgcmpNbsnRsl    string `json:"nmgcmp_nbsn_rsl"`     // 신설합병회사 신설사업부문 최근 사업연도 매출액(원)
	NmgcmpMbsn       string `json:"nmgcmp_mbsn"`         // 신설합병회사(주요사업)
	NmgcmpRlstAtn    string `json:"nmgcmp_rlst_atn"`     // 신설합병회사(재상장신청 여부)
	MgscMgctrd       string `json:"mgsc_mgctrd"`         // 합병일정(합병계약일)
	MgscShddstd      string `json:"mgsc_shddstd"`        // 합병일정(주주확정기준일)
	MgscShclspdBgd   string `json:"mgsc_shclspd_bgd"`    // 합병일정(주주명부 폐쇄기간(시작일))
	MgscShclspdEdd   string `json:"mgsc_shclspd_edd"`    // 합병일정(주주명부 폐쇄기간(종료일))
	MgscMgopRcpdBgd  string `json:"mgsc_mgop_rcpd_bgd"`  // 합병일정(합병반대의사통지 접수기간(시작일))
	MgscMgopRcpdEdd  string `json:"mgsc_mgop_rcpd_edd"`  // 합병일정(합병반대의사통지 접수기간(종료일))
	MgscGmtsckPrd    string `json:"mgsc_gmtsck_prd"`     // 합병일정(주주총회예정일자)
	MgscAprskhExpdBgd string `json:"mgsc_aprskh_expd_bgd"` // 합병일정(주식매수청구권 행사기간(시작일))
	MgscAprskhExpdEdd string `json:"mgsc_aprskh_expd_edd"` // 합병일정(주식매수청구권 행사기간(종료일))
	MgscOsprpdBgd    string `json:"mgsc_osprpd_bgd"`     // 합병일정(구주권 제출기간(시작일))
	MgscOsprpdEdd    string `json:"mgsc_osprpd_edd"`     // 합병일정(구주권 제출기간(종료일))
	MgscTrspprpdBgd  string `json:"mgsc_trspprpd_bgd"`   // 합병일정(매매거래 정지예정기간(시작일))
	MgscTrspprpdEdd  string `json:"mgsc_trspprpd_edd"`   // 합병일정(매매거래 정지예정기간(종료일))
	MgscCdobprpdBgd  string `json:"mgsc_cdobprpd_bgd"`   // 합병일정(채권자이의 제출기간(시작일))
	MgscCdobprpdEdd  string `json:"mgsc_cdobprpd_edd"`   // 합병일정(채권자이의 제출기간(종료일))
	MgscMgdt         string `json:"mgsc_mgdt"`           // 합병일정(합병기일)
	MgscErgmd        string `json:"mgsc_ergmd"`          // 합병일정(종료보고 총회일)
	MgscMgrgsprd     string `json:"mgsc_mgrgsprd"`       // 합병일정(합병등기예정일자)
	MgscNstkdlprd    string `json:"mgsc_nstkdlprd"`      // 합병일정(신주권교부예정일)
	MgscNstklstprd   string `json:"mgsc_nstklstprd"`     // 합병일정(신주의 상장예정일)
	BdlstAtn         string `json:"bdlst_atn"`           // 우회상장 해당 여부
	OtcprBdlstSfAtn  string `json:"otcpr_bdlst_sf_atn"`  // 타법인의 우회상장 요건 충족여부
	AprskhPlnprc     string `json:"aprskh_plnprc"`       // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth string `json:"aprskh_pym_plpd_mth"` // 주식매수청구권(지급예정시기, 지급방법)
	AprskhCtref      string `json:"aprskh_ctref"`        // 주식매수청구권(계약에 미치는 효력)
	Bddd             string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT           string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB           string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn          string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	PoptCtrAtn       string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn        string `json:"popt_ctr_cn"`         // 계약내용
	RsSmAtn          string `json:"rs_sm_atn"`           // 증권신고서 제출대상 여부
	ExSmR            string `json:"ex_sm_r"`             // 제출을 면제받은 경우 그 사유
}

// CompanyDivisionItem 은 회사분할 결정 (cmpDvDecsn) 한 건.
type CompanyDivisionItem struct {
	RceptNo               string `json:"rcept_no"`                  // 접수번호
	CorpCls               string `json:"corp_cls"`                  // 법인구분 (Y/K/N/E)
	CorpCode              string `json:"corp_code"`                 // 고유번호
	CorpName              string `json:"corp_name"`                 // 회사명
	DvMth                 string `json:"dv_mth"`                    // 분할방법
	DvImpef               string `json:"dv_impef"`                  // 분할의 중요영향 및 효과
	DvRt                  string `json:"dv_rt"`                     // 분할비율
	DvTrfbsnprtCn         string `json:"dv_trfbsnprt_cn"`           // 분할로 이전할 사업 및 재산의 내용
	AtdvExcmpCmpnm        string `json:"atdv_excmp_cmpnm"`          // 분할 후 존속회사(회사명)
	AtdvfdtlTast          string `json:"atdvfdtl_tast"`             // 분할 후 존속회사 분할후 재무(원)(자산총계)
	AtdvfdtlTdbt          string `json:"atdvfdtl_tdbt"`             // 분할 후 존속회사 분할후 재무(원)(부채총계)
	AtdvfdtlTeqt          string `json:"atdvfdtl_teqt"`             // 분할 후 존속회사 분할후 재무(원)(자본총계)
	AtdvfdtlCpt           string `json:"atdvfdtl_cpt"`              // 분할 후 존속회사 분할후 재무(원)(자본금)
	AtdvfdtlStd           string `json:"atdvfdtl_std"`              // 분할 후 존속회사 분할후 재무(원)(현재기준)
	AtdvExcmpExbsnRsl     string `json:"atdv_excmp_exbsn_rsl"`      // 분할 후 존속회사 존속사업부문 최근 사업연도매출액(원)
	AtdvExcmpMbsn         string `json:"atdv_excmp_mbsn"`           // 분할 후 존속회사(주요사업)
	AtdvExcmpAtdvLstmnAtn string `json:"atdv_excmp_atdv_lstmn_atn"` // 분할 후 존속회사(분할 후 상장유지 여부)
	DvfcmpCmpnm           string `json:"dvfcmp_cmpnm"`              // 분할설립회사(회사명)
	FfdtlTast             string `json:"ffdtl_tast"`                // 분할설립회사 설립시 재무(원)(자산총계)
	FfdtlTdbt             string `json:"ffdtl_tdbt"`                // 분할설립회사 설립시 재무(원)(부채총계)
	FfdtlTeqt             string `json:"ffdtl_teqt"`                // 분할설립회사 설립시 재무(원)(자본총계)
	FfdtlCpt              string `json:"ffdtl_cpt"`                 // 분할설립회사 설립시 재무(원)(자본금)
	FfdtlStd              string `json:"ffdtl_std"`                 // 분할설립회사 설립시 재무(원)(현재기준)
	DvfcmpNbsnRsl         string `json:"dvfcmp_nbsn_rsl"`           // 분할설립회사 신설사업부문 최근 사업연도 매출액(원)
	DvfcmpMbsn            string `json:"dvfcmp_mbsn"`               // 분할설립회사(주요사업)
	DvfcmpRlstAtn         string `json:"dvfcmp_rlst_atn"`           // 분할설립회사(재상장신청 여부)
	AbcrCrrt              string `json:"abcr_crrt"`                 // 감자에 관한 사항(감자비율(%))
	AbcrOsprpdBgd         string `json:"abcr_osprpd_bgd"`           // 감자에 관한 사항(구주권 제출기간(시작일))
	AbcrOsprpdEdd         string `json:"abcr_osprpd_edd"`           // 감자에 관한 사항(구주권 제출기간(종료일))
	AbcrTrspprpdBgd       string `json:"abcr_trspprpd_bgd"`         // 감자에 관한 사항(매매거래정지 예정기간(시작일))
	AbcrTrspprpdEdd       string `json:"abcr_trspprpd_edd"`         // 감자에 관한 사항(매매거래정지 예정기간(종료일))
	AbcrNstkascnd         string `json:"abcr_nstkascnd"`            // 감자에 관한 사항(신주배정조건)
	AbcrShstkcntRtAtRs    string `json:"abcr_shstkcnt_rt_at_rs"`    // 감자에 관한 사항(주주 주식수 비례여부 및 사유)
	AbcrNstkasstd         string `json:"abcr_nstkasstd"`            // 감자에 관한 사항(신주배정기준일)
	AbcrNstkdlprd         string `json:"abcr_nstkdlprd"`            // 감자에 관한 사항(신주권교부예정일)
	AbcrNstklstprd        string `json:"abcr_nstklstprd"`           // 감자에 관한 사항(신주의 상장예정일)
	GmtsckPrd             string `json:"gmtsck_prd"`                // 주주총회 예정일
	CdobprpdBgd           string `json:"cdobprpd_bgd"`              // 채권자 이의제출기간(시작일)
	CdobprpdEdd           string `json:"cdobprpd_edd"`              // 채권자 이의제출기간(종료일)
	Dvdt                  string `json:"dvdt"`                      // 분할기일
	Dvrgsprd              string `json:"dvrgsprd"`                  // 분할등기 예정일
	Bddd                  string `json:"bddd"`                      // 이사회결의일(결정일)
	OdAAtT                string `json:"od_a_at_t"`                 // 사외이사 참석여부(참석(명))
	OdAAtB                string `json:"od_a_at_b"`                 // 사외이사 참석여부(불참(명))
	AdtAAtn               string `json:"adt_a_atn"`                 // 감사(사외이사가 아닌 감사위원) 참석여부
	PoptCtrAtn            string `json:"popt_ctr_atn"`              // 풋옵션 등 계약 체결여부
	PoptCtrCn             string `json:"popt_ctr_cn"`               // 계약내용
	RsSmAtn               string `json:"rs_sm_atn"`                 // 증권신고서 제출대상 여부
	ExSmR                 string `json:"ex_sm_r"`                   // 제출을 면제받은 경우 그 사유
}

// CompanyDivisionMergerItem 은 회사분할합병 결정 (cmpDvmgDecsn) 한 건.
type CompanyDivisionMergerItem struct {
	RceptNo               string `json:"rcept_no"`                  // 접수번호
	CorpCls               string `json:"corp_cls"`                  // 법인구분 (Y/K/N/E)
	CorpCode              string `json:"corp_code"`                 // 고유번호
	CorpName              string `json:"corp_name"`                 // 회사명
	DvmgMth               string `json:"dvmg_mth"`                  // 분할합병 방법
	DvmgImpef             string `json:"dvmg_impef"`                // 분할합병의 중요영향 및 효과
	DvTrfbsnprtCn         string `json:"dv_trfbsnprt_cn"`           // 분할(분할로 이전할 사업 및 재산의 내용)
	AtdvExcmpCmpnm        string `json:"atdv_excmp_cmpnm"`          // 분할 후 존속회사(회사명)
	AtdvfdtlTast          string `json:"atdvfdtl_tast"`             // 분할 후 존속회사 분할후 재무(원)(자산총계)
	AtdvfdtlTdbt          string `json:"atdvfdtl_tdbt"`             // 분할 후 존속회사 분할후 재무(원)(부채총계)
	AtdvfdtlTeqt          string `json:"atdvfdtl_teqt"`             // 분할 후 존속회사 분할후 재무(원)(자본총계)
	AtdvfdtlCpt           string `json:"atdvfdtl_cpt"`              // 분할 후 존속회사 분할후 재무(원)(자본금)
	AtdvfdtlStd           string `json:"atdvfdtl_std"`              // 분할 후 존속회사 분할후 재무(원)(현재기준)
	AtdvExcmpExbsnRsl     string `json:"atdv_excmp_exbsn_rsl"`      // 분할 후 존속회사 존속사업부문 최근 사업연도매출액(원)
	AtdvExcmpMbsn         string `json:"atdv_excmp_mbsn"`           // 분할 후 존속회사(주요사업)
	AtdvExcmpAtdvLstmnAtn string `json:"atdv_excmp_atdv_lstmn_atn"` // 분할 후 존속회사(분할 후 상장유지 여부)
	DvfcmpCmpnm           string `json:"dvfcmp_cmpnm"`              // 분할설립 회사(회사명)
	FfdtlTast             string `json:"ffdtl_tast"`                // 분할설립 회사 설립시 재무(원)(자산총계)
	FfdtlTdbt             string `json:"ffdtl_tdbt"`                // 분할설립 회사 설립시 재무(원)(부채총계)
	FfdtlTeqt             string `json:"ffdtl_teqt"`                // 분할설립 회사 설립시 재무(원)(자본총계)
	FfdtlCpt              string `json:"ffdtl_cpt"`                 // 분할설립 회사 설립시 재무(원)(자본금)
	FfdtlStd              string `json:"ffdtl_std"`                 // 분할설립 회사 설립시 재무(원)(현재기준)
	DvfcmpNbsnRsl         string `json:"dvfcmp_nbsn_rsl"`           // 분할설립 회사 신설사업부문 최근 사업연도 매출액(원)
	DvfcmpMbsn            string `json:"dvfcmp_mbsn"`               // 분할설립 회사(주요사업)
	DvfcmpAtdvLstmnAt     string `json:"dvfcmp_atdv_lstmn_at"`      // 분할설립 회사(분할후 상장유지여부)
	AbcrCrrt              string `json:"abcr_crrt"`                 // 감자에 관한 사항(감자비율(%))
	AbcrOsprpdBgd         string `json:"abcr_osprpd_bgd"`           // 감자에 관한 사항(구주권 제출기간(시작일))
	AbcrOsprpdEdd         string `json:"abcr_osprpd_edd"`           // 감자에 관한 사항(구주권 제출기간(종료일))
	AbcrTrspprpdBgd       string `json:"abcr_trspprpd_bgd"`         // 감자에 관한 사항(매매거래정지 예정기간(시작일))
	AbcrTrspprpdEdd       string `json:"abcr_trspprpd_edd"`         // 감자에 관한 사항(매매거래정지 예정기간(종료일))
	AbcrNstkascnd         string `json:"abcr_nstkascnd"`            // 감자에 관한 사항(신주배정조건)
	AbcrShstkcntRtAtRs    string `json:"abcr_shstkcnt_rt_at_rs"`    // 감자에 관한 사항(주주 주식수 비례여부 및 사유)
	AbcrNstkasstd         string `json:"abcr_nstkasstd"`            // 감자에 관한 사항(신주배정기준일)
	AbcrNstkdlprd         string `json:"abcr_nstkdlprd"`            // 감자에 관한 사항(신주권교부예정일)
	AbcrNstklstprd        string `json:"abcr_nstklstprd"`           // 감자에 관한 사항(신주의 상장예정일)
	MgStn                 string `json:"mg_stn"`                    // 합병에 관한 사항(합병형태)
	MgptncmpCmpnm         string `json:"mgptncmp_cmpnm"`            // 합병상대 회사(회사명)
	MgptncmpMbsn          string `json:"mgptncmp_mbsn"`             // 합병상대 회사(주요사업)
	MgptncmpRlCmpn        string `json:"mgptncmp_rl_cmpn"`          // 합병상대 회사(회사와의 관계)
	RbsnfdtlTast          string `json:"rbsnfdtl_tast"`             // 합병상대 회사 최근 사업연도 재무(원)(자산총계)
	RbsnfdtlTdbt          string `json:"rbsnfdtl_tdbt"`             // 합병상대 회사 최근 사업연도 재무(원)(부채총계)
	RbsnfdtlTeqt          string `json:"rbsnfdtl_teqt"`             // 합병상대 회사 최근 사업연도 재무(원)(자본총계)
	RbsnfdtlCpt           string `json:"rbsnfdtl_cpt"`              // 합병상대 회사 최근 사업연도 재무(원)(자본금)
	RbsnfdtlSl            string `json:"rbsnfdtl_sl"`               // 합병상대 회사 최근 사업연도 재무(원)(매출액)
	RbsnfdtlNic           string `json:"rbsnfdtl_nic"`              // 합병상대 회사 최근 사업연도 재무(원)(당기순이익)
	EadtatIntn            string `json:"eadtat_intn"`               // 합병상대 회사 외부감사 여부(기관명)
	EadtatOp              string `json:"eadtat_op"`                 // 합병상대 회사 외부감사 여부(감사의견)
	DvmgnstkOstkCnt       string `json:"dvmgnstk_ostk_cnt"`         // 분할합병신주의 종류와 수(주)(보통주식)
	DvmgnstkCstkCnt       string `json:"dvmgnstk_cstk_cnt"`         // 분할합병신주의 종류와 수(주)(종류주식)
	NmgcmpCmpnm           string `json:"nmgcmp_cmpnm"`              // 합병신설 회사(회사명)
	NmgcmpCpt             string `json:"nmgcmp_cpt"`                // 합병신설 회사(자본금(원))
	NmgcmpMbsn            string `json:"nmgcmp_mbsn"`               // 합병신설 회사(주요사업)
	NmgcmpRlstAtn         string `json:"nmgcmp_rlst_atn"`           // 합병신설 회사(재상장신청 여부)
	DvmgRt                string `json:"dvmg_rt"`                   // 분할합병비율
	DvmgRtBs              string `json:"dvmg_rt_bs"`                // 분할합병비율 산출근거
	ExevlAtn              string `json:"exevl_atn"`                 // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs             string `json:"exevl_bs_rs"`               // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn             string `json:"exevl_intn"`                // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd               string `json:"exevl_pd"`                  // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp               string `json:"exevl_op"`                  // 외부평가에 관한 사항(외부평가 의견)
	DvmgscDvmgctrd        string `json:"dvmgsc_dvmgctrd"`           // 분할합병일정(분할합병계약일)
	DvmgscShddstd         string `json:"dvmgsc_shddstd"`            // 분할합병일정(주주확정기준일)
	DvmgscShclspdBgd      string `json:"dvmgsc_shclspd_bgd"`        // 분할합병일정(주주명부 폐쇄기간(시작일))
	DvmgscShclspdEdd      string `json:"dvmgsc_shclspd_edd"`        // 분할합병일정(주주명부 폐쇄기간(종료일))
	DvmgscDvmgopRcpdBgd   string `json:"dvmgsc_dvmgop_rcpd_bgd"`    // 분할합병일정(분할합병반대의사통지 접수기간(시작일))
	DvmgscDvmgopRcpdEdd   string `json:"dvmgsc_dvmgop_rcpd_edd"`    // 분할합병일정(분할합병반대의사통지 접수기간(종료일))
	DvmgscGmtsckPrd       string `json:"dvmgsc_gmtsck_prd"`         // 분할합병일정(주주총회예정일자)
	DvmgscAprskhExpdBgd   string `json:"dvmgsc_aprskh_expd_bgd"`    // 분할합병일정(주식매수청구권 행사기간(시작일))
	DvmgscAprskhExpdEdd   string `json:"dvmgsc_aprskh_expd_edd"`    // 분할합병일정(주식매수청구권 행사기간(종료일))
	DvmgscCdobprpdBgd     string `json:"dvmgsc_cdobprpd_bgd"`       // 분할합병일정(채권자 이의 제출기간(시작일))
	DvmgscCdobprpdEdd     string `json:"dvmgsc_cdobprpd_edd"`       // 분할합병일정(채권자 이의 제출기간(종료일))
	DvmgscDvmgdt          string `json:"dvmgsc_dvmgdt"`             // 분할합병일정(분할합병기일)
	DvmgscErgmd           string `json:"dvmgsc_ergmd"`              // 분할합병일정(종료보고 총회일)
	DvmgscDvmgrgsprd      string `json:"dvmgsc_dvmgrgsprd"`         // 분할합병일정(분할합병등기예정일)
	BdlstAtn              string `json:"bdlst_atn"`                 // 우회상장 해당 여부
	OtcprBdlstSfAtn       string `json:"otcpr_bdlst_sf_atn"`        // 타법인의 우회상장 요건 충족여부
	AprskhExrq            string `json:"aprskh_exrq"`               // 주식매수청구권(행사요건)
	AprskhPlnprc          string `json:"aprskh_plnprc"`             // 주식매수청구권(매수예정가격)
	AprskhExPcMthPdPl     string `json:"aprskh_ex_pc_mth_pd_pl"`    // 주식매수청구권(행사절차, 방법, 기간, 장소)
	AprskhPymPlpdMth      string `json:"aprskh_pym_plpd_mth"`       // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt             string `json:"aprskh_lmt"`                // 주식매수청구권(제한 관련 내용)
	AprskhCtref           string `json:"aprskh_ctref"`              // 주식매수청구권(계약에 미치는 효력)
	Bddd                  string `json:"bddd"`                      // 이사회결의일(결정일)
	OdAAtT                string `json:"od_a_at_t"`                 // 사외이사 참석여부(참석(명))
	OdAAtB                string `json:"od_a_at_b"`                 // 사외이사 참석여부(불참(명))
	AdtAAtn               string `json:"adt_a_atn"`                 // 감사(사외이사가 아닌 감사위원) 참석여부
	PoptCtrAtn            string `json:"popt_ctr_atn"`              // 풋옵션 등 계약 체결여부
	PoptCtrCn             string `json:"popt_ctr_cn"`               // 계약내용
	RsSmAtn               string `json:"rs_sm_atn"`                 // 증권신고서 제출대상 여부
	ExSmR                 string `json:"ex_sm_r"`                   // 제출을 면제받은 경우 그 사유
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/merger_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용(route map 값은 bare 파일명).
- 3개 메서드 각각 fixture 디코딩 → 대표 필드 검증(머리 + 타입 고유 필드 + 외부평가/거버넌스/증권신고서).
- fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).
- `integration_test.go` 에 통합 케이스 1~2개(`//go:build integration`, ErrNoData skip 허용).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "회사합병·분할·분할합병 결정" 추가.

## 비범위 (후속 plan)

- DS005 해외상장(상장 결정·상장·상장폐지 결정·상장폐지 4) — DS005 마지막 그룹.
- DS006 증권신고서(6). DS002 개인별 보수 Ver 2.0 2종.
