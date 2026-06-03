package report

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

func TestXbrlTaxonomy(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/xbrlTaxonomy.json": "xbrlTaxonomy_bs1.json"})
	items, err := c.XbrlTaxonomy(context.Background(), "BS1")
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "BS1", items[0].SjDiv)
	assert.Equal(t, "ifrs_CurrentAssets", items[0].AccountId)
	assert.Equal(t, "유동자산", items[0].LabelKor)
	assert.Equal(t, "X", items[0].DataTp)
}

func TestDownloadXbrl_Binary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "20240312000736", r.URL.Query().Get("rcept_no"))
		assert.Equal(t, "11011", r.URL.Query().Get("reprt_code"))
		w.Write([]byte("PK\x03\x04xbrlzip"))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	b, err := c.DownloadXbrl(context.Background(), "20240312000736", AnnualReport)
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04xbrlzip", string(b))
}

func TestDownloadXbrl_ErrorJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	_, err := c.DownloadXbrl(context.Background(), "x", AnnualReport)
	var apiErr *httpclient.APIError
	assert.ErrorAs(t, err, &apiErr)
}
