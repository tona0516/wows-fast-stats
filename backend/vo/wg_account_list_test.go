package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWGAccountList_AccountIDs(t *testing.T) {
	t.Parallel()
	data := []WGAccountListData{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 0},
		{NickName: "Bob", AccountID: 456},
		{NickName: "Charlie", AccountID: 789},
		{NickName: "John", AccountID: 123},
	}

	w := WGAccountList{
		Data: data,
	}

	expectedIDs := []int{123, 456, 789}
	actualIDs := w.AccountIDs()

	assert.ElementsMatch(t, expectedIDs, actualIDs)
}

func TestWGAccountList_AccountID(t *testing.T) {
	t.Parallel()
	data := []WGAccountListData{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 456},
		{NickName: "Bob", AccountID: 789},
	}

	w := WGAccountList{
		Data: data,
	}

	nickname := "Alice"
	expectedID := 456
	actualID := w.AccountID(nickname)

	assert.Equal(t, expectedID, actualID)
}

func TestWGAccountList_AccountID_NotExist(t *testing.T) {
	t.Parallel()
	data := []WGAccountListData{
		{NickName: "John", AccountID: 123},
		{NickName: "Alice", AccountID: 456},
		{NickName: "Bob", AccountID: 789},
	}

	w := WGAccountList{
		Data: data,
	}

	nickname := "Unknown"
	expectedID := 0
	actualID := w.AccountID(nickname)

	assert.Equal(t, expectedID, actualID)
}
