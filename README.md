# opendart

DART 전자공시시스템 OpenAPI 의 Go 클라이언트 라이브러리.

## 설치

```bash
go get github.com/kenshin579/opendart
```

## 사용

```go
client, _ := opendart.NewClientFromEnv() // OPENDART_API_KEY
ctx := context.Background()

corp, _ := client.ResolveCorpCode(ctx, "005930")        // 종목코드 → corp_code
company, _ := client.Disclosure.GetCompany(ctx, corp)   // 기업개황
res, _ := client.Disclosure.SearchDisclosures(ctx, disclosure.SearchParams{CorpCode: corp})
zip, _ := client.Disclosure.DownloadDocument(ctx, "20240131000326") // 원본 ZIP
```

## 인증

발급받은 API 키를 `OPENDART_API_KEY` 환경변수로 두거나 `opendart.NewClient(apiKey)` 로 전달한다.

## 커버리지

- DS001 공시정보: 기업개황 · 공시검색 · 공시서류원본파일 · 고유번호(corp_code 매핑)
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황
- (예정) DS002 나머지 · DS003~DS006

## 에러 처리

- `errors.Is(err, opendart.ErrNoData)` — 조회 데이터 없음(013)
- `errors.As(err, &apiErr)` — 그 외 OpenDART status (`*opendart.APIError`)

## 문서

API 명세: [`docs/api/`](docs/api/README.md)
