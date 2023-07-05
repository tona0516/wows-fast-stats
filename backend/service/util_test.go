package service

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUtil_makeRange(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []int{1, 2, 3, 4}, makeRange(1, 5))
	assert.Equal(t, []int{-5, -4, -3, -2, -1}, makeRange(-5, 0))
	assert.Equal(t, []int{}, makeRange(0, 0))
	assert.Equal(t, []int{}, makeRange(0, -1))
}

func TestUtil_doParallel_正常系(t *testing.T) {
	t.Parallel()

	// テストデータ
	values := []int{1, 2, 3, 4, 5}

	// テスト用の関数
	fn := func(value int) error {
		return nil
	}

	// テスト実行
	err := doParallel(2, values, fn)

	// 結果の検証
	assert.Nil(t, err)
}

func TestUtil_doParallel_異常系(t *testing.T) {
	t.Parallel()

	// テストデータ
	values := []int{1, 2, 3, 4, 5}

	// テスト用の関数
	fn := func(value int) error {
		if value == 3 {
			return errors.New("error occurred")
		}
		return nil
	}

	// テスト実行
	err := doParallel(2, values, fn)

	// 結果の検証
	assert.EqualError(t, err, "error occurred")
}
