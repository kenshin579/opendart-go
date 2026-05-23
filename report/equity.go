package report

import "context"

// DividendItem 은 배당에 관한 사항 (alotMatter) 한 건.
type DividendItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호 (14자리)
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호 (8자리)
	CorpName string `json:"corp_name"` // 법인명
	Se       string `json:"se"`        // 구분 (주당액면가액, 주당 현금배당금 등)
	StockKnd string `json:"stock_knd"` // 주식 종류 (보통주 등)
	Thstrm   string `json:"thstrm"`    // 당기
	Frmtrm   string `json:"frmtrm"`    // 전기
	Lwfr     string `json:"lwfr"`      // 전전기
	StlmDt   string `json:"stlm_dt"`   // 결산기준일 (YYYY-MM-DD)
}

// Dividend 는 배당에 관한 사항을 조회한다.
func (c *Client) Dividend(ctx context.Context, p ReportParams) ([]DividendItem, error) {
	return getList[DividendItem](ctx, c.http, "/api/alotMatter.json", p)
}

// CapitalChangeItem 은 증자(감자) 현황 (irdsSttus) 한 건.
type CapitalChangeItem struct {
	RceptNo                 string `json:"rcept_no"`                    // 접수번호 (14자리)
	CorpCls                 string `json:"corp_cls"`                    // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                   // 고유번호 (8자리)
	CorpName                string `json:"corp_name"`                   // 법인명
	IsuDcrsDe               string `json:"isu_dcrs_de"`                 // 주식발행 감소일자
	IsuDcrsStle             string `json:"isu_dcrs_stle"`               // 발행 감소 형태
	IsuDcrsStockKnd         string `json:"isu_dcrs_stock_knd"`          // 발행 감소 주식 종류
	IsuDcrsQy               string `json:"isu_dcrs_qy"`                 // 발행 감소 수량
	IsuDcrsMstvdvFvalAmount string `json:"isu_dcrs_mstvdv_fval_amount"` // 발행 감소 주당 액면 가액
	IsuDcrsMstvdvAmount     string `json:"isu_dcrs_mstvdv_amount"`      // 발행 감소 주당 가액
	StlmDt                  string `json:"stlm_dt"`                     // 결산기준일
}

// CapitalChange 는 증자(감자) 현황을 조회한다.
func (c *Client) CapitalChange(ctx context.Context, p ReportParams) ([]CapitalChangeItem, error) {
	return getList[CapitalChangeItem](ctx, c.http, "/api/irdsSttus.json", p)
}

// TreasuryStockItem 은 자기주식 취득 및 처분 현황 (tesstkAcqsDspsSttus) 한 건.
type TreasuryStockItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	StockKnd      string `json:"stock_knd"`       // 주식 종류
	AcqsMth1      string `json:"acqs_mth1"`       // 취득방법 대분류
	AcqsMth2      string `json:"acqs_mth2"`       // 취득방법 중분류
	AcqsMth3      string `json:"acqs_mth3"`       // 취득방법 소분류
	BsisQy        string `json:"bsis_qy"`         // 기초 수량
	ChangeQyAcqs  string `json:"change_qy_acqs"`  // 변동 수량 취득
	ChangeQyDsps  string `json:"change_qy_dsps"`  // 변동 수량 처분
	ChangeQyIncnr string `json:"change_qy_incnr"` // 변동 수량 소각
	TrmendQy      string `json:"trmend_qy"`       // 기말 수량
	Rm            string `json:"rm"`              // 비고
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// TreasuryStock 은 자기주식 취득 및 처분 현황을 조회한다.
func (c *Client) TreasuryStock(ctx context.Context, p ReportParams) ([]TreasuryStockItem, error) {
	return getList[TreasuryStockItem](ctx, c.http, "/api/tesstkAcqsDspsSttus.json", p)
}

// TotalStockItem 은 주식의 총수 현황 (stockTotqySttus) 한 건.
type TotalStockItem struct {
	RceptNo             string `json:"rcept_no"`                // 접수번호
	CorpCls             string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode            string `json:"corp_code"`               // 고유번호
	CorpName            string `json:"corp_name"`               // 회사명
	Se                  string `json:"se"`                      // 구분 (보통주 등)
	IsuStockTotqy       string `json:"isu_stock_totqy"`         // 발행할 주식의 총수
	NowToIsuStockTotqy  string `json:"now_to_isu_stock_totqy"`  // 현재까지 발행한 주식의 총수
	NowToDcrsStockTotqy string `json:"now_to_dcrs_stock_totqy"` // 현재까지 감소한 주식의 총수
	Redc                string `json:"redc"`                    // 감자
	ProfitIncnr         string `json:"profit_incnr"`            // 이익소각
	RdmstkRepy          string `json:"rdmstk_repy"`             // 상환주식의 상환
	Etc                 string `json:"etc"`                     // 기타
	IstcTotqy           string `json:"istc_totqy"`              // 발행주식의 총수
	TesstkCo            string `json:"tesstk_co"`               // 자기주식수
	DistbStockCo        string `json:"distb_stock_co"`          // 유통주식수
	StlmDt              string `json:"stlm_dt"`                 // 결산기준일
}

