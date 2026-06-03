package disclosure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// newTestClient 는 testdata 의 fixture 를 path 별로 서빙하는 disclosure.Client 를 만든다.
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

func TestGetCompany(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/company.json": "company_samsung.json"})
	got, err := c.GetCompany(context.Background(), "00126380")
	require.NoError(t, err)
	assert.Equal(t, "삼성전자(주)", got.CorpName)
	assert.Equal(t, "005930", got.StockCode)
	assert.Equal(t, "Y", got.CorpCls)
	assert.Equal(t, "19690113", got.EstDate)
	assert.Equal(t, "12", got.AccMonth)
}
