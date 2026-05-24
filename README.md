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
- DS002 정기보고서 주요정보: 증자(감자) · 배당 · 자기주식 · 주식총수 · 최대주주 · 최대주주변동 · 소액주주 현황 · 증권 발행실적 · 미상환 잔액(회사채/기업어음/단기사채/신종자본증권/조건부자본증권) · 감사의견 · 감사/비감사용역 · 타법인 출자 · 공모/사모자금 사용내역 · 임원/직원 현황 · 임원·이사·감사 보수현황
- DS003 정기보고서 재무정보: 단일/다중회사 주요계정 · 단일회사 전체 재무제표(개별/연결) · 단일/다중회사 주요 재무지표 · XBRL 택사노미 양식 · 재무제표 원본파일(XBRL)
- DS004 지분공시 종합정보: 대량보유 상황보고(5% 룰) · 임원·주요주주 소유보고
- DS005 주요사항보고서 주요정보: 부도발생 · 영업정지 · 회생절차 개시신청 · 해산사유 발생 · 채권은행 관리절차 개시/중단 · 소송 등의 제기 · 유상/무상/유무상 증자 결정 · 감자 결정 · 사채 발행(전환사채/신주인수권부사채/교환사채/상각형 조건부자본증권 발행결정)
- (예정) DS005 나머지(자기주식/양수도/합병·분할/해외상장) · DS006 · DS002 개인별 보수 Ver2.0

## 에러 처리

- `errors.Is(err, opendart.ErrNoData)` — 조회 데이터 없음(013)
- `errors.As(err, &apiErr)` — 그 외 OpenDART status (`*opendart.APIError`)

## 문서

API 명세: [`docs/api/`](docs/api/README.md)
