package disclosure

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Company 는 기업개황 (company.json) 응답.
type Company struct {
	httpclient.Envelope
	CorpCode    string `json:"corp_code"`     // 고유번호 (8자리)
	CorpName    string `json:"corp_name"`     // 정식명칭
	CorpNameEng string `json:"corp_name_eng"` // 영문명칭
	StockName   string `json:"stock_name"`    // 종목명(상장사)/약식명칭(기타)
	StockCode   string `json:"stock_code"`    // 종목코드 (6자리)
	CeoName     string `json:"ceo_nm"`        // 대표자명
	CorpCls     string `json:"corp_cls"`      // 법인구분 Y(유가)/K(코스닥)/N(코넥스)/E(기타)
	JurirNo     string `json:"jurir_no"`      // 법인등록번호
	BizrNo      string `json:"bizr_no"`       // 사업자등록번호
	Address     string `json:"adres"`         // 주소
	HomeURL     string `json:"hm_url"`        // 홈페이지
	IRURL       string `json:"ir_url"`        // IR 홈페이지
	PhoneNo     string `json:"phn_no"`        // 전화번호
	FaxNo       string `json:"fax_no"`        // 팩스번호
	IndutyCode  string `json:"induty_code"`   // 업종코드
	EstDate     string `json:"est_dt"`        // 설립일 YYYYMMDD
	AccMonth    string `json:"acc_mt"`        // 결산월 MM
}

// GetCompany 는 corp_code(8자리)로 기업개황을 조회한다.
func (c *Client) GetCompany(ctx context.Context, corpCode string) (*Company, error) {
	var out Company
	if err := c.http.GetJSON(ctx, "/api/company.json", map[string]string{"corp_code": corpCode}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
