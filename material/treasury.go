package material

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// TreasuryStockAcquisitionItem 은 자기주식 취득 결정 (tsstkAqDecsn) 한 건.
type TreasuryStockAcquisitionItem struct {
	RceptNo        string `json:"rcept_no"`           // 접수번호
	CorpCls        string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`          // 고유번호
	CorpName       string `json:"corp_name"`          // 회사명
	AqplnStkOstk   string `json:"aqpln_stk_ostk"`     // 취득예정주식(주)(보통주식)
	AqplnStkEstk   string `json:"aqpln_stk_estk"`     // 취득예정주식(주)(기타주식)
	AqplnPrcOstk   string `json:"aqpln_prc_ostk"`     // 취득예정금액(원)(보통주식)
	AqplnPrcEstk   string `json:"aqpln_prc_estk"`     // 취득예정금액(원)(기타주식)
	AqexpdBgd      string `json:"aqexpd_bgd"`         // 취득예상기간(시작일)
	AqexpdEdd      string `json:"aqexpd_edd"`         // 취득예상기간(종료일)
	HdexpdBgd      string `json:"hdexpd_bgd"`         // 보유예상기간(시작일)
	HdexpdEdd      string `json:"hdexpd_edd"`         // 보유예상기간(종료일)
	AqPp           string `json:"aq_pp"`              // 취득목적
	AqMth          string `json:"aq_mth"`             // 취득방법
	CsIvBk         string `json:"cs_iv_bk"`           // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`    // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`    // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 취득 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`           // 취득 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`        // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`           // 취득 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`        // 취득 전 자기주식 보유현황(기타취득(주)(비율%))
	AqDd           string `json:"aq_dd"`              // 취득결정일
	OdAAtT         string `json:"od_a_at_t"`          // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`          // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`          // 감사(사외이사가 아닌 감사위원) 참석여부
	D1ProdlmOstk   string `json:"d1_prodlm_ostk"`     // 1일 매수 주문수량 한도(보통주식)
	D1ProdlmEstk   string `json:"d1_prodlm_estk"`     // 1일 매수 주문수량 한도(기타주식)
}

// TreasuryStockAcquisition 은 자기주식 취득 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockAcquisition(ctx context.Context, p MaterialParams) ([]TreasuryStockAcquisitionItem, error) {
	return httpclient.GetList[TreasuryStockAcquisitionItem](ctx, c.http, "/api/tsstkAqDecsn.json", p.toMap())
}

// TreasuryStockDisposalItem 은 자기주식 처분 결정 (tsstkDpDecsn) 한 건.
type TreasuryStockDisposalItem struct {
	RceptNo        string `json:"rcept_no"`           // 접수번호
	CorpCls        string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`          // 고유번호
	CorpName       string `json:"corp_name"`          // 회사명
	DpplnStkOstk   string `json:"dppln_stk_ostk"`     // 처분예정주식(주)(보통주식)
	DpplnStkEstk   string `json:"dppln_stk_estk"`     // 처분예정주식(주)(기타주식)
	DpstkPrcOstk   string `json:"dpstk_prc_ostk"`     // 처분 대상 주식가격(원)(보통주식)
	DpstkPrcEstk   string `json:"dpstk_prc_estk"`     // 처분 대상 주식가격(원)(기타주식)
	DpplnPrcOstk   string `json:"dppln_prc_ostk"`     // 처분예정금액(원)(보통주식)
	DpplnPrcEstk   string `json:"dppln_prc_estk"`     // 처분예정금액(원)(기타주식)
	DpprpdBgd      string `json:"dpprpd_bgd"`         // 처분예정기간(시작일)
	DpprpdEdd      string `json:"dpprpd_edd"`         // 처분예정기간(종료일)
	DpPp           string `json:"dp_pp"`              // 처분목적
	DpMMkt         string `json:"dp_m_mkt"`           // 처분방법(시장을 통한 매도(주))
	DpMOvtm        string `json:"dp_m_ovtm"`          // 처분방법(시간외대량매매(주))
	DpMOtc         string `json:"dp_m_otc"`           // 처분방법(장외처분(주))
	DpMEtc         string `json:"dp_m_etc"`           // 처분방법(기타(주))
	CsIvBk         string `json:"cs_iv_bk"`           // 위탁투자중개업자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`    // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`    // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 처분 전 자기주식 보유현황(배당가능이익 범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`           // 처분 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`        // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`           // 처분 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`        // 처분 전 자기주식 보유현황(기타취득(주)(비율%))
	DpDd           string `json:"dp_dd"`              // 처분결정일
	OdAAtT         string `json:"od_a_at_t"`          // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`          // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`          // 감사(사외이사가 아닌 감사위원) 참석여부
	D1SlodlmOstk   string `json:"d1_slodlm_ostk"`     // 1일 매도 주문수량 한도(보통주식)
	D1SlodlmEstk   string `json:"d1_slodlm_estk"`     // 1일 매도 주문수량 한도(기타주식)
}

