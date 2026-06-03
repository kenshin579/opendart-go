package material

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// OtherAssetTransferPutbackOptionItem 은 자산양수도(기타), 풋백옵션 (astInhtrfEtcPtbkOpt) 한 건.
type OtherAssetTransferPutbackOptionItem struct {
	RceptNo      string `json:"rcept_no"`       // 접수번호
	CorpCls      string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`      // 고유번호
	CorpName     string `json:"corp_name"`      // 회사명
	RpRsn        string `json:"rp_rsn"`         // 보고 사유
	AstInhtrfPrc string `json:"ast_inhtrf_prc"` // 자산양수ㆍ도 가액
}

// OtherAssetTransferPutbackOption 은 자산양수도(기타), 풋백옵션(주요사항보고서)을 조회한다.
func (c *Client) OtherAssetTransferPutbackOption(ctx context.Context, p MaterialParams) ([]OtherAssetTransferPutbackOptionItem, error) {
	return httpclient.GetList[OtherAssetTransferPutbackOptionItem](ctx, c.http, "/api/astInhtrfEtcPtbkOpt.json", p.toMap())
}

// StockExchangeTransferItem 은 주식교환·이전 결정 (stkExtrDecsn) 한 건.
type StockExchangeTransferItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	ExtrSen             string `json:"extr_sen"`               // 구분
	ExtrStn             string `json:"extr_stn"`               // 교환ㆍ이전 형태
	ExtrTgcmpCmpnm      string `json:"extr_tgcmp_cmpnm"`       // 교환ㆍ이전 대상법인(회사명)
	ExtrTgcmpRp         string `json:"extr_tgcmp_rp"`          // 교환ㆍ이전 대상법인(대표자)
	ExtrTgcmpMbsn       string `json:"extr_tgcmp_mbsn"`        // 교환ㆍ이전 대상법인(주요사업)
	ExtrTgcmpRlCmpn     string `json:"extr_tgcmp_rl_cmpn"`     // 교환ㆍ이전 대상법인(회사와의 관계)
	ExtrTgcmpTisstkOstk string `json:"extr_tgcmp_tisstk_ostk"` // 교환ㆍ이전 대상법인(발행주식총수(주)(보통주식))
	ExtrTgcmpTisstkCstk string `json:"extr_tgcmp_tisstk_cstk"` // 교환ㆍ이전 대상법인(발행주식총수(주)(종류주식))
	RbsnfdtlTast        string `json:"rbsnfdtl_tast"`          // 대상법인 최근 사업연도 요약재무(원)(자산총계)
	RbsnfdtlTdbt        string `json:"rbsnfdtl_tdbt"`          // 대상법인 최근 사업연도 요약재무(원)(부채총계)
	RbsnfdtlTeqt        string `json:"rbsnfdtl_teqt"`          // 대상법인 최근 사업연도 요약재무(원)(자본총계)
	RbsnfdtlCpt         string `json:"rbsnfdtl_cpt"`           // 대상법인 최근 사업연도 요약재무(원)(자본금)
	ExtrRt              string `json:"extr_rt"`                // 교환ㆍ이전 비율
	ExtrRtBs            string `json:"extr_rt_bs"`             // 교환ㆍ이전 비율 산출근거
	ExevlAtn            string `json:"exevl_atn"`              // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs           string `json:"exevl_bs_rs"`            // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn           string `json:"exevl_intn"`             // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd             string `json:"exevl_pd"`               // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp             string `json:"exevl_op"`               // 외부평가에 관한 사항(외부평가 의견)
	ExtrPp              string `json:"extr_pp"`                // 교환ㆍ이전 목적
	ExtrscExtrctrd      string `json:"extrsc_extrctrd"`        // 교환ㆍ이전일정(교환ㆍ이전계약일)
	ExtrscShddstd       string `json:"extrsc_shddstd"`         // 교환ㆍ이전일정(주주확정기준일)
	ExtrscShclspdBgd    string `json:"extrsc_shclspd_bgd"`     // 교환ㆍ이전일정(주주명부 폐쇄기간(시작일))
	ExtrscShclspdEdd    string `json:"extrsc_shclspd_edd"`     // 교환ㆍ이전일정(주주명부 폐쇄기간(종료일))
	ExtrscExtropRcpdBgd string `json:"extrsc_extrop_rcpd_bgd"` // 교환ㆍ이전일정(반대의사 통지접수기간(시작일))
	ExtrscExtropRcpdEdd string `json:"extrsc_extrop_rcpd_edd"` // 교환ㆍ이전일정(반대의사 통지접수기간(종료일))
	ExtrscGmtsckPrd     string `json:"extrsc_gmtsck_prd"`      // 교환ㆍ이전일정(주주총회 예정일자)
	ExtrscAprskhExpdBgd string `json:"extrsc_aprskh_expd_bgd"` // 교환ㆍ이전일정(주식매수청구권 행사기간(시작일))
	ExtrscAprskhExpdEdd string `json:"extrsc_aprskh_expd_edd"` // 교환ㆍ이전일정(주식매수청구권 행사기간(종료일))
	ExtrscOsprpdBgd     string `json:"extrsc_osprpd_bgd"`      // 교환ㆍ이전일정(구주권제출기간(시작일))
	ExtrscOsprpdEdd     string `json:"extrsc_osprpd_edd"`      // 교환ㆍ이전일정(구주권제출기간(종료일))
	ExtrscTrspprpd      string `json:"extrsc_trspprpd"`        // 교환ㆍ이전일정(매매거래정지예정기간)
	ExtrscTrspprpdBgd   string `json:"extrsc_trspprpd_bgd"`    // 교환ㆍ이전일정(매매거래정지예정기간(시작일))
	ExtrscTrspprpdEdd   string `json:"extrsc_trspprpd_edd"`    // 교환ㆍ이전일정(매매거래정지예정기간(종료일))
	ExtrscExtrdt        string `json:"extrsc_extrdt"`          // 교환ㆍ이전일정(교환ㆍ이전일자)
	ExtrscNstkdlprd     string `json:"extrsc_nstkdlprd"`       // 교환ㆍ이전일정(신주권교부예정일)
	ExtrscNstklstprd    string `json:"extrsc_nstklstprd"`      // 교환ㆍ이전일정(신주의 상장예정일)
	AtextrCpcmpnm       string `json:"atextr_cpcmpnm"`         // 교환ㆍ이전 후 완전모회사명
	AprskhPlnprc        string `json:"aprskh_plnprc"`          // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth    string `json:"aprskh_pym_plpd_mth"`    // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt           string `json:"aprskh_lmt"`             // 주식매수청구권(제한 관련 내용)
	AprskhCtref         string `json:"aprskh_ctref"`           // 주식매수청구권(계약에 미치는 효력)
	BdlstAtn            string `json:"bdlst_atn"`              // 우회상장 해당 여부
	OtcprBdlstSfAtn     string `json:"otcpr_bdlst_sf_atn"`     // 타법인의 우회상장 요건 충족 여부
	Bddd                string `json:"bddd"`                   // 이사회결의일(결정일)
	OdAAtT              string `json:"od_a_at_t"`              // 사외이사 참석여부(참석(명))
	OdAAtB              string `json:"od_a_at_b"`              // 사외이사 참석여부(불참(명))
	AdtAAtn             string `json:"adt_a_atn"`              // 감사(사외이사가 아닌 감사위원) 참석여부
	PoptCtrAtn          string `json:"popt_ctr_atn"`           // 풋옵션 등 계약 체결여부
	PoptCtrCn           string `json:"popt_ctr_cn"`            // 계약내용
	RsSmAtn             string `json:"rs_sm_atn"`              // 증권신고서 제출대상 여부
	ExSmR               string `json:"ex_sm_r"`                // 제출을 면제받은 경우 그 사유
}

// StockExchangeTransfer 는 주식교환·이전 결정(주요사항보고서)을 조회한다.
func (c *Client) StockExchangeTransfer(ctx context.Context, p MaterialParams) ([]StockExchangeTransferItem, error) {
	return httpclient.GetList[StockExchangeTransferItem](ctx, c.http, "/api/stkExtrDecsn.json", p.toMap())
}
