package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wfs/internal/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/ratelimit"
)

func TestWargamingApiClient_AccountInfo(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := data.WGAccountInfo{
			WGResponseCommon: data.WGResponseCommon[map[int]data.WGAccountInfoData]{
				Status: "ok",
				Error:  data.WGError{},
				Data:   map[int]data.WGAccountInfoData{},
			},
		}
		server := simpleMockServer(t, 200, expected)
		defer server.Close()

		wargaming := NewWargamingApiClient(
			"",
			*NewApiConfig(
				server.URL,
				0,
				0,
			),
			ratelimit.NewUnlimited(),
		)
		result, err := wargaming.AccountInfo([]int{123, 456})

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系_タイムアウト", func(t *testing.T) {
		t.Parallel()

		timeout := time.Millisecond * 10
		server := timeoutMockServer(
			t,
			http.StatusOK,
			map[string]any{},
			timeout,
		)
		defer server.Close()

		instance := NewWargamingApiClient(
			"",
			*NewApiConfig(
				server.URL,
				0,
				timeout-1,
			),
			ratelimit.NewUnlimited(),
		)
		_, err := instance.AccountInfo([]int{123, 456})

		assert.ErrorIs(t, err, context.DeadlineExceeded)
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
			w.Write([]byte(body))
		}))
		defer server.Close()

		instance := NewWargamingApiClient(
			"",
			*NewApiConfig(
				server.URL,
				0,
				0,
			),
			ratelimit.NewUnlimited(),
		)
		_, err := instance.AccountInfo([]int{123, 456})

		assert.Error(t, err, ErrErrorResponse)
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

				body, err := json.Marshal(data.WGAccountInfo{
					WGResponseCommon: data.WGResponseCommon[map[int]data.WGAccountInfoData]{
						Status: "ok",
						Error:  data.WGError{},
						Data:   map[int]data.WGAccountInfoData{},
					},
				})
				require.NoError(t, err)
				w.Write(body)
			}))
			defer server.Close()

			instance := NewWargamingApiClient(
				"",
				*NewApiConfig(
					server.URL,
					retry,
					0,
				),
				ratelimit.NewUnlimited(),
			)
			_, err := instance.AccountInfo([]int{123, 456})

			assert.NoError(t, err)
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
				w.Write([]byte(body))
			}))
			defer server.Close()

			instance := NewWargamingApiClient(
				"",
				*NewApiConfig(
					server.URL,
					retry,
					0,
				),
				ratelimit.NewUnlimited(),
			)
			_, err := instance.AccountInfo([]int{123, 456})

			assert.Error(t, err, ErrTemporaryUnavaillalble)
			assert.Equal(t, retry+1, calls)
		}
	})
}

func TestWargamingApiClient_AccountListForSearch(t *testing.T) {
	t.Parallel()

	expected := data.WGAccountList{
		WGResponseCommon: data.WGResponseCommon[[]data.WGAccountListData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   []data.WGAccountListData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.AccountListForSearch("player")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	expected := data.WGClansAccountInfo{
		WGResponseCommon: data.WGResponseCommon[map[int]data.WGClansAccountInfoData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int]data.WGClansAccountInfoData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_ClansInfo(t *testing.T) {
	t.Parallel()

	expected := data.WGClansInfo{
		WGResponseCommon: data.WGResponseCommon[map[int]data.WGClansInfoData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int]data.WGClansInfoData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_ShipsStats(t *testing.T) {
	t.Parallel()

	expected := data.WGShipsStats{
		WGResponseCommon: data.WGResponseCommon[map[int][]data.WGShipsStatsData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int][]data.WGShipsStatsData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_EncycShips(t *testing.T) {
	t.Parallel()

	expected := data.WGEncycShips{
		WGResponseCommon: data.WGResponseCommon[map[int]data.WGEncycShipsData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int]data.WGEncycShipsData{},
		},
		Meta: struct {
			PageTotal int `json:"page_total"`
			Page      int `json:"page"`
		}{PageTotal: 5},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_BattleArena(t *testing.T) {
	t.Parallel()

	expected := data.WGBattleArenas{
		WGResponseCommon: data.WGResponseCommon[map[int]data.WGBattleArenasData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int]data.WGBattleArenasData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_BattleTypes(t *testing.T) {
	t.Parallel()

	expected := data.WGBattleTypes{
		WGResponseCommon: data.WGResponseCommon[map[string]data.WGBattleTypesData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[string]data.WGBattleTypesData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingApiClient_ShipsBadges(t *testing.T) {
	t.Parallel()

	expected := data.WGShipsBadges{
		WGResponseCommon: data.WGResponseCommon[map[int][]data.WGShipsBadgesData]{
			Status: "ok",
			Error:  data.WGError{},
			Data:   map[int][]data.WGShipsBadgesData{},
		},
	}
	server := simpleMockServer(t, 200, expected)
	defer server.Close()

	instance := NewWargamingApiClient(
		"",
		*NewApiConfig(
			server.URL,
			0,
			0,
		),
		ratelimit.NewUnlimited(),
	)
	result, err := instance.ShipsBadges(123)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
