// examples/disclosure — DS001 공시정보 사용 예제.
//
// 실행: OPENDART_API_KEY=... go run ./examples/disclosure
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kenshin579/opendart"
	"github.com/kenshin579/opendart/disclosure"
)

func main() {
	client, err := opendart.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	// 종목코드 → corp_code
	corp, err := client.ResolveCorpCode(ctx, "005930")
	if err != nil {
		log.Fatalf("ResolveCorpCode: %v", err)
	}

	// 기업개황
	company, err := client.Disclosure.GetCompany(ctx, corp)
	if err != nil {
		log.Fatalf("GetCompany: %v", err)
	}
	fmt.Printf("회사명: %s (%s) 대표: %s 설립: %s\n",
		company.CorpName, company.StockCode, company.CeoName, company.EstDate)

	// 공시검색 (최근 5건)
	res, err := client.Disclosure.SearchDisclosures(ctx, disclosure.SearchParams{CorpCode: corp, PageCount: 5})
	if err != nil {
		log.Fatalf("SearchDisclosures: %v", err)
	}
	fmt.Printf("총 %d건 중 %d건:\n", res.TotalCount, len(res.List))
	for _, d := range res.List {
		fmt.Printf("  [%s] %s (%s)\n", d.RceptDt, d.ReportNm, d.RceptNo)
	}
}
