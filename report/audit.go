package report

import "context"

// AuditOpinionItem 은 회계감사인의 명칭 및 감사의견 (accnutAdtorNmNdAdtOpinion) 한 건.
type AuditOpinionItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	BsnsYear             string `json:"bsns_year"`               // 사업연도
	Adtor                string `json:"adtor"`                   // 감사인
	AdtOpinion           string `json:"adt_opinion"`             // 감사의견
	AdtReprtSpcmntMatter string `json:"adt_reprt_spcmnt_matter"` // 감사보고서 특기사항
	EmphsMatter          string `json:"emphs_matter"`            // 강조사항 등
	CoreAdtMatter        string `json:"core_adt_matter"`         // 핵심감사사항
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// AuditOpinion 은 회계감사인의 명칭 및 감사의견을 조회한다.
func (c *Client) AuditOpinion(ctx context.Context, p ReportParams) ([]AuditOpinionItem, error) {
	return getList[AuditOpinionItem](ctx, c.http, "/api/accnutAdtorNmNdAdtOpinion.json", p)
}

// AuditServiceContractItem 은 감사용역체결현황 (adtServcCnclsSttus) 한 건.
type AuditServiceContractItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	BsnsYear            string `json:"bsns_year"`              // 사업연도
	Adtor               string `json:"adtor"`                  // 감사인
	Cn                  string `json:"cn"`                     // 내용
	Mendng              string `json:"mendng"`                 // 보수
	TotReqreTime        string `json:"tot_reqre_time"`         // 총소요시간
	AdtCntrctDtlsMendng string `json:"adt_cntrct_dtls_mendng"` // 감사계약내역(보수)
	AdtCntrctDtlsTime   string `json:"adt_cntrct_dtls_time"`   // 감사계약내역(시간)
	RealExcDtlsMendng   string `json:"real_exc_dtls_mendng"`   // 실제수행내역(보수)
	RealExcDtlsTime     string `json:"real_exc_dtls_time"`     // 실제수행내역(시간)
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// AuditServiceContract 는 감사용역체결현황을 조회한다.
func (c *Client) AuditServiceContract(ctx context.Context, p ReportParams) ([]AuditServiceContractItem, error) {
	return getList[AuditServiceContractItem](ctx, c.http, "/api/adtServcCnclsSttus.json", p)
}

// NonAuditServiceContractItem 은 회계감사인과의 비감사용역 계약체결 현황 (accnutAdtorNonAdtServcCnclsSttus) 한 건.
type NonAuditServiceContractItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	BsnsYear      string `json:"bsns_year"`       // 사업연도
	CntrctCnclsDe string `json:"cntrct_cncls_de"` // 계약체결일
	ServcCn       string `json:"servc_cn"`        // 용역내용
	ServcExcPd    string `json:"servc_exc_pd"`    // 용역수행기간
	ServcMendng   string `json:"servc_mendng"`    // 용역보수
	Rm            string `json:"rm"`              // 비고
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// NonAuditServiceContract 는 회계감사인과의 비감사용역 계약체결 현황을 조회한다.
func (c *Client) NonAuditServiceContract(ctx context.Context, p ReportParams) ([]NonAuditServiceContractItem, error) {
	return getList[NonAuditServiceContractItem](ctx, c.http, "/api/accnutAdtorNonAdtServcCnclsSttus.json", p)
}
