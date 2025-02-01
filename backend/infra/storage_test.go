package infra

import (
	"os"
	"path"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
)

const dbPrefix = "storage_test"

func openDB(t *testing.T) *badger.DB {
	t.Helper()

	storagePath := path.Join(dbPrefix, t.Name())
	db, err := badger.Open(badger.DefaultOptions(storagePath))
	assert.NoError(t, err)

	err = db.DropAll()
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.Equal(t, uint(0), actual)

	// 書き込み：正常系
	var expected uint = 10
	err = storage.WriteDataVersion(expected)
	assert.NoError(t, err)

	// 取得：正常系
	actual, err = storage.DataVersion()
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	// 取得：正常系
	actual, err := storage.OwnIGN()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	cleanDB(t, db)
}
