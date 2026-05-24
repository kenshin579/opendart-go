package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
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
