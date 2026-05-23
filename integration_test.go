//go:build integration

package opendart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// 실행: OPENDART_API_KEY=... go test -tags integration -run TestIntegration -v
func TestIntegration_GetCompany(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	require.Equal(t, "00126380", corp)

	company, err := c.Disclosure.GetCompany(context.Background(), corp)
	require.NoError(t, err)
	require.Contains(t, company.CorpName, "삼성전자")
}
