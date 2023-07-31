package infra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWargaming_AccountInfo(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGAccountInfo{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.AccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, domain.WGAccountInfo{}, result)
}

func TestWargaming_AccountList(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGAccountList{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.AccountList([]string{"player_1", "player_2"})

	assert.NoError(t, err)
	assert.Equal(t, domain.WGAccountList{}, result)
}

func TestWargaming_ClansAccountInfo(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGClansAccountInfo{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.ClansAccountInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, domain.WGClansAccountInfo{}, result)
}

func TestWargaming_ClansInfo(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGClansInfo{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.ClansInfo([]int{123, 456})

	assert.NoError(t, err)
	assert.Equal(t, domain.WGClansInfo{}, result)
}

func TestWargaming_ShipsStats(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGShipsStats{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.ShipsStats(123)

	assert.NoError(t, err)
	assert.Equal(t, domain.WGShipsStats{}, result)
}

func TestWargaming_EncycShips(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGEncycShips{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.EncycShips(1)

	assert.NoError(t, err)
	assert.Equal(t, domain.WGEncycShips{}, result)
}

func TestWargaming_BattleArena(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGBattleArenas{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	result, err := wargaming.BattleArenas()

	assert.NoError(t, err)
	assert.Equal(t, domain.WGBattleArenas{}, result)
}

func TestWargaming_BattleTypes(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGBattleTypes{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})
	result, err := wargaming.BattleTypes()

	assert.NoError(t, err)
	assert.Equal(t, domain.WGBattleTypes{}, result)
}

func TestWargaming_Test(t *testing.T) {
	t.Parallel()

	body, err := json.Marshal(domain.WGEncycInfo{})
	assert.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer server.Close()

	wargaming := NewWargaming(RequestConfig{
		URL: server.URL,
	})

	valid, err := wargaming.Test("hoge")
	assert.True(t, valid)
	assert.NoError(t, err)
}

func TestWargaming_AccountInfo_異常系_リトライなし(t *testing.T) {
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
	})

	_, err := wargaming.AccountInfo([]int{123, 456})

	assert.EqualError(t, err, apperr.New(apperr.ErrWGAPI, errors.New(body)).Error())
	assert.Equal(t, 1, calls)
}

func TestWargaming_AccountInfo_正常系_最大リトライ(t *testing.T) {
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
		var calls int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if calls < int(retry+1) {
				_, _ = w.Write([]byte(body))
				return
			}

			body, _ := json.Marshal(domain.WGAccountInfo{})
			_, _ = w.Write(body)
		}))
		defer server.Close()

		wargaming := NewWargaming(RequestConfig{URL: server.URL, Retry: retry})

		_, err := wargaming.AccountInfo([]int{123, 456})

		assert.NoError(t, err)
		assert.Equal(t, int(retry+1), calls)
	}
}

func TestWargaming_AccountInfo_異常系_最大リトライ(t *testing.T) {
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
		var calls int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		}))
		defer server.Close()

		wargaming := NewWargaming(RequestConfig{URL: server.URL, Retry: retry})

		_, err := wargaming.AccountInfo([]int{123, 456})

		assert.EqualError(t, err, apperr.New(
			apperr.ErrWGAPITemporaryUnavaillalble,
			errors.New(body),
		).Error())
		assert.Equal(t, int(retry+1), calls)
	}
}
