package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutives(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/exctvSttus.json": "exctvSttus.json"})
	items, err := c.Executives(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "한종희", items[0].Nm)
	assert.Equal(t, "부회장", items[0].Ofcps)
	assert.Equal(t, "사내이사", items[0].RgistExctvAt)
}

func TestEmployees(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/empSttus.json": "empSttus.json"})
	items, err := c.Employees(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DX", items[0].FoBbm)
	assert.Equal(t, "37,962", items[0].RgllbrCo)
	assert.Equal(t, "38,286", items[0].Sm)
}

func TestUnregisteredExecutiveCompensation(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/unrstExctvMendngSttus.json": "unrstExctvMendngSttus.json"})
	items, err := c.UnregisteredExecutiveCompensation(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "미등기임원", items[0].Se)
	assert.Equal(t, "1,015", items[0].Nmpr)
}

func TestOutsideDirectorChanges(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/outcmpnyDrctrNdChangeSttus.json": "outcmpnyDrctrNdChangeSttus.json"})
	items, err := c.OutsideDirectorChanges(context.Background(), ReportParams{CorpCode: "00126380", BsnsYear: "2023", ReprtCode: AnnualReport})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "11", items[0].DrctrCo)
	assert.Equal(t, "6", items[0].OtcmpDrctrCo)
}
