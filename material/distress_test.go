package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOccurrences(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dfOcr.json": "dfOcr.json"})
	items, err := c.DefaultOccurrences(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DH오토넥스", items[0].CorpName)
	assert.Equal(t, "당사 김제지점 발행 만기어음 부도", items[0].DfCn)
	assert.Equal(t, "48,322,175", items[0].DfAmt)
}
