package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaidInCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/piicDecsn.json": "piicDecsn.json"})
	items, err := c.PaidInCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "남양유업", items[0].CorpName)
	assert.Equal(t, "주주우선공모증자", items[0].IcMthn)
	assert.Equal(t, "7,184,339,000", items[0].FdppOp)
}

func TestFreeCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/fricDecsn.json": "fricDecsn.json"})
	items, err := c.FreeCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00117230", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "14,580,207", items[0].NstkOstkCnt)
	assert.Equal(t, "0.5", items[0].NstkAscntPsOstk)
	assert.Equal(t, "2,500", items[0].FvPs)
}

func TestPaidFreeCapitalIncrease(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/pifricDecsn.json": "pifricDecsn.json"})
	items, err := c.PaidFreeCapitalIncrease(context.Background(), MaterialParams{CorpCode: "00870481", BgnDe: "20220101", EndDe: "20221231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "주주배정후 실권주 일반공모", items[0].PiicIcMthn)
	assert.Equal(t, "23,500,000,000", items[0].PiicFdppDtrp)
	assert.Equal(t, "0.3", items[0].FricNstkAscntPsOstk)
}

func TestCapitalReduction(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/crDecsn.json": "crDecsn.json"})
	items, err := c.CapitalReduction(context.Background(), MaterialParams{CorpCode: "00295857", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "5.00", items[0].CrRtOstk)
	assert.Equal(t, "140,669,318,500", items[0].BfcrCpt)
	assert.Equal(t, "7,033,374,500", items[0].AtcrCpt)
}
