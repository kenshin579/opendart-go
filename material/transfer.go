package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// BusinessAcquisitionItem 은 영업양수 결정 (bsnInhDecsn) 한 건.
type BusinessAcquisitionItem struct {
	RceptNo          string `json:"rcept_no"`            // 접수번호
	CorpCls          string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`           // 고유번호
	CorpName         string `json:"corp_name"`           // 회사명
	InhBsn           string `json:"inh_bsn"`             // 양수영업
	InhBsnMc         string `json:"inh_bsn_mc"`          // 양수영업 주요내용
	InhPrc           string `json:"inh_prc"`             // 양수가액(원)
	AbsnInhAtn       string `json:"absn_inh_atn"`        // 영업전부의 양수 여부
	AstInhBsn        string `json:"ast_inh_bsn"`         // 재무내용(자산액(양수대상 영업부문 A))
	AstCmpAll        string `json:"ast_cmp_all"`         // 재무내용(자산액(당사전체 B))
	AstRt            string `json:"ast_rt"`              // 재무내용(자산액 비중 % A/B)
	SlInhBsn         string `json:"sl_inh_bsn"`          // 재무내용(매출액(양수대상 영업부문 A))
	SlCmpAll         string `json:"sl_cmp_all"`          // 재무내용(매출액(당사전체 B))
	SlRt             string `json:"sl_rt"`               // 재무내용(매출액 비중 % A/B)
	DbtInhBsn        string `json:"dbt_inh_bsn"`         // 재무내용(부채액(양수대상 영업부문 A))
	DbtCmpAll        string `json:"dbt_cmp_all"`         // 재무내용(부채액(당사전체 B))
	DbtRt            string `json:"dbt_rt"`              // 재무내용(부채액 비중 % A/B)
	InhPp            string `json:"inh_pp"`              // 양수목적
	InhAf            string `json:"inh_af"`              // 양수영향
	InhPrdCtrCnsd    string `json:"inh_prd_ctr_cnsd"`    // 양수예정일자(계약체결일)
	InhPrdInhStd     string `json:"inh_prd_inh_std"`     // 양수예정일자(양수기준일)
	DlptnCmpnm       string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt         string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn        string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd       string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn      string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	InhPym           string `json:"inh_pym"`             // 양수대금지급
	ExevlAtn         string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs        string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn        string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd          string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp          string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn     string `json:"gmtsck_spd_atn"`      // 주주총회 특별결의 여부
	GmtsckPrd        string `json:"gmtsck_prd"`          // 주주총회 예정일자
	AprskhPlnprc     string `json:"aprskh_plnprc"`       // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth string `json:"aprskh_pym_plpd_mth"` // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt        string `json:"aprskh_lmt"`          // 주식매수청구권(제한 관련 내용)
	AprskhCtref      string `json:"aprskh_ctref"`        // 주식매수청구권(계약에 미치는 효력)
	Bddd             string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT           string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB           string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn          string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	BdlstAtn         string `json:"bdlst_atn"`           // 우회상장 해당 여부
	N6mTpaiPlann     string `json:"n6m_tpai_plann"`      // 향후 6월이내 제3자배정 증자 등 계획
	OtcprBdlstSfAtn  string `json:"otcpr_bdlst_sf_atn"`  // 타법인의 우회상장 요건 충족여부
	FtcSttAtn        string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn       string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn        string `json:"popt_ctr_cn"`         // 계약내용
}

// BusinessAcquisition 은 영업양수 결정(주요사항보고서)을 조회한다.
func (c *Client) BusinessAcquisition(ctx context.Context, p MaterialParams) ([]BusinessAcquisitionItem, error) {
	return httpclient.GetList[BusinessAcquisitionItem](ctx, c.http, "/api/bsnInhDecsn.json", p.toMap())
}

// BusinessTransferItem 은 영업양도 결정 (bsnTrfDecsn) 한 건.
type BusinessTransferItem struct {
	RceptNo          string `json:"rcept_no"`            // 접수번호
	CorpCls          string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`           // 고유번호
	CorpName         string `json:"corp_name"`           // 회사명
	TrfBsn           string `json:"trf_bsn"`             // 양도영업
	TrfBsnMc         string `json:"trf_bsn_mc"`          // 양도영업 주요내용
	TrfPrc           string `json:"trf_prc"`             // 양도가액(원)
	AstTrfBsn        string `json:"ast_trf_bsn"`         // 재무내용(자산액(양도대상 영업부문 A))
	AstCmpAll        string `json:"ast_cmp_all"`         // 재무내용(자산액(당사전체 B))
	AstRt            string `json:"ast_rt"`              // 재무내용(자산액 비중 % A/B)
	SlTrfBsn         string `json:"sl_trf_bsn"`          // 재무내용(매출액(양도대상 영업부문 A))
	SlCmpAll         string `json:"sl_cmp_all"`          // 재무내용(매출액(당사전체 B))
	SlRt             string `json:"sl_rt"`               // 재무내용(매출액 비중 % A/B)
	TrfPp            string `json:"trf_pp"`              // 양도목적
	TrfAf            string `json:"trf_af"`              // 양도영향
	TrfPrdCtrCnsd    string `json:"trf_prd_ctr_cnsd"`    // 양도예정일자(계약체결일)
	TrfPrdTrfStd     string `json:"trf_prd_trf_std"`     // 양도예정일자(양도기준일)
	DlptnCmpnm       string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt         string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn        string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd       string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn      string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	TrfPym           string `json:"trf_pym"`             // 양도대금지급
	ExevlAtn         string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs        string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn        string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd          string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp          string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn     string `json:"gmtsck_spd_atn"`      // 주주총회 특별결의 여부
	GmtsckPrd        string `json:"gmtsck_prd"`          // 주주총회 예정일자
	AprskhPlnprc     string `json:"aprskh_plnprc"`       // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth string `json:"aprskh_pym_plpd_mth"` // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt        string `json:"aprskh_lmt"`          // 주식매수청구권(제한 관련 내용)
	AprskhCtref      string `json:"aprskh_ctref"`        // 주식매수청구권(계약에 미치는 효력)
	Bddd             string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT           string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB           string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn          string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn        string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn       string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn        string `json:"popt_ctr_cn"`         // 계약내용
}

