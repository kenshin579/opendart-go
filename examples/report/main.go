// examples/report — DS002 정기보고서 주요정보(지분·주식·배당) 사용 예제.
//
// 실행: OPENDART_API_KEY=... go run ./examples/report
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kenshin579/opendart-go"
	"github.com/kenshin579/opendart-go/report"
)

func main() {
	client, err := opendart.NewClientFromEnv()
	if err != nil {
		log.Fatalf("NewClientFromEnv: %v", err)
	}
	ctx := context.Background()

	corp, err := client.ResolveCorpCode(ctx, "005930")
	if err != nil {
		log.Fatalf("ResolveCorpCode: %v", err)
	}

	p := report.ReportParams{CorpCode: corp, BsnsYear: "2023", ReprtCode: report.AnnualReport}

	// 배당에 관한 사항
	dividends, err := client.Report.Dividend(ctx, p)
	if errors.Is(err, opendart.ErrNoData) {
		fmt.Println("배당 데이터 없음")
	} else if err != nil {
		log.Fatalf("Dividend: %v", err)
	} else {
		fmt.Printf("배당 항목 %d건:\n", len(dividends))
		for _, d := range dividends {
			fmt.Printf("  %s (%s): 당기 %s\n", d.Se, d.StockKnd, d.Thstrm)
		}
	}

	// 최대주주 현황
	majors, err := client.Report.MajorShareholders(ctx, p)
	if errors.Is(err, opendart.ErrNoData) {
		fmt.Println("최대주주 데이터 없음")
	} else if err != nil {
		log.Fatalf("MajorShareholders: %v", err)
	} else {
		fmt.Printf("최대주주 %d명:\n", len(majors))
		for _, m := range majors {
			fmt.Printf("  %s (%s): 지분 %s%%\n", m.Nm, m.Relate, m.TrmendPosesnStockQotaRt)
		}
	}
}
