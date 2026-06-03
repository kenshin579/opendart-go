// Package ownership 는 OpenDART DS004 지분공시 종합정보 API sub-client 다.
// opendart.Client.Ownership 로 접근한다.
package ownership

import "github.com/kenshin579/opendart-go/internal/httpclient"

// Client 는 지분공시 종합정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }
