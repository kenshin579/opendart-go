# opendart-go

DART 전자공시시스템 OpenAPI 의 Go 클라이언트 라이브러리.

## 설치

```bash
go get github.com/kenshin579/opendart-go@latest
```

## 사용

```go
client, _ := opendart.NewClientFromEnv() // OPENDART_API_KEY
ctx := context.Background()

corp, _ := client.ResolveCorpCode(ctx, "005930")        // 종목코드 → corp_code
company, _ := client.Disclosure.GetCompany(ctx, corp)   // 기업개황 (DS001)
res, _ := client.Disclosure.SearchDisclosures(ctx, disclosure.SearchParams{CorpCode: corp})
zip, _ := client.Disclosure.DownloadDocument(ctx, "20240131000326") // 원본 ZIP

// 재무정보(DS003)는 client.Report 로 접근
acnt, _ := client.Report.SingleAccount(ctx, report.ReportParams{
    CorpCode: corp, BsnsYear: "2023", ReprtCode: report.AnnualReport,
}) // 단일회사 주요계정
```

## 인증

발급받은 API 키를 `OPENDART_API_KEY` 환경변수로 두거나 `opendart.NewClient(apiKey)` 로 전달한다.

## 커버리지

| 구분 | 영역 | 서비스 |
|------|------|--------|
| DS001 | 공시정보 | `client.Disclosure` |
| DS002 | 정기보고서 주요정보 | `client.Report` |
| DS003 | 정기보고서 재무정보 | `client.Report` |
| DS004 | 지분공시 종합정보 | `client.Ownership` |
| DS005 | 주요사항보고서 주요정보 | `client.Material` |
| DS006 | 증권신고서 주요정보 | `client.Registration` |

DS001~DS006 의 세부 엔드포인트를 거의 모두 제공한다. 각 항목별 요청·응답 명세는 아래 문서를 참고한다.

- 공식 명세: [OpenDART 개발가이드](https://opendart.fss.or.kr/intro/main.do)
- 로컬 명세(크롤링본): [`docs/api/`](docs/api/README.md)

> (예정) DS002 개인별 보수 Ver2.0

## 에러 처리

- `errors.Is(err, opendart.ErrNoData)` — 조회 데이터 없음(013)
- `errors.As(err, &apiErr)` — 그 외 OpenDART status (`*opendart.APIError`)

## 예제

실행 가능한 예제: [`examples/disclosure`](examples/disclosure/main.go) · [`examples/report`](examples/report/main.go)

```bash
OPENDART_API_KEY=... go run ./examples/disclosure
```

## License

Released under the [MIT License](LICENSE).
