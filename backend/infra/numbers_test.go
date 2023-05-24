package infra

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"changeme/backend/vo"
)

func TestNumbers_ExpectedStats(t *testing.T) {
	t.Parallel()

	// テスト用のレスポンスデータ
	mockResponse := []byte(`{
        "time": 1621699200,
        "data": {
            "1234": {"average_damage_dealt": 50000, "average_frags": 1.2, "win_rate": 52.3},
            "5678": {"average_damage_dealt": 60000, "average_frags": 1.5, "win_rate": 56.8},
            "9012": []
        }
    }`)

	// テスト用の HTTP サーバーを作成し、モックのレスポンスを返すように設定
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockResponse)
	}))
	defer server.Close()

	// テスト対象のインスタンスを作成
	numbers := NewNumbers(server.URL)

	// ExpectedStats メソッドを実行して結果を取得
	expectedStats, err := numbers.ExpectedStats()
	assert.NoError(t, err)

	// 期待される結果
	expectedResult := vo.NSExpectedStats{
		Time: 1621699200,
		Data: map[int]vo.NSExpectedStatsData{
			1234: {
				AverageDamageDealt: 50000,
				AverageFrags:       1.2,
				WinRate:            52.3,
			},
			5678: {
				AverageDamageDealt: 60000,
				AverageFrags:       1.5,
				WinRate:            56.8,
			},
		},
	}

	// 結果の比較
	assert.Equal(t, expectedResult, expectedStats)
}
