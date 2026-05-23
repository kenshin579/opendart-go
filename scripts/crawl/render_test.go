package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMarkdown(t *testing.T) {
	ref := APIRef{
		GrpCd: "DS001", Category: "공시정보", APIID: "2019002",
		Name: "기업개황", Desc: "DART에 등록되어있는 기업의 개황정보를 제공합니다.",
	}
	spec := APISpec{
		BasicInfo: Table{
			Headers: []string{"메서드", "요청URL", "인코딩", "출력포멧"},
			Rows:    [][]string{{"GET", "https://opendart.fss.or.kr/api/company.json", "UTF-8", "JSON"}},
		},
		Request: Table{
			Headers: []string{"요청키", "명칭", "타입", "필수여부", "값설명"},
			Rows:    [][]string{{"crtfc_key", "API 인증키", "STRING(40)", "Y", "발급받은 인증키(40자리)"}},
		},
		Response: Table{
			Headers: []string{"응답키", "명칭", "출력설명"},
			Rows:    [][]string{{"corp_name", "정식명칭", "정식회사 명칭"}},
		},
	}
	assert.Equal(t, readFixture(t, "company.golden.md"), renderMarkdown(ref, spec))
}

func TestRenderTableEscapesPipe(t *testing.T) {
	out := renderTable(Table{Headers: []string{"a"}, Rows: [][]string{{"x|y"}}})
	assert.Contains(t, out, `x\|y`)
}

func TestRenderIndex(t *testing.T) {
	refs := []APIRef{
		{GrpCd: "DS001", Category: "공시정보", APIID: "2019002", Name: "기업개황", Desc: "개황정보"},
		{GrpCd: "DS001", Category: "공시정보", APIID: "2019001", Name: "공시검색", Desc: "검색"},
	}
	out := renderIndex(refs)
	assert.Contains(t, out, "## 공시정보 (DS001)")
	// apiId 오름차순 정렬: 공시검색(2019001) 이 기업개황(2019002) 보다 먼저
	assert.Less(t, strings.Index(out, "공시검색"), strings.Index(out, "기업개황"))
	assert.Contains(t, out, "[기업개황](<공시정보/기업개황.md>)")
}

func TestSanitizeReplacesPathChars(t *testing.T) {
	assert.Equal(t, "a_b", sanitize("a/b"))
	assert.Equal(t, "공시정보", sanitize("공시정보"))
}
