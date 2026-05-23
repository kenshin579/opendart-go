package disclosure

import (
	"context"
	"strconv"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// SearchParams 는 공시검색 (list.json) 요청 인자. 빈 값/0 은 쿼리에서 생략한다.
type SearchParams struct {
	CorpCode       string // 고유번호 (8자리)
	BgnDe          string // 검색 시작 접수일자 YYYYMMDD
	EndDe          string // 검색 종료 접수일자 YYYYMMDD
	LastReprtAt    string // 최종보고서만 검색 여부 (Y/N, 기본 N)
	PblntfTy       string // 공시유형 (A 정기/B 주요사항/C 발행/D 지분/E 기타/F 외부감사/G 펀드/H 자산유동화/I 거래소/J 공정위)
	PblntfDetailTy string // 공시상세유형 (4자리 코드)
	CorpCls        string // 법인구분 (Y/K/N/E)
	Sort           string // 정렬 (date/crp/rpt, 기본 date)
	SortMth        string // 정렬방법 (asc/desc, 기본 desc)
	PageNo         int    // 페이지 번호 (1~n, 0이면 생략)
	PageCount      int    // 페이지당 건수 (1~100, 0이면 생략)
}

// toMap 은 빈 값/0 을 제외한 쿼리 파라미터 맵을 만든다.
func (p SearchParams) toMap() map[string]string {
	m := map[string]string{}
	put := func(k, v string) {
		if v != "" {
			m[k] = v
		}
	}
	put("corp_code", p.CorpCode)
	put("bgn_de", p.BgnDe)
	put("end_de", p.EndDe)
	put("last_reprt_at", p.LastReprtAt)
	put("pblntf_ty", p.PblntfTy)
	put("pblntf_detail_ty", p.PblntfDetailTy)
	put("corp_cls", p.CorpCls)
	put("sort", p.Sort)
	put("sort_mth", p.SortMth)
	if p.PageNo > 0 {
		m["page_no"] = strconv.Itoa(p.PageNo)
	}
	if p.PageCount > 0 {
		m["page_count"] = strconv.Itoa(p.PageCount)
	}
	return m
}

// DisclosureItem 은 공시검색 결과 한 건.
type DisclosureItem struct {
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpName  string `json:"corp_name"`  // 종목명(법인명)
	CorpCode  string `json:"corp_code"`  // 고유번호 (8자리)
	StockCode string `json:"stock_code"` // 종목코드 (6자리)
	ReportNm  string `json:"report_nm"`  // 보고서명
	RceptNo   string `json:"rcept_no"`   // 접수번호 (DownloadDocument 인자)
	FlrNm     string `json:"flr_nm"`     // 공시 제출인명
	RceptDt   string `json:"rcept_dt"`   // 접수일자 YYYYMMDD
	Rm        string `json:"rm"`         // 비고
}

// SearchResult 는 공시검색 응답 (페이지네이션 + 목록).
type SearchResult struct {
	httpclient.Envelope
	PageNo     int              `json:"page_no"`     // 페이지 번호
	PageCount  int              `json:"page_count"`  // 페이지당 건수
	TotalCount int              `json:"total_count"` // 총 건수
	TotalPage  int              `json:"total_page"`  // 총 페이지 수
	List       []DisclosureItem `json:"list"`        // 공시 목록
}

// SearchDisclosures 는 조건별 공시보고서를 검색한다.
func (c *Client) SearchDisclosures(ctx context.Context, params SearchParams) (*SearchResult, error) {
	var out SearchResult
	if err := c.http.GetJSON(ctx, "/api/list.json", params.toMap(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}
