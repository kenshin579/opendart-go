# OpenDART 문서 크롤러 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** OpenDART 개발가이드의 전체 API 명세를 크롤링해 `docs/api/{한글카테고리}/{API명}.md` 로 변환하는 재사용 가능한 Go 크롤러를 만든다.

**Architecture:** `scripts/crawl` 의 `package main`. 순수 함수(`parseList`, `parseDetail`, `renderMarkdown`, `renderIndex`)와 I/O(`httpGet`, `writeDoc`, `writeIndex`, `run`)를 분리해 순수 함수를 fixture 기반 golden 테스트로 검증한다. 6개 카테고리 목록 페이지에서 API 링크를 동적 수집 → 각 상세 페이지를 `goquery` 로 `<caption>` 기준 표 추출 → markdown 렌더링.

**Tech Stack:** Go 1.25+, `github.com/PuerkitoBio/goquery` (HTML 파싱), `github.com/stretchr/testify` (test).

**Spec:** `docs/superpowers/specs/2026-05-23-opendart-doc-crawler-design.md`

**검증된 HTML 구조 (사실):**
- 목록 페이지 `guide/main.do?apiGrpCd=DS00X`: 표 행 셀 = `[번호, API명, 상세기능(설명), 개발가이드(detail링크)]`. 링크 href = `/guide/detail.do?apiGrpCd=DS001&apiId=2019002`.
- 상세 페이지 `guide/detail.do?...`: `<table>` 을 `<caption>` 으로 식별 — `기본 정보`(헤더: 메서드/요청URL/인코딩/출력포멧), `요청 인자`(요청키/명칭/타입/필수여부/값설명), `응답 결과`(응답키/명칭/출력설명, **타입 컬럼 없음**), `OpenAPI 테스트`(스킵), `메시지 설명`(전 API 공통 상태코드 — per-API md 에서 제외). 각 표는 헤더 1행이 `<th>`, 데이터 행이 `<td>`.

---

## File Structure

```
opendart/
  go.mod                       # module github.com/kenshin579/opendart
  go.sum
  scripts/crawl/
    model.go                   # 타입 + 카테고리 상수
    list.go                    # parseList, apiIDFromHref
    detail.go                  # parseDetail, extractTable, cellText
    render.go                  # renderMarkdown, renderIndex, renderTable + 헬퍼
    fetch.go                   # httpGet, 상수
    main.go                    # run 오케스트레이션, writeDoc, writeIndex, sanitize
    list_test.go
    detail_test.go
    render_test.go
    testdata/
      ds001_list.html          # 캡처한 DS001 목록 페이지
      company_detail.html      # 캡처한 기업개황 상세 페이지
      company.golden.md        # renderMarkdown 기대 출력
  docs/api/                    # (크롤러 실행 산출물)
```

---

### Task 1: 모듈 초기화 · 의존성 · 테스트 fixture

**Files:**
- Create: `go.mod`, `go.sum` (생성됨)
- Create: `scripts/crawl/testdata/ds001_list.html`
- Create: `scripts/crawl/testdata/company_detail.html`

- [ ] **Step 1: Go 모듈 초기화 및 의존성 추가**

Run (레포 루트 `/Users/user/src/workspace_moneyflow/opendart` 에서):
```bash
go mod init github.com/kenshin579/opendart
go get github.com/PuerkitoBio/goquery@latest
go get github.com/stretchr/testify@latest
```
Expected: `go.mod` 에 `module github.com/kenshin579/opendart`, `go 1.25` (또는 설치된 버전), goquery/testify require 추가.

- [ ] **Step 2: 테스트 fixture 캡처 (실제 OpenDART 페이지)**

Run:
```bash
mkdir -p scripts/crawl/testdata
curl -s -A "opendart-doc-crawler" \
  "https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS001" \
  -o scripts/crawl/testdata/ds001_list.html
curl -s -A "opendart-doc-crawler" \
  "https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS001&apiId=2019002" \
  -o scripts/crawl/testdata/company_detail.html
```

- [ ] **Step 3: fixture 검증**

Run:
```bash
grep -c "detail.do?apiGrpCd=DS001" scripts/crawl/testdata/ds001_list.html
grep -c "기업개황\|company.json" scripts/crawl/testdata/company_detail.html
```
Expected: 두 grep 모두 1 이상 (페이지가 비어있지 않음).

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum scripts/crawl/testdata
git commit -m "chore: init opendart module + crawler test fixtures"
```

---

### Task 2: 데이터 모델 + 카테고리 상수

**Files:**
- Create: `scripts/crawl/model.go`

- [ ] **Step 1: model.go 작성**

```go
package main

