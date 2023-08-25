package response

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_FieldQuery(t *testing.T) {
	t.Parallel()

	// テストデータ
	type TestData struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Detail struct {
			Hoge string `json:"hoge"`
			Fuga string `json:"fuga"`
		} `json:"detail"`
	}

	// テスト対象のデータ型を取得
	dataType := reflect.TypeOf(TestData{})

	// fieldQuery 関数を実行して結果を取得
	result := fieldQuery(dataType)

	// 期待される結果
	expectedResult := "id,name,detail.hoge,detail.fuga"

	// 結果の比較
	assert.Equal(t, expectedResult, result)
}

func TestUtil_ToSnakeCase(t *testing.T) {
	t.Parallel()

	// テストデータ
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "camelCase", expected: "camel_case"},
		{input: "PascalCase", expected: "pascal_case"},
		{input: "snake_case", expected: "snake_case"},
		{input: "lowercase", expected: "lowercase"},
		{input: "UPPERCASE", expected: "uppercase"},
		{input: "mixed_Case", expected: "mixed_case"},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		// toSnakeCase 関数を実行して結果を取得
		result := toSnakeCase(tc.input)

		// 結果の比較
		assert.Equal(t, tc.expected, result)
	}
}
