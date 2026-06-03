package opendart

import (
	"context"

	"github.com/kenshin579/opendart-go/internal/corpcode"
)

// CorpCodeEntry 는 corp_code 매핑의 한 회사. (internal 타입 별칭 — 외부에서 opendart.CorpCodeEntry 로 사용)
type CorpCodeEntry = corpcode.Entry

// ResolveCorpCode 는 종목코드(6자리)를 corp_code(8자리)로 변환한다.
// 매핑이 없으면 ErrCorpCodeNotFound.
func (c *Client) ResolveCorpCode(ctx context.Context, stockCode string) (string, error) {
	e, ok, err := c.corp.ByStockCode(ctx, stockCode)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrCorpCodeNotFound
	}
	return e.CorpCode, nil
}

// LookupCorpCode 는 corp_code(8자리)로 회사 엔트리를 조회한다.
func (c *Client) LookupCorpCode(ctx context.Context, corpCode string) (CorpCodeEntry, error) {
	e, ok, err := c.corp.ByCorpCode(ctx, corpCode)
	if err != nil {
		return CorpCodeEntry{}, err
	}
	if !ok {
		return CorpCodeEntry{}, ErrCorpCodeNotFound
	}
	return e, nil
}

// CorpCodes 는 전체 회사 엔트리를 반환한다 (사용자 직접 필터용).
func (c *Client) CorpCodes(ctx context.Context) ([]CorpCodeEntry, error) {
	return c.corp.Entries(ctx)
}

// RefreshCorpCodes 는 TTL 을 무시하고 corp_code 매핑을 강제 재다운로드한다.
func (c *Client) RefreshCorpCodes(ctx context.Context) error {
	return c.corp.Refresh(ctx)
}
