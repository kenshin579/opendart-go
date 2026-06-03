package registration

import "github.com/kenshin579/opendart-go/internal/httpclient"

// Client 는 DS006 증권신고서 주요정보 API 클라이언트.
type Client struct{ http *httpclient.Client }

// New 는 registration.Client 를 만든다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// Params 는 DS006 증권신고서 공통 요청 파라미터(corp_code + 기간).
type Params struct {
	CorpCode string // 고유번호(8자리)
	BgnDe    string // 검색시작 접수일자(YYYYMMDD)
	EndDe    string // 검색종료 접수일자(YYYYMMDD)
}

func (p Params) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" {
		m["bgn_de"] = p.BgnDe
	}
	if p.EndDe != "" {
		m["end_de"] = p.EndDe
	}
	return m
}
