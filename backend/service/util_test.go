package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_DoParallel_Success(t *testing.T) {
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

func TestUtil_DoParallel_Failure(t *testing.T) {
	t.Parallel()

	// テストデータ
	values := []int{1, 2, 3, 4, 5}

	// テスト用の関数
	fn := func(value int) error {
		if value == 3 {
			//nolint:goerr113
			return errors.New("error occurred")
		}
		return nil
	}

	// テスト実行
	err := doParallel(2, values, fn)

	// 結果の検証
	assert.EqualError(t, err, "error occurred")
}
