package material

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 testdata fixture 를 path 별로 서빙하는 material.Client 를 만든다.
func newTestClient(t *testing.T, routes map[string]string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture, ok := routes[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		b, err := os.ReadFile(filepath.Join("testdata", fixture))
		require.NoError(t, err)
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	return New(hc)
}

func TestMaterialParams_toMap(t *testing.T) {
	m := MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"}.toMap()
	assert.Equal(t, "00126089", m["corp_code"])
	assert.Equal(t, "20230101", m["bgn_de"])
	assert.Equal(t, "20231231", m["end_de"])

	only := MaterialParams{CorpCode: "00126089"}.toMap()
	_, hasBgn := only["bgn_de"]
	assert.False(t, hasBgn)
	assert.Equal(t, "00126089", only["corp_code"])
}
