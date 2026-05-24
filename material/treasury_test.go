package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreasuryStockAcquisition(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqDecsn.json": "tsstkAqDecsn.json",
	})

	items, err := c.TreasuryStockAcquisition(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20241118000328", got.RceptNo)
	assert.Equal(t, "50,144,628", got.AqplnStkOstk)
	assert.Equal(t, "유가증권시장을 통한 장내 매수", got.AqMth)
	assert.Equal(t, "-", got.AqWtnDivOstkRt)
	assert.Equal(t, "6,525,123", got.D1ProdlmOstk)
}

func TestTreasuryStockDisposal(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkDpDecsn.json": "tsstkDpDecsn.json",
	})

	items, err := c.TreasuryStockDisposal(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230610000222", got.RceptNo)
	assert.Equal(t, "500,000", got.DpplnStkOstk)
	assert.Equal(t, "37,500,000,000", got.DpplnPrcOstk)
	assert.Equal(t, "임직원 상여 지급", got.DpPp)
	assert.Equal(t, "500,000", got.DpMMkt)
	assert.Equal(t, "125,000", got.D1SlodlmOstk)
}

func TestTreasuryStockTrustContract(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/tsstkAqTrctrCnsDecsn.json": "tsstkAqTrctrCnsDecsn.json",
	})

	items, err := c.TreasuryStockTrustContract(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230310000333", got.RceptNo)
	assert.Equal(t, "50,000,000,000", got.CtrPrc)
	assert.Equal(t, "한국투자증권", got.CtrCnsInt)
	assert.Equal(t, "2023년 03월 10일", got.Bddd)
	assert.Equal(t, "주가 안정 및 주주가치 제고", got.CtrPp)
}
