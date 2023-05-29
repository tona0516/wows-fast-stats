// TODO backoff tests

package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// テスト用のデータ型.
type TestData struct {
	Name string `json:"name"`
}

func TestAPIClient_GetRequest_NoQuery(t *testing.T) {
	t.Parallel()

	// テスト用のデータとレスポンスを準備
	testData := TestData{Name: "test"}
	//nolint:errchkjson
	responseData, _ := json.Marshal(testData)

	// テスト用の HTTP サーバーを作成し、モックのレスポンスを返すように設定
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseData)
	}))
	defer server.Close()

	// APIClient のインスタンスを作成
	client := NewAPIClient[TestData](server.URL)

	// GetRequest を呼び出してレスポンスを取得
	response, err := client.GetRequest(map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, testData, response)
}

func TestAPIClient_GetRequest_InvalidResponseBody(t *testing.T) {
	t.Parallel()

	// テスト用の無効な JSON データ
	invalidJSON := []byte("invalid-json")

	// テスト用の HTTP サーバーを作成し、無効な JSON レスポンスを返すように設定
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(invalidJSON)
	}))
	defer server.Close()

	// APIClient のインスタンスを作成
	client := NewAPIClient[TestData](server.URL)

	// GetRequest を呼び出してエラーを確認
	response, err := client.GetRequest(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, TestData{}, response)
}

func TestAPIClient_GetRequest_JSONUnmarshalError(t *testing.T) {
	t.Parallel()

	// テスト用の無効な JSON データ
	invalidJSON := []byte(`{"name":}`)

	// テスト用の HTTP サーバーを作成し、無効な JSON レスポンスを返すように設定
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(invalidJSON)
	}))
	defer server.Close()

	// APIClient のインスタンスを作成
	client := NewAPIClient[TestData](server.URL)

	// GetRequest を呼び出してエラーを確認
	response, err := client.GetRequest(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, TestData{}, response)
}

func TestAPIClient_MaxRetry_Failure(t *testing.T) {
	t.Parallel()

	// テスト用の HTTP サーバーを作成し、モックのレスポンスを返すように設定
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint:forcetypeassert
		c, _, _ := w.(http.Hijacker).Hijack()
		defer c.Close()
		calls++
	}))
	defer server.Close()

	// APIClient のインスタンスを作成
	client := NewAPIClient[TestData](server.URL)

	// GetRequest を呼び出してレスポンスを取得
	_, err := client.GetRequest(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, 4, calls)
}

func TestAPIClient_MaxRetry_Success(t *testing.T) {
	t.Parallel()

	// テスト用のデータとレスポンスを準備
	testData := TestData{Name: "test"}
	//nolint:errchkjson
	responseData, _ := json.Marshal(testData)

	// テスト用の HTTP サーバーを作成し、モックのレスポンスを返すように設定
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 4 {
			//nolint:forcetypeassert
			c, _, _ := w.(http.Hijacker).Hijack()
			defer c.Close()
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(responseData)
		}
	}))
	defer server.Close()

	// APIClient のインスタンスを作成
	client := NewAPIClient[TestData](server.URL)

	// GetRequest を呼び出してレスポンスを取得
	_, err := client.GetRequest(map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 4, calls)
}
