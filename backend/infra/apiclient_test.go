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

func TestAPIClient_GetRequest(t *testing.T) {
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
	client := &APIClient[TestData]{}

	// GetRequest を呼び出してレスポンスを取得
	response, err := client.GetRequest(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, testData, response)
}

func TestAPIClient_GetRequest_Error(t *testing.T) {
	t.Parallel()

	// APIClient のインスタンスを作成
	client := &APIClient[TestData]{}

	// テスト用の無効な URL
	invalidURL := "invalid-url"

	// GetRequest を呼び出してエラーを確認
	response, err := client.GetRequest(invalidURL)
	assert.Error(t, err)
	assert.Equal(t, TestData{}, response)
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
	client := &APIClient[TestData]{}

	// GetRequest を呼び出してエラーを確認
	response, err := client.GetRequest(server.URL)
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
	client := &APIClient[TestData]{}

	// GetRequest を呼び出してエラーを確認
	response, err := client.GetRequest(server.URL)
	assert.Error(t, err)
	assert.Equal(t, TestData{}, response)
}
