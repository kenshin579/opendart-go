package report

import "context"

// AccountItem 은 단일/다중회사 주요계정 (fnlttSinglAcnt / fnlttMultiAcnt) 한 건.
type AccountItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	CorpCode        string `json:"corp_code"`         // 고유번호
	StockCode       string `json:"stock_code"`        // 종목 코드
	FsDiv           string `json:"fs_div"`            // 개별/연결구분
	FsNm            string `json:"fs_nm"`             // 개별/연결명
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	AccountNm       string `json:"account_nm"`        // 계정명
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmDt        string `json:"thstrm_dt"`         // 당기일자
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmDt        string `json:"frmtrm_dt"`         // 전기일자
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmDt     string `json:"bfefrmtrm_dt"`      // 전전기일자
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}

// SingleAccount 는 단일회사 주요계정을 조회한다.
func (c *Client) SingleAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttSinglAcnt.json", p)
}

// MultiAccount 는 다중회사 주요계정을 조회한다. p.CorpCode 는 콤마로 여러 고유번호를 전달한다.
func (c *Client) MultiAccount(ctx context.Context, p ReportParams) ([]AccountItem, error) {
	return getList[AccountItem](ctx, c.http, "/api/fnlttMultiAcnt.json", p)
}

// FsDiv 는 개별/연결 구분.
type FsDiv string

const (
	FsDivSeparate     FsDiv = "OFS" // 재무제표(개별)
	FsDivConsolidated FsDiv = "CFS" // 연결재무제표
)

// IndexClass 는 재무지표 분류 코드.
type IndexClass string

const (
	IndexProfitability IndexClass = "M210000" // 수익성지표
	IndexStability     IndexClass = "M220000" // 안정성지표
	IndexGrowth        IndexClass = "M230000" // 성장성지표
	IndexActivity      IndexClass = "M240000" // 활동성지표
)

// FinancialStatementParams 는 전체 재무제표 요청 인자.
type FinancialStatementParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
	FsDiv     FsDiv      // 개별(OFS)/연결(CFS)
}

func (p FinancialStatementParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
		"fs_div":     string(p.FsDiv),
	}
}

// FinancialIndexParams 는 주요 재무지표 요청 인자. 다중회사는 CorpCode 를 콤마로 구분한다.
type FinancialIndexParams struct {
	CorpCode  string     // 고유번호 (8자리; 다중은 콤마 구분)
	BsnsYear  string     // 사업연도 (4자리)
	ReprtCode ReportCode // 보고서 코드
	IdxClCode IndexClass // 지표분류코드
}

func (p FinancialIndexParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":   p.CorpCode,
		"bsns_year":   p.BsnsYear,
		"reprt_code":  string(p.ReprtCode),
		"idx_cl_code": string(p.IdxClCode),
	}
}

// FullStatementItem 은 단일회사 전체 재무제표 (fnlttSinglAcntAll) 한 건.
type FullStatementItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	ReprtCode       string `json:"reprt_code"`        // 보고서 코드
	BsnsYear        string `json:"bsns_year"`         // 사업 연도
	CorpCode        string `json:"corp_code"`         // 고유번호
	SjDiv           string `json:"sj_div"`            // 재무제표구분
	SjNm            string `json:"sj_nm"`             // 재무제표명
	AccountId       string `json:"account_id"`        // 계정ID
	AccountNm       string `json:"account_nm"`        // 계정명
	AccountDetail   string `json:"account_detail"`    // 계정상세
	ThstrmNm        string `json:"thstrm_nm"`         // 당기명
	ThstrmAmount    string `json:"thstrm_amount"`     // 당기금액
	ThstrmAddAmount string `json:"thstrm_add_amount"` // 당기누적금액
	FrmtrmNm        string `json:"frmtrm_nm"`         // 전기명
	FrmtrmAmount    string `json:"frmtrm_amount"`     // 전기금액
	FrmtrmQNm       string `json:"frmtrm_q_nm"`       // 전기명(분/반기)
	FrmtrmQAmount   string `json:"frmtrm_q_amount"`   // 전기금액(분/반기)
	FrmtrmAddAmount string `json:"frmtrm_add_amount"` // 전기누적금액
	BfefrmtrmNm     string `json:"bfefrmtrm_nm"`      // 전전기명
	BfefrmtrmAmount string `json:"bfefrmtrm_amount"`  // 전전기금액
	Ord             string `json:"ord"`               // 계정과목 정렬순서
	Currency        string `json:"currency"`          // 통화 단위
}

// SingleFullStatement 는 단일회사 전체 재무제표를 조회한다.
func (c *Client) SingleFullStatement(ctx context.Context, p FinancialStatementParams) ([]FullStatementItem, error) {
	return getListParams[FullStatementItem](ctx, c.http, "/api/fnlttSinglAcntAll.json", p.toMap())
}

// FinancialIndexItem 은 단일/다중회사 주요 재무지표 (fnlttSinglIndx / fnlttCmpnyIndx) 한 건.
type FinancialIndexItem struct {
	ReprtCode string `json:"reprt_code"`  // 보고서 코드
	BsnsYear  string `json:"bsns_year"`   // 사업 연도
	CorpCode  string `json:"corp_code"`   // 고유번호
	StockCode string `json:"stock_code"`  // 종목 코드
	StlmDt    string `json:"stlm_dt"`     // 결산기준일
	IdxClCode string `json:"idx_cl_code"` // 지표분류코드
	IdxClNm   string `json:"idx_cl_nm"`   // 지표분류명
	IdxCode   string `json:"idx_code"`    // 지표코드
	IdxNm     string `json:"idx_nm"`      // 지표명
	IdxVal    string `json:"idx_val"`     // 지표값
}

// SingleIndex 는 단일회사 주요 재무지표를 조회한다.
func (c *Client) SingleIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttSinglIndx.json", p.toMap())
}

// MultiIndex 는 다중회사 주요 재무지표를 조회한다. p.CorpCode 는 콤마로 여러 고유번호를 전달한다.
func (c *Client) MultiIndex(ctx context.Context, p FinancialIndexParams) ([]FinancialIndexItem, error) {
	return getListParams[FinancialIndexItem](ctx, c.http, "/api/fnlttCmpnyIndx.json", p.toMap())
}
