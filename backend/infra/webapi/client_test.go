package webapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"wfs/backend/application/vo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のデータ型.
type TestData struct {
	Name string `json:"name"`
}

func TestGetRequest_正常系_クエリなし(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write(expected.ByteBody)
	}))
	defer server.Close()
	expected.FullURL = server.URL

	actual, err := GetRequest[TestData](server.URL, 1*time.Second)

	assert.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestGetRequest_正常系_クエリあり(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write(expected.ByteBody)
	}))
	defer server.Close()
	expected.FullURL = server.URL + "?hoge=fuga"

	actual, err := GetRequest[TestData](server.URL, 1*time.Second, vo.NewPair("hoge", "fuga"))

	assert.Equal(t, expected, actual)
	require.NoError(t, err)
}

func TestGetRequest_異常系_タイムアウト(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(normalResponse().ByteBody)
	}))
	defer server.Close()

	_, err := GetRequest[TestData](server.URL, 100*time.Millisecond)

	require.Error(t, err)
}

func TestGetRequest_異常系_不正なレスポンス(t *testing.T) {
	t.Parallel()

	responses := []struct {
		body string
	}{
		{body: "<html></html>"},
		{body: `{"name":}`},
	}

	for _, res := range responses {
		res := res

		t.Run("", func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(res.body))
			}))
			defer server.Close()

			_, err := GetRequest[TestData](server.URL, 1*time.Second)

			require.Error(t, err)
		})
	}
}

func TestPostMultipartFormData_正常系(t *testing.T) {
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
