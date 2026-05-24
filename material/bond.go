package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// ConvertibleBondItem 은 전환사채권 발행결정 (cvbdIsDecsn) 한 건.
type ConvertibleBondItem struct {
	RceptNo                    string `json:"rcept_no"`                       // 접수번호
	CorpCls                    string `json:"corp_cls"`                       // 법인구분 (Y/K/N/E)
	CorpCode                   string `json:"corp_code"`                      // 고유번호
	CorpName                   string `json:"corp_name"`                      // 회사명
	BdTm                       string `json:"bd_tm"`                          // 사채의 종류(회차)
	BdKnd                      string `json:"bd_knd"`                         // 사채의 종류(종류)
	BdFta                      string `json:"bd_fta"`                         // 사채의 권면(전자등록)총액 (원)
	AtcscRmislmt               string `json:"atcsc_rmislmt"`                  // 정관상 잔여 발행한도 (원)
	OvisFta                    string `json:"ovis_fta"`                       // 해외발행(권면(전자등록)총액)
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
	CvRt                       string `json:"cv_rt"`                          // 전환비율 (%)
	CvPrc                      string `json:"cv_prc"`                         // 전환가액 (원/주)
	CvisstkKnd                 string `json:"cvisstk_knd"`                    // 전환에 따라 발행할 주식(종류)
	CvisstkCnt                 string `json:"cvisstk_cnt"`                    // 전환에 따라 발행할 주식(주식수)
	CvisstkTisstkVs            string `json:"cvisstk_tisstk_vs"`              // 전환에 따라 발행할 주식(주식총수 대비 %)
	CvrqpdBgd                  string `json:"cvrqpd_bgd"`                     // 전환청구기간(시작일)
	CvrqpdEdd                  string `json:"cvrqpd_edd"`                     // 전환청구기간(종료일)
	ActMktprcflCvprcLwtrsprc   string `json:"act_mktprcfl_cvprc_lwtrsprc"`    // 시가하락 전환가액 조정(최저 조정가액 원)
	ActMktprcflCvprcLwtrsprcBs string `json:"act_mktprcfl_cvprc_lwtrsprc_bs"` // 시가하락 전환가액 조정(최저 조정가액 근거)
	RmislmtLt70p               string `json:"rmislmt_lt70p"`                  // 시가하락 조정(전환가 70% 미만 조정가능 잔여한도 원)
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

// ConvertibleBondIssuance 는 전환사채권 발행결정(주요사항보고서)을 조회한다.
func (c *Client) ConvertibleBondIssuance(ctx context.Context, p MaterialParams) ([]ConvertibleBondItem, error) {
	return httpclient.GetList[ConvertibleBondItem](ctx, c.http, "/api/cvbdIsDecsn.json", p.toMap())
}