// TreasuryStockDisposal 은 자기주식 처분 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockDisposal(ctx context.Context, p MaterialParams) ([]TreasuryStockDisposalItem, error) {
	return httpclient.GetList[TreasuryStockDisposalItem](ctx, c.http, "/api/tsstkDpDecsn.json", p.toMap())
}

// TreasuryStockTrustContractItem 은 자기주식취득 신탁계약 체결 결정 (tsstkAqTrctrCnsDecsn) 한 건.
type TreasuryStockTrustContractItem struct {
	RceptNo        string `json:"rcept_no"`           // 접수번호
	CorpCls        string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`          // 고유번호
	CorpName       string `json:"corp_name"`          // 회사명
	CtrPrc         string `json:"ctr_prc"`            // 계약금액(원)
	CtrPdBgd       string `json:"ctr_pd_bgd"`         // 계약기간(시작일)
	CtrPdEdd       string `json:"ctr_pd_edd"`         // 계약기간(종료일)
	CtrPp          string `json:"ctr_pp"`             // 계약목적
	CtrCnsInt      string `json:"ctr_cns_int"`        // 계약체결기관
	CtrCnsPrd      string `json:"ctr_cns_prd"`        // 계약체결 예정일자
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`    // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`    // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 계약 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`           // 계약 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`        // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`           // 계약 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`        // 계약 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`          // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`          // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`          // 감사(사외이사가 아닌 감사위원) 참석여부
	CsIvBk         string `json:"cs_iv_bk"`           // 위탁투자중개업자
}

// TreasuryStockTrustContract 는 자기주식취득 신탁계약 체결 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockTrustContract(ctx context.Context, p MaterialParams) ([]TreasuryStockTrustContractItem, error) {
	return httpclient.GetList[TreasuryStockTrustContractItem](ctx, c.http, "/api/tsstkAqTrctrCnsDecsn.json", p.toMap())
}

// TreasuryStockTrustCancellationItem 은 자기주식취득 신탁계약 해지 결정 (tsstkAqTrctrCcDecsn) 한 건.
type TreasuryStockTrustCancellationItem struct {
	RceptNo        string `json:"rcept_no"`           // 접수번호
	CorpCls        string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`          // 고유번호
	CorpName       string `json:"corp_name"`          // 회사명
	CtrPrcBfcc     string `json:"ctr_prc_bfcc"`       // 계약금액(원)(해지 전)
	CtrPrcAtcc     string `json:"ctr_prc_atcc"`       // 계약금액(원)(해지 후)
	CtrPdBfccBgd   string `json:"ctr_pd_bfcc_bgd"`    // 해지 전 계약기간(시작일)
	CtrPdBfccEdd   string `json:"ctr_pd_bfcc_edd"`    // 해지 전 계약기간(종료일)
	CcPp           string `json:"cc_pp"`              // 해지목적
	CcInt          string `json:"cc_int"`             // 해지기관
	CcPrd          string `json:"cc_prd"`             // 해지예정일자
	TpRmAtcc       string `json:"tp_rm_atcc"`         // 해지후 신탁재산의 반환방법
	AqWtnDivOstk   string `json:"aq_wtn_div_ostk"`    // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(보통주식))
	AqWtnDivOstkRt string `json:"aq_wtn_div_ostk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	AqWtnDivEstk   string `json:"aq_wtn_div_estk"`    // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(기타주식))
	AqWtnDivEstkRt string `json:"aq_wtn_div_estk_rt"` // 해지 전 자기주식 보유현황(배당가능범위 내 취득(주)(비율%))
	EaqOstk        string `json:"eaq_ostk"`           // 해지 전 자기주식 보유현황(기타취득(주)(보통주식))
	EaqOstkRt      string `json:"eaq_ostk_rt"`        // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	EaqEstk        string `json:"eaq_estk"`           // 해지 전 자기주식 보유현황(기타취득(주)(기타주식))
	EaqEstkRt      string `json:"eaq_estk_rt"`        // 해지 전 자기주식 보유현황(기타취득(주)(비율%))
	Bddd           string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT         string `json:"od_a_at_t"`          // 사외이사 참석여부(참석(명))
	OdAAtB         string `json:"od_a_at_b"`          // 사외이사 참석여부(불참(명))
	AdtAAtn        string `json:"adt_a_atn"`          // 감사(사외이사가 아닌 감사위원) 참석여부
}

// TreasuryStockTrustCancellation 은 자기주식취득 신탁계약 해지 결정(주요사항보고서)을 조회한다.
func (c *Client) TreasuryStockTrustCancellation(ctx context.Context, p MaterialParams) ([]TreasuryStockTrustCancellationItem, error) {
	return httpclient.GetList[TreasuryStockTrustCancellationItem](ctx, c.http, "/api/tsstkAqTrctrCcDecsn.json", p.toMap())
}
