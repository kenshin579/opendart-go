package corpcode

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<result>
	<list><corp_code>00126380</corp_code><corp_name>삼성전자(주)</corp_name><corp_eng_name>SAMSUNG</corp_eng_name><stock_code>005930</stock_code><modify_date>20240101</modify_date></list>
	<list><corp_code>00434003</corp_code><corp_name>다코</corp_name><corp_eng_name>Daco</corp_eng_name><stock_code> </stock_code><modify_date>20170630</modify_date></list>
</result>`

func makeZip(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	require.NoError(t, err)
	_, err = w.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func TestCache_IndexAndLookup(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	c := New(t.TempDir(), time.Hour, func(ctx context.Context) ([]byte, error) { return zipBytes, nil })

	e, ok, err := c.ByStockCode(context.Background(), "005930")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "00126380", e.CorpCode)

	_, ok, err = c.ByStockCode(context.Background(), "")
	require.NoError(t, err)
	assert.False(t, ok)

	e2, ok, err := c.ByCorpCode(context.Background(), "00434003")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "다코", e2.CorpName)
	assert.Equal(t, "", e2.StockCode)

	all, err := c.Entries(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestCache_DiskCacheAvoidsRefetch(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	dir := t.TempDir()
	calls := 0
	fetch := func(ctx context.Context) ([]byte, error) { calls++; return zipBytes, nil }

	c1 := New(dir, time.Hour, fetch)
	_, _, _ = c1.ByStockCode(context.Background(), "005930")
	require.Equal(t, 1, calls)

	c2 := New(dir, time.Hour, fetch)
	_, _, _ = c2.ByStockCode(context.Background(), "005930")
	assert.Equal(t, 1, calls)
}

func TestCache_StaleFallbackOnFetchError(t *testing.T) {
	zipBytes := makeZip(t, "CORPCODE.xml", sampleXML)
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, cacheFileName), zipBytes, 0o644))
	old := time.Now().Add(-48 * time.Hour)
	require.NoError(t, os.Chtimes(filepath.Join(dir, cacheFileName), old, old))

	c := New(dir, time.Hour, func(ctx context.Context) ([]byte, error) {
		return nil, errors.New("network down")
	})
	e, ok, err := c.ByCorpCode(context.Background(), "00126380")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "삼성전자(주)", e.CorpName)
}
