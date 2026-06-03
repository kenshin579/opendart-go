package material

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

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

// PaidInCapitalIncrease 는 유상증자 결정을 조회한다.
func (c *Client) PaidInCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidInCapitalIncreaseItem, error) {
	return httpclient.GetList[PaidInCapitalIncreaseItem](ctx, c.http, "/api/piicDecsn.json", p.toMap())
}

// FreeCapitalIncreaseItem 은 무상증자 결정 (fricDecsn) 한 건.
type FreeCapitalIncreaseItem struct {
	RceptNo         string `json:"rcept_no"`           // 접수번호
	CorpCls         string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`          // 고유번호
	CorpName        string `json:"corp_name"`          // 회사명
	NstkOstkCnt     string `json:"nstk_ostk_cnt"`      // 신주의 종류와 수(보통주식)
	NstkEstkCnt     string `json:"nstk_estk_cnt"`      // 신주의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`              // 1주당 액면가액 (원)
	BficTisstkOstk  string `json:"bfic_tisstk_ostk"`   // 증자전 발행주식총수(보통주식)
	BficTisstkEstk  string `json:"bfic_tisstk_estk"`   // 증자전 발행주식총수(기타주식)
	NstkAsstd       string `json:"nstk_asstd"`         // 신주배정기준일
	NstkAscntPsOstk string `json:"nstk_ascnt_ps_ostk"` // 1주당 신주배정 주식수(보통주식)
	NstkAscntPsEstk string `json:"nstk_ascnt_ps_estk"` // 1주당 신주배정 주식수(기타주식)
	NstkDividrk     string `json:"nstk_dividrk"`       // 신주의 배당기산일
	NstkDlprd       string `json:"nstk_dlprd"`         // 신주권교부예정일
	NstkLstprd      string `json:"nstk_lstprd"`        // 신주의 상장 예정일
	Bddd            string `json:"bddd"`               // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`          // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`          // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`          // 감사(감사위원) 참석여부
}

// FreeCapitalIncrease 는 무상증자 결정을 조회한다.
func (c *Client) FreeCapitalIncrease(ctx context.Context, p MaterialParams) ([]FreeCapitalIncreaseItem, error) {
	return httpclient.GetList[FreeCapitalIncreaseItem](ctx, c.http, "/api/fricDecsn.json", p.toMap())
}

