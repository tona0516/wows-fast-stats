package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs/backend/application/vo"

	"github.com/stretchr/testify/assert"
)

// テスト用のデータ型.
type TestData struct {
	Name string `json:"name"`
}

func TestUtil_getRequest_正常系_クエリなし(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write(expected.BodyByte)
	}))
	defer server.Close()

	actual, err := getRequest[TestData](server.URL)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestUtil_getRequest_正常系_クエリあり(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write(expected.BodyByte)
	}))
	defer server.Close()

	actual, err := getRequest[TestData](server.URL, vo.NewPair("hoge", "fuga"))

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestUtil_getRequest_異常系_不正なレスポンス(t *testing.T) {
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

			expected := APIResponse[TestData]{
				StatusCode: http.StatusOK,
				Body:       TestData{},
				BodyByte:   []byte(res.body),
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(expected.StatusCode)
				_, _ = w.Write(expected.BodyByte)
			}))
			defer server.Close()

			actual, err := getRequest[TestData](server.URL)

			assert.Error(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func normalResponse() APIResponse[TestData] {
	testData := TestData{Name: "test"}

	//nolint:errchkjson
	body, _ := json.Marshal(testData)
	return APIResponse[TestData]{
		StatusCode: http.StatusOK,
		Body:       testData,
		BodyByte:   body,
	}
}
