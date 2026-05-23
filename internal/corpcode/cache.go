// Package corpcode 는 OpenDART 고유번호(corpCode.xml) 의 디스크 캐시 + 인덱스다.
// 다운로드 비용이 큰 전체 회사 매핑 ZIP 을 TTL 단위로 재사용하고, 다운로드 실패 시
// 오래된 캐시로 fallback 한다 (KIS mastercache 모델).
package corpcode

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Entry 는 CORPCODE.xml 의 한 회사.
type Entry struct {
	CorpCode    string `xml:"corp_code"`     // 고유번호 (8자리)
	CorpName    string `xml:"corp_name"`     // 정식회사명칭
	CorpEngName string `xml:"corp_eng_name"` // 영문 정식명칭
	StockCode   string `xml:"stock_code"`    // 종목코드 (6자리, 비상장은 빈 값)
	ModifyDate  string `xml:"modify_date"`   // 최종변경일자 YYYYMMDD
}

// FetchFunc 는 corpCode.xml ZIP 바이트를 가져온다 (캐시 미스/갱신 시 호출).
type FetchFunc func(ctx context.Context) ([]byte, error)

const cacheFileName = "corpcode.zip"

// Cache 는 corp_code 매핑의 디스크 캐시 + 인메모리 인덱스.
type Cache struct {
	dir   string
	ttl   time.Duration
	fetch FetchFunc

	mu      sync.Mutex
	loaded  bool
	entries []Entry
	byStock map[string]Entry
	byCorp  map[string]Entry
}

// New 는 Cache 를 만든다. dir 는 캐시 디렉토리, ttl 은 재다운로드 주기.
func New(dir string, ttl time.Duration, fetch FetchFunc) *Cache {
	return &Cache{
		dir: dir, ttl: ttl, fetch: fetch,
		byStock: map[string]Entry{}, byCorp: map[string]Entry{},
	}
}

// ensure 는 인덱스가 준비되도록 한다 (lazy 1회). force=true 면 강제 재다운로드.
func (c *Cache) ensure(ctx context.Context, force bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loaded && !force {
		return nil
	}
	zipBytes, err := c.obtain(ctx, force)
	if err != nil {
		return err
	}
	entries, err := parseZip(zipBytes)
	if err != nil {
		return err
	}
	c.index(entries)
	c.loaded = true
	return nil
}

// obtain 은 신선한 캐시가 있으면 디스크에서, 아니면 fetch 후 저장한다.
// fetch 실패 시 오래된 캐시라도 있으면 fallback 한다(가용성 우선). 이 경우 프로세스
// 수명 동안 stale 데이터를 쓰며, 네트워크 회복 시 RefreshCorpCodes 로만 갱신된다.
func (c *Cache) obtain(ctx context.Context, force bool) ([]byte, error) {
	path := filepath.Join(c.dir, cacheFileName)
	if !force {
		if b, ok := readFresh(path, c.ttl); ok {
			return b, nil
		}
	}
	b, err := c.fetch(ctx)
	if err != nil {
		if old, rerr := os.ReadFile(path); rerr == nil {
			return old, nil // stale fallback
		}
		return nil, err
	}
	if mkErr := os.MkdirAll(c.dir, 0o755); mkErr == nil {
		_ = os.WriteFile(path, b, 0o644)
	}
	return b, nil
}

// readFresh 는 TTL 내 캐시 파일을 읽는다.
func readFresh(path string, ttl time.Duration) ([]byte, bool) {
	fi, err := os.Stat(path)
	if err != nil || time.Since(fi.ModTime()) > ttl {
		return nil, false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return b, true
}

// index 는 엔트리로 두 맵을 구축한다. 빈 stock_code(비상장)는 stockCode 맵에서 제외.
func (c *Cache) index(entries []Entry) {
	c.entries = entries
	c.byStock = make(map[string]Entry, len(entries))
	c.byCorp = make(map[string]Entry, len(entries))
	for _, e := range entries {
		c.byCorp[e.CorpCode] = e
		if e.StockCode != "" {
			c.byStock[e.StockCode] = e
		}
	}
}

// parseZip 은 ZIP 안 첫 .xml 을 찾아 엔트리로 파싱하고 stock_code 를 trim 한다.
func parseZip(zipBytes []byte) ([]Entry, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("corpcode: open zip: %w", err)
	}
	for _, f := range zr.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".xml") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		var doc struct {
			List []Entry `xml:"list"`
		}
		if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
			return nil, fmt.Errorf("corpcode: parse xml: %w", err)
		}
		for i := range doc.List {
			doc.List[i].StockCode = strings.TrimSpace(doc.List[i].StockCode)
		}
		return doc.List, nil
	}
	return nil, fmt.Errorf("corpcode: no xml in zip")
}

// Entries 는 전체 엔트리의 복사본을 반환한다 (호출자 변경이 내부 인덱스에 영향 없도록).
func (c *Cache) Entries(ctx context.Context) ([]Entry, error) {
	if err := c.ensure(ctx, false); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Entry, len(c.entries))
	copy(out, c.entries)
	return out, nil
}

// ByStockCode 는 종목코드로 엔트리를 찾는다.
func (c *Cache) ByStockCode(ctx context.Context, stockCode string) (Entry, bool, error) {
	if err := c.ensure(ctx, false); err != nil {
		return Entry{}, false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.byStock[strings.TrimSpace(stockCode)]
	return e, ok, nil
}

// ByCorpCode 는 고유번호로 엔트리를 찾는다.
func (c *Cache) ByCorpCode(ctx context.Context, corpCode string) (Entry, bool, error) {
	if err := c.ensure(ctx, false); err != nil {
		return Entry{}, false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.byCorp[strings.TrimSpace(corpCode)]
	return e, ok, nil
}

// Refresh 는 TTL 을 무시하고 강제 재다운로드한다.
func (c *Cache) Refresh(ctx context.Context) error {
	return c.ensure(ctx, true)
}
