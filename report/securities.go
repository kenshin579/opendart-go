package report

import "context"

// DebtSecuritiesIssuanceItem 은 채무증권 발행실적 (detScritsIsuAcmslt) 한 건.
type DebtSecuritiesIssuanceItem struct {
	RceptNo       string `json:"rcept_no"`       // 접수번호
	CorpCls       string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`      // 고유번호
	CorpName      string `json:"corp_name"`      // 회사명
	IsuCmpny      string `json:"isu_cmpny"`      // 발행회사
	ScritsKndNm   string `json:"scrits_knd_nm"`  // 증권종류
	IsuMthNm      string `json:"isu_mth_nm"`     // 발행방법
	IsuDe         string `json:"isu_de"`         // 발행일자
	FacvaluTotamt string `json:"facvalu_totamt"` // 권면(전자등록)총액
	Intrt         string `json:"intrt"`          // 이자율
	EvlGradInstt  string `json:"evl_grad_instt"` // 평가등급(평가기관)
	Mtd           string `json:"mtd"`            // 만기일
	RepyAt        string `json:"repy_at"`        // 상환여부
	MngtCmpny     string `json:"mngt_cmpny"`     // 주관회사
	StlmDt        string `json:"stlm_dt"`        // 결산기준일
}

// DebtSecuritiesIssuance 는 채무증권 발행실적을 조회한다.
func (c *Client) DebtSecuritiesIssuance(ctx context.Context, p ReportParams) ([]DebtSecuritiesIssuanceItem, error) {
	return getList[DebtSecuritiesIssuanceItem](ctx, c.http, "/api/detScritsIsuAcmslt.json", p)
}

// CorporateBondBalanceItem 은 회사채 미상환 잔액 (cprndNrdmpBlce) 한 건.
type CorporateBondBalanceItem struct {
	RceptNo            string `json:"rcept_no"`              // 접수번호
	CorpCls            string `json:"corp_cls"`              // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`             // 고유번호
	CorpName           string `json:"corp_name"`             // 회사명
	RemndrExprtn1      string `json:"remndr_exprtn1"`        // 잔여만기 (구분1)
	RemndrExprtn2      string `json:"remndr_exprtn2"`        // 잔여만기 (구분2)
	Yy1Below           string `json:"yy1_below"`             // 1년 이하
	Yy1ExcessYy2Below  string `json:"yy1_excess_yy2_below"`  // 1년초과 2년이하
	Yy2ExcessYy3Below  string `json:"yy2_excess_yy3_below"`  // 2년초과 3년이하
	Yy3ExcessYy4Below  string `json:"yy3_excess_yy4_below"`  // 3년초과 4년이하
	Yy4ExcessYy5Below  string `json:"yy4_excess_yy5_below"`  // 4년초과 5년이하
	Yy5ExcessYy10Below string `json:"yy5_excess_yy10_below"` // 5년초과 10년이하
	Yy10Excess         string `json:"yy10_excess"`           // 10년초과
	Sm                 string `json:"sm"`                    // 합계
	StlmDt             string `json:"stlm_dt"`               // 결산기준일
}

// CorporateBondBalance 는 회사채 미상환 잔액을 조회한다.
func (c *Client) CorporateBondBalance(ctx context.Context, p ReportParams) ([]CorporateBondBalanceItem, error) {
	return getList[CorporateBondBalanceItem](ctx, c.http, "/api/cprndNrdmpBlce.json", p)
}

// CommercialPaperBalanceItem 은 기업어음증권 미상환 잔액 (entrprsBilScritsNrdmpBlce) 한 건.
type CommercialPaperBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Yy1ExcessYy2Below    string `json:"yy1_excess_yy2_below"`    // 1년초과 2년이하
	Yy2ExcessYy3Below    string `json:"yy2_excess_yy3_below"`    // 2년초과 3년이하
	Yy3Excess            string `json:"yy3_excess"`              // 3년 초과
	Sm                   string `json:"sm"`                      // 합계
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// CommercialPaperBalance 는 기업어음증권 미상환 잔액을 조회한다.
func (c *Client) CommercialPaperBalance(ctx context.Context, p ReportParams) ([]CommercialPaperBalanceItem, error) {
	return getList[CommercialPaperBalanceItem](ctx, c.http, "/api/entrprsBilScritsNrdmpBlce.json", p)
}

