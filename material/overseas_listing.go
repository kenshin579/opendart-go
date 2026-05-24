package material

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// OverseasListingDecisionItem 은 해외 증권시장 주권등 상장 결정 (ovLstDecsn) 한 건.
type OverseasListingDecisionItem struct {
	RceptNo         string `json:"rcept_no"`          // 접수번호
	CorpCls         string `json:"corp_cls"`          // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`         // 고유번호
	CorpName        string `json:"corp_name"`         // 회사명
	LstprstkOstkCnt string `json:"lstprstk_ostk_cnt"` // 상장예정주식 종류ㆍ수(주)(보통주식)
	LstprstkEstkCnt string `json:"lstprstk_estk_cnt"` // 상장예정주식 종류ㆍ수(주)(기타주식)
	TisstkOstk      string `json:"tisstk_ostk"`       // 발행주식 총수(주)(보통주식)
	TisstkEstk      string `json:"tisstk_estk"`       // 발행주식 총수(주)(기타주식)
	PsmthNstkSl     string `json:"psmth_nstk_sl"`     // 공모방법(신주발행 (주))
	PsmthOstkSl     string `json:"psmth_ostk_sl"`     // 공모방법(구주매출 (주))
	Fdpp            string `json:"fdpp"`              // 자금조달(신주발행) 목적
	LststkOrlst     string `json:"lststk_orlst"`      // 상장증권(원주상장 (주))
	LststkDrlst     string `json:"lststk_drlst"`      // 상장증권(DR상장 (주))
	LstexNt         string `json:"lstex_nt"`          // 상장거래소(소재국가)
	Lstpp           string `json:"lstpp"`             // 해외상장목적
	Lstprd          string `json:"lstprd"`            // 상장예정일자
	Bddd            string `json:"bddd"`              // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`         // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`         // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`         // 감사(감사위원) 참석여부
}

// OverseasListingDecision 은 해외 증권시장 주권등 상장 결정(주요사항보고서)을 조회한다.
func (c *Client) OverseasListingDecision(ctx context.Context, p MaterialParams) ([]OverseasListingDecisionItem, error) {
	return httpclient.GetList[OverseasListingDecisionItem](ctx, c.http, "/api/ovLstDecsn.json", p.toMap())
}

// OverseasListingItem 은 해외 증권시장 주권등 상장 (ovLst) 한 건.
type OverseasListingItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	LststkOstkCnt string `json:"lststk_ostk_cnt"` // 상장주식 종류 및 수(보통주식(주))
	LststkEstkCnt string `json:"lststk_estk_cnt"` // 상장주식 종류 및 수(기타주식(주))
	LstexNt       string `json:"lstex_nt"`        // 상장거래소(소재국가)
	StkCd         string `json:"stk_cd"`          // 종목 명 (code)
	Lstd          string `json:"lstd"`            // 상장일자
	Cfd           string `json:"cfd"`             // 확인일자
}

// OverseasListing 은 해외 증권시장 주권등 상장(주요사항보고서)을 조회한다.
func (c *Client) OverseasListing(ctx context.Context, p MaterialParams) ([]OverseasListingItem, error) {
	return httpclient.GetList[OverseasListingItem](ctx, c.http, "/api/ovLst.json", p.toMap())
}
