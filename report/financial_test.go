package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglAcnt.json": "fnlttSinglAcnt.json"})
	items, err := c.SingleAccount(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "연결재무제표", items[0].FsNm)
	assert.Equal(t, "195,936,557,000,000", items[0].ThstrmAmount)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestMultiAccount(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttMultiAcnt.json": "fnlttMultiAcnt.json"})
	items, err := c.MultiAccount(context.Background(), ReportParams{CorpCode: "00126380,00164779", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "유동자산", items[0].AccountNm)
	assert.Equal(t, "00126380", items[0].CorpCode)
}

func TestFsDivAndIndexParams_toMap(t *testing.T) {
	sm := FinancialStatementParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, FsDiv: FsDivSeparate}.toMap()
	assert.Equal(t, "OFS", sm["fs_div"])
	assert.Equal(t, "11011", sm["reprt_code"])

	im := FinancialIndexParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability}.toMap()
	assert.Equal(t, "M210000", im["idx_cl_code"])
	assert.Equal(t, "00126380", im["corp_code"])
}

func TestSingleFullStatement(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglAcntAll.json": "fnlttSinglAcntAll.json"})
	items, err := c.SingleFullStatement(context.Background(), FinancialStatementParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, FsDiv: FsDivSeparate})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "ifrs-full_Assets", items[0].AccountId)
	assert.Equal(t, "자산총계", items[0].AccountNm)
	assert.Equal(t, "296857289000000", items[0].ThstrmAmount)
}

func TestSingleIndex(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttSinglIndx.json": "fnlttSinglIndx.json"})
	items, err := c.SingleIndex(context.Background(), FinancialIndexParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "수익성지표", items[0].IdxClNm)
	assert.Equal(t, "세전계속사업이익률", items[0].IdxNm)
	assert.Equal(t, "5.981", items[0].IdxVal)
}

func TestMultiIndex(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fnlttCmpnyIndx.json": "fnlttCmpnyIndx.json"})
	items, err := c.MultiIndex(context.Background(), FinancialIndexParams{CorpCode: "00126380,00164779", BsnsYear: "2023", ReprtCode: AnnualReport, IdxClCode: IndexProfitability})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "00126380", items[0].CorpCode)
	assert.Equal(t, "5.981", items[0].IdxVal)
}
