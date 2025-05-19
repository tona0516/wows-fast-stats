package infra

import (
	"testing"
	"wfs/backend/data"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openDB(t *testing.T) *badger.DB {
	t.Helper()

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(t, err)

	return db
}

func TestStorage_DataVersion(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	// 取得：保存されていない場合0を返却する
	actual, err := storage.DataVersion()
	require.NoError(t, err)
	assert.Equal(t, uint(0), actual)

	// 書き込み：正常系
	var expected uint = 10
	err = storage.WriteDataVersion(expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err = storage.DataVersion()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_UserConfig(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	err := delete(storage.db, userConfigKey)
	require.NoError(t, err)

	// 取得：保存されていない場合はデフォルト値を返却する
	actual, err := storage.UserConfig()
	require.NoError(t, err)
	assert.Equal(t, data.DefaultUserConfig(), actual)
	assert.False(t, storage.IsExistUserConfig())

	// 書き込み：正常系
	expected := data.UserConfig{
		FontSize: "large",
		Displays: data.Displays{
			Ship:    data.Ship{PR: true},
			Overall: data.Overall{PR: false},
		},
	}
	err = storage.WriteUserConfig(expected)
	require.NoError(t, err)
	assert.True(t, storage.IsExistUserConfig())

	// 取得：正常系
	actual, err = storage.UserConfig()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_UserConfigV2(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	err := delete(storage.db, userConfigKey)
	require.NoError(t, err)

	// 取得：保存されていない場合はデフォルト値を返却する
	actual, err := storage.UserConfigV2()
	require.NoError(t, err)
	assert.Equal(t, data.DefaultUserConfigV2(), actual)

	// 書き込み：正常系
	expected := data.UserConfigV2{
		FontSize: "large",
		Display: data.UCDisplay{
			Ship:    data.UCDisplayShip{PR: true},
			Overall: data.UCDisplayOverall{PR: false},
		},
	}
	err = storage.WriteUserConfigV2(expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err = storage.UserConfigV2()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_AlertPlayers(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	assertEmpty := func() {
		actual, err := storage.AlertPlayers()
		require.NoError(t, err)
		assert.Equal(t, []data.AlertPlayer{}, actual)
		assert.False(t, storage.IsExistAlertPlayers())
	}

	// 取得：保存されていない場合空のスライスを返却する
	assertEmpty()

	// 書き込み：空配列を書き込もうとした場合はキーごと削除される
	err := storage.WriteAlertPlayers([]data.AlertPlayer{})
	require.NoError(t, err)
	assertEmpty()

	// 書き込み：正常系
	expected := []data.AlertPlayer{
		{
			AccountID: 100,
			Name:      "test",
			Pattern:   "bi-check-circle-fill",
			Message:   "hello",
		},
		{
			AccountID: 200,
			Name:      "hoge",
			Pattern:   "bi-check-circle-fill",
			Message:   "memo",
		},
	}
	err = storage.WriteAlertPlayers(expected)
	require.NoError(t, err)
	assert.True(t, storage.IsExistAlertPlayers())

	// 取得：正常系
	actual, err := storage.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_ExpectedStats(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	expected := data.ExpectedStats{
		1: data.ExpectedValues{
			AverageDamageDealt: 123,
			AverageFrags:       456,
			WinRate:            789,
		},
		10: data.ExpectedValues{
			AverageDamageDealt: 1,
			AverageFrags:       2,
			WinRate:            3,
		},
	}

	// 書き込み：正常系
	err := storage.WriteExpectedStats(expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := storage.ExpectedStats()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_OwnIGN(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	expected := "tonango"
	// 書き込み：正常系
	err := storage.WriteOwnIGN(expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := storage.OwnIGN()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
