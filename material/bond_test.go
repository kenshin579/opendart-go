package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertibleBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/cvbdIsDecsn.json": "cvbdIsDecsn.json",
	})

	items, err := c.ConvertibleBondIssuance(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230815000123", got.RceptNo)
	assert.Equal(t, "30,000,000,000", got.BdFta)
	assert.Equal(t, "15,000", got.CvPrc)
	assert.Equal(t, "2024년 08월 16일", got.CvrqpdBgd)
	assert.Equal(t, "미제출", got.RsSmAtn)
	assert.Equal(t, "기명식 무보증 사모 전환사채", got.BdKnd)
	assert.Equal(t, "10,500", got.ActMktprcflCvprcLwtrsprc)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
