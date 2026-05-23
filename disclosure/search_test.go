package disclosure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchDisclosures(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/list.json": "list_samsung.json"})
	res, err := c.SearchDisclosures(context.Background(), SearchParams{CorpCode: "00126380", PageCount: 2})
	require.NoError(t, err)
	assert.Equal(t, 1, res.PageNo)
	assert.Equal(t, 15, res.TotalCount)
	assert.Equal(t, 8, res.TotalPage)
	require.Len(t, res.List, 2)
	assert.Equal(t, "20240131000326", res.List[0].RceptNo)
	assert.Equal(t, "특수관계인과의내부거래", res.List[0].ReportNm)
}

func TestSearchParams_toMap(t *testing.T) {
	m := SearchParams{CorpCode: "00126380", PblntfTy: "A", PageNo: 2, PageCount: 50}.toMap()
	assert.Equal(t, "00126380", m["corp_code"])
	assert.Equal(t, "A", m["pblntf_ty"])
	assert.Equal(t, "2", m["page_no"])
	assert.Equal(t, "50", m["page_count"])
	_, hasBgn := m["bgn_de"]
	assert.False(t, hasBgn)

	empty := SearchParams{}.toMap()
	_, hasPage := empty["page_no"]
	assert.False(t, hasPage)
}
