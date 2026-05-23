package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return string(b)
}

func TestParseList(t *testing.T) {
	refs, err := parseList(readFixture(t, "ds001_list.html"))
	require.NoError(t, err)
	require.Len(t, refs, 4) // DS001 = 공시검색/기업개황/공시서류원본파일/고유번호

	var company APIRef
	for _, r := range refs {
		if r.APIID == "2019002" {
			company = r
		}
	}
	assert.Equal(t, "기업개황", company.Name)
	assert.Contains(t, company.Desc, "개황정보")
}

func TestAPIIDFromHref(t *testing.T) {
	assert.Equal(t, "2019002", apiIDFromHref("/guide/detail.do?apiGrpCd=DS001&apiId=2019002"))
	assert.Equal(t, "", apiIDFromHref("/guide/main.do?apiGrpCd=DS001"))
}
