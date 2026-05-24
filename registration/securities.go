package registration

import (
	"context"
	"encoding/json"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// RsGeneralItem 은 증권신고서 일반사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsGeneralItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Sbd      string `json:"sbd"`       // 청약기일
	Pymd     string `json:"pymd"`      // 납입기일
	Sband    string `json:"sband"`     // 청약공고일
	Asand    string `json:"asand"`     // 배정공고일
	Asstd    string `json:"asstd"`     // 배정기준일
	Exstk    string `json:"exstk"`     // 신주인수권에 관한 사항(행사대상증권)
	Exprc    string `json:"exprc"`     // 신주인수권에 관한 사항(행사가격)
	Expd     string `json:"expd"`      // 신주인수권에 관한 사항(행사기간)
	RptRcpn  string `json:"rpt_rcpn"`  // 주요사항보고서(접수번호)
}

// RsSecurityTypeItem 은 증권신고서 증권의종류 그룹 항목(지분증권/증권예탁증권 공통).
type RsSecurityTypeItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Stksen   string `json:"stksen"`    // 증권의종류
	Stkcnt   string `json:"stkcnt"`    // 증권수량
	Fv       string `json:"fv"`        // 액면가액
	Slprc    string `json:"slprc"`     // 모집(매출)가액
	Slta     string `json:"slta"`      // 모집(매출)총액
	Slmthn   string `json:"slmthn"`    // 모집(매출)방법
}

// RsUnderwriterItem 은 증권신고서 인수인정보 그룹 항목(지분증권/증권예탁증권 공통).
type RsUnderwriterItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Actsen   string `json:"actsen"`    // 인수인구분
	Actnmn   string `json:"actnmn"`    // 인수인명
	Stksen   string `json:"stksen"`    // 증권의종류
	Udtcnt   string `json:"udtcnt"`    // 인수수량
	Udtamt   string `json:"udtamt"`    // 인수금액
	Udtprc   string `json:"udtprc"`    // 인수대가
	Udtmth   string `json:"udtmth"`    // 인수방법
}

// RsFundUsageItem 은 증권신고서 자금의사용목적 그룹 항목(지분증권/증권예탁증권 공통).
type RsFundUsageItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Se       string `json:"se"`        // 구분
	Amt      string `json:"amt"`       // 금액
}

// RsSellerItem 은 증권신고서 매출인에관한사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsSellerItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Hdr       string `json:"hdr"`        // 보유자
	RlCmp     string `json:"rl_cmp"`     // 회사와의관계
	BfslHdstk string `json:"bfsl_hdstk"` // 매출전보유증권수
	Slstk     string `json:"slstk"`      // 매출증권수
	AtslHdstk string `json:"atsl_hdstk"` // 매출후보유증권수
}

// EquityRetailPutbackOptionItem 은 지분증권 일반청약자환매청구권 그룹 항목(지분증권 전용).
type EquityRetailPutbackOptionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Grtrs    string `json:"grtrs"`     // 부여사유
	Exavivr  string `json:"exavivr"`   // 행사가능 투자자
	Grtcnt   string `json:"grtcnt"`    // 부여수량
	Expd     string `json:"expd"`      // 행사기간
	Exprc    string `json:"exprc"`     // 행사가격
}

// EquitySecuritiesRegistration 은 지분증권 증권신고서(estkRs)의 그룹별 항목.
type EquitySecuritiesRegistration struct {
	General             []RsGeneralItem                 // 일반사항
	SecurityTypes       []RsSecurityTypeItem            // 증권의종류
	Underwriters        []RsUnderwriterItem             // 인수인정보
	FundUsage           []RsFundUsageItem               // 자금의사용목적
	Sellers             []RsSellerItem                  // 매출인에관한사항
	RetailPutbackOption []EquityRetailPutbackOptionItem // 일반청약자환매청구권
}

