package main

// categories 는 OpenDART 개발가이드 6개 그룹 (apiGrpCd → 한글명).
var categories = []struct {
	Code string
	Name string
}{
	{"DS001", "공시정보"},
	{"DS002", "정기보고서 주요정보"},
	{"DS003", "정기보고서 재무정보"},
	{"DS004", "지분공시 종합정보"},
	{"DS005", "주요사항보고서 주요정보"},
	{"DS006", "증권신고서 주요정보"},
}

// APIRef 는 카테고리 목록 페이지에서 추출한 개별 API 식별 정보.
type APIRef struct {
	GrpCd    string // DS001
	Category string // 공시정보
	APIID    string // 2019002
	Name     string // 기업개황
	Desc     string // 상세기능 설명
}

// Table 은 상세 페이지의 한 표 (헤더 행 + 데이터 행들).
type Table struct {
	Headers []string
	Rows    [][]string
}

// APISpec 은 detail 페이지에서 추출한 명세. 메시지 설명(공통 상태코드)은 제외.
type APISpec struct {
	BasicInfo Table // 기본 정보 (메서드/요청URL/인코딩/출력포멧)
	Request   Table // 요청 인자
	Response  Table // 응답 결과
}
