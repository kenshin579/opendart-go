package disclosure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/internal/httpclient"
)

func TestDownloadDocument_Binary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "20240131000326", r.URL.Query().Get("rcept_no"))
		w.Write([]byte("PK\x03\x04docdata"))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	b, err := c.DownloadDocument(context.Background(), "20240131000326")
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04docdata", string(b))
}

func TestDownloadDocument_ErrorJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	}))
	t.Cleanup(srv.Close)
	c := New(httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()}))

	_, err := c.DownloadDocument(context.Background(), "x")
	var apiErr *httpclient.APIError
	assert.ErrorAs(t, err, &apiErr)
}
