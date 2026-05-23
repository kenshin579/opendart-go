# OpenDART 라이브러리 본체 설계 — Foundation + DS001 (공시정보)

- 작성일: 2026-05-23
- 모듈: `github.com/kenshin579/opendart`
- 범위: **라이브러리 기반(Client/인증/HTTP/에러/corp_code 인프라) + DS001 공시정보 카테고리 수직 슬라이스**

## 배경 & 목표

`korea-investment-stock`(KIS) 처럼 다른 개발자도 사용할 수 있는 공개 Go 라이브러리.
Phase 1 에서 전체 API 명세를 크롤링해 `docs/api/` 에 확보했다(PR #1, main merged). 이 spec 은
그 docs 를 source of truth 로 삼아 **라이브러리 본체의 기반**을 세우고, **DS001 공시정보**
카테고리를 첫 구현 슬라이스로 끝까지 완성한다. 나머지 5개 카테고리(DS002~DS006)는 동일
패턴을 재사용하는 후속 plan 으로 진행한다.

설계 접근: **KIS 구조를 따르되 단순화** — OpenDART 는 OAuth/토큰/계좌가 없고 단일 API key
(`crtfc_key`) 쿼리 파라미터, 전부 GET, JSON(+일부 ZIP) 응답이므로 KIS 의 token/ratelimit
계층을 들어낸다.

## OpenDART API 표면 요약 (docs 기반 사실)

- 인증: 모든 요청에 `crtfc_key`(40자리) 쿼리 파라미터 1개.
- 베이스 URL: `https://opendart.fss.or.kr`. 각 API 는 `.json` / `.xml` 두 엔드포인트 — **우리는 `.json` 만 사용**.
- 공통 envelope: 모든 JSON 응답에 `status`(`"000"`=정상) + `message`.
- 응답 3종: ① 단일 객체(기업개황) ② list+페이지네이션(공시검색) ③ 바이너리 ZIP(corpCode.xml, 공시서류원본파일 등).
- 회사별 API 는 종목코드가 아니라 **corp_code(8자리)** 를 요구. 고유번호 API 가 전체 매핑을 ZIP(→CORPCODE.xml)으로 제공.
- 호출 한도: **일일 20,000건**(초당 아님). status `020`=요청제한초과.
- 서버 quirk: **TLS1.2 + RSA 키교환 cipher 전용**. Go 기본 HTTP 클라이언트는 handshake 실패 → cipher 명시 필요(Phase 1 크롤러에서 검증).

## 모듈 레이아웃

```
opendart/
  client.go            # Client + NewClient + sub-client 와이어링 (package opendart)
  config.go            # functional options
  from_env.go          # NewClientFromEnv (OPENDART_API_KEY)
  errors.go            # APIError + ErrNoData + ErrCorpCodeNotFound + status→error 매핑
  corpcode.go          # root client 메서드: ResolveCorpCode/LookupCorpCode/CorpCodes/RefreshCorpCodes
  disclosure/          # DS001 공시정보 sub-client
    client.go          #   sub-client 구조체 + New
    company.go         #   기업개황 GetCompany
    search.go          #   공시검색 SearchDisclosures
    document.go        #   공시서류원본파일 DownloadDocument
    *_test.go
    testdata/          #   응답 JSON / ZIP fixture
  internal/
    httpclient/        # GET + crtfc_key 주입 + JSON 디코드 + envelope status 검사 + TLS RSA cipher
      *_test.go
    corpcode/          # ZIP 다운로드→unzip→XML 파싱→디스크 캐시(TTL)→인덱스
      testdata/
      *_test.go
  examples/disclosure/ # DS001 사용 예제
  docs/api/...         # 크롤링된 명세 (구현 레퍼런스)
```

**호출 스타일:**
```go
client, _ := opendart.NewClientFromEnv()                    // OPENDART_API_KEY
corp, _ := client.ResolveCorpCode(ctx, "005930")            // 종목코드 → corp_code
company, _ := client.Disclosure.GetCompany(ctx, corp)       // 기업개황
res, _ := client.Disclosure.SearchDisclosures(ctx, params)  // 공시검색 (list)
zip, _ := client.Disclosure.DownloadDocument(ctx, rceptNo)  // 원본 ZIP []byte
```

## Client 생성 & 옵션

```go
func NewClient(apiKey string, opts ...Option) (*Client, error)   // apiKey 빈 값이면 에러
func NewClientFromEnv(opts ...Option) (*Client, error)           // OPENDART_API_KEY 읽음

// functional options
WithHTTPClient(*http.Client)          // 사용자 정의 클라이언트 (기본: TLS RSA cipher 내장 클라이언트)
WithBaseURL(string)                   // 기본 https://opendart.fss.or.kr (테스트 override 용)
WithTimeout(time.Duration)            // 기본 30s
WithCorpCodeCacheDir(string)          // 기본 os.UserCacheDir()/opendart, 실패 시 os.TempDir()/opendart
WithCorpCodeCacheTTL(time.Duration)   // 기본 24h
```

```go
type Client struct {
    apiKey     string
    http       *httpclient.Client
    corp       *corpcode.Cache
    Disclosure *disclosure.Client
    // 후속: Report, Financial, Ownership, Material, Registration
}
```

## HTTP / envelope 계층 (`internal/httpclient`)

모든 GET 요청의 단일 통로:

- `crtfc_key` 를 쿼리 파라미터로 자동 주입. base URL + `.json` 경로.
- 내장 `*http.Client` 는 **표준 `net/http`** 기반(resty 미사용 — 단순 GET 뿐이라 YAGNI). TLS 설정:
  `MinVersion: tls.VersionTLS12`, `CipherSuites: {TLS_RSA_WITH_AES_128_GCM_SHA256, TLS_RSA_WITH_AES_256_GCM_SHA384}`.
  `WithHTTPClient` 로 교체 가능.
- 응답 JSON 디코드 후 **envelope `status` 검사**: `000`→정상 / `013`→`ErrNoData` / 그 외→`*APIError{Status, Message}`.
- envelope 디코딩: 공용 `statusEnvelope` 를 각 응답 타입에 임베드. 임베드를 통해 작은
  인터페이스를 만족시켜 httpclient 가 status/message 를 읽는다(타입별 재디코드 불필요):
```go
type statusEnvelope struct {
    Status  string `json:"status"`  // 000=정상
    Message string `json:"message"`
}
func (e statusEnvelope) apiStatus() (status, message string) { return e.Status, e.Message }

type statusReader interface{ apiStatus() (string, string) } // 임베드로 자동 구현

// GetJSON 은 out 으로 디코드 후 out.apiStatus() 로 status 검사 → 에러 변환.
func (c *Client) GetJSON(ctx context.Context, path string, params map[string]string, out statusReader) error
func (c *Client) GetBytes(ctx context.Context, path string, params map[string]string) ([]byte, error)
```
`statusEnvelope`/`statusReader` 는 `internal/httpclient` 패키지에 정의(임베드용으로 export, 예:
`httpclient.Envelope`)하고, 응답 타입(`disclosure.Company`, `disclosure.SearchResult`)이 이를
임베드해 `GetJSON` 에 전달한다. status→error 변환(`013`→`ErrNoData`, 그 외→`*APIError`)은
httpclient 가 수행하되, sentinel/`APIError` 타입은 root `opendart` 패키지에 두므로 httpclient 는
변환 콜백(또는 `errors.go` 가 주입한 매핑 함수)을 통해 root 패키지의 에러를 생성한다 — internal→root
import 순환을 피하기 위함.

**rate limiting 없음 / 응답 캐싱 없음** — 일일 쿼터·응답 캐싱은 호출자 책임(KIS 와 동일). 무거운 정적
파일(corp_code)만 디스크 캐시.

## corp_code 인프라 (`internal/corpcode` + corpcode.go)

KIS `internal/mastercache` 모델.

`internal/corpcode.Cache`:
- 다운로드: `GET /api/corpCode.xml?crtfc_key=...` → ZIP 바이트.
- 압축해제: `archive/zip` 로 ZIP 안의 `CORPCODE.xml` 추출.
- 파싱: `encoding/xml` 로 `<list>` 엔트리 → `[]Entry`.
- 디스크 캐시: 원본 ZIP 저장, TTL(기본 24h) 경과 시 재다운로드. **다운로드 실패 시 옛 캐시 fallback**.
- 인메모리 인덱스: 최초 로드 시 `stockCode→Entry`, `corpCode→Entry` 맵 구축(`sync.Once`/mutex, lazy 1회). 빈 stock_code(비상장)는 stockCode 맵에서 제외.
- 주입형 `FetchFunc` 로 다운로드를 분리 → 네트워크 없이 테스트.

```go
type Entry struct {
    CorpCode    string `xml:"corp_code"`     // 고유번호 (8자리)
    CorpName    string `xml:"corp_name"`     // 정식회사명칭
    CorpEngName string `xml:"corp_eng_name"` // 영문 정식명칭
    StockCode   string `xml:"stock_code"`    // 종목코드 (6자리, 비상장은 빈 값)
    ModifyDate  string `xml:"modify_date"`   // 최종변경일자 YYYYMMDD
}
```

root `Client` 메서드(corpcode.go):
```go
func (c *Client) ResolveCorpCode(ctx context.Context, stockCode string) (string, error) // 없으면 ErrCorpCodeNotFound
func (c *Client) LookupCorpCode(ctx context.Context, corpCode string) (corpcode.Entry, error)
func (c *Client) CorpCodes(ctx context.Context) ([]corpcode.Entry, error)               // 전체 (사용자 직접 필터)
func (c *Client) RefreshCorpCodes(ctx context.Context) error                            // TTL 무시 강제 재다운로드
```

## DS001 공시정보 sub-client (`disclosure/`)

### ① 기업개황 — `GetCompany(ctx, corpCode string) (*Company, error)`

`GET /api/company.json` (params: corp_code).

```go
type Company struct {
    CorpName    string `json:"corp_name"`     // 정식명칭
    CorpNameEng string `json:"corp_name_eng"` // 영문명칭
    StockName   string `json:"stock_name"`    // 종목명(상장사)/약식명칭(기타)
    StockCode   string `json:"stock_code"`    // 종목코드 (6자리)
    CeoName     string `json:"ceo_nm"`        // 대표자명
    CorpCls     string `json:"corp_cls"`      // 법인구분 Y(유가)/K(코스닥)/N(코넥스)/E(기타)
    JurirNo     string `json:"jurir_no"`      // 법인등록번호
    BizrNo      string `json:"bizr_no"`       // 사업자등록번호
    Address     string `json:"adres"`         // 주소
    HomeURL     string `json:"hm_url"`        // 홈페이지
    IRURL       string `json:"ir_url"`        // IR 홈페이지
    PhoneNo     string `json:"phn_no"`        // 전화번호
    FaxNo       string `json:"fax_no"`        // 팩스번호
    IndutyCode  string `json:"induty_code"`   // 업종코드
    EstDate     string `json:"est_dt"`        // 설립일 YYYYMMDD
    AccMonth    string `json:"acc_mt"`        // 결산월 MM
}
```

### ② 공시검색 — `SearchDisclosures(ctx, params SearchParams) (*SearchResult, error)`

`GET /api/list.json`. 빈 값/0 은 쿼리에서 생략(API 기본값 사용).

```go
type SearchParams struct {
    CorpCode       string // 고유번호 (8자리)
    BgnDe          string // 검색 시작 접수일자 YYYYMMDD
    EndDe          string // 검색 종료 접수일자 YYYYMMDD
    LastReprtAt    string // 최종보고서만 검색 여부 (Y/N, 기본 N)
    PblntfTy       string // 공시유형 (A 정기/B 주요사항/C 발행/D 지분/E 기타/F 외부감사/G 펀드/H 자산유동화/I 거래소/J 공정위)
    PblntfDetailTy string // 공시상세유형 (4자리 코드)
    CorpCls        string // 법인구분 (Y/K/N/E)
    Sort           string // 정렬 (date/crp/rpt, 기본 date)
    SortMth        string // 정렬방법 (asc/desc, 기본 desc)
    PageNo         int    // 페이지 번호 (1~n, 0이면 생략→기본 1)
    PageCount      int    // 페이지당 건수 (1~100, 0이면 생략→기본 10)
}

type SearchResult struct {
    PageNo     int              `json:"page_no"`     // 페이지 번호
    PageCount  int              `json:"page_count"`  // 페이지당 건수
    TotalCount int              `json:"total_count"` // 총 건수
    TotalPage  int              `json:"total_page"`  // 총 페이지 수
    List       []DisclosureItem `json:"list"`        // 공시 목록
}

type DisclosureItem struct {
    CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
    CorpName  string `json:"corp_name"`  // 종목명(법인명)
    CorpCode  string `json:"corp_code"`  // 고유번호 (8자리)
    StockCode string `json:"stock_code"` // 종목코드 (6자리)
    ReportNm  string `json:"report_nm"`  // 보고서명
    RceptNo   string `json:"rcept_no"`   // 접수번호 (DownloadDocument 인자)
    FlrNm     string `json:"flr_nm"`     // 공시 제출인명
    RceptDt   string `json:"rcept_dt"`   // 접수일자 YYYYMMDD
    Rm        string `json:"rm"`         // 비고
}
```

> 구현 메모: `list.json` 의 numeric envelope 필드(page_no/total_count 등)가 JSON number 인지
> string 인지 fixture 로 확정한다. string 이면 `,string` 태그를 붙인다.

### ③ 공시서류원본파일 — `DownloadDocument(ctx, rceptNo string) ([]byte, error)`

`GET /api/document.xml` (params: rcept_no) → ZIP 원본 바이트 그대로 반환. 압축해제·파싱은 호출자
몫(임의 공시 원본은 형태가 다양한 바이너리). `GetBytes` 사용. 단, 바이너리 응답이 아니라 에러
JSON(status≠000)일 수 있으므로, 응답이 JSON 에러 envelope 면 `*APIError` 로 변환한다.

## 에러 처리 (errors.go)

```go
type APIError struct {
    Status  string // OpenDART status 코드 (000 외)
    Message string // OpenDART message
}
func (e *APIError) Error() string // 예: "opendart: [020] 요청 제한을 초과하였습니다."

var ErrNoData = errors.New("opendart: no data (013)")               // status 013 — httpclient 가 변환
var ErrCorpCodeNotFound = errors.New("opendart: corp_code not found") // 매핑 실패
```

- 호출자: `errors.Is(err, opendart.ErrNoData)` / `errors.As(err, &apiErr)` 로 코드 접근.

## 컨벤션

- 모든 응답 struct 필드에 한글 설명 코멘트 (docs/api 명칭·설명 반영).
- 도메인 코드 주석은 한국어, godoc 노출 함수는 한 줄 요약(한국어).
- 표준 `net/http`(httpclient 내부). resty 미사용.
- 파라미터 맵 빌드는 내부 헬퍼로 빈 값 omit.
- 모든 파일 UTF-8.

## 테스트 전략

- **httpclient**: `httptest`/`httpmock` 로 200 정상·013(ErrNoData)·기타 status(APIError)·비-200 검증.
- **disclosure 각 메서드**: 실제 응답 JSON 을 `testdata/` fixture 로 두고 디코딩·필드 매핑·파라미터 빌드 검증(네트워크 없음).
- **corp_code**: 작은 CORPCODE.xml ZIP fixture + 주입형 `FetchFunc` 로 파싱·인덱스·TTL·stale fallback.
- **통합 테스트(선택)**: 실제 `OPENDART_API_KEY` 필요 테스트는 build tag `//go:build integration` 으로 분리 → 기본 `go test` 에서 제외.

## 산출물 (이 슬라이스)

Client/options/from_env/errors + `internal/httpclient` + `internal/corpcode` + corp_code root 메서드
+ `disclosure`(기업개황/공시검색/원본파일) + 단위 테스트 + `examples/disclosure` + README.

## 비범위 (후속 plan)

DS002~DS006 카테고리, 숫자-문자열 coercion 헬퍼(재무정보에서 도입), XML 출력, 응답 캐싱,
rps rate limiter.
