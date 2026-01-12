package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
	"wfs/backend/config"
	"wfs/backend/data"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWargamingClient_AccountInfo(t *testing.T) {
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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			WargamingClient: config.WargamingConfig{
				URL:          server.URL,
				RetryCount:   0,
				Timeout:      0,
				RateLimitRPS: 1,
				AppID:        "",
			},
		})
		wargaming, err := NewWargamingClient(injector)
		require.NoError(t, err)
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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			WargamingClient: config.WargamingConfig{
				URL:          server.URL,
				RetryCount:   0,
				Timeout:      timeout - 1,
				RateLimitRPS: 1,
				AppID:        "",
			},
		})
		instance, err := NewWargamingClient(injector)
		require.NoError(t, err)
		_, err = instance.AccountInfo([]int{123, 456})

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

		injector := do.New()
		do.ProvideValue(injector, config.Config{
			WargamingClient: config.WargamingConfig{
				URL:          server.URL,
				RetryCount:   0,
				Timeout:      0,
				RateLimitRPS: 1,
				AppID:        "",
			},
		})
		instance, err := NewWargamingClient(injector)
		require.NoError(t, err)
		_, err = instance.AccountInfo([]int{123, 456})

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

			injector := do.New()
			do.ProvideValue(injector, config.Config{
				WargamingClient: config.WargamingConfig{
					URL:          server.URL,
					RetryCount:   retry,
					Timeout:      0,
					RateLimitRPS: 1,
					AppID:        "",
				},
			})
			instance, err := NewWargamingClient(injector)
			require.NoError(t, err)
			_, err = instance.AccountInfo([]int{123, 456})

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

			injector := do.New()
			do.ProvideValue(injector, config.Config{
				WargamingClient: config.WargamingConfig{
					URL:          server.URL,
					RetryCount:   retry,
					Timeout:      0,
					RateLimitRPS: 1,
					AppID:        "",
				},
			})
			instance, err := NewWargamingClient(injector)
			require.NoError(t, err)
			_, err = instance.AccountInfo([]int{123, 456})

			assert.Error(t, err, ErrTemporaryUnavaillalble)
			assert.Equal(t, retry+1, calls)
		}
	})
}

func TestWargamingClient_ClansAccountInfo(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_ClansInfo(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_ShipsStats(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_EncycShips(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_BattleArena(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_BattleTypes(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_ShipsBadges(t *testing.T) {
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

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		WargamingClient: config.WargamingConfig{
			URL:          server.URL,
			RetryCount:   0,
			Timeout:      0,
			RateLimitRPS: 1,
			AppID:        "",
		},
	})
	instance, err := NewWargamingClient(injector)
	require.NoError(t, err)
	result, err := instance.ShipsBadges(123)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargamingClient_fieldQuery(t *testing.T) {
	t.Parallel()

	// テストデータ
	type TestData struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Detail struct {
			Hoge string `json:"hoge"`
			Fuga string `json:"fuga"`
		} `json:"detail"`
	}

	// テスト対象のデータ型を取得
	dataType := reflect.TypeOf(TestData{})

	// fieldQuery 関数を実行して結果を取得
	instance := WargamingClient{}
	result := instance.fieldQuery(dataType)

	// 期待される結果
	expectedResult := "id,name,detail.hoge,detail.fuga"

	// 結果の比較
	assert.Equal(t, expectedResult, result)
}

func TestWargamingClient_toSnakeCase(t *testing.T) {
	t.Parallel()

	// テストデータ
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "camelCase", expected: "camel_case"},
		{input: "PascalCase", expected: "pascal_case"},
		{input: "snake_case", expected: "snake_case"},
		{input: "lowercase", expected: "lowercase"},
		{input: "UPPERCASE", expected: "uppercase"},
		{input: "mixed_Case", expected: "mixed_case"},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		// toSnakeCase 関数を実行して結果を取得
		instance := WargamingClient{}
		result := instance.toSnakeCase(tc.input)

		// 結果の比較
		assert.Equal(t, tc.expected, result)
	}
}
