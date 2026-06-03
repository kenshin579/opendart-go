package registration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart-go/internal/httpclient"
)

// newTestClient 는 지정한 testdata fixture 를 모든 요청에 서빙하는 registration.Client 를 만든다.
func newTestClient(t *testing.T, fixture string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(filepath.Join("testdata", fixture))
		require.NoError(t, err)
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	hc := httpclient.New(httpclient.Config{APIKey: "KEY", BaseURL: srv.URL, HTTPClient: srv.Client()})
	return New(hc)
}
