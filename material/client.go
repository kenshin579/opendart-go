// Package material 는 OpenDART DS005 주요사항보고서 주요정보 API sub-client 다.
// opendart.Client.Material 로 접근한다.
package material

import "github.com/kenshin579/opendart/internal/httpclient"

// Client 는 주요사항보고서 주요정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// MaterialParams 는 DS005 공통 요청 인자 (날짜 범위). 빈 값은 쿼리에서 생략한다(OpenDART 기본값 적용).
type MaterialParams struct {
	CorpCode string // 고유번호 (8자리)
	BgnDe    string // 시작일 YYYYMMDD
	EndDe    string // 종료일 YYYYMMDD
}

func (p MaterialParams) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" {
		m["bgn_de"] = p.BgnDe
	}
	if p.EndDe != "" {
		m["end_de"] = p.EndDe
	}
	return m
}
