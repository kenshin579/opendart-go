package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResp struct {
	Envelope
	CorpName string `json:"corp_name"`
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
}

func TestGetJSON_SuccessAndKeyInjection(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "KEY", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
		w.Write([]byte(`{"status":"000","message":"정상","corp_name":"삼성전자(주)"}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", map[string]string{"corp_code": "00126380"}, &out)
	require.NoError(t, err)
	assert.Equal(t, "삼성전자(주)", out.CorpName)
}

func TestGetJSON_NoData(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", nil, &out)
	assert.ErrorIs(t, err, ErrNoData)
}

func TestGetJSON_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/api/company.json", nil, &out)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "010", apiErr.Status)
}

func TestGetBytes_Binary(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PK\x03\x04zipdata"))
	})
	b, err := c.GetBytes(context.Background(), "/api/corpCode.xml", nil)
	require.NoError(t, err)
	assert.Equal(t, "PK\x03\x04zipdata", string(b))
}

func TestGetBytes_ErrorJSON(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	})
	_, err := c.GetBytes(context.Background(), "/api/document.xml", nil)
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
}

func TestGetJSON_HTTPError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	var out testResp
	err := c.GetJSON(context.Background(), "/x", nil, &out)
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNoData)
}
