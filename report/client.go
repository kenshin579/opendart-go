// Package report 는 OpenDART DS002 정기보고서 주요정보 API sub-client 다.
// opendart.Client.Report 로 접근한다.
package report

import (
	"context"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// Client 는 정기보고서 주요정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// ReportCode 는 정기보고서 종류 코드.
type ReportCode string

const (
	Q1Report     ReportCode = "11013" // 1분기보고서
	HalfReport   ReportCode = "11012" // 반기보고서
	Q3Report     ReportCode = "11014" // 3분기보고서
	AnnualReport ReportCode = "11011" // 사업보고서
)

// ReportParams 는 DS002 공통 요청 인자 (전부 필수).
type ReportParams struct {
	CorpCode  string     // 고유번호 (8자리)
	BsnsYear  string     // 사업연도 (4자리, 2015 이후)
	ReprtCode ReportCode // 보고서 코드
}

func (p ReportParams) toMap() map[string]string {
	return map[string]string{
		"corp_code":  p.CorpCode,
		"bsns_year":  p.BsnsYear,
		"reprt_code": string(p.ReprtCode),
	}
}

// listResponse 는 DS002 공통 list 응답 envelope.
type listResponse[T any] struct {
	httpclient.Envelope
	List []T `json:"list"`
}

// getList 는 공통 list 조회 헬퍼. GetJSON 의 status 검사를 거친 뒤 list 만 반환한다.
// 조회 데이터 없음(013)은 httpclient 가 ErrNoData 로 변환한다.
func getList[T any](ctx context.Context, hc *httpclient.Client, path string, p ReportParams) ([]T, error) {
	var resp listResponse[T]
	if err := hc.GetJSON(ctx, path, p.toMap(), &resp); err != nil {
		return nil, err
	}
	return resp.List, nil
}
