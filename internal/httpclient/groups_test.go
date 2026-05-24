package httpclient

import (
	"context"
	"encoding/json"
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGroups_Success(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "KEY", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "00107598", r.URL.Query().Get("corp_code"))
		w.Write([]byte(`{"status":"000","message":"정상","group":[` +
			`{"title":"일반사항","list":[{"rcept_no":"20230515002454","corp_name":"남양유업"}]},` +
			`{"title":"증권의종류","list":[{"stksen":"우선주"},{"stksen":"보통주"}]}` +
			`]}`))
	})

	groups, err := GetGroups(context.Background(), c, "/api/estkRs.json", map[string]string{"corp_code": "00107598"})
	require.NoError(t, err)
	require.Len(t, groups, 2)

	assert.Equal(t, "일반사항", groups[0].Title)
	assert.Equal(t, "증권의종류", groups[1].Title)

	var general []struct {
		RceptNo  string `json:"rcept_no"`
		CorpName string `json:"corp_name"`
	}
	require.NoError(t, json.Unmarshal(groups[0].List, &general))
	require.Len(t, general, 1)
	assert.Equal(t, "20230515002454", general[0].RceptNo)
	assert.Equal(t, "남양유업", general[0].CorpName)

	var types []struct {
		Stksen string `json:"stksen"`
	}
	require.NoError(t, json.Unmarshal(groups[1].List, &types))
	require.Len(t, types, 2)
	assert.Equal(t, "우선주", types[0].Stksen)
}

func TestGetGroups_NoData(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	})
	_, err := GetGroups(context.Background(), c, "/api/estkRs.json", nil)
	assert.ErrorIs(t, err, ErrNoData)
}

func TestGetGroups_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"010","message":"등록되지 않은 인증키입니다."}`))
	})
	_, err := GetGroups(context.Background(), c, "/api/estkRs.json", nil)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, "010", apiErr.Status)
}
