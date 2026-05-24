# OpenDART DS005 주요사항보고서 주요정보 — 부실·법적 이벤트 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 부실·법적 이벤트 7개 API** — 신규 `material` 패키지 + 공통 `MaterialParams`

## 배경 & 목표

DS001~DS004 는 main 에 머지됨. 이 spec 은 DS005 주요사항보고서 주요정보(총 36개)의 첫 그룹 —
**부실·법적 이벤트 7개**다. DS005 36개는 모두 동일 요청(`corp_code`+`bgn_de`+`end_de` 날짜 범위)과
list 응답을 가지므로, 신규 `material` 패키지에 공통 `MaterialParams` + 기존 `httpclient.GetList[T]`
재사용으로 구현한다(DS002 report 패턴과 동형). 나머지 그룹(증자·감자/사채발행/자기주식/양수도/
합병·분할/해외상장 ~29개)은 동일 패턴 후속 plan.

## API 표면 (docs 기반 사실)

- 36개 전부 동일 요청 파라미터: `crtfc_key`(자동 주입) + `corp_code` + `bgn_de` + `end_de`. JSON `list[]`.
- 부실·법적 이벤트 7개: 부도발생(dfOcr) · 영업정지(bsnSp) · 회생절차 개시신청(ctrcvsBgrq) ·
  해산사유 발생(dsRsOcr) · 채권은행 등의 관리절차 개시(bnkMngtPcbg) · 채권은행 등의 관리절차 중단
  (bnkMngtPcsp) · 소송 등의 제기(lwstLg).
- 공통 머리: `rcept_no`/`corp_cls`/`corp_code`/`corp_name`. 값은 문자열.

## 아키텍처

신규 sub-package `material/`. `client.Material` 로 노출(root `client.go` 와이어링 1줄). 기존
`httpclient.GetList[T]` 재사용, 새 추상화 없음.

```
material/
  client.go         # Client + New + MaterialParams(+toMap)
  distress.go       # 7개 메서드 + item struct
  client_test.go    # newTestClient + MaterialParams.toMap 테스트
  distress_test.go  # 7개 fixture 테스트
  testdata/         # 7개 실 응답 fixture
client.go           # (수정) Material 필드 + 와이어링
README.md           # (수정) DS005 커버리지
integration_test.go # (수정) DefaultOccurrences 통합 케이스
```

### 공통 (material/client.go)

```go
// Package material 는 OpenDART DS005 주요사항보고서 주요정보 API sub-client 다.
package material

import "github.com/kenshin579/opendart/internal/httpclient"

type Client struct {
	http *httpclient.Client
}

func New(http *httpclient.Client) *Client { return &Client{http: http} }

// MaterialParams 는 DS005 공통 요청 인자 (날짜 범위). 빈 값은 쿼리에서 생략(OpenDART 기본값 적용).
type MaterialParams struct {
	CorpCode string // 고유번호 (8자리)
	BgnDe    string // 시작일 YYYYMMDD
	EndDe    string // 종료일 YYYYMMDD
}

func (p MaterialParams) toMap() map[string]string {
	m := map[string]string{"corp_code": p.CorpCode}
	if p.BgnDe != "" {
		m["bgn_de"] = p.BgnDe
	}
	if p.EndDe != "" {
		m["end_de"] = p.EndDe
	}
	return m
}
```

## 7개 메서드 (material/distress.go)

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `DefaultOccurrences` | 부도발생 | `/api/dfOcr.json` |
| `BusinessSuspensions` | 영업정지 | `/api/bsnSp.json` |
| `RehabilitationApplications` | 회생절차 개시신청 | `/api/ctrcvsBgrq.json` |
| `DissolutionCauses` | 해산사유 발생 | `/api/dsRsOcr.json` |
| `CreditorBankManagementStart` | 채권은행 등의 관리절차 개시 | `/api/bnkMngtPcbg.json` |
| `CreditorBankManagementStop` | 채권은행 등의 관리절차 중단 | `/api/bnkMngtPcsp.json` |
| `Lawsuits` | 소송 등의 제기 | `/api/lwstLg.json` |

