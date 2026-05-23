# OpenDART 문서 크롤러 설계 (Phase 1)

- 작성일: 2026-05-23
- 모듈: `github.com/kenshin579/opendart`
- 범위: **OpenDART 개발가이드 전체 API 명세를 크롤링하여 markdown 으로 변환**

## 배경 & 목표

`korea-investment-stock`(KIS) 라이브러리처럼, 다른 개발자도 사용할 수 있는 공개 Go
라이브러리 `opendart` 를 개발한다. 최종 목표는 OpenDART(전자공시시스템 DART OpenAPI)의
**전체 API** 를 커버하는 라이브러리이지만, 작업 순서상 **문서를 먼저 크롤링해 md 로 확보**
하고, 그 md 를 source of truth 로 삼아 라이브러리 구조를 이후 별도 세션에서 설계·구현한다.

이 spec 은 그 첫 단계인 **문서 크롤러(Phase 1)** 만 다룬다.

## 비범위 (Out of Scope)

라이브러리 본체는 이 세션에서 다루지 않는다. docs 확보 후 별도 brainstorming 으로 진행:

- `Client` / sub-client (`Disclosure`, `Financial` 등) / 응답 타입
- 단일 API key 인증 래퍼 (`OPENDART_API_KEY`)
- corp_code 인프라 (고유번호 ZIP→CORPCODE.xml 다운로드·캐시 + 종목코드↔corp_code 매핑 헬퍼)
- examples

> 참고용 north star (확정 아님): KIS 와 동일하게 `client.Disclosure.GetCompany(ctx, corpCode)`
> 같은 1-level 그룹핑 호출 스타일, `internal/httpclient`, `internal/corpcode`, 도메인별
> sub-package (`disclosure/`, `report/`, `financial/`, `ownership/`, `material/`,
> `registration/`) 구조를 따를 가능성이 높다. 단, 실제 설계는 크롤링된 docs 를 보고 결정한다.

## OpenDART 문서 구조 (크롤링 대상)

개발가이드는 6개 카테고리(apiGrpCd)로 구성된다:

| apiGrpCd | 카테고리 (한글) |
|----------|-----------------|
| DS001 | 공시정보 |
| DS002 | 정기보고서 주요정보 |
| DS003 | 정기보고서 재무정보 |
| DS004 | 지분공시 종합정보 |
| DS005 | 주요사항보고서 주요정보 |
| DS006 | 증권신고서 주요정보 |

- **카테고리 목록 페이지**: `https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS00X`
  — 각 API 의 detail 링크(`apiId`)와 한글명을 포함.
- **API 상세 페이지**: `https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS00X&apiId=YYYYYYY`
  — **서버 렌더링 HTML**(검증 완료). 다음을 포함:
  - 제목 / 개요 설명
  - 요청 URL (`.json` / `.xml` 두 엔드포인트)
  - 요청 파라미터 표 (파라미터·한글명·타입·필수·설명)
  - 응답 결과 필드 표 (필드·한글명·타입·설명)
  - 응답 예시는 페이지에 없는 경우가 많음(테스트 도구만 제공)

## 산출물 디렉토리/포맷

```
opendart/
  go.mod                         # module github.com/kenshin579/opendart
  scripts/
    crawl/
      main.go                    # 크롤러 진입점
      ...                        # 파싱/렌더링 로직
      testdata/                  # 대표 detail.do HTML fixture + golden md
  docs/
    api/
      README.md                  # 전체 API 인덱스 (자동 생성)
      공시정보/
        기업개황.md
        공시검색.md
        ...
      정기보고서 주요정보/
        ...
      정기보고서 재무정보/
        ...
      지분공시 종합정보/
        ...
      주요사항보고서 주요정보/
        ...
      증권신고서 주요정보/
        ...
    superpowers/specs/           # 본 문서 등
```

- **디렉토리 네이밍**: 한글 카테고리명 (사람이 보기 직관적, KIS `docs/api` 와 동일한 결).
- **파일명**: API 한글명 기준 (`기업개황.md`).
- 모든 파일 UTF-8.

### 각 API md 포맷 (① ~ ⑤)

```markdown
# 기업개황

> apiGrpCd: DS001 · apiId: 2019002
> 엔드포인트: `GET https://opendart.fss.or.kr/api/company.json` (`.xml` 도 제공)

DART 에 등록되어있는 기업의 개황정보를 제공합니다.   # ① 개요

## 요청 파라미터                                       # ③

| 파라미터 | 한글명 | 타입 | 필수 | 설명 |
|----------|--------|------|------|------|
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |
| corp_code | 고유번호 | STRING(8) | Y | 공시대상회사의 고유번호(8자리) |

## 응답 필드                                           # ④

| 필드 | 한글명 | 타입 | 설명 |
|------|--------|------|------|
| status | 에러 및 정보 코드 | STRING | 정상: 000, 그 외 에러 |
| corp_name | 정식명칭 | STRING | 정식회사 명칭 |
| ... | ... | ... | ... |

## 응답 예시                                           # ⑤ (있을 때만)
...
```

## 크롤러 동작

`go run ./scripts/crawl` 으로 실행하는 단일 Go 프로그램.

1. **API 목록 수집** — 6개 카테고리 목록 페이지를 받아 각 API 의 `(apiGrpCd, apiId, 한글명)`
   을 동적 추출. (하드코딩 seed 대신 목록에서 추출 → 신규 API 자동 포함)
2. **상세 페이지 파싱** — 각 `detail.do` 를 받아 `goquery`(`github.com/PuerkitoBio/goquery`)
   로 ①~⑤ 추출.
3. **md 렌더링** — 깔끔한 GitHub markdown 표로 작성. 상단 메타데이터(apiGrpCd, apiId,
   엔드포인트) 명시.
4. **파일 쓰기** — `docs/api/{한글카테고리}/{API한글명}.md`. 카테고리 코드→한글명 매핑은
   크롤러 내 상수. 재실행 시 멱등하게 재생성.
5. **인덱스 생성** — `docs/api/README.md` 에 전체 API 목록 표(카테고리/이름/엔드포인트/링크)
   자동 생성.

### 예의 / 견고성

- 요청 간 짧은 delay(약 300ms), 명시적 User-Agent.
- 개별 페이지 파싱 실패 시 **전체 중단하지 않고** 스킵 → 실행 끝에 실패 목록 리포트.

## 테스트 & 검증

- **단위 테스트**: 대표 detail.do 페이지 HTML 을 `scripts/crawl/testdata/` 에 fixture 로
  저장 → 파싱 결과를 golden md 와 비교. 네트워크 없이 `go test ./...` 로 회귀 검증.
- **수동 검증**: 크롤러 1회 full 실행 후
  1. 6개 카테고리 디렉토리가 모두 생성되었는지
  2. 기업개황 등 핵심 API md 를 실제 포털과 눈으로 대조
  3. 실패 리포트가 비었는지

## 의존성

- Go 1.25+
- `github.com/PuerkitoBio/goquery` (HTML 파싱)
- `github.com/stretchr/testify` (test)
