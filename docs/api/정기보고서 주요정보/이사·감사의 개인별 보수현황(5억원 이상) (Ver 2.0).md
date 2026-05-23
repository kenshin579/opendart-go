# 이사·감사의 개인별 보수현황(5억원 이상) (Ver 2.0)

> apiGrpCd: DS002 · apiId: 2026001
>
> `GET https://opendart.fss.or.kr/api/hmvAuditIndvdlBySttusV2.json`

정기보고서(사업, 분기, 반기보고서) 내에 이사·감사의 개인별 보수현황(5억원 이상)을 제공합니다. ※ 2026년 5월 이후부터 제출된 보고서해당

## 기본 정보

| 메서드 | 요청URL | 인코딩 | 출력포멧 |
| --- | --- | --- | --- |
| GET | https://opendart.fss.or.kr/api/hmvAuditIndvdlBySttusV2.json | UTF-8 | JSON |
| GET | https://opendart.fss.or.kr/api/hmvAuditIndvdlBySttusV2.xml | UTF-8 | XML |

## 요청 인자

| 요청키 | 명칭 | 타입 | 필수여부 | 값설명 |
| --- | --- | --- | --- | --- |
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |
| corp_code | 고유번호 | STRING(8) | Y | 공시대상회사의 고유번호(8자리)※ 개발가이드 > 공시정보 > 고유번호 참고 |
| bsns_year | 사업연도 | STRING(4) | Y | 사업연도(4자리)※ 2026년 5월 이후 부터 정보제공 |
| reprt_code | 보고서 코드 | STRING(5) | Y | 1분기보고서 : 11013반기보고서 : 110123분기보고서 : 11014사업보고서 : 11011 |

## 응답 결과

| 응답키 | 명칭 | 출력설명 |
| --- | --- | --- |
| result |  |  |
| status | 에러 및 정보 코드 | (※메시지 설명 참조) |
| message | 에러 및 정보 메시지 | (※메시지 설명 참조) |
| rcept_no | 접수번호 | 접수번호(14자리) ※ 공시뷰어 연결에 이용예시- PC용 : https://dart.fss.or.kr/dsaf001/main.do?rcpNo=접수번호 |
| corp_cls | 법인구분 | 법인구분 : Y(유가), K(코스닥), N(코넥스), E(기타) |
| corp_code | 고유번호 | 공시대상회사의 고유번호(8자리) |
| corp_name | 회사명 | 공시대상회사명 |
| stlm_dt | 결산기준일 | YYYY-MM-DD |
| group |  |  |
| nm | 이름 | 홍길동 |
| fscl_year | 사업연도 | 당기,전기,전전기 |
| ofcps | 직위 | 이사, 대표이사 등 |
| mendng_totamt | 보수 총액 | 9,999,999,999 |
| list |  |  |
| stk_bsd_pd_mendng_totamt_knd | 보수총액 중 주식기준보상 지급액-종류 |  |
| stk_bsd_pd_mendng_totamt_qty | 보수총액 중 주식기준보상 지급액-수량 | 9,999,999,999 |
| stk_bsd_pd_mendng_totamt_amt | 보수총액 중 주식기준보상 지급액-금액 | 9,999,999,999 |
| stk_opt_exrcsbl_qty | 주식매수선택권 행사가능수량 | 9,999,999,999 |
| stk_opt_unexrcsbl_qty | 주식매수선택권 행사불가수량 | 9,999,999,999 |
| stk_opt_exrc_pr | 주식매수선택권 행사가격 | 9,999,999,999 |
| stk_opt_rmn_blce | 주식매수선택권 잔여금액 | 9,999,999,999 |
| othr_stk_bsd_cmpn_unpyd_qty | 그 외 주식기준 보상 미지급수량 | 9,999,999,999 |
| othr_stk_bsd_cmpn_mkt_vl | 그 외 주식기준 보상 시장가치 | 9,999,999,999 |
| rm | 비고 |  |
