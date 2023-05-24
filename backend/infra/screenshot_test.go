package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenshot_Save_Success(t *testing.T) {
	t.Parallel()

	// テストデータの作成
	rawData := "Hello, world!"
	base64Data := "SGVsbG8sIHdvcmxkIQ=="
	path := "screenshot_test/screenshot.png"

	// Saveメソッドの呼び出し
	screenshot := Screenshot{}
	err := screenshot.Save(path, base64Data)
	defer os.RemoveAll(filepath.Dir(path))

	// エラーの確認
	assert.NoError(t, err)
	assert.FileExists(t, path)

	// ファイルの内容の確認
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, content, []byte(rawData))

	// ファイルの削除
	os.Remove("screenshot.png")
}
