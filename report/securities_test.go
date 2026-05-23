package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebtSecuritiesIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/detScritsIsuAcmslt.json": "detScritsIsuAcmslt.json"})
	items, err := c.DebtSecuritiesIssuance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "회사채", items[0].ScritsKndNm)
	assert.Equal(t, "128,940,000,000", items[0].FacvaluTotamt)
	assert.Equal(t, "2027.10.01", items[0].Mtd)
}

func TestCorporateBondBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cprndNrdmpBlce.json": "cprndNrdmpBlce.json"})
	items, err := c.CorporateBondBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "541,548,000,000", items[0].Sm)
	assert.Equal(t, "522,207,000,000", items[0].Yy1ExcessYy2Below)
}

func TestCommercialPaperBalance(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/entrprsBilScritsNrdmpBlce.json": "entrprsBilScritsNrdmpBlce.json"})
	items, err := c.CommercialPaperBalance(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미상환잔액", items[0].RemndrExprtn1)
	assert.Equal(t, "-", items[0].Yy3Excess)
}
