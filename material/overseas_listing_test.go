package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverseasListingDecision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLstDecsn.json": "ovLstDecsn.json"})
	items, err := c.OverseasListingDecision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20241031000536", got.RceptNo)
	assert.Equal(t, "38,685,850", got.LstprstkEstkCnt)
	assert.Equal(t, "런던증권거래소(영국)", got.LstexNt)
	assert.Equal(t, "2025년 년 월 1일일", got.Lstprd)
}

func TestOverseasListing(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLst.json": "ovLst.json"})
	items, err := c.OverseasListing(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231010000222", got.RceptNo)
	assert.Equal(t, "1,000,000", got.LststkOstkCnt)
	assert.Equal(t, "TEST", got.StkCd)
	assert.Equal(t, "2023년 10월 02일", got.Lstd)
}
