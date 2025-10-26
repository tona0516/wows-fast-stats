package util

import (
	"reflect"
	"testing"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	result := FieldQuery(dataType)

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
		result := ToSnakeCase(tc.input)

		// 結果の比較
		assert.Equal(t, tc.expected, result)
	}
}

func TestUtil_MakeRange(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []int{1, 2, 3, 4}, MakeRange(1, 5))
	assert.Equal(t, []int{-5, -4, -3, -2, -1}, MakeRange(-5, 0))
	assert.Equal(t, []int{}, MakeRange(0, 0))
	assert.Equal(t, []int{}, MakeRange(0, -1))
}

func TestUtil_DoParallel(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		values := MakeRange(1, 5)

		var calls int
		err := DoParallel(values, func(value int) error {
			calls++
			return nil
		})

		require.NoError(t, err)
		assert.Len(t, values, calls)
	})
	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		values := MakeRange(1, 5)

		expected := apperr.HTTPRequestError
		err := DoParallel(values, func(value int) error {
			if value == values[len(values)-1] {
				return failure.New(expected)
			}
			return nil
		})

		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, expected, code)
	})
}
