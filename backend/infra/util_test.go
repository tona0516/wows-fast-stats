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

func TestUtil_getRequest_正常系_クエリなし(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write([]byte(expected.BodyString))
	}))
	defer server.Close()

	actual, err := getRequest[TestData](server.URL, map[string]string{}, 0)

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestUtil_getRequest_正常系_最終リトライで成功(t *testing.T) {
	t.Parallel()

	expected := normalResponse()

	retry := 1
	maxCalls := retry + 1
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++

		if calls < maxCalls {
			//nolint:forcetypeassert
			c, _, _ := w.(http.Hijacker).Hijack()
			defer c.Close()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(expected.StatusCode)
		_, _ = w.Write([]byte(expected.BodyString))
	}))
	defer server.Close()

	actual, err := getRequest[TestData](server.URL, map[string]string{}, uint64(retry))

	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
	assert.Equal(t, maxCalls, calls)
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
				BodyString: res.body,
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(expected.StatusCode)
				_, _ = w.Write([]byte(expected.BodyString))
			}))
			defer server.Close()

			actual, err := getRequest[TestData](server.URL, map[string]string{}, 0)

			assert.Error(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestUtil_getRequest_異常系_すべてのリトライが失敗(t *testing.T) {
	t.Parallel()

	retry := 1
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint:forcetypeassert
		c, _, _ := w.(http.Hijacker).Hijack()
		defer c.Close()
		calls++
	}))
	defer server.Close()

	_, err := getRequest[TestData](server.URL, map[string]string{}, uint64(retry))

	assert.Error(t, err)
	assert.Equal(t, retry+1, calls)
}

func normalResponse() APIResponse[TestData] {
	testData := TestData{Name: "test"}

	//nolint:errchkjson
	body, _ := json.Marshal(testData)
	return APIResponse[TestData]{
		StatusCode: http.StatusOK,
		Body:       testData,
		BodyString: string(body),
	}
}
