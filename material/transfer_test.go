package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBusinessAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnInhDecsn.json": "bsnInhDecsn.json"})
	items, err := c.BusinessAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230410000111", got.RceptNo)
	assert.Equal(t, "반도체 사업부문", got.InhBsn)
	assert.Equal(t, "500,000,000,000", got.InhPrc)
	assert.Equal(t, "삼일회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestBusinessTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnTrfDecsn.json": "bsnTrfDecsn.json"})
	items, err := c.BusinessTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230510000222", got.RceptNo)
	assert.Equal(t, "디스플레이 사업부문", got.TrfBsn)
	assert.Equal(t, "300,000,000,000", got.TrfPrc)
	assert.Equal(t, "안진회계법인", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestTangibleAssetAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tgastInhDecsn.json": "tgastInhDecsn.json"})
	items, err := c.TangibleAssetAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230610000333", got.RceptNo)
	assert.Equal(t, "토지 및 건물", got.AstSen)
	assert.Equal(t, "150,000,000,000", got.InhdtlInhprc)
	assert.Equal(t, "한국감정원", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestTangibleAssetTransfer(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tgastTrfDecsn.json": "tgastTrfDecsn.json"})
	items, err := c.TangibleAssetTransfer(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230710000444", got.RceptNo)
	assert.Equal(t, "토지", got.AstSen)
	assert.Equal(t, "80,000,000,000", got.TrfdtlTrfprc)
	assert.Equal(t, "한국감정원", got.ExevlIntn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
