// Package disclosure 는 OpenDART DS001 공시정보 API sub-client 다.
// opendart.Client.Disclosure 로 접근한다.
package disclosure

import "github.com/kenshin579/opendart/internal/httpclient"

// Client 는 공시정보 sub-client.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root opendart.NewClient 가 호출한다.
func New(http *httpclient.Client) *Client { return &Client{http: http} }