```go
// DefaultItem 은 부도발생 (dfOcr) 한 건.
type DefaultItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DfCn     string `json:"df_cn"`     // 부도내용
	DfAmt    string `json:"df_amt"`    // 부도금액
	DfBnk    string `json:"df_bnk"`    // 부도발생은행
	Dfd      string `json:"dfd"`       // 최종부도(당좌거래정지)일자
	DfRs     string `json:"df_rs"`     // 부도사유 및 경위
}

// BusinessSuspensionItem 은 영업정지 (bsnSp) 한 건.
type BusinessSuspensionItem struct {
	RceptNo   string `json:"rcept_no"`    // 접수번호
	CorpCls   string `json:"corp_cls"`    // 법인구분 (Y/K/N/E)
	CorpCode  string `json:"corp_code"`   // 고유번호
	CorpName  string `json:"corp_name"`   // 회사명
	BsnspRm   string `json:"bsnsp_rm"`    // 영업정지 분야
	BsnspAmt  string `json:"bsnsp_amt"`   // 영업정지 내역(영업정지금액)
	Rsl       string `json:"rsl"`         // 영업정지 내역(최근매출총액)
	SlVs      string `json:"sl_vs"`       // 영업정지 내역(매출액 대비)
	LsAtn     string `json:"ls_atn"`      // 영업정지 내역(대규모법인여부)
	KrxSttAtn string `json:"krx_stt_atn"` // 영업정지 내역(거래소 의무공시 해당 여부)
	BsnspCn   string `json:"bsnsp_cn"`    // 영업정지 내용
	BsnspRs   string `json:"bsnsp_rs"`    // 영업정지사유
	FtCtp     string `json:"ft_ctp"`      // 향후대책
	BsnspAf   string `json:"bsnsp_af"`    // 영업정지영향
	Bsnspd    string `json:"bsnspd"`      // 영업정지일자
	Bddd      string `json:"bddd"`        // 이사회결의일(결정일)
	OdAAtT    string `json:"od_a_at_t"`   // 사외이사 참석여부(참석)
	OdAAtB    string `json:"od_a_at_b"`   // 사외이사 참석여부(불참)
	AdtAAtn   string `json:"adt_a_atn"`   // 감사(감사위원) 참석여부
}

// RehabilitationItem 은 회생절차 개시신청 (ctrcvsBgrq) 한 건.
type RehabilitationItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Apcnt    string `json:"apcnt"`     // 신청인 (회사와의 관계)
	Cpct     string `json:"cpct"`      // 관할법원
	RqRs     string `json:"rq_rs"`     // 신청사유
	Rqd      string `json:"rqd"`       // 신청일자
	FtCtpSc  string `json:"ft_ctp_sc"` // 향후대책 및 일정
}

// DissolutionItem 은 해산사유 발생 (dsRsOcr) 한 건.
type DissolutionItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	DsRs     string `json:"ds_rs"`     // 해산사유
	DsRsd    string `json:"ds_rsd"`    // 해산사유발생일(결정일)
	OdAAtT   string `json:"od_a_at_t"` // 사외이사 참석여부(참석)
	OdAAtB   string `json:"od_a_at_b"` // 사외이사 참석여부(불참)
	AdtAAtn  string `json:"adt_a_atn"` // 감사(감사위원) 참석여부
}

// CreditorBankMgmtStartItem 은 채권은행 등의 관리절차 개시 (bnkMngtPcbg) 한 건.
type CreditorBankMgmtStartItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcbgDd string `json:"mngt_pcbg_dd"` // 관리절차개시 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	MngtPd     string `json:"mngt_pd"`      // 관리기간
	MngtRs     string `json:"mngt_rs"`      // 관리사유
	Cfd        string `json:"cfd"`          // 확인일자
}

// CreditorBankMgmtStopItem 은 채권은행 등의 관리절차 중단 (bnkMngtPcsp) 한 건.
type CreditorBankMgmtStopItem struct {
	RceptNo    string `json:"rcept_no"`     // 접수번호
	CorpCls    string `json:"corp_cls"`     // 법인구분 (Y/K/N/E)
	CorpCode   string `json:"corp_code"`    // 고유번호
	CorpName   string `json:"corp_name"`    // 회사명
	MngtPcspDd string `json:"mngt_pcsp_dd"` // 관리절차중단 결정일자
	MngtInt    string `json:"mngt_int"`     // 관리기관
	SpRs       string `json:"sp_rs"`        // 중단사유
	FtCtp      string `json:"ft_ctp"`       // 향후대책
	Cfd        string `json:"cfd"`          // 확인일자
}

// LawsuitItem 은 소송 등의 제기 (lwstLg) 한 건.
type LawsuitItem struct {
	RceptNo  string `json:"rcept_no"`  // 접수번호
	CorpCls  string `json:"corp_cls"`  // 법인구분 (Y/K/N/E)
	CorpCode string `json:"corp_code"` // 고유번호
	CorpName string `json:"corp_name"` // 회사명
	Icnm     string `json:"icnm"`      // 사건의 명칭
	AcAp     string `json:"ac_ap"`     // 원고·신청인
	RqCn     string `json:"rq_cn"`     // 청구내용
	Cpct     string `json:"cpct"`      // 관할법원
	FtCtp    string `json:"ft_ctp"`    // 향후대책
	Lgd      string `json:"lgd"`       // 제기일자
	Cfd      string `json:"cfd"`       // 확인일자
}
```

각 메서드는 위 패턴(`httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap())`)으로 작성한다.

## root 와이어링 (client.go)

```go
type Client struct {
	http *httpclient.Client
	corp *corpcode.Cache

	Disclosure *disclosure.Client // DS001
	Report     *report.Client     // DS002 + DS003
	Ownership  *ownership.Client  // DS004
	Material   *material.Client   // DS005 주요사항보고서 주요정보
}
// NewClient 내부: c.Material = material.New(hc)
```

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/client_test.go`: 기존 패턴 `newTestClient` + `MaterialParams.toMap` 테스트(corp_code 항상, bgn_de/end_de omit 확인).
- `material/distress_test.go`: 7개 메서드 각각 실 응답 fixture 디코딩 → 대표 필드 검증.
- fixture 는 실 API 로 캡처해 임베드(데이터 있는 종목/기간; 데이터 없는 이벤트는 발생 회사로 캡처).
- `integration_test.go` 에 `DefaultOccurrences` 통합 케이스(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" 에 DS005 주요사항보고서(부실·법적 이벤트) 추가.

## 비범위 (후속 plan)

- DS005 나머지 그룹: 증자·감자 / 사채 발행 / 자기주식 / 영업·자산 양수도 / 합병·분할 / 해외상장 (~29개).
- DS006 증권신고서 주요정보.
- DS002 개인별 보수 Ver 2.0 2종.
