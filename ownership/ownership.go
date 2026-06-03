package ownership

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// MajorStockItem 은 대량보유 상황보고 (majorstock) 한 건.
type MajorStockItem struct {
	RceptNo    string `json:"rcept_no"`    // 접수번호
	RceptDt    string `json:"rcept_dt"`    // 접수일자
	CorpCode   string `json:"corp_code"`   // 고유번호
	CorpName   string `json:"corp_name"`   // 회사명
	ReportTp   string `json:"report_tp"`   // 보고구분
	Repror     string `json:"repror"`      // 대표보고자
	Stkqy      string `json:"stkqy"`       // 보유주식등의 수
	StkqyIrds  string `json:"stkqy_irds"`  // 보유주식등의 증감
	Stkrt      string `json:"stkrt"`       // 보유비율
	StkrtIrds  string `json:"stkrt_irds"`  // 보유비율 증감
	CtrStkqy   string `json:"ctr_stkqy"`   // 주요체결 주식등의 수
	CtrStkrt   string `json:"ctr_stkrt"`   // 주요체결 보유비율
	ReportResn string `json:"report_resn"` // 보고사유
}

// MajorStockReports 는 대량보유 상황보고(5% 룰)를 조회한다.
func (c *Client) MajorStockReports(ctx context.Context, corpCode string) ([]MajorStockItem, error) {
	return httpclient.GetList[MajorStockItem](ctx, c.http, "/api/majorstock.json", map[string]string{"corp_code": corpCode})
}

// ExecutiveStockItem 은 임원·주요주주 소유보고 (elestock) 한 건.
type ExecutiveStockItem struct {
	RceptNo            string `json:"rcept_no"`               // 접수번호
	RceptDt            string `json:"rcept_dt"`               // 접수일자
	CorpCode           string `json:"corp_code"`              // 고유번호
	CorpName           string `json:"corp_name"`              // 회사명
	Repror             string `json:"repror"`                 // 보고자
	IsuExctvRgistAt    string `json:"isu_exctv_rgist_at"`     // 발행회사 관계 임원(등기여부)
	IsuExctvOfcps      string `json:"isu_exctv_ofcps"`        // 발행회사 관계 임원 직위
	IsuMainShrholdr    string `json:"isu_main_shrholdr"`      // 발행회사 관계 주요주주
	SpStockLmpCnt      string `json:"sp_stock_lmp_cnt"`       // 특정증권등 소유 수
	SpStockLmpIrdsCnt  string `json:"sp_stock_lmp_irds_cnt"`  // 특정증권등 소유 증감 수
	SpStockLmpRate     string `json:"sp_stock_lmp_rate"`      // 특정증권등 소유 비율
	SpStockLmpIrdsRate string `json:"sp_stock_lmp_irds_rate"` // 특정증권등 소유 증감 비율
}

// ExecutiveStockReports 는 임원·주요주주 소유보고를 조회한다.
func (c *Client) ExecutiveStockReports(ctx context.Context, corpCode string) ([]ExecutiveStockItem, error) {
	return httpclient.GetList[ExecutiveStockItem](ctx, c.http, "/api/elestock.json", map[string]string{"corp_code": corpCode})
}