// BusinessTransfer 는 영업양도 결정(주요사항보고서)을 조회한다.
func (c *Client) BusinessTransfer(ctx context.Context, p MaterialParams) ([]BusinessTransferItem, error) {
	return httpclient.GetList[BusinessTransferItem](ctx, c.http, "/api/bsnTrfDecsn.json", p.toMap())
}

// TangibleAssetAcquisitionItem 은 유형자산 양수 결정 (tgastInhDecsn) 한 건.
type TangibleAssetAcquisitionItem struct {
	RceptNo           string `json:"rcept_no"`               // 접수번호
	CorpCls           string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode          string `json:"corp_code"`              // 고유번호
	CorpName          string `json:"corp_name"`              // 회사명
	AstSen            string `json:"ast_sen"`                // 자산구분
	AstNm             string `json:"ast_nm"`                 // 자산명
	InhdtlInhprc      string `json:"inhdtl_inhprc"`          // 양수내역(양수금액(원))
	InhdtlTast        string `json:"inhdtl_tast"`            // 양수내역(자산총액(원))
	InhdtlTastVs      string `json:"inhdtl_tast_vs"`         // 양수내역(자산총액대비(%))
	InhPp             string `json:"inh_pp"`                 // 양수목적
	InhAf             string `json:"inh_af"`                 // 양수영향
	InhPrdCtrCnsd     string `json:"inh_prd_ctr_cnsd"`       // 양수예정일자(계약체결일)
	InhPrdInhStd      string `json:"inh_prd_inh_std"`        // 양수예정일자(양수기준일)
	InhPrdRgsPrd      string `json:"inh_prd_rgs_prd"`        // 양수예정일자(등기예정일)
	DlptnCmpnm        string `json:"dlptn_cmpnm"`            // 거래상대방(회사명/성명)
	DlptnCpt          string `json:"dlptn_cpt"`              // 거래상대방(자본금(원))
	DlptnMbsn         string `json:"dlptn_mbsn"`             // 거래상대방(주요사업)
	DlptnHoadd        string `json:"dlptn_hoadd"`            // 거래상대방(본점소재지(주소))
	DlptnRlCmpn       string `json:"dlptn_rl_cmpn"`          // 거래상대방(회사와의 관계)
	DlPym             string `json:"dl_pym"`                 // 거래대금지급
	ExevlAtn          string `json:"exevl_atn"`              // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs         string `json:"exevl_bs_rs"`            // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn         string `json:"exevl_intn"`             // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd           string `json:"exevl_pd"`               // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp           string `json:"exevl_op"`               // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn      string `json:"gmtsck_spd_atn"`         // 주주총회 특별결의 여부
	GmtsckPrd         string `json:"gmtsck_prd"`             // 주주총회 예정일자
	AprskhExrq        string `json:"aprskh_exrq"`            // 주식매수청구권(행사요건)
	AprskhPlnprc      string `json:"aprskh_plnprc"`          // 주식매수청구권(매수예정가격)
	AprskhExPcMthPdPl string `json:"aprskh_ex_pc_mth_pd_pl"` // 주식매수청구권(행사절차, 방법, 기간, 장소)
	AprskhPymPlpdMth  string `json:"aprskh_pym_plpd_mth"`    // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt         string `json:"aprskh_lmt"`             // 주식매수청구권(제한 관련 내용)
	AprskhCtref       string `json:"aprskh_ctref"`           // 주식매수청구권(계약에 미치는 효력)
	Bddd              string `json:"bddd"`                   // 이사회결의일(결정일)
	OdAAtT            string `json:"od_a_at_t"`              // 사외이사 참석여부(참석(명))
	OdAAtB            string `json:"od_a_at_b"`              // 사외이사 참석여부(불참(명))
	AdtAAtn           string `json:"adt_a_atn"`              // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn         string `json:"ftc_stt_atn"`            // 공정거래위원회 신고대상 여부
	PoptCtrAtn        string `json:"popt_ctr_atn"`           // 풋옵션 등 계약 체결여부
	PoptCtrCn         string `json:"popt_ctr_cn"`            // 계약내용
}

// TangibleAssetAcquisition 은 유형자산 양수 결정(주요사항보고서)을 조회한다.
func (c *Client) TangibleAssetAcquisition(ctx context.Context, p MaterialParams) ([]TangibleAssetAcquisitionItem, error) {
	return httpclient.GetList[TangibleAssetAcquisitionItem](ctx, c.http, "/api/tgastInhDecsn.json", p.toMap())
}
