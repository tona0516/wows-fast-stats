package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGAccountList_AccountIDs(t *testing.T) {
	t.Parallel()

	w := WGAccountList{}
	w.Data = []WGAccountListData{
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

	w := WGAccountList{}
	w.Data = []WGAccountListData{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 456},
		{NickName: "Bob", AccountID: 789},
	}

	nickname := "Alice"
	expectedID := 456
	actualID := w.AccountID(nickname)

	assert.Equal(t, expectedID, actualID)
}

func TestWGAccountList_AccountID_NotExist(t *testing.T) {
	t.Parallel()

	w := WGAccountList{}
	w.Data = []WGAccountListData{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 456},
		{NickName: "Bob", AccountID: 789},
	}

	nickname := "Unknown"
	expectedID := 0
	actualID := w.AccountID(nickname)

	assert.Equal(t, expectedID, actualID)
}