// TotalStock 은 주식의 총수 현황을 조회한다.
func (c *Client) TotalStock(ctx context.Context, p ReportParams) ([]TotalStockItem, error) {
	return getList[TotalStockItem](ctx, c.http, "/api/stockTotqySttus.json", p)
}

// MajorShareholderItem 은 최대주주 현황 (hyslrSttus) 한 건.
type MajorShareholderItem struct {
	RceptNo                 string `json:"rcept_no"`                    // 접수번호
	CorpCls                 string `json:"corp_cls"`                    // 법인구분 (Y/K/N/E)
	CorpCode                string `json:"corp_code"`                   // 고유번호
	CorpName                string `json:"corp_name"`                   // 법인명
	Nm                      string `json:"nm"`                          // 성명
	Relate                  string `json:"relate"`                      // 관계
	StockKnd                string `json:"stock_knd"`                   // 주식 종류
	BsisPosesnStockCo       string `json:"bsis_posesn_stock_co"`        // 기초 소유 주식 수
	BsisPosesnStockQotaRt   string `json:"bsis_posesn_stock_qota_rt"`   // 기초 소유 주식 지분율
	TrmendPosesnStockCo     string `json:"trmend_posesn_stock_co"`      // 기말 소유 주식 수
	TrmendPosesnStockQotaRt string `json:"trmend_posesn_stock_qota_rt"` // 기말 소유 주식 지분율
	Rm                      string `json:"rm"`                          // 비고
	StlmDt                  string `json:"stlm_dt"`                     // 결산기준일
}

// MajorShareholders 는 최대주주 현황을 조회한다.
func (c *Client) MajorShareholders(ctx context.Context, p ReportParams) ([]MajorShareholderItem, error) {
	return getList[MajorShareholderItem](ctx, c.http, "/api/hyslrSttus.json", p)
}

// MajorShareholderChangeItem 은 최대주주 변동현황 (hyslrChgSttus) 한 건.
type MajorShareholderChangeItem struct {
	RceptNo        string `json:"rcept_no"`         // 접수번호
	CorpCls        string `json:"corp_cls"`         // 법인구분 (Y/K/N/E)
	CorpCode       string `json:"corp_code"`        // 고유번호
	CorpName       string `json:"corp_name"`        // 법인명
	ChangeOn       string `json:"change_on"`        // 변동 일
	MxmmShrholdrNm string `json:"mxmm_shrholdr_nm"` // 최대 주주 명
	PosesnStockCo  string `json:"posesn_stock_co"`  // 소유 주식 수
	QotaRt         string `json:"qota_rt"`          // 지분율
	ChangeCause    string `json:"change_cause"`     // 변동 원인
	Rm             string `json:"rm"`               // 비고
	StlmDt         string `json:"stlm_dt"`          // 결산기준일
}

// MajorShareholderChanges 는 최대주주 변동현황을 조회한다.
func (c *Client) MajorShareholderChanges(ctx context.Context, p ReportParams) ([]MajorShareholderChangeItem, error) {
	return getList[MajorShareholderChangeItem](ctx, c.http, "/api/hyslrChgSttus.json", p)
}

// MinorityShareholderItem 은 소액주주 현황 (mrhlSttus) 한 건.
type MinorityShareholderItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 법인명
	Se            string `json:"se"`              // 구분
	ShrholdrCo    string `json:"shrholdr_co"`     // 주주수
	ShrholdrTotCo string `json:"shrholdr_tot_co"` // 전체 주주수
	ShrholdrRate  string `json:"shrholdr_rate"`   // 주주 비율
	HoldStockCo   string `json:"hold_stock_co"`   // 보유 주식수
	StockTotCo    string `json:"stock_tot_co"`    // 총발행 주식수
	HoldStockRate string `json:"hold_stock_rate"` // 보유 주식 비율
	StlmDt        string `json:"stlm_dt"`         // 결산기준일
}

// MinorityShareholders 는 소액주주 현황을 조회한다.
func (c *Client) MinorityShareholders(ctx context.Context, p ReportParams) ([]MinorityShareholderItem, error) {
	return getList[MinorityShareholderItem](ctx, c.http, "/api/mrhlSttus.json", p)
}
