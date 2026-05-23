package report

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// newTestClient 는 testdata fixture 를 path 별로 서빙하는 report.Client 를 만든다.
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

func TestReportParams_toMap(t *testing.T) {
	m := ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport}.toMap()
	assert.Equal(t, "00126380", m["corp_code"])
	assert.Equal(t, "2023", m["bsns_year"])
	assert.Equal(t, "11011", m["reprt_code"])
}

func TestGetList_NoData(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	c := New(hc)
	_, err := c.Dividend(context.Background(), ReportParams{CorpCode: "x", BsnsYear: "2023", ReprtCode: AnnualReport})
	assert.ErrorIs(t, err, httpclient.ErrNoData)
}
