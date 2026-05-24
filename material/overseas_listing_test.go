package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverseasListingDecision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLstDecsn.json": "ovLstDecsn.json"})
	items, err := c.OverseasListingDecision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20241031000536", got.RceptNo)
	assert.Equal(t, "38,685,850", got.LstprstkEstkCnt)
	assert.Equal(t, "런던증권거래소(영국)", got.LstexNt)
	assert.Equal(t, "2025년 년 월 1일일", got.Lstprd)
}

func TestOverseasListing(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovLst.json": "ovLst.json"})
	items, err := c.OverseasListing(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231010000222", got.RceptNo)
	assert.Equal(t, "1,000,000", got.LststkOstkCnt)
	assert.Equal(t, "TEST", got.StkCd)
	assert.Equal(t, "2023년 10월 02일", got.Lstd)
}

func TestOverseasDelistingDecision(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovDlstDecsn.json": "ovDlstDecsn.json"})
	items, err := c.OverseasDelistingDecision(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20200101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20241031000508", got.RceptNo)
	assert.Equal(t, "38,685,850", got.DlststkEstkCnt)
	assert.Equal(t, "룩셈부르크증권거래소에 상장된 주식예탁증서(DR) 우선주의 거래량 미미", got.DlstRs)
	assert.Equal(t, "2024년 10월 31일", got.Bddd)
}

func TestOverseasDelisting(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ovDlst.json": "ovDlst.json"})
	items, err := c.OverseasDelisting(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	got := items[0]
	assert.Equal(t, "20231215000444", got.RceptNo)
	assert.Equal(t, "미국 (NASDAQ)", got.LstexNt)
	assert.Equal(t, "2023년 12월 14일", got.Tredd)
	assert.Equal(t, "거래량 부족", got.DlstRs)
}
