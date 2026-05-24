# OpenDART DS006 증권신고서 주요정보 — Sub-1 «증권» 설계

- 작성일: 2026-05-25
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS006 증권신고서 Sub-1 3개 API** (신규 `registration` 패키지)

## 배경 & 목표

DS005 완료(36/36). DS006 증권신고서는 신규 카테고리(6개). 응답이 DS005 의 평면 `list[]` 와 달리
**그룹형**(`{status, message, group:[{title, list:[...]}]}`, group 은 JSON 배열 — 실 API 확인됨,
남양유업 estkRs 2023). 24개 그룹 스키마라 2 sub-group 분할: 이 spec 은 **Sub-1 «증권» 3개**
(지분증권/채무증권/증권예탁증권). Sub-2 «신고»(합병/분할/교환이전)는 후속. 신규 패키지 `registration`,
root `client.Registration` 와이어링.

## 새 인프라

### internal/httpclient.GetGroups (신규)
DS005 의 `GetList[T]` 와 별개로, 그룹형 응답 디코더를 추가한다(status 엔벨로프 재사용 → 013→ErrNoData).
```go
// Group 은 증권신고서(DS006) 응답의 그룹 하나(title + raw list).
type Group struct {
	Title string          `json:"title"` // 그룹명칭(예: 일반사항)
	List  json.RawMessage `json:"list"`  // 그룹 항목 배열(타입은 호출측에서 결정)
}

type groupEnvelope struct {
	Envelope        // status/message + APIStatus() (013→ErrNoData, 그 외→APIError)
	Group []Group `json:"group"`
}

// GetGroups 는 그룹형(DS006) 응답을 디코딩해 그룹 목록을 반환한다.
func GetGroups(ctx context.Context, c *Client, path string, params map[string]string) ([]Group, error) {
	var env groupEnvelope
	if err := c.GetJSON(ctx, path, params, &env); err != nil {
		return nil, err
	}
	return env.Group, nil
}
```
(`Envelope`/`GetJSON`/`StatusReader` 는 기존 그대로 재사용. `groupEnvelope` 가 `APIStatus()` 를 임베드로 만족.)

### registration 패키지 (신규)
```
registration/
  client.go             # Client{http *httpclient.Client}, New, Params + toMap, 그룹 디코드 헬퍼
  securities.go         # Sub-1 3개 메서드 + wrapper/item struct
  client_test.go        # newTestClient (httptest, testdata fixture)
  securities_test.go    # 3개 fixture 테스트
  testdata/             # 3개 fixture
```
root `client.go`: `Registration *registration.Client` 필드 + `c.Registration = registration.New(hc)`.

`registration/client.go`:
```go
package registration

type Client struct{ http *httpclient.Client }
func New(http *httpclient.Client) *Client { return &Client{http: http} }

// Params 는 DS006 증권신고서 공통 요청 파라미터(corp_code + 기간).
type Params struct {
	CorpCode string // 고유번호(8자리)
	BgnDe    string // 검색시작 접수일자(YYYYMMDD)
	EndDe    string // 검색종료 접수일자(YYYYMMDD)
}
func (p Params) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" { m["bgn_de"] = p.BgnDe }
	if p.EndDe != "" { m["end_de"] = p.EndDe }
	return m
}

// unmarshalGroup 은 title 이 일치하는 그룹의 list 를 dst(타입 슬라이스 포인터)에 언마샬한다.
func unmarshalGroup(groups []httpclient.Group, title string, dst any) error {
	for _, g := range groups {
		if g.Title == title {
			return json.Unmarshal(g.List, dst)
		}
	}
	return nil // 그룹 없으면 빈 슬라이스 유지
}
```

## API 표면 (registration/securities.go)

| 메서드 | 한글 | 엔드포인트 | 반환 |
|--------|------|-----------|------|
| `EquitySecurities` | 지분증권 | `/api/estkRs.json` | `*EquitySecuritiesRegistration` |
| `DebtSecurities` | 채무증권 | `/api/bdRs.json` | `*DebtSecuritiesRegistration` |
| `DepositaryReceipts` | 증권예탁증권 | `/api/stkdpRs.json` | `*DepositaryReceiptsRegistration` |

