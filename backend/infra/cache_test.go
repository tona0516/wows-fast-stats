package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func TestCache_Serialize_Deserialize(t *testing.T) {
	// テスト用のデータ
	data := "test data"

	// シリアライズとデシリアライズを行うためのキャッシュインスタンスを作成
	cache := NewCache[string]("test", "cache")
	defer os.RemoveAll(cache.Dir)

	// シリアライズしてキャッシュに保存
	err := cache.Serialize(data)
	assert.NoError(t, err)

	// デシリアライズしてデータを取得
	loadedData, err := cache.Deserialize()
	assert.NoError(t, err)
	assert.Equal(t, data, loadedData)
}

//nolint:paralleltest
func TestCache_Deserialize_FileNotFound(t *testing.T) {
	// 存在しないキャッシュファイルを指定してデシリアライズを行うためのキャッシュインスタンスを作成
	cache := NewCache[string]("nonexistent", "cache")
	defer os.RemoveAll(cache.Dir)

	// デシリアライズを実行してエラーを確認
	_, err := cache.Deserialize()
	assert.Error(t, err)
	assert.EqualError(t, err, "I301 Deserialize open cache/nonexistent.bin: no such file or directory")
}

//nolint:paralleltest
func TestCache_Deserialize_DecodeError(t *testing.T) {
	// キャッシュファイルのデコードエラーをシミュレートするためのキャッシュインスタンスを作成
	cache := NewCache[string]("invalid", "cache")
	defer os.RemoveAll(cache.Dir)

	// シリアライズされた無効なデータを書き込む
	err := os.Mkdir(cache.Dir, 0o755)
	assert.NoError(t, err)
	file, err := os.Create(filepath.Join(cache.Dir, cache.Name+".bin"))
	assert.NoError(t, err)
	defer file.Close()
	_, err = file.Write([]byte("invalid data"))
	assert.NoError(t, err)

	// デシリアライズを実行してエラーを確認
	_, err = cache.Deserialize()
	assert.Error(t, err)
	assert.EqualError(t, err, "I301 Deserialize unexpected EOF")
}
