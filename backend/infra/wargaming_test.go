package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra/response"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/ratelimit"
)

const testAppID = "test_appid"

//nolint:gochecknoglobals
var limiter = ratelimit.NewUnlimited()

func NewWGResponse(
	status string,
	errCode int,
	errMsg string,
) *response.WGResponse[any] {
	return &response.WGResponse[any]{
		Status: status,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    errCode,
			Message: errMsg,
		},
	}
}

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()
	successResponse := response.WGResponse[data.WGAccountInfo]{
		Data: data.WGAccountInfo{},
	}
	retryableMessages := []string{
		"REQUEST_LIMIT_EXCEEDED",
		"SOURCE_NOT_AVAILABLE",
	}

	t.Run("データを取得できる", func(t *testing.T) {
		t.Parallel()
		server := simpleMockServer(200, successResponse)
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			1*time.Second,
			0,
		), limiter)
		result, err := wargaming.AccountInfo(testAppID, []int{123, 456})

		require.NoError(t, err)
		assert.Equal(t, data.WGAccountInfo{}, result)
	})

	t.Run("最大回数リトライしてデータを取得できる", func(t *testing.T) {
		t.Parallel()
		retry := 2

		for _, message := range retryableMessages {
			body := NewWGResponse("error", 407, message)

			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				if calls < retry+1 {
					jsonStr, _ := json.Marshal(body)
					_, _ = w.Write(jsonStr)
					return
				}

				body, _ := json.Marshal(successResponse)
				_, _ = w.Write(body)
			}))
			defer server.Close()

			wargaming := NewWargaming(*NewAPIConfig(
				server.URL,
				1*time.Second,
				retry,
			), limiter)

			actual, err := wargaming.AccountInfo(testAppID, []int{123, 456})

			require.NoError(t, err)
			assert.Equal(t, data.WGAccountInfo{}, actual)
			assert.Equal(t, retry+1, calls)
		}
	})

	t.Run("タイムアウト", func(t *testing.T) {
		t.Parallel()
		timeout := 1 * time.Second

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(timeout * 2)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			jsonStr, _ := json.Marshal(successResponse)
			_, _ = w.Write(jsonStr)
		}))
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			timeout,
			0,
		), limiter)

		_, err := wargaming.AccountInfo(testAppID, []int{123, 456})
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.WGAPIError, code)
	})

	t.Run("リトライせずにエラーを返却する", func(t *testing.T) {
		t.Parallel()
		retry := 2

		body := NewWGResponse("error", 407, "INVALID_APPLICATION_ID")

		var calls int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			jsonStr, _ := json.Marshal(body)
			_, _ = w.Write(jsonStr)
		}))
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			1*time.Second,
			retry,
		), limiter)

		_, err := wargaming.AccountInfo(testAppID, []int{123, 456})
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.WGAPIError, code)
		assert.Equal(t, 1, calls)
	})

	t.Run("最大回数リトライしてエラーを返却する", func(t *testing.T) {
		t.Parallel()
		retry := 2

		for _, message := range retryableMessages {
			body := NewWGResponse("error", 407, message)

			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				jsonStr, _ := json.Marshal(body)
				_, _ = w.Write(jsonStr)
			}))
			defer server.Close()

			wargaming := NewWargaming(*NewAPIConfig(
				server.URL,
				1*time.Second,
				retry,
			), limiter)

			_, err := wargaming.AccountInfo(testAppID, []int{123, 456})

			code, ok := failure.CodeOf(err)
			assert.True(t, ok)
			assert.Equal(t, apperr.WGAPITemporaryUnavaillalble, code)
			assert.Equal(t, 3, calls)
		}
	})
}