**공유 item 타입**(지분증권·증권예탁증권의 5개 공통 그룹은 스키마 동일 → 공유):
```go
// RsGeneralItem 은 증권신고서 일반사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsGeneralItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Sbd      string `json:"sbd"`       // 청약기일
	Pymd     string `json:"pymd"`      // 납입기일
	Sband    string `json:"sband"`     // 청약공고일
	Asand    string `json:"asand"`     // 배정공고일
	Asstd    string `json:"asstd"`     // 배정기준일
	Exstk    string `json:"exstk"`     // 신주인수권에 관한 사항(행사대상증권)
	Exprc    string `json:"exprc"`     // 신주인수권에 관한 사항(행사가격)
	Expd     string `json:"expd"`      // 신주인수권에 관한 사항(행사기간)
	RptRcpn  string `json:"rpt_rcpn"`  // 주요사항보고서(접수번호)
}

// RsSecurityTypeItem 은 증권신고서 증권의종류 그룹 항목(지분증권/증권예탁증권 공통).
type RsSecurityTypeItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Stksen   string `json:"stksen"`    // 증권의종류
	Stkcnt   string `json:"stkcnt"`    // 증권수량
	Fv       string `json:"fv"`        // 액면가액
	Slprc    string `json:"slprc"`     // 모집(매출)가액
	Slta     string `json:"slta"`      // 모집(매출)총액
	Slmthn   string `json:"slmthn"`    // 모집(매출)방법
}

// RsUnderwriterItem 은 증권신고서 인수인정보 그룹 항목(지분증권/증권예탁증권 공통).
type RsUnderwriterItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Actsen   string `json:"actsen"`    // 인수인구분
	Actnmn   string `json:"actnmn"`    // 인수인명
	Stksen   string `json:"stksen"`    // 증권의종류
	Udtcnt   string `json:"udtcnt"`    // 인수수량
	Udtamt   string `json:"udtamt"`    // 인수금액
	Udtprc   string `json:"udtprc"`    // 인수대가
	Udtmth   string `json:"udtmth"`    // 인수방법
}

// RsFundUsageItem 은 증권신고서 자금의사용목적 그룹 항목(지분증권/증권예탁증권 공통).
type RsFundUsageItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Se       string `json:"se"`        // 구분
	Amt      string `json:"amt"`       // 금액
}

// RsSellerItem 은 증권신고서 매출인에관한사항 그룹 항목(지분증권/증권예탁증권 공통).
type RsSellerItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Hdr       string `json:"hdr"`        // 보유자
	RlCmp     string `json:"rl_cmp"`     // 회사와의관계
	BfslHdstk string `json:"bfsl_hdstk"` // 매출전보유증권수
	Slstk     string `json:"slstk"`      // 매출증권수
	AtslHdstk string `json:"atsl_hdstk"` // 매출후보유증권수
}

// EquityRetailPutbackOptionItem 은 지분증권 일반청약자환매청구권 그룹 항목(지분증권 전용).
type EquityRetailPutbackOptionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Grtrs    string `json:"grtrs"`     // 부여사유
	Exavivr  string `json:"exavivr"`   // 행사가능 투자자
	Grtcnt   string `json:"grtcnt"`    // 부여수량
	Expd     string `json:"expd"`      // 행사기간
	Exprc    string `json:"exprc"`     // 행사가격
}
```

