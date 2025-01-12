//nolint:paralleltest
package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testInstallPath = "testdata"

func TestLocalFile_SaveScreenshot(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// テストデータの作成
		rawData := "Hello, world!"
		base64Data := "SGVsbG8sIHdvcmxkIQ=="
		path := "screenshot_test/screenshot.png"

		// テストで生成したディレクトリを削除
		defer os.RemoveAll(filepath.Dir(path))

		// テスト
		instance := LocalFile{}
		err := instance.SaveScreenshot(path, base64Data)

		// アサーション
		assert.NoError(t, err)
		assert.FileExists(t, path)
		content, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, content, []byte(rawData))
	})
}
