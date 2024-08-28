package infra

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
	"wfs/backend/data"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dbPrefix = "storage_test"

func openDB(t *testing.T) *badger.DB {
	t.Helper()

	storagePath := path.Join(dbPrefix, t.Name())
	db, err := badger.Open(badger.DefaultOptions(storagePath))
	require.NoError(t, err)

	err = db.DropAll()
	require.NoError(t, err)

	return db
}

func cleanDB(t *testing.T, db *badger.DB) {
	t.Helper()

	_ = db.DropAll()
	_ = db.Close()
	_ = os.RemoveAll(path.Join(dbPrefix, t.Name()))
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

	cleanDB(t, db)
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

	cleanDB(t, db)
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

	cleanDB(t, db)
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

	cleanDB(t, db)
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

	cleanDB(t, db)
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

	cleanDB(t, db)
}

func TestStorage_BattleHistory(t *testing.T) {
	t.Parallel()

	db := openDB(t)
	storage := NewStorage(db)

	expected := data.Battle{
		Meta: data.Meta{
			Unixtime: 1723585938,
			Arena:    "test_arena",
			Type:     "test_type",
			OwnShip:  "test_ownship",
		},
		Teams: []data.Team{
			{
				Players: []data.Player{
					{
						PlayerInfo: data.PlayerInfo{},
						ShipInfo:   data.ShipInfo{},
						PvPSolo:    data.PlayerStats{},
						PvPAll:     data.PlayerStats{},
					},
				},
			},
		},
	}
	expectedKey := fmt.Sprintf(
		"battle_%s_%s_%s_%s",
		time.Unix(expected.Meta.Unixtime, 0).Format(time.DateTime),
		expected.Meta.Type,
		expected.Meta.OwnShip,
		expected.Meta.Arena,
	)

	// 書き込み：正常系
	err := storage.WriteBattleHistory(expected)

	require.NoError(t, err)

	// 取得：正常系
	actual, err := storage.BattleHistory(expectedKey)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)

	// キー取得：正常系
	actualKeys, err := storage.BattleHistoryKeys()

	require.NoError(t, err)
	assert.Len(t, actualKeys, 1)
	assert.Equal(t, expectedKey, actualKeys[0])

	// 削除：正常系
	err = storage.DeleteBattleHistory(expectedKey)

	require.NoError(t, err)
	_, err = storage.BattleHistory(expectedKey)
	require.Error(t, err)
	actualKeys, err = storage.BattleHistoryKeys()
	require.NoError(t, err)
	assert.Empty(t, actualKeys)

	cleanDB(t, db)
}
