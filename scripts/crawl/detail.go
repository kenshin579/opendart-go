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
//
// 가정: 표는 flat 구조 — 헤더는 첫 행의 <th>, 데이터는 각 행의 <td> 이며
// colspan/rowspan(병합 셀)은 없다. OpenDART 개발가이드의 현재 표들은 모두
// 이 구조다(fixture 로 확인). 만약 향후 병합 셀이 도입되면 헤더 수와 행 셀 수가
// 어긋나 renderTable 이 컬럼을 잘못 채울 수 있으므로 그때 보강이 필요하다.
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