// categories 는 OpenDART 개발가이드 6개 그룹 (apiGrpCd → 한글명).
var categories = []struct {
	Code string
	Name string
}{
	{"DS001", "공시정보"},
	{"DS002", "정기보고서 주요정보"},
	{"DS003", "정기보고서 재무정보"},
	{"DS004", "지분공시 종합정보"},
	{"DS005", "주요사항보고서 주요정보"},
	{"DS006", "증권신고서 주요정보"},
}

// APIRef 는 카테고리 목록 페이지에서 추출한 개별 API 식별 정보.
type APIRef struct {
	GrpCd    string // DS001
	Category string // 공시정보
	APIID    string // 2019002
	Name     string // 기업개황
	Desc     string // 상세기능 설명
}

// Table 은 상세 페이지의 한 표 (헤더 행 + 데이터 행들).
type Table struct {
	Headers []string
	Rows    [][]string
}

// APISpec 은 detail 페이지에서 추출한 명세. 메시지 설명(공통 상태코드)은 제외.
type APISpec struct {
	BasicInfo Table // 기본 정보 (메서드/요청URL/인코딩/출력포멧)
	Request   Table // 요청 인자
	Response  Table // 응답 결과
}
```

- [ ] **Step 2: 빌드 확인 (타입 선언은 동작이 없어 테스트 대신 빌드로 검증)**

Run: `go build ./scripts/crawl/`
Expected: 에러 없음 (미사용 타입 경고 없음 — Go 는 미사용 패키지 레벨 선언 허용).

- [ ] **Step 3: Commit**

```bash
git add scripts/crawl/model.go
git commit -m "feat(crawl): add data model and category constants"
```

---

### Task 3: 목록 페이지 파서 (parseList)

**Files:**
- Create: `scripts/crawl/list.go`
- Test: `scripts/crawl/list_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`scripts/crawl/list_test.go`:
```go
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
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./scripts/crawl/ -run 'TestParseList|TestAPIIDFromHref' -v`
Expected: FAIL — `undefined: parseList`, `undefined: apiIDFromHref`.

- [ ] **Step 3: list.go 구현**

```go
package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseList 는 카테고리 목록 페이지 HTML 에서 detail.do 링크를 가진 모든 행을 찾아
// APIRef 리스트를 만든다. GrpCd/Category 는 호출자가 채운다 (여기선 APIID/Name/Desc 만).
func parseList(html string) ([]APIRef, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var refs []APIRef
	doc.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		link := tr.Find(`a[href*="detail.do"]`)
		if link.Length() == 0 {
			return
		}
		href, _ := link.Attr("href")
		apiID := apiIDFromHref(href)
		if apiID == "" {
			return
		}
		cells := tr.Find("td") // [번호, API명, 상세기능, 개발가이드(링크)]
		refs = append(refs, APIRef{
			APIID: apiID,
			Name:  strings.TrimSpace(cells.Eq(1).Text()),
			Desc:  strings.Join(strings.Fields(cells.Eq(2).Text()), " "),
		})
	})
	return refs, nil
}

// apiIDFromHref 는 "...&apiId=2019002" 에서 "2019002" 를 추출한다. 없으면 "".
func apiIDFromHref(href string) string {
	const key = "apiId="
	i := strings.Index(href, key)
	if i < 0 {
		return ""
	}
	id := href[i+len(key):]
	if j := strings.IndexAny(id, "&\""); j >= 0 {
		id = id[:j]
	}
	return strings.TrimSpace(id)
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./scripts/crawl/ -run 'TestParseList|TestAPIIDFromHref' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add scripts/crawl/list.go scripts/crawl/list_test.go
git commit -m "feat(crawl): parse category list page into APIRefs"
```

---

### Task 4: 상세 페이지 파서 (parseDetail / extractTable)

