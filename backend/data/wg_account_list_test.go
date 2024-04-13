package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGAccountList_AccountIDs(t *testing.T) {
	t.Parallel()

	w := WGAccountList{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 0},
		{NickName: "Bob", AccountID: 456},
		{NickName: "Charlie", AccountID: 789},
		{NickName: "John", AccountID: 123},
	}

	expectedIDs := []int{123, 456, 789}
	actualIDs := w.AccountIDs()

	assert.ElementsMatch(t, expectedIDs, actualIDs)
}

func TestWGAccountList_AccountID(t *testing.T) {
	t.Parallel()

	w := WGAccountList{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 456},
		{NickName: "Bob", AccountID: 789},
	}

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		expectedID := 456
		actualID := w.AccountID("Alice")

		assert.Equal(t, expectedID, actualID)
	})

	t.Run("異常系_存在しない場合は0を返す", func(t *testing.T) {
		t.Parallel()

		expectedID := 0
		actualID := w.AccountID("Unknown")

		assert.Equal(t, expectedID, actualID)
	})
}
