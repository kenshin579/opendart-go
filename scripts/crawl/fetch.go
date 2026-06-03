package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent   = "opendart-doc-crawler (github.com/kenshin579/opendart-go)"
	politeDelay = 300 * time.Millisecond
)

// newHTTPClient 는 OpenDART 서버에 맞춘 HTTP 클라이언트를 만든다.
//
// 서버(opendart.fss.or.kr)는 TLS 1.2 의 RSA 키교환 cipher
// (AES128-GCM-SHA256 등)만 지원한다. Go 기본 ClientHello 는 forward secrecy
// 가 없는 RSA 키교환 cipher 를 제외하므로 handshake 가 실패한다. 해당 cipher 를
// 명시해 curl 과 동일하게 협상하도록 한다.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
		},
	}
}

// httpGet 은 User-Agent 를 붙여 URL 본문을 문자열로 가져온다.
func httpGet(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
