# 재무제표 원본파일(XBRL)

> apiGrpCd: DS003 · apiId: 2019019
>
> `GET https://opendart.fss.or.kr/api/fnlttXbrl.xml`

상장법인(유가증권, 코스닥) 및 주요 비상장법인(사업보고서 제출대상 & IFRS 적용)이 제출한 정기보고서 내에 XBRL재무제표의 원본파일(XBRL)을 제공합니다.

## 기본 정보

| 메서드 | 요청URL | 인코딩 | 출력포멧 |
| --- | --- | --- | --- |
| GET | https://opendart.fss.or.kr/api/fnlttXbrl.xml | UTF-8 | Zip FILE (binary) |

## 요청 인자

| 요청키 | 명칭 | 타입 | 필수여부 | 값설명 |
| --- | --- | --- | --- | --- |
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |
| rcept_no | 접수번호 | STRING(8) | Y | 접수번호 ※ 조회방법 : 공시검색API 호출 > 응답요청 값 rcept_no 추출 |
| reprt_code | 보고서 코드 | STRING(5) | Y | 1분기보고서 : 11013반기보고서 : 110123분기보고서 : 11014사업보고서 : 11011 |

## 응답 결과

| 응답키 | 명칭 | 출력설명 |
| --- | --- | --- |
| result |  |  |
| status | 에러 및 정보 코드 | (※메시지 설명 참조) |
| message | 에러 및 정보 메시지 | (※메시지 설명 참조) |
