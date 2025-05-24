package infra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		server := simpleMockServer(200, response.WGAccountInfo{
			WGResponseCommon: response.WGResponseCommon[data.WGAccountInfo]{
				Status: "",
				Error:  response.WGError{},
				Data:   map[int]data.WGAccountInfoData{},
			},
		})
		defer server.Close()

		wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
		result, err := wargaming.AccountInfo([]int{123, 456})

		require.NoError(t, err)
		assert.Equal(t, data.WGAccountInfo{}, result)
	})

	t.Run("異常系_リトライなし", func(t *testing.T) {
		t.Parallel()
		body := `{
            "status":"error",
            "error":{
                "field":null,
                "message":"INVALID_APPLICATION_ID",
                "code":407,
                "value":null
            }
        }`

		var calls int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		}))
		defer server.Close()

		wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
		_, err := wargaming.AccountInfo([]int{123, 456})

		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.WGAPIError, code)
		assert.Equal(t, 1, calls)
	})

	t.Run("正常系_最大リトライ", func(t *testing.T) {
		t.Parallel()
		messages := []string{
			"REQUEST_LIMIT_EXCEEDED",
			"SOURCE_NOT_AVAILABLE",
		}

		for _, message := range messages {
			body := fmt.Sprintf(`{
                "status":"error",
                "error":{
                    "field":null,
                    "message":"%s",
                    "code":407,
                    "value":null
                }
            }`, message)

			var retry int = 1
			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				if calls < retry+1 {
					_, _ = w.Write([]byte(body))
					return
				}

				body, _ := json.Marshal(response.WGAccountInfo{})
				_, _ = w.Write(body)
			}))
			defer server.Close()

			wargaming := NewWargaming(server.URL, retry, 0, 0, 1, "")
			_, err := wargaming.AccountInfo([]int{123, 456})

			require.NoError(t, err)
			assert.Equal(t, retry+1, calls)
		}
	})

	t.Run("異常系_最大リトライ", func(t *testing.T) {
		t.Parallel()
		messages := []string{
			"REQUEST_LIMIT_EXCEEDED",
			"SOURCE_NOT_AVAILABLE",
		}

		for _, message := range messages {
			body := fmt.Sprintf(`{
                "status":"error",
                "error":{
                    "field":null,
                    "message":"%s",
                    "code":407,
                    "value":null
                }
            }`, message)

			var retry int = 1
			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
			}))
			defer server.Close()

			wargaming := NewWargaming(server.URL, retry, 0, 0, 1, "")
			_, err := wargaming.AccountInfo([]int{123, 456})

			code, ok := failure.CodeOf(err)
			assert.True(t, ok)
			assert.Equal(t, apperr.WGAPITemporaryUnavaillalble, code)
			assert.Equal(t, retry+1, calls)
		}
	})
}

func TestWargaming_AccountListForSearch(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGAccountList{
		WGResponseCommon: response.WGResponseCommon[data.WGAccountList]{
			Status: "",
			Error:  response.WGError{},
			Data:   []data.WGAccountListData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.AccountListForSearch("player")

	require.NoError(t, err)
	assert.Equal(t, data.WGAccountList{}, result)
}

func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGClansAccountInfo{
		WGResponseCommon: response.WGResponseCommon[data.WGClansAccountInfo]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int]data.WGClansAccountInfoData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.ClansAccountInfo([]int{123, 456})

	require.NoError(t, err)
	assert.Equal(t, data.WGClansAccountInfo{}, result)
}

func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGClansInfo{
		WGResponseCommon: response.WGResponseCommon[data.WGClansInfo]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int]data.WGClansInfoData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.ClansInfo([]int{123, 456})

	require.NoError(t, err)
	assert.Equal(t, data.WGClansInfo{}, result)
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGShipsStats{
		WGResponseCommon: response.WGResponseCommon[data.WGShipsStats]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int][]data.WGShipsStatsData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.ShipsStats(123)

	require.NoError(t, err)
	assert.Equal(t, data.WGShipsStats{}, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	expectedPageTotal := 5
	server := simpleMockServer(200, response.WGEncycShips{
		WGResponseCommon: response.WGResponseCommon[data.WGEncycShips]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int]data.WGEncycShipsData{},
		},
		Meta: struct {
			PageTotal int `json:"page_total"`
			Page      int `json:"page"`
		}{PageTotal: expectedPageTotal},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, pageTotal, err := wargaming.EncycShips(1)

	require.NoError(t, err)
	assert.Equal(t, data.WGEncycShips{}, result)
	assert.Equal(t, expectedPageTotal, pageTotal)
}

func TestWargaming_EncycInfo(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGEncycInfo{})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.EncycInfo()

	require.NoError(t, err)
	assert.Equal(t, data.WGEncycInfoData{}, result)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGBattleArenas{
		WGResponseCommon: response.WGResponseCommon[data.WGBattleArenas]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int]data.WGBattleArenasData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.BattleArenas()

	require.NoError(t, err)
	assert.Equal(t, data.WGBattleArenas{}, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGBattleTypes{
		WGResponseCommon: response.WGResponseCommon[data.WGBattleTypes]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[string]data.WGBattleTypesData{},
		},
	})
	defer server.Close()

	wargaming := NewWargaming(server.URL, 0, 0, 0, 1, "")
	result, err := wargaming.BattleTypes()

	require.NoError(t, err)
	assert.Equal(t, data.WGBattleTypes{}, result)
}
