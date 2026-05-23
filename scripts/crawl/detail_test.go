package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDetail(t *testing.T) {
	spec, err := parseDetail(readFixture(t, "company_detail.html"))
	require.NoError(t, err)

	// 기본 정보
	assert.Equal(t, []string{"메서드", "요청URL", "인코딩", "출력포멧"}, spec.BasicInfo.Headers)
	require.GreaterOrEqual(t, len(spec.BasicInfo.Rows), 1)
	assert.Equal(t, "GET", spec.BasicInfo.Rows[0][0])
	assert.Contains(t, spec.BasicInfo.Rows[0][1], "company.json")

	// 요청 인자
	assert.Equal(t, []string{"요청키", "명칭", "타입", "필수여부", "값설명"}, spec.Request.Headers)
	assert.Equal(t, "crtfc_key", spec.Request.Rows[0][0])

	// 응답 결과 (타입 컬럼 없음 — 3컬럼)
	assert.Equal(t, []string{"응답키", "명칭", "출력설명"}, spec.Response.Headers)
	found := false
	for _, r := range spec.Response.Rows {
		if len(r) > 0 && r[0] == "corp_name" {
			found = true
		}
	}
	assert.True(t, found, "응답 결과에 corp_name 필드가 있어야 함")
}
