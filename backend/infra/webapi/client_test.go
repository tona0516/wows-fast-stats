package webapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のデータ型.
type TestRequestBody struct {
	ID string `json:"ID"`
}

type TestResponseBody struct {
	Name string `json:"name"`
}

func mockServer[T, U any](response Response[T, U], responseTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(responseTime)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(response.StatusCode)
		_, _ = w.Write(response.BodyByte)
	}))
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	t.Run("正常系_クエリなし", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[any, TestResponseBody]{
			Request: Request[any]{
				Method: http.MethodGet,
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 0)
		defer server.Close()
		expected.Request.URL = server.URL

		actual, err := GetRequest[TestResponseBody](server.URL, 1*time.Second, nil, nil)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("正常系_クエリあり", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[any, TestResponseBody]{
			Request: Request[any]{
				Method: http.MethodGet,
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 0)
		defer server.Close()
		expected.Request.URL = server.URL + "?hoge=fuga"

		actual, err := GetRequest[TestResponseBody](server.URL, 1*time.Second, map[string]string{"hoge": "fuga"}, nil)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("異常系_タイムアウト", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[any, TestResponseBody]{
			Request: Request[any]{
				Method: http.MethodGet,
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 1*time.Second)
		defer server.Close()
		expected.Request.URL = server.URL

		_, err := GetRequest[TestResponseBody](server.URL, 100*time.Millisecond, nil, nil)

		require.Error(t, err)
	})

	t.Run("異常系_不正なレスポンス", func(t *testing.T) {
		t.Parallel()

		responses := []struct {
			name string
			body string
		}{
			{name: "HTML", body: "<html></html>"},
			{name: "不正なJSON", body: `{"name":}`},
		}

		for _, res := range responses {
			res := res

			t.Run(res.name, func(t *testing.T) {
				t.Parallel()

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(res.body))
				}))
				defer server.Close()

				_, err := GetRequest[TestResponseBody](server.URL, 1*time.Second, nil, nil)

				require.Error(t, err)
			})
		}
	})
}

func TestPostRequestJSON(t *testing.T) {
	t.Parallel()

	t.Run("正常系_ボディなし", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[any, TestResponseBody]{
			Request: Request[any]{
				Method: http.MethodPost,
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 0)
		defer server.Close()
		expected.Request.URL = server.URL

		actual, err := PostRequestJSON[any, TestResponseBody](server.URL, 1*time.Second, nil, nil)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("正常系_ボディあり", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[TestRequestBody, TestResponseBody]{
			Request: Request[TestRequestBody]{
				Method: http.MethodPost,
				Body:   TestRequestBody{ID: "test_id"},
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 0)
		defer server.Close()
		expected.Request.URL = server.URL

		actual, err := PostRequestJSON[TestRequestBody, TestResponseBody](
			server.URL,
			1*time.Second,
			TestRequestBody{ID: "test_id"},
			nil,
		)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("異常系_タイムアウト", func(t *testing.T) {
		t.Parallel()

		responseBody := TestResponseBody{Name: "test_name"}
		responseBodyByte, _ := json.Marshal(responseBody)
		expected := Response[TestRequestBody, TestResponseBody]{
			Request: Request[TestRequestBody]{
				Method: http.MethodPost,
			},
			StatusCode: http.StatusOK,
			Body:       responseBody,
			BodyByte:   responseBodyByte,
		}

		server := mockServer(expected, 1*time.Second)
		defer server.Close()
		expected.Request.URL = server.URL

		_, err := PostRequestJSON[TestRequestBody, TestResponseBody](server.URL, 100*time.Millisecond, TestRequestBody{}, nil)

		require.Error(t, err)
	})

	t.Run("異常系_不正なレスポンス", func(t *testing.T) {
		t.Parallel()

		responses := []struct {
			name string
			body string
		}{
			{name: "HTML", body: "<html></html>"},
			{name: "不正なJSON", body: `{"name":}`},
		}

		for _, res := range responses {
			res := res

			t.Run(res.name, func(t *testing.T) {
				t.Parallel()

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(res.body))
				}))
				defer server.Close()

				_, err := PostRequestJSON[TestRequestBody, TestResponseBody](server.URL, 1*time.Second, TestRequestBody{}, nil)

				require.Error(t, err)
			})
		}
	})
}