// PaidFreeCapitalIncreaseItem 은 유무상증자 결정 (pifricDecsn) 한 건. piic_* 유상 / fric_* 무상.
type PaidFreeCapitalIncreaseItem struct {
	RceptNo             string `json:"rcept_no"`                // 접수번호
	CorpCls             string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`               // 고유번호
	CorpName            string `json:"corp_name"`               // 회사명
	PiicNstkOstkCnt     string `json:"piic_nstk_ostk_cnt"`      // 유상증자 신주수(보통주식)
	PiicNstkEstkCnt     string `json:"piic_nstk_estk_cnt"`      // 유상증자 신주수(기타주식)
	PiicFvPs            string `json:"piic_fv_ps"`              // 유상증자 1주당 액면가액
	PiicBficTisstkOstk  string `json:"piic_bfic_tisstk_ostk"`   // 유상증자 증자전 발행총수(보통주식)
	PiicBficTisstkEstk  string `json:"piic_bfic_tisstk_estk"`   // 유상증자 증자전 발행총수(기타주식)
	PiicFdppFclt        string `json:"piic_fdpp_fclt"`          // 유상증자 자금조달목적(시설자금)
	PiicFdppBsninh      string `json:"piic_fdpp_bsninh"`        // 유상증자 자금조달목적(영업양수자금)
	PiicFdppOp          string `json:"piic_fdpp_op"`            // 유상증자 자금조달목적(운영자금)
	PiicFdppDtrp        string `json:"piic_fdpp_dtrp"`          // 유상증자 자금조달목적(채무상환자금)
	PiicFdppOcsa        string `json:"piic_fdpp_ocsa"`          // 유상증자 자금조달목적(타법인 증권 취득자금)
	PiicFdppEtc         string `json:"piic_fdpp_etc"`           // 유상증자 자금조달목적(기타자금)
	PiicIcMthn          string `json:"piic_ic_mthn"`            // 유상증자 증자방식
	FricNstkOstkCnt     string `json:"fric_nstk_ostk_cnt"`      // 무상증자 신주수(보통주식)
	FricNstkEstkCnt     string `json:"fric_nstk_estk_cnt"`      // 무상증자 신주수(기타주식)
	FricFvPs            string `json:"fric_fv_ps"`              // 무상증자 1주당 액면가액
	FricBficTisstkOstk  string `json:"fric_bfic_tisstk_ostk"`   // 무상증자 증자전 발행총수(보통주식)
	FricBficTisstkEstk  string `json:"fric_bfic_tisstk_estk"`   // 무상증자 증자전 발행총수(기타주식)
	FricNstkAsstd       string `json:"fric_nstk_asstd"`         // 무상증자 신주배정기준일
	FricNstkAscntPsOstk string `json:"fric_nstk_ascnt_ps_ostk"` // 무상증자 1주당 신주배정수(보통주식)
	FricNstkAscntPsEstk string `json:"fric_nstk_ascnt_ps_estk"` // 무상증자 1주당 신주배정수(기타주식)
	FricNstkDividrk     string `json:"fric_nstk_dividrk"`       // 무상증자 신주 배당기산일
	FricNstkDlprd       string `json:"fric_nstk_dlprd"`         // 무상증자 신주권교부예정일
	FricNstkLstprd      string `json:"fric_nstk_lstprd"`        // 무상증자 신주 상장예정일
	FricBddd            string `json:"fric_bddd"`               // 무상증자 이사회결의일(결정일)
	FricOdAAtT          string `json:"fric_od_a_at_t"`          // 무상증자 사외이사 참석(참석)
	FricOdAAtB          string `json:"fric_od_a_at_b"`          // 무상증자 사외이사 참석(불참)
	FricAdtAAtn         string `json:"fric_adt_a_atn"`          // 무상증자 감사 참석여부
	SslAt               string `json:"ssl_at"`                  // 공매도 해당여부
	SslBgd              string `json:"ssl_bgd"`                 // 공매도 시작일
	SslEdd              string `json:"ssl_edd"`                 // 공매도 종료일
}

// PaidFreeCapitalIncrease 는 유무상증자 결정을 조회한다.
func (c *Client) PaidFreeCapitalIncrease(ctx context.Context, p MaterialParams) ([]PaidFreeCapitalIncreaseItem, error) {
	return httpclient.GetList[PaidFreeCapitalIncreaseItem](ctx, c.http, "/api/pifricDecsn.json", p.toMap())
}

// CapitalReductionItem 은 감자 결정 (crDecsn) 한 건.
type CapitalReductionItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	CorpCls         string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`         // 고유번호
	CorpName        string `json:"corp_name"`         // 회사명
	CrstkOstkCnt    string `json:"crstk_ostk_cnt"`    // 감자주식의 종류와 수(보통주식)
	CrstkEstkCnt    string `json:"crstk_estk_cnt"`    // 감자주식의 종류와 수(기타주식)
	FvPs            string `json:"fv_ps"`             // 1주당 액면가액 (원)
	BfcrCpt         string `json:"bfcr_cpt"`          // 감자전 자본금 (원)
	AtcrCpt         string `json:"atcr_cpt"`          // 감자후 자본금 (원)
	BfcrTisstkOstk  string `json:"bfcr_tisstk_ostk"`  // 감자전 발행주식수(보통주식)
	AtcrTisstkOstk  string `json:"atcr_tisstk_ostk"`  // 감자후 발행주식수(보통주식)
	BfcrTisstkEstk  string `json:"bfcr_tisstk_estk"`  // 감자전 발행주식수(기타주식)
	AtcrTisstkEstk  string `json:"atcr_tisstk_estk"`  // 감자후 발행주식수(기타주식)
	CrRtOstk        string `json:"cr_rt_ostk"`        // 감자비율(보통주식 %)
	CrRtEstk        string `json:"cr_rt_estk"`        // 감자비율(기타주식 %)
	CrStd           string `json:"cr_std"`            // 감자기준일
	CrMth           string `json:"cr_mth"`            // 감자방법
	CrRs            string `json:"cr_rs"`             // 감자사유
	CrscGmtsckPrd   string `json:"crsc_gmtsck_prd"`   // 감자일정(주주총회 예정일)
	CrscTrnmsppd    string `json:"crsc_trnmsppd"`     // 감자일정(명의개서정지기간)
	CrscOsprpd      string `json:"crsc_osprpd"`       // 감자일정(구주권 제출기간)
	CrscTrspprpd    string `json:"crsc_trspprpd"`     // 감자일정(매매거래 정지예정기간)
	CrscOsprpdBgd   string `json:"crsc_osprpd_bgd"`   // 감자일정(구주권 제출기간 시작일)
	CrscOsprpdEdd   string `json:"crsc_osprpd_edd"`   // 감자일정(구주권 제출기간 종료일)
	CrscTrspprpdBgd string `json:"crsc_trspprpd_bgd"` // 감자일정(매매거래 정지예정기간 시작일)
	CrscTrspprpdEdd string `json:"crsc_trspprpd_edd"` // 감자일정(매매거래 정지예정기간 종료일)
	CrscNstkdlprd   string `json:"crsc_nstkdlprd"`    // 감자일정(신주권교부예정일)
	CrscNstklstprd  string `json:"crsc_nstklstprd"`   // 감자일정(신주상장예정일)
	CdobprpdBgd     string `json:"cdobprpd_bgd"`      // 채권자 이의제출기간(시작일)
	CdobprpdEdd     string `json:"cdobprpd_edd"`      // 채권자 이의제출기간(종료일)
	OsprNstkdlPl    string `json:"ospr_nstkdl_pl"`    // 구주권제출 및 신주권교부장소
	Bddd            string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`         // 사외이사 참석여부(참석)
	OdAAtB          string `json:"od_a_at_b"`         // 사외이사 참석여부(불참)
	AdtAAtn         string `json:"adt_a_atn"`         // 감사(감사위원) 참석여부
	FtcSttAtn       string `json:"ftc_stt_atn"`       // 공정거래위원회 신고대상 여부
}

// CapitalReduction 은 감자 결정을 조회한다.
func (c *Client) CapitalReduction(ctx context.Context, p MaterialParams) ([]CapitalReductionItem, error) {
	return httpclient.GetList[CapitalReductionItem](ctx, c.http, "/api/crDecsn.json", p.toMap())
}
