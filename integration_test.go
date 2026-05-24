//go:build integration

package opendart

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kenshin579/opendart/material"
	"github.com/kenshin579/opendart/report"
)

// 실행: OPENDART_API_KEY=... go test -tags integration -run TestIntegration -v
func TestIntegration_GetCompany(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	require.Equal(t, "00126380", corp)

	company, err := c.Disclosure.GetCompany(context.Background(), corp)
	require.NoError(t, err)
	require.Contains(t, company.CorpName, "삼성전자")
}

func TestIntegration_Dividend(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.Dividend(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_DebtSecuritiesIssuance(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.DebtSecuritiesIssuance(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_AuditOpinion(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.AuditOpinion(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_Executives(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.Executives(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_SingleAccount(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Report.SingleAccount(context.Background(), report.ReportParams{
		CorpCode:  corp,
		BsnsYear:  "2023",
		ReprtCode: report.AnnualReport,
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_XbrlTaxonomy(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Report.XbrlTaxonomy(context.Background(), "BS1")
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_MajorStockReports(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Ownership.MajorStockReports(context.Background(), corp)
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_DefaultOccurrences(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Material.DefaultOccurrences(context.Background(), material.MaterialParams{
		CorpCode: "00126089", // DH오토넥스 (실제 부도 사례)
		BgnDe:    "20230101",
		EndDe:    "20231231",
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_PaidInCapitalIncrease(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Material.PaidInCapitalIncrease(context.Background(), material.MaterialParams{
		CorpCode: "00107598", // 남양유업 (실제 유상증자 사례)
		BgnDe:    "20230101",
		EndDe:    "20231231",
	})
	require.NoError(t, err)
	require.NotEmpty(t, items)
}

func TestIntegration_ConvertibleBondIssuance(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	items, err := c.Material.ConvertibleBondIssuance(context.Background(), material.MaterialParams{
		CorpCode: "00126380", // 삼성전자 (전환사채 발행 사례 없을 수 있음 → ErrNoData 허용)
		BgnDe:    "20200101",
		EndDe:    "20241231",
	})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 전환사채 발행 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_TreasuryStockAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)

	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)

	items, err := c.Material.TreasuryStockAcquisition(context.Background(), material.MaterialParams{
		CorpCode: corp,
		BgnDe:    "20200101",
		EndDe:    "20241231",
	})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 자기주식 취득 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_BusinessAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.BusinessAcquisition(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 영업양수 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_OtherCorpStockAcquisition(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OtherCorpStockAcquisition(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 타법인주식 양수 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_OtherAssetTransferPutbackOption(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.OtherAssetTransferPutbackOption(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 자산양수도(기타) 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}

func TestIntegration_StockExchangeTransfer(t *testing.T) {
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	corp, err := c.ResolveCorpCode(context.Background(), "005930")
	require.NoError(t, err)
	items, err := c.Material.StockExchangeTransfer(context.Background(), material.MaterialParams{CorpCode: corp, BgnDe: "20200101", EndDe: "20241231"})
	if errors.Is(err, ErrNoData) {
		t.Skip("해당 기간 주식교환·이전 데이터 없음")
	}
	require.NoError(t, err)
	for _, it := range items {
		require.NotEmpty(t, it.RceptNo)
	}
}