**Files:**
- Create: `scripts/crawl/detail.go`
- Test: `scripts/crawl/detail_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`scripts/crawl/detail_test.go`:
```go
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
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./scripts/crawl/ -run TestParseDetail -v`
Expected: FAIL — `undefined: parseDetail`.

- [ ] **Step 3: detail.go 구현**

```go
package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseDetail 은 detail 페이지 HTML 에서 caption 으로 식별되는 핵심 3개 표를 추출한다.
func parseDetail(html string) (APISpec, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return APISpec{}, err
	}
	return APISpec{
		BasicInfo: extractTable(doc, "기본 정보"),
		Request:   extractTable(doc, "요청 인자"),
		Response:  extractTable(doc, "응답 결과"),
	}, nil
}

// extractTable 은 <caption> 텍스트에 captionContains 를 포함하는 첫 <table> 을 찾아
// 헤더(<th>)와 데이터 행(<td> 들)을 반환한다. 없으면 빈 Table.
func extractTable(doc *goquery.Document, captionContains string) Table {
	var out Table
	doc.Find("table").EachWithBreak(func(_ int, t *goquery.Selection) bool {
		caption := strings.TrimSpace(t.Find("caption").First().Text())
		if !strings.Contains(caption, captionContains) {
			return true // 계속
		}
		t.Find("th").Each(func(_ int, th *goquery.Selection) {
			out.Headers = append(out.Headers, cellText(th))
		})
		t.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			if tds.Length() == 0 {
				return // 헤더 행(th 만) 스킵
			}
			row := make([]string, 0, tds.Length())
			tds.Each(func(_ int, td *goquery.Selection) {
				row = append(row, cellText(td))
			})
			out.Rows = append(out.Rows, row)
		})
		return false // 첫 매치 후 중단
	})
	return out
}