func TestWargaming_AccountList(t *testing.T) {
	t.Parallel()

	t.Run("データを取得できる", func(t *testing.T) {
		t.Parallel()
		server := simpleMockServer(200, response.WGResponse[data.WGAccountList]{
			Data: []data.WGAccountListData{},
		})
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.AccountList(testAppID, []string{"A", "B"})

		require.NoError(t, err)
		assert.Equal(t, data.WGAccountList{}, result)
	})

	t.Run("要素0の場合リクエストせずにデータを返却する", func(t *testing.T) {
		t.Parallel()
		wargaming := NewWargaming(*NewAPIConfig(
			"",
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.AccountList(testAppID, []string{})

		require.NoError(t, err)
		assert.Equal(t, data.WGAccountList{}, result)
	})
}

func TestWargaming_AccountListForSearch(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGResponse[data.WGAccountList]{
		Data: []data.WGAccountListData{},
	})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	result, err := wargaming.AccountListForSearch(testAppID, "player")

	require.NoError(t, err)
	assert.Equal(t, data.WGAccountList{}, result)
}

//nolint:dupl
func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	t.Run("データを取得できる", func(t *testing.T) {
		t.Parallel()
		server := simpleMockServer(200, response.WGResponse[data.WGClansAccountInfo]{
			Data: map[int]data.WGClansAccountInfoData{},
		})
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.ClansAccountInfo(testAppID, []int{123, 456})

		require.NoError(t, err)
		assert.Equal(t, data.WGClansAccountInfo{}, result)
	})

	t.Run("要素0の場合リクエストせずにデータを返却する", func(t *testing.T) {
		t.Parallel()
		wargaming := NewWargaming(*NewAPIConfig(
			"",
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.ClansAccountInfo(testAppID, []int{})

		require.NoError(t, err)
		assert.Equal(t, data.WGClansAccountInfo{}, result)
	})
}

//nolint:dupl
func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	t.Run("データを取得できる", func(t *testing.T) {
		t.Parallel()
		server := simpleMockServer(200, response.WGResponse[data.WGClansInfo]{
			Data: map[int]data.WGClansInfoData{},
		})
		defer server.Close()

		wargaming := NewWargaming(*NewAPIConfig(
			server.URL,
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.ClansInfo(testAppID, []int{123, 456})

		require.NoError(t, err)
		assert.Equal(t, data.WGClansInfo{}, result)
	})

	t.Run("要素0の場合リクエストせずにデータを返却する", func(t *testing.T) {
		t.Parallel()

		wargaming := NewWargaming(*NewAPIConfig(
			"",
			1*time.Second,
			0,
		), limiter)

		result, err := wargaming.ClansInfo(testAppID, []int{})

		require.NoError(t, err)
		assert.Equal(t, data.WGClansInfo{}, result)
	})
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGResponse[data.WGShipsStats]{
		Data: map[int][]data.WGShipsStatsData{},
	})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	result, err := wargaming.ShipsStats(testAppID, 123)

	require.NoError(t, err)
	assert.Equal(t, data.WGShipsStats{}, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	expectedPageTotal := 5
	server := simpleMockServer(200, response.WGEncycShips{
		WGResponse: response.WGResponse[data.WGEncycShips]{
			Data: map[int]data.WGEncycShipsData{},
		},
		Meta: struct {
			PageTotal int `json:"page_total"`
			Page      int `json:"page"`
		}{PageTotal: expectedPageTotal},
	})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	result, pageTotal, err := wargaming.EncycShips(testAppID, 1)

	require.NoError(t, err)
	assert.Equal(t, data.WGEncycShips{}, result)
	assert.Equal(t, expectedPageTotal, pageTotal)
}

func TestWargaming_EncycInfo(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGEncycInfo{})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	result, err := wargaming.EncycInfo(testAppID)

	require.NoError(t, err)
	assert.Equal(t, data.WGEncycInfoData{}, result)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGResponse[data.WGBattleArenas]{
		Data: map[int]data.WGBattleArenasData{},
	})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	result, err := wargaming.BattleArenas(testAppID)

	require.NoError(t, err)
	assert.Equal(t, data.WGBattleArenas{}, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGResponse[data.WGBattleTypes]{
		Data: map[string]data.WGBattleTypesData{},
	})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)
	result, err := wargaming.BattleTypes(testAppID)

	require.NoError(t, err)
	assert.Equal(t, data.WGBattleTypes{}, result)
}

func TestWargaming_Test(t *testing.T) {
	t.Parallel()

	server := simpleMockServer(200, response.WGEncycInfo{})
	defer server.Close()

	wargaming := NewWargaming(*NewAPIConfig(
		server.URL,
		1*time.Second,
		0,
	), limiter)

	valid := wargaming.Test("hoge")
	assert.True(t, valid)
}