**채무증권 전용 item 타입**(모든 그룹에 tm(회차) 포함 → 별도):
```go
// DebtGeneralItem 은 채무증권 증권신고서 일반사항 그룹 항목.
type DebtGeneralItem struct {
	RceptNo     string `json:"rcept_no"`      // 접수번호
	CorpCls     string `json:"corp_cls"`      // 법인구분 (Y/K/N/E)
	CorpCode    string `json:"corp_code"`     // 고유번호
	CorpName    string `json:"corp_name"`     // 회사명
	Tm          string `json:"tm"`            // 회차
	Bdnmn       string `json:"bdnmn"`         // 채무증권 명칭
	Slmth       string `json:"slmth"`         // 모집(매출)방법
	Fta         string `json:"fta"`           // 권면(전자등록)총액
	Slta        string `json:"slta"`          // 모집(매출)총액
	Isprc       string `json:"isprc"`         // 발행가액
	Intr        string `json:"intr"`          // 이자율
	Isrr        string `json:"isrr"`          // 발행수익률
	Rpd         string `json:"rpd"`           // 상환기일
	PrintPymint string `json:"print_pymint"`  // 원리금지급대행기관
	MngtCmp     string `json:"mngt_cmp"`      // (사채)관리회사
	CdrtInt     string `json:"cdrt_int"`      // 신용등급(신용평가기관)
	Sbd         string `json:"sbd"`           // 청약기일
	Pymd        string `json:"pymd"`          // 납입기일
	Sband       string `json:"sband"`         // 청약공고일
	Asand       string `json:"asand"`         // 배정공고일
	Asstd       string `json:"asstd"`         // 배정기준일
	Dpcrn       string `json:"dpcrn"`         // 표시통화
	DpcrAmt     string `json:"dpcr_amt"`      // 표시통화기준발행규모
	Usarn       string `json:"usarn"`         // 사용지역
	Usntn       string `json:"usntn"`         // 사용국가
	WnexplAt    string `json:"wnexpl_at"`     // 원화 교환 예정 여부
	Udtintnm    string `json:"udtintnm"`      // 인수기관명
	GrtInt      string `json:"grt_int"`       // 보증을 받은 경우(보증기관)
	GrtAmt      string `json:"grt_amt"`       // 보증을 받은 경우(보증금액)
	IcmgMgknd   string `json:"icmg_mgknd"`    // 담보 제공의 경우(담보의 종류)
	IcmgMgamt   string `json:"icmg_mgamt"`    // 담보 제공의 경우(담보금액)
	EstkExstk   string `json:"estk_exstk"`    // 지분증권과 연계된 경우(행사대상증권)
	EstkExrt    string `json:"estk_exrt"`     // 지분증권과 연계된 경우(권리행사비율)
	EstkExprc   string `json:"estk_exprc"`    // 지분증권과 연계된 경우(권리행사가격)
	EstkExpd    string `json:"estk_expd"`     // 지분증권과 연계된 경우(권리행사기간)
	RptRcpn     string `json:"rpt_rcpn"`      // 주요사항보고서(접수번호)
	DrcbAt      string `json:"drcb_at"`       // 파생결합사채해당여부
	DrcbUast    string `json:"drcb_uast"`     // 파생결합사채(기초자산)
	DrcbOptknd  string `json:"drcb_optknd"`   // 파생결합사채(옵션종류)
	DrcbMtd     string `json:"drcb_mtd"`      // 파생결합사채(만기일)
}

// DebtUnderwriterItem 은 채무증권 증권신고서 인수인정보 그룹 항목.
type DebtUnderwriterItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Tm       string `json:"tm"`        // 회차
	Actsen   string `json:"actsen"`    // 인수인구분
	Actnmn   string `json:"actnmn"`    // 인수인명
	Stksen   string `json:"stksen"`    // 증권의종류
	Udtcnt   string `json:"udtcnt"`    // 인수수량
	Udtamt   string `json:"udtamt"`    // 인수금액
	Udtprc   string `json:"udtprc"`    // 인수대가
	Udtmth   string `json:"udtmth"`    // 인수방법
}

// DebtFundUsageItem 은 채무증권 증권신고서 자금의사용목적 그룹 항목.
type DebtFundUsageItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Tm       string `json:"tm"`        // 회차
	Se       string `json:"se"`        // 구분
	Amt      string `json:"amt"`       // 금액
}

// DebtSellerItem 은 채무증권 증권신고서 매출인에관한사항 그룹 항목.
type DebtSellerItem struct {
	RceptNo   string `json:"rcept_no"`   // 접수번호
	CorpCls   string `json:"corp_cls"`   // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`  // 고유번호
	CorpName  string `json:"corp_name"`  // 회사명
	Tm        string `json:"tm"`         // 회차
	Hdr       string `json:"hdr"`        // 보유자
	RlCmp     string `json:"rl_cmp"`     // 회사와의관계
	BfslHdstk string `json:"bfsl_hdstk"` // 매출전보유증권수
	Slstk     string `json:"slstk"`      // 매출증권수
	AtslHdstk string `json:"atsl_hdstk"` // 매출후보유증권수
}
```

**Wrapper 타입 + 메서드:**
```go
// EquitySecuritiesRegistration 은 지분증권 증권신고서(estkRs)의 그룹별 항목.
type EquitySecuritiesRegistration struct {
	General             []RsGeneralItem                 // 일반사항
	SecurityTypes       []RsSecurityTypeItem            // 증권의종류
	Underwriters        []RsUnderwriterItem             // 인수인정보
	FundUsage           []RsFundUsageItem               // 자금의사용목적
	Sellers             []RsSellerItem                  // 매출인에관한사항
	RetailPutbackOption []EquityRetailPutbackOptionItem // 일반청약자환매청구권
}

