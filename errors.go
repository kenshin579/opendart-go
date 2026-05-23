package opendart

import (
	"errors"

	"github.com/kenshin579/opendart/internal/httpclient"
)

// APIError 는 OpenDART status != "000" 응답. errors.As 로 Status/Message 접근.
type APIError = httpclient.APIError

// ErrNoData 는 status 013 (조회된 데이터 없음).
var ErrNoData = httpclient.ErrNoData

// ErrCorpCodeNotFound 는 종목코드/고유번호 매핑 실패.
var ErrCorpCodeNotFound = errors.New("opendart: corp_code not found")
