package webapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GET_Success(t *testing.T) {
	t.Parallel()

	// モックサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test-path", r.URL.Path)
		assert.Equal(t, "value", r.URL.Query().Get("key"))
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"GET success"}`))
	}))
	defer mockServer.Close()

	// クライアントの作成
	client := NewClient(mockServer.URL,
		WithPath("test-path"),
		WithQuery(map[string]string{"key": "value"}),
		WithHeaders(map[string]string{"Authorization": "Bearer token"}),
	)

	// GETリクエストの実行
	res, body, err := client.GET()

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{"message":"GET success"}`, string(body))
}

func TestClient_POST_Success(t *testing.T) {
	t.Parallel()

	// モックサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test-path", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"key":"value"}`, string(body))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"message":"POST success"}`))
	}))
	defer mockServer.Close()

	// クライアントの作成
	client := NewClient(mockServer.URL,
		WithPath("test-path"),
		WithBody(map[string]string{"key": "value"}),
		WithHeaders(map[string]string{"Content-Type": "application/json"}),
	)

	// POSTリクエストの実行
	res, body, err := client.POST()

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, `{"message":"POST success"}`, string(body))
}

func TestClient_GET_ErrorResponse(t *testing.T) {
	t.Parallel()

	// モックサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal server error"}`))
	}))
	defer mockServer.Close()

	// クライアントの作成
	client := NewClient(mockServer.URL)

	// GETリクエストの実行
	res, body, err := client.GET()

	// 結果の検証
	require.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.JSONEq(t, `{"error":"Internal server error"}`, string(body))
}

func TestClient_Timeout(t *testing.T) {
	t.Parallel()

	// モックサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// クライアントの作成
	client := NewClient(mockServer.URL, WithTimeout(1*time.Second))

	// GETリクエストの実行
	_, _, err := client.GET()

	// 結果の検証
	assert.Error(t, err)
}

func TestClient_InsecureTLS(t *testing.T) {
	t.Parallel()

	// HTTPSモックサーバーの作成
	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Insecure TLS success"}`))
	}))
	defer mockServer.Close()

	// クライアントの作成
	client := NewClient(mockServer.URL,
		WithIsInsecure(true),
	)

	// GETリクエストの実行
	res, body, err := client.GET()

	// 結果の検証
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{"message":"Insecure TLS success"}`, string(body))
}

func TestClient_MissingBaseURL(t *testing.T) {
	t.Parallel()

	// クライアントの作成
	client := NewClient("")

	// GETリクエストの実行
	_, _, err := client.GET()

	// 結果の検証
	assert.Error(t, err)
}