// cellText 는 셀 내부 텍스트를 공백 정규화해 한 줄로 만든다 (md 표 안전).
func cellText(s *goquery.Selection) string {
	return strings.Join(strings.Fields(s.Text()), " ")
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./scripts/crawl/ -run TestParseDetail -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add scripts/crawl/detail.go scripts/crawl/detail_test.go
git commit -m "feat(crawl): parse detail page tables by caption"
```

---

### Task 5: Markdown 렌더러 (renderMarkdown / renderTable)

**Files:**
- Create: `scripts/crawl/render.go`
- Test: `scripts/crawl/render_test.go`
- Create: `scripts/crawl/testdata/company.golden.md`

- [ ] **Step 1: golden fixture 작성**

`scripts/crawl/testdata/company.golden.md` (파일 끝에 단일 개행 1개):
```markdown
# 기업개황

> apiGrpCd: DS001 · apiId: 2019002
>
> `GET https://opendart.fss.or.kr/api/company.json`

DART에 등록되어있는 기업의 개황정보를 제공합니다.

## 기본 정보

| 메서드 | 요청URL | 인코딩 | 출력포멧 |
| --- | --- | --- | --- |
| GET | https://opendart.fss.or.kr/api/company.json | UTF-8 | JSON |

## 요청 인자

| 요청키 | 명칭 | 타입 | 필수여부 | 값설명 |
| --- | --- | --- | --- | --- |
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |

## 응답 결과

| 응답키 | 명칭 | 출력설명 |
| --- | --- | --- |
| corp_name | 정식명칭 | 정식회사 명칭 |
```

- [ ] **Step 2: 실패하는 테스트 작성**

`scripts/crawl/render_test.go`:
```go
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
```

- [ ] **Step 3: 테스트 실패 확인**

Run: `go test ./scripts/crawl/ -run 'TestRenderMarkdown|TestRenderTableEscapesPipe' -v`
Expected: FAIL — `undefined: renderMarkdown`, `undefined: renderTable`.

- [ ] **Step 4: render.go 구현**

```go
package main

import (
	"fmt"
	"strings"
)

// renderMarkdown 은 한 API 의 md 문서를 생성한다 (파일 끝 개행 1개).
func renderMarkdown(ref APIRef, spec APISpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", ref.Name)
	fmt.Fprintf(&b, "> apiGrpCd: %s · apiId: %s\n", ref.GrpCd, ref.APIID)
	if ep := endpointSummary(spec.BasicInfo); ep != "" {
		fmt.Fprintf(&b, ">\n> %s\n", ep)
	}
	b.WriteString("\n")
	if ref.Desc != "" {
		fmt.Fprintf(&b, "%s\n\n", ref.Desc)
	}
	writeSection(&b, "기본 정보", spec.BasicInfo)
	writeSection(&b, "요청 인자", spec.Request)
	writeSection(&b, "응답 결과", spec.Response)
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// endpointSummary 는 기본 정보 표 첫 행에서 "`GET https://...`" 요약을 만든다.
func endpointSummary(t Table) string {
	if len(t.Rows) == 0 || len(t.Rows[0]) < 2 {
		return ""
	}
	return fmt.Sprintf("`%s %s`", t.Rows[0][0], t.Rows[0][1])
}

// writeSection 은 표가 비어있지 않으면 ## 헤더 + md 표를 출력한다.
func writeSection(b *strings.Builder, title string, t Table) {
	if len(t.Rows) == 0 {
		return
	}
	fmt.Fprintf(b, "## %s\n\n", title)
	b.WriteString(renderTable(t))
	b.WriteString("\n")
}

// renderTable 은 Table 을 GitHub markdown 표 문자열로 만든다 (각 행 끝 개행).
func renderTable(t Table) string {
	ncol := len(t.Headers)
	for _, r := range t.Rows {
		if len(r) > ncol {
			ncol = len(r)
		}
	}
	if ncol == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("| " + strings.Join(padRow(t.Headers, ncol), " | ") + " |\n")
	b.WriteString("|" + strings.Repeat(" --- |", ncol) + "\n")
	for _, r := range t.Rows {
		b.WriteString("| " + strings.Join(padRow(r, ncol), " | ") + " |\n")
	}
	return b.String()
}

// padRow 는 셀의 파이프를 이스케이프하고 행 길이를 ncol 로 맞춘다.
func padRow(cells []string, ncol int) []string {
	out := make([]string, ncol)
	for i := 0; i < ncol; i++ {
		if i < len(cells) {
			out[i] = strings.ReplaceAll(cells[i], "|", `\|`)
		}
	}
	return out
}
```

- [ ] **Step 5: 테스트 통과 확인**

Run: `go test ./scripts/crawl/ -run 'TestRenderMarkdown|TestRenderTableEscapesPipe' -v`
Expected: PASS. (실패 시 golden 파일의 공백/개행을 출력과 대조 — 파일 끝 개행 1개 확인.)

- [ ] **Step 6: Commit**

```bash
git add scripts/crawl/render.go scripts/crawl/render_test.go scripts/crawl/testdata/company.golden.md
git commit -m "feat(crawl): render API spec to markdown"
```

---

### Task 6: 인덱스 렌더러 (renderIndex)

**Files:**
- Modify: `scripts/crawl/render.go` (renderIndex 추가)
- Modify: `scripts/crawl/render_test.go` (TestRenderIndex 추가)

- [ ] **Step 1: 실패하는 테스트 추가**

`scripts/crawl/render_test.go` 에 추가:
```go
func TestRenderIndex(t *testing.T) {
	refs := []APIRef{
		{GrpCd: "DS001", Category: "공시정보", APIID: "2019002", Name: "기업개황", Desc: "개황정보"},
		{GrpCd: "DS001", Category: "공시정보", APIID: "2019001", Name: "공시검색", Desc: "검색"},
	}
	out := renderIndex(refs)
	assert.Contains(t, out, "## 공시정보 (DS001)")
	// apiId 오름차순 정렬: 공시검색(2019001) 이 기업개황(2019002) 보다 먼저
	assert.Less(t, strings.Index(out, "공시검색"), strings.Index(out, "기업개황"))
	assert.Contains(t, out, "[기업개황](공시정보/기업개황.md)")
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./scripts/crawl/ -run TestRenderIndex -v`
Expected: FAIL — `undefined: renderIndex`.

- [ ] **Step 3: render.go 에 renderIndex 추가**

```go
import "sort" // 파일 상단 import 블록에 추가

// renderIndex 는 docs/api/README.md 본문을 생성한다. 카테고리(GrpCd) → apiId 순 정렬.
func renderIndex(refs []APIRef) string {
	sorted := append([]APIRef(nil), refs...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].GrpCd != sorted[j].GrpCd {
			return sorted[i].GrpCd < sorted[j].GrpCd
		}
		return sorted[i].APIID < sorted[j].APIID
	})

	var b strings.Builder
	b.WriteString("# OpenDART API 문서\n\n")
	b.WriteString("OpenDART 개발가이드(https://opendart.fss.or.kr)에서 크롤링한 API 명세입니다.\n\n")
	b.WriteString("`go run ./scripts/crawl` 으로 재생성됩니다.\n")

	curGrp := ""
	for _, r := range sorted {
		if r.GrpCd != curGrp {
			fmt.Fprintf(&b, "\n## %s (%s)\n\n", r.Category, r.GrpCd)
			b.WriteString("| API | 설명 | 문서 |\n| --- | --- | --- |\n")
			curGrp = r.GrpCd
		}
		link := fmt.Sprintf("%s/%s.md", sanitize(r.Category), sanitize(r.Name))
		fmt.Fprintf(&b, "| %s | %s | [%s](%s) |\n",
			r.Name, strings.ReplaceAll(r.Desc, "|", `\|`), r.Name, link)
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}
```

> 주: `sanitize` 는 Task 7 의 main.go 에서 정의된다. 같은 `package main` 이라 컴파일 시점에 함께 보이므로 Task 7 완료 전에는 빌드가 `undefined: sanitize` 로 실패할 수 있다. 이 테스트만 먼저 통과시키려면 Step 3 직후 임시로 `func sanitize(s string) string { return s }` 를 render.go 하단에 두고, Task 7 Step 3 에서 main.go 로 옮기며 제거한다. (또는 Task 7 을 먼저 끝낸 뒤 Task 6 테스트를 실행한다.)

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./scripts/crawl/ -run TestRenderIndex -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add scripts/crawl/render.go scripts/crawl/render_test.go
git commit -m "feat(crawl): render docs/api index"
```

---

### Task 7: HTTP fetch + 오케스트레이션 (fetch.go / main.go)

**Files:**
- Create: `scripts/crawl/fetch.go`
- Create: `scripts/crawl/main.go`
- Modify: `scripts/crawl/render.go` (Task 6 의 임시 `sanitize` 가 있었다면 제거)

- [ ] **Step 1: fetch.go 작성**

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent   = "opendart-doc-crawler (github.com/kenshin579/opendart)"
	politeDelay = 300 * time.Millisecond
)

// httpGet 은 User-Agent 를 붙여 URL 본문을 문자열로 가져온다.
func httpGet(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
```

- [ ] **Step 2: main.go 작성**

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	listURLFmt   = "https://opendart.fss.or.kr/guide/main.do?apiGrpCd=%s"
	detailURLFmt = "https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=%s&apiId=%s"
	docsRoot     = "docs/api"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	client := &http.Client{Timeout: 30 * time.Second}
	var done []APIRef
	var failures []string

	for _, cat := range categories {
		listHTML, err := httpGet(client, fmt.Sprintf(listURLFmt, cat.Code))
		if err != nil {
			failures = append(failures, fmt.Sprintf("list %s: %v", cat.Code, err))
			continue
		}
		refs, err := parseList(listHTML)
		if err != nil {
			failures = append(failures, fmt.Sprintf("parse list %s: %v", cat.Code, err))
			continue
		}
		for i := range refs {
			refs[i].GrpCd = cat.Code
			refs[i].Category = cat.Name
		}
		time.Sleep(politeDelay)

		for _, ref := range refs {
			detailHTML, err := httpGet(client, fmt.Sprintf(detailURLFmt, ref.GrpCd, ref.APIID))
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s fetch: %v", ref.Category, ref.Name, err))
				continue
			}
			spec, err := parseDetail(detailHTML)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s parse: %v", ref.Category, ref.Name, err))
				continue
			}
			if err := writeDoc(ref, spec); err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s write: %v", ref.Category, ref.Name, err))
				continue
			}
			done = append(done, ref)
			fmt.Printf("✓ %s / %s\n", ref.Category, ref.Name)
			time.Sleep(politeDelay)
		}
	}

	if err := writeIndex(done); err != nil {
		return err
	}

	fmt.Printf("\n완료: %d개 API 문서 생성\n", len(done))
	if len(failures) > 0 {
		fmt.Printf("\n실패 %d건:\n", len(failures))
		for _, f := range failures {
			fmt.Println("  -", f)
		}
	}
	return nil
}

// writeDoc 은 한 API md 를 docs/api/{카테고리}/{이름}.md 로 쓴다.
func writeDoc(ref APIRef, spec APISpec) error {
	dir := filepath.Join(docsRoot, sanitize(ref.Category))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(
		filepath.Join(dir, sanitize(ref.Name)+".md"),
		[]byte(renderMarkdown(ref, spec)), 0o644,
	)
}

// writeIndex 는 docs/api/README.md 를 쓴다.
func writeIndex(refs []APIRef) error {
	if err := os.MkdirAll(docsRoot, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(docsRoot, "README.md"), []byte(renderIndex(refs)), 0o644)
}

// sanitize 는 파일/디렉토리명에 부적합한 문자를 치환한다.
func sanitize(name string) string {
	r := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", `"`, "_", "<", "_", ">", "_", "|", "_",
	)
	return strings.TrimSpace(r.Replace(name))
}
```

- [ ] **Step 3: Task 6 임시 sanitize 제거 (있을 경우)**

render.go 하단에 임시 `func sanitize` 를 넣었다면 제거한다 (이제 main.go 에 정의됨).

- [ ] **Step 4: 빌드 · vet · 전체 테스트**

Run:
```bash
go build ./scripts/crawl/
go vet ./scripts/crawl/
go test ./scripts/crawl/ -v
```
Expected: 빌드/vet 에러 없음, 모든 테스트 PASS (`sanitize` 중복 정의 없음).

- [ ] **Step 5: Commit**

```bash
git add scripts/crawl/fetch.go scripts/crawl/main.go scripts/crawl/render.go
git commit -m "feat(crawl): wire fetch + orchestration"
```

---

### Task 8: 전체 크롤링 실행 + 수동 검증 + 문서 커밋

**Files:**
- Create: `docs/api/**/*.md` (크롤러 산출물)
- Create: `docs/api/README.md`

- [ ] **Step 1: 전체 크롤링 실행**

Run (레포 루트에서): `go run ./scripts/crawl`
Expected: `✓ 공시정보 / 기업개황` 등 진행 로그, 마지막에 `완료: N개 API 문서 생성`. **실패 건이 0 이어야 함** (실패 목록이 출력되면 해당 페이지 구조를 점검).

- [ ] **Step 2: 산출물 구조 확인**

Run:
```bash
ls docs/api
find docs/api -name '*.md' | wc -l
cat docs/api/공시정보/기업개황.md
```
Expected: 6개 한글 카테고리 디렉토리 + `README.md`. 기업개황.md 가 golden 과 동일한 구조(제목/메타/요청 인자/응답 결과 표).

- [ ] **Step 3: 실제 포털과 대조 (수동)**

`docs/api/공시정보/기업개황.md` 를 https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS001&apiId=2019002 와 눈으로 대조: 요청 인자(crtfc_key, corp_code)와 응답 필드(corp_name, ceo_nm, est_dt 등)가 누락 없이 들어갔는지 확인.

- [ ] **Step 4: 인코딩 확인 (한글)**

Run: `file -I docs/api/공시정보/기업개황.md`
Expected: `charset=utf-8`.

- [ ] **Step 5: Commit (생성된 문서)**

```bash
git add docs/api
git commit -m "docs: crawl OpenDART API specs to markdown"
```

---

## Self-Review Notes

- **Spec coverage:** 디렉토리 네이밍(한글) = Task 7 sanitize/writeDoc · md 포맷 ①~④ = Task 5 · 목록 동적 수집 = Task 3 · goquery 파싱 = Task 4 · 인덱스 = Task 6 · 예의(delay)/실패 리포트 = Task 7 · golden 테스트 = Task 3~6 · full 실행 수동 검증 = Task 8. 모두 매핑됨. (메시지 설명/응답 예시 ⑤ 는 spec 에서 "있을 때만"으로 선택사항 — 전 API 공통 상태코드라 per-API 에서 제외, 후속 단일 문서로 가능.)
- **Type consistency:** `APIRef`/`Table`/`APISpec` 필드명, `parseList`/`parseDetail`/`extractTable`/`cellText`/`renderMarkdown`/`renderTable`/`padRow`/`endpointSummary`/`writeSection`/`renderIndex`/`httpGet`/`writeDoc`/`writeIndex`/`sanitize` 시그니처가 정의 태스크와 사용 태스크에서 일치.
- **sanitize 순환 주의:** Task 6 의 renderIndex 가 sanitize(main.go) 를 참조 → 같은 패키지지만 작성 순서상 임시 stub 처리 가이드를 Task 6 Step 3 주석 + Task 7 Step 3 에 명시.
