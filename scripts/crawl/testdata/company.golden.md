# 기업개황

> apiGrpCd: DS001 · apiId: 2019002
>
> `GET https://opendart.fss.or.kr/api/company.json`

DART에 등록되어있는 기업의 개황정보를 제공합니다.

## 기본 정보

| 메서드 | 요청URL | 인코딩 | 출력포멧 |
| --- | --- | --- | --- |
| GET | https://opendart.fss.or.kr/api/company.json | UTF-8 | JSON |

## 요청 인자

| 요청키 | 명칭 | 타입 | 필수여부 | 값설명 |
| --- | --- | --- | --- | --- |
| crtfc_key | API 인증키 | STRING(40) | Y | 발급받은 인증키(40자리) |

## 응답 결과

| 응답키 | 명칭 | 출력설명 |
| --- | --- | --- |
| corp_name | 정식명칭 | 정식회사 명칭 |
