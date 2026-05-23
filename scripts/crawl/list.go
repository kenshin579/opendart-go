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
