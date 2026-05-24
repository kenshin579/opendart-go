package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanyMerger(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cmpMgDecsn.json": "cmpMgDecsn.json"})
	items, err := c.CompanyMerger(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230410000111", got.RceptNo)
	assert.Equal(t, "흡수합병", got.MgMth)
	assert.Equal(t, "1:0.5", got.MgRt)
	assert.Equal(t, "합병상대회사", got.MgptncmpCmpnm)
	assert.Equal(t, "2023년 06월 10일", got.MgscGmtsckPrd)
	assert.Equal(t, "제출", got.RsSmAtn)
}

func TestCompanyDivision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/cmpDvDecsn.json": "cmpDvDecsn.json"})
	items, err := c.CompanyDivision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20230510000222", got.RceptNo)
	assert.Equal(t, "인적분할", got.DvMth)
	assert.Equal(t, "배터리 사업부문", got.DvTrfbsnprtCn)
	assert.Equal(t, "30.0", got.AbcrCrrt)
	assert.Equal(t, "2023년 07월 01일", got.Dvdt)
	assert.Equal(t, "제출", got.RsSmAtn)
}
