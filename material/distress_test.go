package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOccurrences(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dfOcr.json": "dfOcr.json"})
	items, err := c.DefaultOccurrences(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "DH오토넥스", items[0].CorpName)
	assert.Equal(t, "당사 김제지점 발행 만기어음 부도", items[0].DfCn)
	assert.Equal(t, "48,322,175", items[0].DfAmt)
}

func TestBusinessSuspensions(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bsnSp.json": "bsnSp.json"})
	items, err := c.BusinessSuspensions(context.Background(), MaterialParams{CorpCode: "00153393", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "방적 사업", items[0].BsnspRm)
	assert.Equal(t, "97,254,982,693", items[0].BsnspAmt)
	assert.Equal(t, "2023년 08월 31일", items[0].Bsnspd)
}

func TestRehabilitationApplications(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/ctrcvsBgrq.json": "ctrcvsBgrq.json"})
	items, err := c.RehabilitationApplications(context.Background(), MaterialParams{CorpCode: "00126089", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "서울회생법원", items[0].Cpct)
	assert.Equal(t, "주식회사 대유플러스", items[0].Apcnt)
}

func TestDissolutionCauses(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/dsRsOcr.json": "dsRsOcr.json"})
	items, err := c.DissolutionCauses(context.Background(), MaterialParams{CorpCode: "00580603", BgnDe: "20200101", EndDe: "20201231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "존립기간의 만료", items[0].DsRs)
	assert.Equal(t, "2020년 03월 27일", items[0].DsRsd)
}

func TestCreditorBankManagementStart(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bnkMngtPcbg.json": "bnkMngtPcbg.json"})
	items, err := c.CreditorBankManagementStart(context.Background(), MaterialParams{CorpCode: "00153861", BgnDe: "20240101", EndDe: "20251231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "태영건설", items[0].CorpName)
	assert.Equal(t, "경영정상화", items[0].MngtRs)
	assert.Equal(t, "2024년 01월 11일", items[0].MngtPcbgDd)
}

func TestCreditorBankManagementStop(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/bnkMngtPcsp.json": "bnkMngtPcsp.json"})
	items, err := c.CreditorBankManagementStop(context.Background(), MaterialParams{CorpCode: "00245481", BgnDe: "20200101", EndDe: "20201231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "경영정상화 완료", items[0].SpRs)
	assert.Equal(t, "2020년 04월 10일", items[0].MngtPcspDd)
}

func TestLawsuits(t *testing.T) {
	c := newTestClient(t, map[string]string{"/api/lwstLg.json": "lwstLg.json"})
	items, err := c.Lawsuits(context.Background(), MaterialParams{CorpCode: "01070149", BgnDe: "20240101", EndDe: "20241231"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "신주발행금지 등 임시의 지위를 구하는 가처분", items[0].Icnm)
	assert.Equal(t, "수원지방법원", items[0].Cpct)
}
