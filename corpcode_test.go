package opendart

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const corpXML = `<?xml version="1.0" encoding="UTF-8"?>
<result>
	<list><corp_code>00126380</corp_code><corp_name>삼성전자(주)</corp_name><corp_eng_name>SAMSUNG</corp_eng_name><stock_code>005930</stock_code><modify_date>20240101</modify_date></list>
	<list><corp_code>00434003</corp_code><corp_name>다코</corp_name><corp_eng_name>Daco</corp_eng_name><stock_code> </stock_code><modify_date>20170630</modify_date></list>
</result>`

func corpCodeServer(t *testing.T) *httptest.Server {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("CORPCODE.xml")
	require.NoError(t, err)
	_, err = w.Write([]byte(corpXML))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	zipBytes := buf.Bytes()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(zipBytes)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newCorpTestClient(t *testing.T) *Client {
	srv := corpCodeServer(t)
	c, err := NewClient("KEY",
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithCorpCodeCacheDir(t.TempDir()),
	)
	require.NoError(t, err)
	return c
}

func TestResolveCorpCode(t *testing.T) {
	c := newCorpTestClient(t)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	assert.Equal(t, "00126380", corp)

	_, err = c.ResolveCorpCode(context.Background(), "999999")
	assert.ErrorIs(t, err, ErrCorpCodeNotFound)
}

func TestLookupAndList(t *testing.T) {
	c := newCorpTestClient(t)
	e, err := c.LookupCorpCode(context.Background(), "00434003")
	require.NoError(t, err)
	assert.Equal(t, "다코", e.CorpName)

	all, err := c.CorpCodes(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 2)
}