// ShortTermBondBalanceItem 은 단기사채 미상환 잔액 (srtpdPsndbtNrdmpBlce) 한 건.
type ShortTermBondBalanceItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 회사명
	RemndrExprtn1        string `json:"remndr_exprtn1"`          // 잔여만기 (구분1)
	RemndrExprtn2        string `json:"remndr_exprtn2"`          // 잔여만기 (구분2)
	De10Below            string `json:"de10_below"`              // 10일 이하
	De10ExcessDe30Below  string `json:"de10_excess_de30_below"`  // 10일초과 30일이하
	De30ExcessDe90Below  string `json:"de30_excess_de90_below"`  // 30일초과 90일이하
	De90ExcessDe180Below string `json:"de90_excess_de180_below"` // 90일초과 180일이하
	De180ExcessYy1Below  string `json:"de180_excess_yy1_below"`  // 180일초과 1년이하
	Sm                   string `json:"sm"`                      // 합계
	IsuLmt               string `json:"isu_lmt"`                 // 발행 한도
	RemndrLmt            string `json:"remndr_lmt"`              // 잔여 한도
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// ShortTermBondBalance 는 단기사채 미상환 잔액을 조회한다.
func (c *Client) ShortTermBondBalance(ctx context.Context, p ReportParams) ([]ShortTermBondBalanceItem, error) {
	return getList[ShortTermBondBalanceItem](ctx, c.http, "/api/srtpdPsndbtNrdmpBlce.json", p)
}

// HybridSecuritiesBalanceItem 은 신종자본증권 미상환 잔액 (newCaplScritsNrdmpBlce) 한 건.
type HybridSecuritiesBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy5Below   string `json:"yy1_excess_yy5_below"`   // 1년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy15Below string `json:"yy10_excess_yy15_below"` // 10년초과 15년이하
	Yy15ExcessYy20Below string `json:"yy15_excess_yy20_below"` // 15년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// HybridSecuritiesBalance 는 신종자본증권 미상환 잔액을 조회한다.
func (c *Client) HybridSecuritiesBalance(ctx context.Context, p ReportParams) ([]HybridSecuritiesBalanceItem, error) {
	return getList[HybridSecuritiesBalanceItem](ctx, c.http, "/api/newCaplScritsNrdmpBlce.json", p)
}

// ContingentCapitalBalanceItem 은 조건부 자본증권 미상환 잔액 (cndlCaplScritsNrdmpBlce) 한 건.
type ContingentCapitalBalanceItem struct {
	RceptNo             string `json:"rcept_no"`               // 접수번호
	CorpCls             string `json:"corp_cls"`               // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`              // 고유번호
	CorpName            string `json:"corp_name"`              // 회사명
	RemndrExprtn1       string `json:"remndr_exprtn1"`         // 잔여만기 (구분1)
	RemndrExprtn2       string `json:"remndr_exprtn2"`         // 잔여만기 (구분2)
	Yy1Below            string `json:"yy1_below"`              // 1년 이하
	Yy1ExcessYy2Below   string `json:"yy1_excess_yy2_below"`   // 1년초과 2년이하
	Yy2ExcessYy3Below   string `json:"yy2_excess_yy3_below"`   // 2년초과 3년이하
	Yy3ExcessYy4Below   string `json:"yy3_excess_yy4_below"`   // 3년초과 4년이하
	Yy4ExcessYy5Below   string `json:"yy4_excess_yy5_below"`   // 4년초과 5년이하
	Yy5ExcessYy10Below  string `json:"yy5_excess_yy10_below"`  // 5년초과 10년이하
	Yy10ExcessYy20Below string `json:"yy10_excess_yy20_below"` // 10년초과 20년이하
	Yy20ExcessYy30Below string `json:"yy20_excess_yy30_below"` // 20년초과 30년이하
	Yy30Excess          string `json:"yy30_excess"`            // 30년초과
	Sm                  string `json:"sm"`                     // 합계
	StlmDt              string `json:"stlm_dt"`                // 결산기준일
}

// ContingentCapitalBalance 는 조건부 자본증권 미상환 잔액을 조회한다.
func (c *Client) ContingentCapitalBalance(ctx context.Context, p ReportParams) ([]ContingentCapitalBalanceItem, error) {
	return getList[ContingentCapitalBalanceItem](ctx, c.http, "/api/cndlCaplScritsNrdmpBlce.json", p)
}
