package report

import "context"

// TaxonomyItem 은 XBRL 택사노미 재무제표양식 (xbrlTaxonomy) 한 건.
type TaxonomyItem struct {
	SjDiv     string `json:"sj_div"`     // 재무제표구분
	AccountId string `json:"account_id"` // 계정ID
	AccountNm string `json:"account_nm"` // 계정명
	BsnsDe    string `json:"bsns_de"`    // 기준일 (YYYYMMDD)
	LabelKor  string `json:"label_kor"`  // 한글 출력명
	LabelEng  string `json:"label_eng"`  // 영문 출력명
	DataTp    string `json:"data_tp"`    // 데이터 유형 (일부 행에는 없음)
	IfrsRef   string `json:"ifrs_ref"`   // IFRS Reference
}

// XbrlTaxonomy 는 표준 XBRL 재무제표 택사노미 양식을 조회한다.
// sjDiv 는 재무제표구분 코드(BS1~4 재무상태표 / IS 손익계산서 / CIS 포괄손익 / CF 현금흐름표 /
// SCE 자본변동표 등). 잘못된 코드는 *opendart.APIError 로 반환된다.
func (c *Client) XbrlTaxonomy(ctx context.Context, sjDiv string) ([]TaxonomyItem, error) {
	return getListParams[TaxonomyItem](ctx, c.http, "/api/xbrlTaxonomy.json", map[string]string{"sj_div": sjDiv})
}

// DownloadXbrl 은 접수번호(rceptNo)+보고서코드로 재무제표 원본 XBRL(ZIP) 을 그대로 반환한다.
// 압축 해제·파싱은 호출자 몫.
func (c *Client) DownloadXbrl(ctx context.Context, rceptNo string, reprtCode ReportCode) ([]byte, error) {
	return c.http.GetBytes(ctx, "/api/fnlttXbrl.xml", map[string]string{
		"rcept_no":   rceptNo,
		"reprt_code": string(reprtCode),
	})
}
