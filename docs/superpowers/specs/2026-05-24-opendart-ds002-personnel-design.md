# OpenDART DS002 정기보고서 주요정보 — 임원·직원·보수 그룹 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS002 임원·직원·보수 9개 API** (`report` 패키지 확장). 개인별 보수 Ver 2.0 2종은 보류.

## 배경 & 목표

DS002 공통 추상화(`report` 패키지: `getList`/`ReportParams`)와 지분·주식·배당 7개(PR #3),
증권 발행·미상환 6개(PR #4), 감사·자금·출자 6개(PR #5)는 main 에 머지됨. 이 spec 은 DS002 마지막
그룹인 **임원·직원·보수**다. 표준 요청(`corp_code`+`bsns_year`+`reprt_code`)과 list 응답이라 기존
`getList[T]` 패턴을 재사용한다.

**Ver 2.0 보류 결정:** 개인별 보수 2종에는 enriched Ver 2.0 변종(`indvdlByPayV2`,
`hmvAuditIndvdlBySttusV2`)이 있으나, **현재 어떤 회사·연도로도 데이터가 조회되지 않는다**(삼성·SK하이닉스·
현대차·NAVER·카카오 2023~2024 전부 status 013). 원본 보수 API docs 의 "2026년 4월까지 제출된
보고서 해당" 표기로 보아 V2 는 2026년 4월 이후 제출 보고서부터 적용되는 미래 포맷이다. 실 응답을
캡처할 수 없고 응답 구조(docs 의 `group` 행 → 중첩 `group[]` 가능성)도 확정 불가하므로, 이 프로젝트의
"실 응답으로 검증" 원칙에 따라 **V2 2종은 데이터·구조 확인 후 별도 작업으로 보류**한다. 이 그룹은 원본
9개만 구현한다 (완료 시 DS002 28/30).

## API 표면 (docs 기반 사실)

- 요청: 전 API 공통 `crtfc_key`(자동 주입) + `corp_code` + `bsns_year` + `reprt_code` (추가 파라미터 없음).
- 응답: 공통 envelope(`status`/`message`) + `list[]`. 숫자는 콤마 문자열, 빈 값 "-".
- 이사·감사 전체 보수 3종은 응답 list 항목에 `fscl_year`(사업연도) + 주식기준보상 관련 필드 다수 포함.

## 아키텍처

`report` 패키지에 새 파일 `personnel.go` 추가 (기존 `equity.go`/`securities.go`/`audit.go` 와 동일하게
그룹당 1파일). `client.Report` 는 이미 root 에 와이어링됨 → **root 변경 불필요**. 메서드는 기존
`getList[T]`/`ReportParams` 재사용.

```
report/
  personnel.go        # 9개 메서드 + item struct (신규)
  personnel_test.go   # 9개 fixture 테스트 (신규, newTestClient 재사용)
  testdata/           # 9개 실 응답 JSON fixture 추가
README.md             # (수정) DS002 커버리지 + V2 보류 명시
integration_test.go   # (수정) Executives 통합 케이스
```

## 9개 엔드포인트 (report/personnel.go)

각 메서드: `func (c *Client) X(ctx, p ReportParams) ([]XItem, error) { return getList[XItem](ctx, c.http, "<path>", p) }`.

| 메서드 | 한글 | 엔드포인트 |
|--------|------|-----------|
| `Executives` | 임원 현황 | `/api/exctvSttus.json` |
| `Employees` | 직원 현황 | `/api/empSttus.json` |
| `UnregisteredExecutiveCompensation` | 미등기임원 보수현황 | `/api/unrstExctvMendngSttus.json` |
| `OutsideDirectorChanges` | 사외이사 및 그 변동현황 | `/api/outcmpnyDrctrNdChangeSttus.json` |
| `DirectorAuditorApprovedCompensation` | 이사·감사 전체 보수(주총 승인금액) | `/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json` |
| `DirectorAuditorTotalCompensation` | 이사·감사 전체 보수(지급금액-전체) | `/api/hmvAuditAllSttus.json` |
| `DirectorAuditorCompensationByType` | 이사·감사 전체 보수(지급금액-유형별) | `/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json` |
| `IndividualDirectorAuditorCompensation` | 이사·감사 개인별 보수(5억 이상) | `/api/hmvAuditIndvdlBySttus.json` |
| `IndividualTop5Compensation` | 개인별 보수지급(5억이상 상위5인) | `/api/indvdlByPay.json` |

```go
// ExecutiveItem 은 임원 현황 (exctvSttus) 한 건.
type ExecutiveItem struct {
	RceptNo            string `json:"rcept_no"`            // 접수번호
	CorpCls            string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode           string `json:"corp_code"`           // 고유번호
	CorpName           string `json:"corp_name"`           // 법인명
	Nm                 string `json:"nm"`                  // 성명
	Sexdstn            string `json:"sexdstn"`             // 성별
	BirthYm            string `json:"birth_ym"`            // 출생 년월
	Ofcps              string `json:"ofcps"`               // 직위
	RgistExctvAt       string `json:"rgist_exctv_at"`      // 등기 임원 여부
	FteAt              string `json:"fte_at"`              // 상근 여부
	ChrgJob            string `json:"chrg_job"`            // 담당 업무
	MainCareer         string `json:"main_career"`         // 주요 경력
	MxmmShrholdrRelate string `json:"mxmm_shrholdr_relate"` // 최대 주주 관계
	HffcPd             string `json:"hffc_pd"`             // 재직 기간
	TenureEndOn        string `json:"tenure_end_on"`       // 임기 만료 일
	StlmDt             string `json:"stlm_dt"`             // 결산기준일
}

// EmployeeItem 은 직원 현황 (empSttus) 한 건.
type EmployeeItem struct {
	RceptNo              string `json:"rcept_no"`                // 접수번호
	CorpCls              string `json:"corp_cls"`                // 법인구분 (Y/K/N/E)
	CorpCode             string `json:"corp_code"`               // 고유번호
	CorpName             string `json:"corp_name"`               // 법인명
	FoBbm                string `json:"fo_bbm"`                  // 사업부문
	Sexdstn              string `json:"sexdstn"`                 // 성별
	ReformBfeEmpCoRgllbr string `json:"reform_bfe_emp_co_rgllbr"` // 개정 전 직원 수 정규직
	ReformBfeEmpCoCnttk  string `json:"reform_bfe_emp_co_cnttk"`  // 개정 전 직원 수 계약직
	ReformBfeEmpCoEtc    string `json:"reform_bfe_emp_co_etc"`    // 개정 전 직원 수 기타
	RgllbrCo             string `json:"rgllbr_co"`               // 정규직 수
	RgllbrAbacptLabrrCo  string `json:"rgllbr_abacpt_labrr_co"`   // 정규직 단시간 근로자 수
	CnttkCo              string `json:"cnttk_co"`                // 계약직 수
	CnttkAbacptLabrrCo   string `json:"cnttk_abacpt_labrr_co"`    // 계약직 단시간 근로자 수
	Sm                   string `json:"sm"`                      // 합계
	AvrgCnwkSdytrn       string `json:"avrg_cnwk_sdytrn"`         // 평균 근속 연수
	FyerSalaryTotamt     string `json:"fyer_salary_totamt"`       // 연간 급여 총액
	JanSalaryAm          string `json:"jan_salary_am"`            // 1인평균 급여 액
	Rm                   string `json:"rm"`                      // 비고
	StlmDt               string `json:"stlm_dt"`                 // 결산기준일
}

// UnregisteredExecutiveCompensationItem 은 미등기임원 보수현황 (unrstExctvMendngSttus) 한 건.
type UnregisteredExecutiveCompensationItem struct {
	RceptNo          string `json:"rcept_no"`           // 접수번호
	CorpCls          string `json:"corp_cls"`           // 법인구분 (Y/K/N/E)
	CorpCode         string `json:"corp_code"`          // 고유번호
	CorpName         string `json:"corp_name"`          // 회사명
	Se               string `json:"se"`                 // 구분
	Nmpr             string `json:"nmpr"`               // 인원수
	FyerSalaryTotamt string `json:"fyer_salary_totamt"`  // 연간급여 총액
	JanSalaryAm      string `json:"jan_salary_am"`       // 1인평균 급여액
	Rm               string `json:"rm"`                 // 비고
	StlmDt           string `json:"stlm_dt"`            // 결산기준일
}

// OutsideDirectorChangeItem 은 사외이사 및 그 변동현황 (outcmpnyDrctrNdChangeSttus) 한 건.
type OutsideDirectorChangeItem struct {
	RceptNo      string `json:"rcept_no"`      // 접수번호
	CorpCls      string `json:"corp_cls"`      // 법인구분 (Y/K/N/E)
	CorpCode     string `json:"corp_code"`     // 고유번호
	CorpName     string `json:"corp_name"`     // 회사명
	DrctrCo      string `json:"drctr_co"`      // 이사의 수
	OtcmpDrctrCo string `json:"otcmp_drctr_co"` // 사외이사 수
	Apnt         string `json:"apnt"`          // 사외이사 변동현황(선임)
	Rlsofc       string `json:"rlsofc"`        // 사외이사 변동현황(해임)
	MdstrmResig  string `json:"mdstrm_resig"`  // 사외이사 변동현황(중도퇴임)
	StlmDt       string `json:"stlm_dt"`       // 결산기준일
}

// DirectorAuditorApprovedCompensationItem 은 이사·감사 전체의 보수현황(주주총회 승인금액)
// (drctrAdtAllMendngSttusGmtsckConfmAmount) 한 건.
type DirectorAuditorApprovedCompensationItem struct {
	RceptNo           string `json:"rcept_no"`            // 접수번호
	CorpCls           string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode          string `json:"corp_code"`           // 고유번호
	CorpName          string `json:"corp_name"`           // 회사명
	Se                string `json:"se"`                  // 구분
	Nmpr              string `json:"nmpr"`                // 인원수
	GmtsckConfmAmount string `json:"gmtsck_confm_amount"`  // 주주총회 승인금액
	Rm                string `json:"rm"`                  // 비고
	StlmDt            string `json:"stlm_dt"`             // 결산기준일
	FsclYear          string `json:"fscl_year"`           // 사업연도
}

// DirectorAuditorTotalCompensationItem 은 이사·감사 전체의 보수현황(보수지급금액 - 이사·감사 전체)
// (hmvAuditAllSttus) 한 건.
type DirectorAuditorTotalCompensationItem struct {
	RceptNo               string `json:"rcept_no"`                  // 접수번호
	CorpCls               string `json:"corp_cls"`                  // 법인구분 (Y/K/N/E)
	CorpCode              string `json:"corp_code"`                 // 고유번호
	CorpName              string `json:"corp_name"`                 // 법인명
	Nmpr                  string `json:"nmpr"`                      // 인원수
	MendngTotamt          string `json:"mendng_totamt"`             // 보수 총액
	JanAvrgMendngAm       string `json:"jan_avrg_mendng_am"`        // 1인 평균 보수 액
	Rm                    string `json:"rm"`                        // 비고
	StlmDt                string `json:"stlm_dt"`                   // 결산기준일
	FsclYear              string `json:"fscl_year"`                 // 사업연도
	StkBsdPdMendngTotamt  string `json:"stk_bsd_pd_mendng_totamt"`  // 보수총액 중 주식기준보상 지급액
	StkOptExrcsblQty      string `json:"stk_opt_exrcsbl_qty"`       // 주식매수선택권 행사가능수량
	StkOptUnexrcsblQty    string `json:"stk_opt_unexrcsbl_qty"`     // 주식매수선택권 행사불가수량
	StkOptRmnBlce         string `json:"stk_opt_rmn_blce"`          // 주식매수선택권 잔여금액
	OthrStkBsdCmpnUnpydQty string `json:"othr_stk_bsd_cmpn_unpyd_qty"` // 그 외 주식기준 보상 미지급수량
	OthrStkBsdCmpnMktVl   string `json:"othr_stk_bsd_cmpn_mkt_vl"`  // 그 외 주식기준 보상 시장가치
}

// DirectorAuditorCompensationByTypeItem 은 이사·감사 전체의 보수현황(보수지급금액 - 유형별)
// (drctrAdtAllMendngSttusMendngPymntamtTyCl) 한 건.
type DirectorAuditorCompensationByTypeItem struct {
	RceptNo               string `json:"rcept_no"`                  // 접수번호
	CorpCls               string `json:"corp_cls"`                  // 법인구분 (Y/K/N/E)
	CorpCode              string `json:"corp_code"`                 // 고유번호
	CorpName              string `json:"corp_name"`                 // 회사명
	Se                    string `json:"se"`                        // 구분 (등기이사/사외이사/감사위원회 위원 등)
	Nmpr                  string `json:"nmpr"`                      // 인원수
	PymntTotamt           string `json:"pymnt_totamt"`              // 보수총액
	Psn1AvrgPymntamt      string `json:"psn1_avrg_pymntamt"`        // 1인당 평균보수액
	Rm                    string `json:"rm"`                        // 비고
	StlmDt                string `json:"stlm_dt"`                   // 결산기준일
	FsclYear              string `json:"fscl_year"`                 // 사업연도
	StkBsdPdMendngTotamt  string `json:"stk_bsd_pd_mendng_totamt"`  // 보수총액 중 주식기준보상 지급액
	StkOptExrcsblQty      string `json:"stk_opt_exrcsbl_qty"`       // 주식매수선택권 행사가능수량
	StkOptUnexrcsblQty    string `json:"stk_opt_unexrcsbl_qty"`     // 주식매수선택권 행사불가수량
	StkOptRmnBlce         string `json:"stk_opt_rmn_blce"`          // 주식매수선택권 잔여금액
	OthrStkBsdCmpnUnpydQty string `json:"othr_stk_bsd_cmpn_unpyd_qty"` // 그 외 주식기준 보상 미지급수량
	OthrStkBsdCmpnMktVl   string `json:"othr_stk_bsd_cmpn_mkt_vl"`  // 그 외 주식기준 보상 시장가치
}

// IndividualDirectorAuditorCompensationItem 은 이사·감사의 개인별 보수현황(5억원 이상)
// (hmvAuditIndvdlBySttus) 한 건.
type IndividualDirectorAuditorCompensationItem struct {
	RceptNo                   string `json:"rcept_no"`                      // 접수번호
	CorpCls                   string `json:"corp_cls"`                      // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                     // 고유번호
	CorpName                  string `json:"corp_name"`                     // 법인명
	Nm                        string `json:"nm"`                            // 이름
	Ofcps                     string `json:"ofcps"`                         // 직위
	MendngTotamt              string `json:"mendng_totamt"`                 // 보수 총액
	MendngTotamtCtInclsMendng string `json:"mendng_totamt_ct_incls_mendng"` // 보수 총액 비 포함 보수
	StlmDt                    string `json:"stlm_dt"`                       // 결산기준일
}

// IndividualTop5CompensationItem 은 개인별 보수지급 금액(5억이상 상위5인) (indvdlByPay) 한 건.
type IndividualTop5CompensationItem struct {
	RceptNo                   string `json:"rcept_no"`                      // 접수번호
	CorpCls                   string `json:"corp_cls"`                      // 법인구분 (Y/K/N/E)
	CorpCode                  string `json:"corp_code"`                     // 고유번호
	CorpName                  string `json:"corp_name"`                     // 법인명
	Nm                        string `json:"nm"`                            // 이름
	Ofcps                     string `json:"ofcps"`                         // 직위
	MendngTotamt              string `json:"mendng_totamt"`                 // 보수 총액
	MendngTotamtCtInclsMendng string `json:"mendng_totamt_ct_incls_mendng"` // 보수 총액 비 포함 보수
	StlmDt                    string `json:"stlm_dt"`                       // 결산기준일
}
```

각 메서드는 위 패턴(`getList[XItem](ctx, c.http, "<path>", p)`)으로 작성한다.

> 참고: `IndividualDirectorAuditorCompensationItem` 과 `IndividualTop5CompensationItem` 은 현재 응답
> 필드가 동일하지만, 의미(전체 5억 이상 vs 상위 5인)와 엔드포인트가 다르므로 별도 타입으로 둔다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `report/personnel_test.go`: 기존 `report/client_test.go` 의 `newTestClient` 재사용.
- 9개 메서드 각각: 실 응답 JSON fixture 디코딩 → 대표 필드 매핑 검증.
- fixture 는 실 API 로 캡처해 임베드(계획 작성 단계, 삼성전자 2023 사업보고서 — 9개 모두 데이터 존재 확인).
- `integration_test.go` 에 `Executives` 통합 케이스 추가(`//go:build integration`).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, 숫자 coercion 없음(콤마 string 유지), UTF-8.
- README "커버리지" DS002 줄에 "임원/직원 현황·보수현황" 추가, "(예정)" 줄을 "개인별 보수 Ver2.0 2종 · DS003~DS006" 으로 갱신.

## 비범위 (후속 plan)

- DS002 개인별 보수 Ver 2.0 2종(`indvdlByPayV2`, `hmvAuditIndvdlBySttusV2`) — 실데이터 확보 +
  응답 구조(중첩 `group[]` 여부) 확인 후 구현.
- DS003~DS006 카테고리.
- 신규 예제(기존 `examples/report` 로 충분).
