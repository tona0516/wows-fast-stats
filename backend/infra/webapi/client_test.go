package webapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のデータ型.
type TestData struct {
	Name string `json:"name"`
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	t.Run("正常系_クエリなし", func(t *testing.T) {
		t.Parallel()

		expected := normalResponse()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(expected.StatusCode)
			_, _ = w.Write(expected.ByteBody)
		}))
		defer server.Close()
		expected.FullURL = server.URL

		actual, err := GetRequest[TestData](server.URL, 1*time.Second, nil)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("正常系_クエリあり", func(t *testing.T) {
		t.Parallel()

		expected := normalResponse()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(expected.StatusCode)
			_, _ = w.Write(expected.ByteBody)
		}))
		defer server.Close()
		expected.FullURL = server.URL + "?hoge=fuga"

		actual, err := GetRequest[TestData](server.URL, 1*time.Second, map[string]string{"hoge": "fuga"})

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("異常系_タイムアウト", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(normalResponse().ByteBody)
		}))
		defer server.Close()

		_, err := GetRequest[TestData](server.URL, 100*time.Millisecond, nil)

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
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(res.body))
				}))
				defer server.Close()

				_, err := GetRequest[TestData](server.URL, 1*time.Second, nil)

				require.Error(t, err)
			})
		}
	})
}

func TestPostMultipartFormData(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expected := normalResponse()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(expected.StatusCode)
			_, _ = w.Write(expected.ByteBody)
		}))
		defer server.Close()
		expected.FullURL = server.URL

		filename := "testfile.txt"
		_, err := os.Create(filename)
		defer os.Remove(filename)
		require.NoError(t, err)

		actual, err := PostMultipartFormData[TestData](
			server.URL,
			1*time.Second,
			[]Form{
				NewForm("content_name", "content_value", false),
				NewForm("file", filename, true),
			},
		)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})
}

func normalResponse() Response[TestData] {
	testData := TestData{Name: "test"}
	//nolint:errchkjson
	byteBody, _ := json.Marshal(testData)

	return Response[TestData]{
		StatusCode: http.StatusOK,
		Body:       testData,
		ByteBody:   byteBody,
	}
}
