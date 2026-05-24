package material

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertibleBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/cvbdIsDecsn.json": "cvbdIsDecsn.json",
	})

	items, err := c.ConvertibleBondIssuance(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230815000123", got.RceptNo)
	assert.Equal(t, "30,000,000,000", got.BdFta)
	assert.Equal(t, "15,000", got.CvPrc)
	assert.Equal(t, "2024년 08월 16일", got.CvrqpdBgd)
	assert.Equal(t, "미제출", got.RsSmAtn)
	assert.Equal(t, "기명식 무보증 사모 전환사채", got.BdKnd)
	assert.Equal(t, "10,500", got.ActMktprcflCvprcLwtrsprc)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestBondWithWarrantIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/bdwtIsDecsn.json": "bdwtIsDecsn.json",
	})

	items, err := c.BondWithWarrantIssuance(context.Background(), MaterialParams{CorpCode: "00164779", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230910000222", got.RceptNo)
	assert.Equal(t, "8,000", got.ExPrc)
	assert.Equal(t, "분리형", got.BdwtDivAtn)
	assert.Equal(t, "2024년 09월 12일", got.ExpdBgd)
	assert.Equal(t, "2,500,000", got.NstkIsstkCnt)
	assert.Equal(t, "20,000,000,000", got.BdFta)
	assert.Equal(t, "미제출", got.RsSmAtn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestExchangeableBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/exbdIsDecsn.json": "exbdIsDecsn.json",
	})

	items, err := c.ExchangeableBondIssuance(context.Background(), MaterialParams{CorpCode: "00126380", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230705000333", got.RceptNo)
	assert.Equal(t, "50,000,000,000", got.BdFta)
	assert.Equal(t, "25,000", got.ExPrc)
	assert.Equal(t, "자기주식(기명식 보통주)", got.Extg)
	assert.Equal(t, "2024년 07월 06일", got.ExrqpdBgd)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}

func TestContingentConvertibleBondIssuance(t *testing.T) {
	c := newTestClient(t, map[string]string{
		"/api/wdCocobdIsDecsn.json": "wdCocobdIsDecsn.json",
	})

	items, err := c.ContingentConvertibleBondIssuance(context.Background(), MaterialParams{CorpCode: "00164779", BgnDe: "20230101", EndDe: "20231231"})
	require.NoError(t, err)
	require.Len(t, items, 1)

	got := items[0]
	assert.Equal(t, "20230601000444", got.RceptNo)
	assert.Equal(t, "100,000,000,000", got.BdFta)
	assert.Equal(t, "4.5", got.BdIntrSf) // 표면이자율 (이 struct 에선 sf/ex 라벨이 다른 3종과 반대)
	assert.Equal(t, "4.5", got.BdIntrEx) // 만기이자율
	assert.Contains(t, got.DbtrsSc, "부실금융기관")
	assert.Equal(t, "제출", got.RsSmAtn)
	assert.Equal(t, "미해당", got.FtcSttAtn)
}