// DebtGeneralItem 은 채무증권 증권신고서 일반사항 그룹 항목.
type DebtGeneralItem struct {
	RceptNo     string `json:"rcept_no"`     // 접수번호
	CorpCls     string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode    string `json:"corp_code"`    // 고유번호
	CorpName    string `json:"corp_name"`    // 회사명
	Tm          string `json:"tm"`           // 회차
	Bdnmn       string `json:"bdnmn"`        // 채무증권 명칭
	Slmth       string `json:"slmth"`        // 모집(매출)방법
	Fta         string `json:"fta"`          // 권면(전자등록)총액
	Slta        string `json:"slta"`         // 모집(매출)총액
	Isprc       string `json:"isprc"`        // 발행가액
	Intr        string `json:"intr"`         // 이자율
	Isrr        string `json:"isrr"`         // 발행수익률
	Rpd         string `json:"rpd"`          // 상환기일
	PrintPymint string `json:"print_pymint"` // 원리금지급대행기관
	MngtCmp     string `json:"mngt_cmp"`     // (사채)관리회사
	CdrtInt     string `json:"cdrt_int"`     // 신용등급(신용평가기관)
	Sbd         string `json:"sbd"`          // 청약기일
	Pymd        string `json:"pymd"`         // 납입기일
	Sband       string `json:"sband"`        // 청약공고일
	Asand       string `json:"asand"`        // 배정공고일
	Asstd       string `json:"asstd"`        // 배정기준일
	Dpcrn       string `json:"dpcrn"`        // 표시통화
	DpcrAmt     string `json:"dpcr_amt"`     // 표시통화기준발행규모
	Usarn       string `json:"usarn"`        // 사용지역
	Usntn       string `json:"usntn"`        // 사용국가
	WnexplAt    string `json:"wnexpl_at"`    // 원화 교환 예정 여부
	Udtintnm    string `json:"udtintnm"`     // 인수기관명
	GrtInt      string `json:"grt_int"`      // 보증을 받은 경우(보증기관)
	GrtAmt      string `json:"grt_amt"`      // 보증을 받은 경우(보증금액)
	IcmgMgknd   string `json:"icmg_mgknd"`   // 담보 제공의 경우(담보의 종류)
	IcmgMgamt   string `json:"icmg_mgamt"`   // 담보 제공의 경우(담보금액)
	EstkExstk   string `json:"estk_exstk"`   // 지분증권과 연계된 경우(행사대상증권)
	EstkExrt    string `json:"estk_exrt"`    // 지분증권과 연계된 경우(권리행사비율)
	EstkExprc   string `json:"estk_exprc"`   // 지분증권과 연계된 경우(권리행사가격)
	EstkExpd    string `json:"estk_expd"`    // 지분증권과 연계된 경우(권리행사기간)
	RptRcpn     string `json:"rpt_rcpn"`     // 주요사항보고서(접수번호)
	DrcbAt      string `json:"drcb_at"`      // 파생결합사채해당여부
	DrcbUast    string `json:"drcb_uast"`    // 파생결합사채(기초자산)
	DrcbOptknd  string `json:"drcb_optknd"`  // 파생결합사채(옵션종류)
	DrcbMtd     string `json:"drcb_mtd"`     // 파생결합사채(만기일)
}

// DebtUnderwriterItem 은 채무증권 증권신고서 인수인정보 그룹 항목.
type DebtUnderwriterItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Tm       string `json:"tm"`        // 회차
	Actsen   string `json:"actsen"`    // 인수인구분
	Actnmn   string `json:"actnmn"`    // 인수인명
	Stksen   string `json:"stksen"`    // 증권의종류
	Udtcnt   string `json:"udtcnt"`    // 인수수량
	Udtamt   string `json:"udtamt"`    // 인수금액
	Udtprc   string `json:"udtprc"`    // 인수대가
	Udtmth   string `json:"udtmth"`    // 인수방법
}

// DebtFundUsageItem 은 채무증권 증권신고서 자금의사용목적 그룹 항목.
type DebtFundUsageItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Tm       string `json:"tm"`        // 회차
	Se       string `json:"se"`        // 구분
	Amt      string `json:"amt"`       // 금액
}

// DebtSellerItem 은 채무증권 증권신고서 매출인에관한사항 그룹 항목.
type DebtSellerItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Tm        string `json:"tm"`         // 회차
	Hdr       string `json:"hdr"`        // 보유자
	RlCmp     string `json:"rl_cmp"`     // 회사와의관계
	BfslHdstk string `json:"bfsl_hdstk"` // 매출전보유증권수
	Slstk     string `json:"slstk"`      // 매출증권수
	AtslHdstk string `json:"atsl_hdstk"` // 매출후보유증권수
}

// DebtSecuritiesRegistration 은 채무증권 증권신고서(bdRs)의 그룹별 항목.
type DebtSecuritiesRegistration struct {
	General      []DebtGeneralItem     // 일반사항
	Underwriters []DebtUnderwriterItem // 인수인정보
	FundUsage    []DebtFundUsageItem   // 자금의사용목적
	Sellers      []DebtSellerItem      // 매출인에관한사항
}

// DebtSecurities 는 채무증권 증권신고서(DS006)를 조회한다.
func (c *Client) DebtSecurities(ctx context.Context, p Params) (*DebtSecuritiesRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/bdRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DebtSecuritiesRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// DepositaryReceiptsRegistration 은 증권예탁증권 증권신고서(stkdpRs)의 그룹별 항목.
type DepositaryReceiptsRegistration struct {
	General       []RsGeneralItem      // 일반사항
	SecurityTypes []RsSecurityTypeItem // 증권의종류
	Underwriters  []RsUnderwriterItem  // 인수인정보
	FundUsage     []RsFundUsageItem    // 자금의사용목적
	Sellers       []RsSellerItem       // 매출인에관한사항
}

// DepositaryReceipts 는 증권예탁증권 증권신고서(DS006)를 조회한다.
func (c *Client) DepositaryReceipts(ctx context.Context, p Params) (*DepositaryReceiptsRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/stkdpRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DepositaryReceiptsRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "증권의종류":
			derr = json.Unmarshal(g.List, &out.SecurityTypes)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// EquitySecurities 는 지분증권 증권신고서(DS006)를 조회한다.
func (c *Client) EquitySecurities(ctx context.Context, p Params) (*EquitySecuritiesRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/estkRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &EquitySecuritiesRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "증권의종류":
			derr = json.Unmarshal(g.List, &out.SecurityTypes)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		case "일반청약자환매청구권":
			derr = json.Unmarshal(g.List, &out.RetailPutbackOption)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
