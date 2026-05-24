package report

import "context"

// ExecutiveItem 은 임원 현황 (exctvSttus) 한 건.
type ExecutiveItem struct {
	RceptNo            string `json:"rcept_no"`             // 접수번호
	CorpCls            string `json:"corp_cls"`             // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`            // 고유번호
	CorpName           string `json:"corp_name"`            // 법인명
	Nm                 string `json:"nm"`                   // 성명
	Sexdstn            string `json:"sexdstn"`              // 성별
	BirthYm            string `json:"birth_ym"`             // 출생 년월
	Ofcps              string `json:"ofcps"`                // 직위
	RgistExctvAt       string `json:"rgist_exctv_at"`       // 등기 임원 여부
	FteAt              string `json:"fte_at"`               // 상근 여부
	ChrgJob            string `json:"chrg_job"`             // 담당 업무
	MainCareer         string `json:"main_career"`          // 주요 경력
	MxmmShrholdrRelate string `json:"mxmm_shrholdr_relate"` // 최대 주주 관계
	HffcPd             string `json:"hffc_pd"`              // 재직 기간
	TenureEndOn        string `json:"tenure_end_on"`        // 임기 만료 일
	StlmDt             string `json:"stlm_dt"`              // 결산기준일
}

// Executives 는 임원 현황을 조회한다.
func (c *Client) Executives(ctx context.Context, p ReportParams) ([]ExecutiveItem, error) {
	return getList[ExecutiveItem](ctx, c.http, "/api/exctvSttus.json", p)
}

// EmployeeItem 은 직원 현황 (empSttus) 한 건.
type EmployeeItem struct {
	RceptNo              string `json:"rcept_no"`                 // 접수번호
	CorpCls              string `json:"corp_cls"`                 // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`                // 고유번호
	CorpName             string `json:"corp_name"`                // 법인명
	FoBbm                string `json:"fo_bbm"`                   // 사업부문
	Sexdstn              string `json:"sexdstn"`                  // 성별
	ReformBfeEmpCoRgllbr string `json:"reform_bfe_emp_co_rgllbr"` // 개정 전 직원 수 정규직
	ReformBfeEmpCoCnttk  string `json:"reform_bfe_emp_co_cnttk"`  // 개정 전 직원 수 계약직
	ReformBfeEmpCoEtc    string `json:"reform_bfe_emp_co_etc"`    // 개정 전 직원 수 기타
	RgllbrCo             string `json:"rgllbr_co"`                // 정규직 수
	RgllbrAbacptLabrrCo  string `json:"rgllbr_abacpt_labrr_co"`   // 정규직 단시간 근로자 수
	CnttkCo              string `json:"cnttk_co"`                 // 계약직 수
	CnttkAbacptLabrrCo   string `json:"cnttk_abacpt_labrr_co"`    // 계약직 단시간 근로자 수
	Sm                   string `json:"sm"`                       // 합계
	AvrgCnwkSdytrn       string `json:"avrg_cnwk_sdytrn"`         // 평균 근속 연수
	FyerSalaryTotamt     string `json:"fyer_salary_totamt"`       // 연간 급여 총액
	JanSalaryAm          string `json:"jan_salary_am"`            // 1인평균 급여 액
	Rm                   string `json:"rm"`                       // 비고
	StlmDt               string `json:"stlm_dt"`                  // 결산기준일
}

// Employees 는 직원 현황을 조회한다.
func (c *Client) Employees(ctx context.Context, p ReportParams) ([]EmployeeItem, error) {
	return getList[EmployeeItem](ctx, c.http, "/api/empSttus.json", p)
}

// UnregisteredExecutiveCompensationItem 은 미등기임원 보수현황 (unrstExctvMendngSttus) 한 건.
type UnregisteredExecutiveCompensationItem struct {
	RceptNo          string `json:"rcept_no"`            // 접수번호
	CorpCls          string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`           // 고유번호
	CorpName         string `json:"corp_name"`           // 회사명
	Se               string `json:"se"`                  // 구분
	FyerSalaryTotamt string `json:"fyer_salary_totamt"`  // 연간급여 총액
	JanSalaryAm      string `json:"jan_salary_am"`       // 1인평균 급여액
	Nmpr             string `json:"nmpr"`                // 인원수
	Rm               string `json:"rm"`                  // 비고
	StlmDt           string `json:"stlm_dt"`             // 결산기준일
}

// UnregisteredExecutiveCompensation 은 미등기임원 보수현황을 조회한다.
func (c *Client) UnregisteredExecutiveCompensation(ctx context.Context, p ReportParams) ([]UnregisteredExecutiveCompensationItem, error) {
	return getList[UnregisteredExecutiveCompensationItem](ctx, c.http, "/api/unrstExctvMendngSttus.json", p)
}

// OutsideDirectorChangeItem 은 사외이사 및 그 변동현황 (outcmpnyDrctrNdChangeSttus) 한 건.
type OutsideDirectorChangeItem struct {
	RceptNo      string `json:"rcept_no"`        // 접수번호
	CorpCls      string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`       // 고유번호
	CorpName     string `json:"corp_name"`       // 회사명
	DrctrCo      string `json:"drctr_co"`        // 이사의 수
	OtcmpDrctrCo string `json:"otcmp_drctr_co"`  // 사외이사 수
	Apnt         string `json:"apnt"`            // 사외이사 변동현황(선임)
	Rlsofc       string `json:"rlsofc"`          // 사외이사 변동현황(해임)
	MdstrmResig  string `json:"mdstrm_resig"`    // 사외이사 변동현황(중도퇴임)
	StlmDt       string `json:"stlm_dt"`         // 결산기준일
}

// OutsideDirectorChanges 는 사외이사 및 그 변동현황을 조회한다.
func (c *Client) OutsideDirectorChanges(ctx context.Context, p ReportParams) ([]OutsideDirectorChangeItem, error) {
	return getList[OutsideDirectorChangeItem](ctx, c.http, "/api/outcmpnyDrctrNdChangeSttus.json", p)
}
