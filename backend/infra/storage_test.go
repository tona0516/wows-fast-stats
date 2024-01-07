package infra

import (
	"os"
	"testing"
	"wfs/backend/domain"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var storage *Storage

func TestMain(m *testing.M) {
	// before all
	storagePath := "./unit_test_storage"
	db, _ := badger.Open(badger.DefaultOptions(storagePath))
	storage = NewStorage(db)

	code := m.Run()

	// after all
	os.RemoveAll(storagePath)

	os.Exit(code)
}

func TestStorage_DataVersion(t *testing.T) {
	t.Parallel()

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

	// 取得：保存されていない場合はデフォルト値を返却する
	actual, err := storage.UserConfig()
	require.NoError(t, err)
	assert.Equal(t, domain.DefaultUserConfig, actual)
	assert.False(t, storage.IsExistUserConfig())

	// 書き込み：正常系
	expected := domain.UserConfig{
		FontSize: "large",
		Displays: domain.Displays{
			Ship:    domain.Ship{PR: true},
			Overall: domain.Overall{PR: false},
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

func TestStorage_AlertPlayers(t *testing.T) {
	t.Parallel()

	// 取得：保存されていない場合空のスライスを返却する
	actual, err := storage.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, []domain.AlertPlayer{}, actual)
	assert.False(t, storage.IsExistAlertPlayers())

	// 書き込み：正常系
	expected := []domain.AlertPlayer{
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
	actual, err = storage.AlertPlayers()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestStorage_ExpectedStats(t *testing.T) {
	t.Parallel()

	expected := domain.ExpectedStats{
		1: domain.ExpectedValues{
			AverageDamageDealt: 123,
			AverageFrags:       456,
			WinRate:            789,
		},
		10: domain.ExpectedValues{
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

	expected := "tonango"
	// 書き込み：正常系
	err := storage.WriteOwnIGN(expected)
	require.NoError(t, err)

	// 取得：正常系
	actual, err := storage.OwnIGN()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