// EquitySecurities 는 지분증권 증권신고서(DS006)를 조회한다.
func (c *Client) EquitySecurities(ctx context.Context, p Params) (*EquitySecuritiesRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/estkRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &EquitySecuritiesRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "증권의종류":
			derr = json.Unmarshal(g.List, &out.SecurityTypes)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		case "일반청약자환매청구권":
			derr = json.Unmarshal(g.List, &out.RetailPutbackOption)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// DebtSecuritiesRegistration 은 채무증권 증권신고서(bdRs)의 그룹별 항목.
type DebtSecuritiesRegistration struct {
	General      []DebtGeneralItem     // 일반사항
	Underwriters []DebtUnderwriterItem // 인수인정보
	FundUsage    []DebtFundUsageItem   // 자금의사용목적
	Sellers      []DebtSellerItem      // 매출인에관한사항
}

// DebtSecurities 는 채무증권 증권신고서(DS006)를 조회한다.
func (c *Client) DebtSecurities(ctx context.Context, p Params) (*DebtSecuritiesRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/bdRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DebtSecuritiesRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}

// DepositaryReceiptsRegistration 은 증권예탁증권 증권신고서(stkdpRs)의 그룹별 항목.
type DepositaryReceiptsRegistration struct {
	General       []RsGeneralItem      // 일반사항
	SecurityTypes []RsSecurityTypeItem // 증권의종류
	Underwriters  []RsUnderwriterItem  // 인수인정보
	FundUsage     []RsFundUsageItem    // 자금의사용목적
	Sellers       []RsSellerItem       // 매출인에관한사항
}

// DepositaryReceipts 는 증권예탁증권 증권신고서(DS006)를 조회한다.
func (c *Client) DepositaryReceipts(ctx context.Context, p Params) (*DepositaryReceiptsRegistration, error) {
	groups, err := httpclient.GetGroups(ctx, c.http, "/api/stkdpRs.json", p.toMap())
	if err != nil {
		return nil, err
	}
	out := &DepositaryReceiptsRegistration{}
	for _, g := range groups {
		var derr error
		switch g.Title {
		case "일반사항":
			derr = json.Unmarshal(g.List, &out.General)
		case "증권의종류":
			derr = json.Unmarshal(g.List, &out.SecurityTypes)
		case "인수인정보":
			derr = json.Unmarshal(g.List, &out.Underwriters)
		case "자금의사용목적":
			derr = json.Unmarshal(g.List, &out.FundUsage)
		case "매출인에관한사항":
			derr = json.Unmarshal(g.List, &out.Sellers)
		}
		if derr != nil {
			return nil, derr
		}
	}
	return out, nil
}
```
(위 switch 디스패치는 `unmarshalGroup` 헬퍼로 대체 가능하나, 명시적 switch 가 컴파일타임 타입 안전성이 높아 채택. 구현 시 선택.)

## 에러 처리

`GetGroups` 가 status 엔벨로프 처리: 013→`opendart.ErrNoData`, 그 외→`*opendart.APIError`. 데이터 있으면 그룹 일부가 비어도 해당 슬라이스는 nil/empty.

## 테스트 전략

- `registration/client_test.go`: `newTestClient(t, fixture string) *Client` (httptest, 단일 응답 서빙) — DS006 는 엔드포인트당 단일 호출이라 path별 라우팅보다 fixture 직접 지정이 단순.
- `registration/securities_test.go`: 3개 메서드 각각 fixture 디코딩 → 각 그룹 슬라이스 len + 대표 필드 검증. 빈 그룹/누락 그룹 처리도 1케이스.
- fixture 는 실 API 캡처(남양유업 estkRs 확인됨; bdRs/stkdpRs 는 발행 종목 탐색, 불가 시 docs 스키마 일치 샘플 — 그룹 배열 형태 유지).
- `integration_test.go` 에 `EquitySecurities`(남양유업 00107598) 통합 케이스(`//go:build integration`); 추가로 채무/예탁 1개(ErrNoData skip).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지"에 DS006 줄 신설: "DS006 증권신고서 주요정보: 지분증권 · 채무증권 · 증권예탁증권".

## 비범위 (후속 plan)

- DS006 Sub-2 «신고» 3개(합병 mgRs/분할 dvRs/주식의포괄적교환·이전 extrRs — 일반사항/발행증권/당사회사 3그룹 구조).
- DS002 개인별 보수 Ver 2.0 2종(데이터 확보 시).
