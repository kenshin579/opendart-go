# XBRL택사노미재무제표양식

> apiGrpCd: DS003 · apiId: 2020001
>
> `GET https://opendart.fss.or.kr/api/xbrlTaxonomy.json`

금융감독원 회계포탈에서 제공하는 IFRS 기반 XBRL 재무제표 공시용 표준계정과목체계(계정과목) 을 제공합니다.

## 기본 정보

| 메서드 | 요청URL | 인코딩 | 출력포멧 |
| --- | --- | --- | --- |
| GET | https://opendart.fss.or.kr/api/xbrlTaxonomy.json | UTF-8 | JSON |
| GET | https://opendart.fss.or.kr/api/xbrlTaxonomy.xml | UTF-8 | XML |

## 요청 인자

| 요청키 | 명칭 | 타입 | 필수여부 | 값설명 |
| --- | --- | --- | --- | --- |
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |
| sj_div | 재무제표구분 | STRING(5) | Y | (※재무제표구분 참조) |

## 응답 결과

| 응답키 | 명칭 | 출력설명 |
| --- | --- | --- |
| result |  |  |
| status | 에러 및 정보 코드 | (※메시지 설명 참조) |
| message | 에러 및 정보 메시지 | (※메시지 설명 참조) |
| list |  |  |
| sj_div | 재무제표구분 | 재무제표구분 |
| account_id | 계정ID | 계정 고유명칭 |
| account_nm | 계정명 | 계정명 |
| bsns_de | 기준일 | 적용 기준일 |
| label_kor | 한글 출력명 | 한글 출력명 |
| label_eng | 영문 출력명 | 영문 출력명 |
| data_tp | 데이터 유형 | ※ 데이타 유형설명 - text block : 제목 - Text : Text - yyyy-mm-dd : Date - X : Monetary Value - (X): Monetary Value(Negative) - X.XX : Decimalized Value - Shares : Number of shares (주식 수) - For each : 공시된 항목이 전후로 반복적으로 공시될 경우 사용 - 공란 : 입력 필요 없음 |
| ifrs_ref | IFRS Reference | IFRS Reference ※ 출력예시K-IFRS 1001 문단 54 (9),K-IFRS 1007 문단 45 |
