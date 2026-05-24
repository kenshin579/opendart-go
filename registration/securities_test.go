package registration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebtSecurities(t *testing.T) {
	c := newTestClient(t, "bdRs.json")
	res, err := c.DebtSecurities(context.Background(), Params{CorpCode: "00164779", BgnDe: "20180101", EndDe: "20241231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 22)
	require.Len(t, res.Underwriters, 100)
	require.Len(t, res.FundUsage, 48)
	require.Len(t, res.Sellers, 23)
	assert.Equal(t, "20180312000588", res.General[0].RceptNo)
	assert.Equal(t, "218", res.General[0].Tm)
	assert.Equal(t, "무보증사채", res.General[0].Bdnmn)
	assert.Equal(t, "SK하이닉스", res.General[0].CorpName)
	assert.Equal(t, "운영자금", res.FundUsage[0].Se)
}

func TestEquitySecurities(t *testing.T) {
	c := newTestClient(t, "estkRs.json")
	res, err := c.EquitySecurities(context.Background(), Params{CorpCode: "00107598", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.SecurityTypes, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 2)
	require.Len(t, res.Sellers, 1)
	require.Len(t, res.RetailPutbackOption, 1)
	assert.Equal(t, "20230515002454", res.General[0].RceptNo)
	assert.Equal(t, "남양유업", res.General[0].CorpName)
	assert.Equal(t, "우선주", res.SecurityTypes[0].Stksen)
	assert.Equal(t, "NH투자증권", res.Underwriters[0].Actnmn)
	assert.Equal(t, "운영자금", res.FundUsage[0].Se)
	assert.Equal(t, "발행제비용", res.FundUsage[1].Se)
}

func TestDepositaryReceipts(t *testing.T) {
	c := newTestClient(t, "stkdpRs.json")
	res, err := c.DepositaryReceipts(context.Background(), Params{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.General, 1)
	require.Len(t, res.SecurityTypes, 1)
	require.Len(t, res.Underwriters, 1)
	require.Len(t, res.FundUsage, 1)
	require.Len(t, res.Sellers, 1)
	assert.Equal(t, "20230701000222", res.General[0].RceptNo)
	assert.Equal(t, "주식예탁증서(DR)", res.SecurityTypes[0].Stksen)
	assert.Equal(t, "시설자금", res.FundUsage[0].Se)
}
