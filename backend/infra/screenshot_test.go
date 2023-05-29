package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenshot_Save_正常系(t *testing.T) {
	t.Parallel()

	// テストデータの作成
	rawData := "Hello, world!"
	base64Data := "SGVsbG8sIHdvcmxkIQ=="
	path := "screenshot_test/screenshot.png"

	// テストで生成したディレクトリを削除
	defer os.RemoveAll(filepath.Dir(path))

	// テスト
	screenshot := Screenshot{}
	err := screenshot.Save(path, base64Data)

	// アサーション
	assert.NoError(t, err)
	assert.FileExists(t, path)
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, content, []byte(rawData))
}
