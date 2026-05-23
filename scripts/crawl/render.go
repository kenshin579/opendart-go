package main

import (
	"fmt"
	"sort"
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
		// 링크 대상은 <> 로 감싼다 — 카테고리/파일명에 공백·괄호가 있어도
		// (예: "정기보고서 주요정보/증자(감자) 현황.md") GitHub/CommonMark 에서
		// 링크가 깨지지 않도록 한다.
		fmt.Fprintf(&b, "| %s | %s | [%s](<%s>) |\n",
			r.Name, strings.ReplaceAll(r.Desc, "|", `\|`), r.Name, docRelPath(r))
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// docRelPath 는 docs/api 기준 상대 문서 경로 "{카테고리}/{이름}.md" 를 반환한다
// (항상 "/" 구분자 — markdown 링크와 파일 쓰기가 같은 규칙을 공유하도록).
func docRelPath(ref APIRef) string {
	return sanitize(ref.Category) + "/" + sanitize(ref.Name) + ".md"
}

// sanitize 는 파일/디렉토리명에 부적합한 문자를 치환한다.
func sanitize(name string) string {
	r := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", `"`, "_", "<", "_", ">", "_", "|", "_",
	)
	return strings.TrimSpace(r.Replace(name))
}
