package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDividend(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/alotMatter.json": "alotMatter.json"})
	items, err := c.Dividend(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "주당액면가액(원)", items[0].Se)
	assert.Equal(t, "100", items[0].Thstrm)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
}

func TestCapitalChange(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/irdsSttus.json": "irdsSttus.json"})
	items, err := c.CapitalChange(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성전자", items[0].CorpName)
	assert.Equal(t, "2023-12-31", items[0].StlmDt)
}

func TestTreasuryStock(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/tesstkAcqsDspsSttus.json": "tesstkAcqsDspsSttus.json"})
	items, err := c.TreasuryStock(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestTotalStock(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/stockTotqySttus.json": "stockTotqySttus.json"})
	items, err := c.TotalStock(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "보통주", items[0].Se)
	assert.Equal(t, "20,000,000,000", items[0].IsuStockTotqy)
	assert.Equal(t, "5,969,782,550", items[0].DistbStockCo)
}

func TestMajorShareholders(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hyslrSttus.json": "hyslrSttus.json"})
	items, err := c.MajorShareholders(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성생명보험㈜", items[0].Nm)
	assert.Equal(t, "최대주주 본인", items[0].Relate)
	assert.Equal(t, "508,157,148", items[0].TrmendPosesnStockCo)
}

func TestMajorShareholderChanges(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/hyslrChgSttus.json": "hyslrChgSttus.json"})
	items, err := c.MajorShareholderChanges(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "삼성생명보험㈜", items[0].MxmmShrholdrNm)
	assert.Equal(t, "변동전 최대주주의 피상속", items[0].ChangeCause)
}

func TestMinorityShareholders(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/mrhlSttus.json": "mrhlSttus.json"})
	items, err := c.MinorityShareholders(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "소액주주", items[0].Se)
	assert.Equal(t, "4,672,039", items[0].ShrholdrCo)
	assert.Equal(t, "67.30%", items[0].HoldStockRate)
}
