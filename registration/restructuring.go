package registration

import (
	"context"
	"encoding/json"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// RestructuringGeneralItem 은 증권신고서(합병/분할/교환이전) 일반사항 그룹 항목.
type RestructuringGeneralItem struct {
	RceptNo       string `json:"rcept_no"`       // 접수번호
	CorpCls       string `json:"corp_cls"`       // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`      // 고유번호
	CorpName      string `json:"corp_name"`      // 회사명
	Stn           string `json:"stn"`            // 형태
	Bddd          string `json:"bddd"`           // 이사회 결의일
	Ctrd          string `json:"ctrd"`           // 계약일
	GmtsckShddstd string `json:"gmtsck_shddstd"` // 주주총회를 위한 주주확정일
	ApGmtsck      string `json:"ap_gmtsck"`      // 승인을 위한 주주총회일
	AprskhPdBgd   string `json:"aprskh_pd_bgd"`  // 주식매수청구권 행사 기간 및 가격(시작일)
	AprskhPdEdd   string `json:"aprskh_pd_edd"`  // 주식매수청구권 행사 기간 및 가격(종료일)
	AprskhPrc     string `json:"aprskh_prc"`     // 주식매수청구권 행사 기간 및 가격(주식매수청구가격-회사제시)
	MgdtEtc       string `json:"mgdt_etc"`       // 합병기일등
	RtVl          string `json:"rt_vl"`          // 비율 또는 가액
	ExevlInt      string `json:"exevl_int"`      // 외부평가기관
	GrtmnEtc      string `json:"grtmn_etc"`      // 지급 교부금 등
	RptRcpn       string `json:"rpt_rcpn"`       // 주요사항보고서(접수번호)
}

// RestructuringIssuedSecurityItem 은 증권신고서(합병/분할/교환이전) 발행증권 그룹 항목.
type RestructuringIssuedSecurityItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Kndn     string `json:"kndn"`      // 종류
	Cnt      string `json:"cnt"`       // 수량
	Fv       string `json:"fv"`        // 액면가액
	Slprc    string `json:"slprc"`     // 모집(매출)가액
	Slta     string `json:"slta"`      // 모집(매출)총액
}

// RestructuringPartyCompanyItem 은 증권신고서(합병/분할/교환이전) 당사회사에관한사항 그룹 항목.
type RestructuringPartyCompanyItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Cmpnm    string `json:"cmpnm"`     // 회사명
	Sen      string `json:"sen"`       // 구분
	Tast     string `json:"tast"`      // 총자산
	Cpt      string `json:"cpt"`       // 자본금
	IsstkKnd string `json:"isstk_knd"` // 발행주식수(주식의종류)
	IsstkCnt string `json:"isstk_cnt"` // 발행주식수(주식수)
}

// MergerRegistration 은 합병 증권신고서(mgRs)의 그룹별 항목.
type MergerRegistration struct {
	General          []RestructuringGeneralItem        // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem   // 당사회사에관한사항
}

// DivisionRegistration 은 분할 증권신고서(dvRs)의 그룹별 항목.
type DivisionRegistration struct {
	General          []RestructuringGeneralItem        // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem   // 당사회사에관한사항
}

// Division 은 분할 증권신고서(DS006)를 조회한다.
func (c *Client) Division(ctx context.Context, p Params) (*DivisionRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/dvRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DivisionRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// StockExchangeTransferRegistration 은 주식의포괄적교환·이전 증권신고서(extrRs)의 그룹별 항목.
type StockExchangeTransferRegistration struct {
	General          []RestructuringGeneralItem        // 일반사항
	IssuedSecurities []RestructuringIssuedSecurityItem // 발행증권
	PartyCompanies   []RestructuringPartyCompanyItem   // 당사회사에관한사항
}

// StockExchangeTransfer 는 주식의포괄적교환·이전 증권신고서(DS006)를 조회한다.
func (c *Client) StockExchangeTransfer(ctx context.Context, p Params) (*StockExchangeTransferRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/extrRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &StockExchangeTransferRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// Merger 는 합병 증권신고서(DS006)를 조회한다.
func (c *Client) Merger(ctx context.Context, p Params) (*MergerRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/mgRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &MergerRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "발행증권":
			derr = json.Unmarshal(g.List, &out.IssuedSecurities)
		case "당사회사에관한사항":
			derr = json.Unmarshal(g.List, &out.PartyCompanies)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
