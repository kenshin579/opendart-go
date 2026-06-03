package material

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// DefaultItem 은 부도발생 (dfOcr) 한 건.
type DefaultItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DfCn     string `json:"df_cn"`     // 부도내용
	DfAmt    string `json:"df_amt"`    // 부도금액
	DfBnk    string `json:"df_bnk"`    // 부도발생은행
	Dfd      string `json:"dfd"`       // 최종부도(당좌거래정지)일자
	DfRs     string `json:"df_rs"`     // 부도사유 및 경위
}

// DefaultOccurrences 는 부도발생을 조회한다.
func (c *Client) DefaultOccurrences(ctx context.Context, p MaterialParams) ([]DefaultItem, error) {
	return httpclient.GetList[DefaultItem](ctx, c.http, "/api/dfOcr.json", p.toMap())
}

// BusinessSuspensionItem 은 영업정지 (bsnSp) 한 건.
type BusinessSuspensionItem struct {
	RceptNo   string `json:"rcept_no"`    // 접수번호
	CorpCls   string `json:"corp_cls"`    // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`   // 고유번호
	CorpName  string `json:"corp_name"`   // 회사명
	BsnspRm   string `json:"bsnsp_rm"`    // 영업정지 분야
	BsnspAmt  string `json:"bsnsp_amt"`   // 영업정지 내역(영업정지금액)
	Rsl       string `json:"rsl"`         // 영업정지 내역(최근매출총액)
	SlVs      string `json:"sl_vs"`       // 영업정지 내역(매출액 대비)
	LsAtn     string `json:"ls_atn"`      // 영업정지 내역(대규모법인여부)
	KrxSttAtn string `json:"krx_stt_atn"` // 영업정지 내역(거래소 의무공시 해당 여부)
	BsnspCn   string `json:"bsnsp_cn"`    // 영업정지 내용
	BsnspRs   string `json:"bsnsp_rs"`    // 영업정지사유
	FtCtp     string `json:"ft_ctp"`      // 향후대책
	BsnspAf   string `json:"bsnsp_af"`    // 영업정지영향
	Bsnspd    string `json:"bsnspd"`      // 영업정지일자
	Bddd      string `json:"bddd"`        // 이사회결의일(결정일)
	OdAAtT    string `json:"od_a_at_t"`   // 사외이사 참석여부(참석)
	OdAAtB    string `json:"od_a_at_b"`   // 사외이사 참석여부(불참)
	AdtAAtn   string `json:"adt_a_atn"`   // 감사(감사위원) 참석여부
}

// BusinessSuspensions 는 영업정지를 조회한다.
func (c *Client) BusinessSuspensions(ctx context.Context, p MaterialParams) ([]BusinessSuspensionItem, error) {
	return httpclient.GetList[BusinessSuspensionItem](ctx, c.http, "/api/bsnSp.json", p.toMap())
}

// RehabilitationItem 은 회생절차 개시신청 (ctrcvsBgrq) 한 건.
type RehabilitationItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Apcnt    string `json:"apcnt"`     // 신청인 (회사와의 관계)
	Cpct     string `json:"cpct"`      // 관할법원
	RqRs     string `json:"rq_rs"`     // 신청사유
	Rqd      string `json:"rqd"`       // 신청일자
	FtCtpSc  string `json:"ft_ctp_sc"` // 향후대책 및 일정
}

// RehabilitationApplications 는 회생절차 개시신청을 조회한다.
func (c *Client) RehabilitationApplications(ctx context.Context, p MaterialParams) ([]RehabilitationItem, error) {
	return httpclient.GetList[RehabilitationItem](ctx, c.http, "/api/ctrcvsBgrq.json", p.toMap())
}

// DissolutionItem 은 해산사유 발생 (dsRsOcr) 한 건.
type DissolutionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DsRs     string `json:"ds_rs"`     // 해산사유
	DsRsd    string `json:"ds_rsd"`    // 해산사유발생일(결정일)
	OdAAtT   string `json:"od_a_at_t"` // 사외이사 참석여부(참석)
	OdAAtB   string `json:"od_a_at_b"` // 사외이사 참석여부(불참)
	AdtAAtn  string `json:"adt_a_atn"` // 감사(감사위원) 참석여부
}

// DissolutionCauses 는 해산사유 발생을 조회한다.
func (c *Client) DissolutionCauses(ctx context.Context, p MaterialParams) ([]DissolutionItem, error) {
	return httpclient.GetList[DissolutionItem](ctx, c.http, "/api/dsRsOcr.json", p.toMap())
}

// CreditorBankMgmtStartItem 은 채권은행 등의 관리절차 개시 (bnkMngtPcbg) 한 건.
type CreditorBankMgmtStartItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcbgDd string `json:"mngt_pcbg_dd"` // 관리절차개시 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	MngtPd     string `json:"mngt_pd"`      // 관리기간
	MngtRs     string `json:"mngt_rs"`      // 관리사유
	Cfd        string `json:"cfd"`          // 확인일자
}

// CreditorBankManagementStart 는 채권은행 등의 관리절차 개시를 조회한다.
func (c *Client) CreditorBankManagementStart(ctx context.Context, p MaterialParams) ([]CreditorBankMgmtStartItem, error) {
	return httpclient.GetList[CreditorBankMgmtStartItem](ctx, c.http, "/api/bnkMngtPcbg.json", p.toMap())
}

// CreditorBankMgmtStopItem 은 채권은행 등의 관리절차 중단 (bnkMngtPcsp) 한 건.
type CreditorBankMgmtStopItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcspDd string `json:"mngt_pcsp_dd"` // 관리절차중단 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	SpRs       string `json:"sp_rs"`        // 중단사유
	FtCtp      string `json:"ft_ctp"`       // 향후대책
	Cfd        string `json:"cfd"`          // 확인일자
}

// CreditorBankManagementStop 은 채권은행 등의 관리절차 중단을 조회한다.
func (c *Client) CreditorBankManagementStop(ctx context.Context, p MaterialParams) ([]CreditorBankMgmtStopItem, error) {
	return httpclient.GetList[CreditorBankMgmtStopItem](ctx, c.http, "/api/bnkMngtPcsp.json", p.toMap())
}

// LawsuitItem 은 소송 등의 제기 (lwstLg) 한 건.
type LawsuitItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Icnm     string `json:"icnm"`      // 사건의 명칭
	AcAp     string `json:"ac_ap"`     // 원고·신청인
	RqCn     string `json:"rq_cn"`     // 청구내용
	Cpct     string `json:"cpct"`      // 관할법원
	FtCtp    string `json:"ft_ctp"`    // 향후대책
	Lgd      string `json:"lgd"`       // 제기일자
	Cfd      string `json:"cfd"`       // 확인일자
}

// Lawsuits 는 소송 등의 제기를 조회한다.
func (c *Client) Lawsuits(ctx context.Context, p MaterialParams) ([]LawsuitItem, error) {
	return httpclient.GetList[LawsuitItem](ctx, c.http, "/api/lwstLg.json", p.toMap())
}
