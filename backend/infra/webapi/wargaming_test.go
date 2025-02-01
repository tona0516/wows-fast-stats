package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/infra/response"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"

	"go.uber.org/ratelimit"
)

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := response.WGAccountInfo{
			114: {
				HiddenProfile: true,
			},
		}
		server := simpleMockServer(200, response.WGAccountInfoResponse{
			WGResponseCommon: response.WGResponseCommon[response.WGAccountInfo]{
				Status: "",
				Error:  response.WGError{},
				Data:   expected,
			},
		})
		defer server.Close()

		wargaming := NewWargaming(RequestConfig{
			URL: server.URL,
		}, ratelimit.New(10), "")

		result, err := wargaming.AccountInfo([]int{123, 456})
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
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

		wargaming := NewWargaming(RequestConfig{
			URL: server.URL,
		}, ratelimit.New(10), "")

		_, err := wargaming.AccountInfo([]int{123, 456})
		assert.True(t, failure.Is(err, apperr.WGAPIError))
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

			var retry uint64 = 1
			var calls uint64
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

			wargaming := NewWargaming(RequestConfig{URL: server.URL, Retry: retry}, ratelimit.New(10), "")

			_, err := wargaming.AccountInfo([]int{123, 456})

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

			var retry uint64 = 1
			var calls uint64
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
			}))
			defer server.Close()

			wargaming := NewWargaming(RequestConfig{URL: server.URL, Retry: retry}, ratelimit.New(10), "")

			_, err := wargaming.AccountInfo([]int{123, 456})

			assert.True(t, failure.Is(err, apperr.WGAPITemporaryUnavaillalble))
			assert.Equal(t, retry+1, calls)
		}
	})
}

func TestWargaming_AccountList(t *testing.T) {
	t.Parallel()

	expected := response.WGAccountList{
		{
			NickName:  "1919",
			AccountID: 114514,
		},
	}

	server := simpleMockServer(200, response.WGAccountListResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGAccountList]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	})
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.AccountList([]string{"player"})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_AccountListForSearch(t *testing.T) {
	t.Parallel()

	expected := response.WGAccountList{
		{
			NickName:  "1919",
			AccountID: 114514,
		},
	}

	server := simpleMockServer(200, response.WGAccountListResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGAccountList]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	})
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.AccountListForSearch("player")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	expected := response.WGClansAccountInfo{
		114514: response.WGClansAccountInfoData{ClanID: 1919},
	}
	res := response.WGClansAccountInfoResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGClansAccountInfo]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	}

	server := simpleMockServer(200, res)
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	expected := response.WGClansInfo{
		114: {
			Tag:         "514",
			Description: "1919",
		},
	}
	res := response.WGClansInfoResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGClansInfo]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	}

	server := simpleMockServer(200, res)
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	expected := response.WGShipsStats{
		114: {
			{
				ShipID: 514,
			},
		},
	}
	server := simpleMockServer(200, response.WGShipsStatsResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGShipsStats]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	})
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	expectedPageTotal := 5
	expected := response.WGEncycShips{
		WGResponseCommon: response.WGResponseCommon[map[int]response.WGEncycShipsData]{
			Status: "",
			Error:  response.WGError{},
			Data:   map[int]response.WGEncycShipsData{},
		},
		Meta: struct {
			PageTotal int `json:"page_total"`
			Page      int `json:"page"`
		}{PageTotal: expectedPageTotal},
	}

	server := simpleMockServer(200, expected)
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, pageTotal, err := wargaming.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, expectedPageTotal, pageTotal)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	expected := response.WGBattleArenas{
		114: {Name: "いつもの"},
	}

	server := simpleMockServer(200, response.WGBattleArenasResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGBattleArenas]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	})
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")

	result, err := wargaming.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	expected := response.WGBattleTypes{
		"PVP": {Name: "ランダム戦"},
	}

	server := simpleMockServer(200, response.WGBattleTypesResponse{
		WGResponseCommon: response.WGResponseCommon[response.WGBattleTypes]{
			Status: "",
			Error:  response.WGError{},
			Data:   expected,
		},
	})
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	}, ratelimit.New(10), "")
	result, err := wargaming.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
