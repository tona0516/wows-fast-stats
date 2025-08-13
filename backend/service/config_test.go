package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/infra"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validInstallPath = "install_path_test"

//nolint:paralleltest
func TestConfig_UpdateInstallPath(t *testing.T) {
	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		fileStore := infra.NewFileStore(t.TempDir())

		// テスト
		c := NewConfig(nil, nil, fileStore, nil)
		err := c.UpdateInstallPath(validInstallPath)

		// アサーション
		assert.NoError(t, err)
		actual, err := fileStore.Get(data.InstallPathKey)
		assert.NoError(t, err)
		assert.Equal(t, validInstallPath, actual)
	})

	t.Run("異常系", func(t *testing.T) {
		params := map[string]failure.StringCode{
			"":             apperr.EmptyInstallPath,   // 空文字
			"invalid/path": apperr.InvalidInstallPath, // 配下にWorldOfWarships.exeが存在しないパス
		}

		for path, expected := range params {
			// テスト
			c := NewConfig(nil, nil, nil, nil)
			err := c.UpdateInstallPath(path)

			// アサーション
			assert.EqualError(t, apperr.Unwrap(err), expected.ErrorCode())
		}
	})
}

func TestConfig_AlertPlayers(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
	})
}

func TestConfig_UpdateAlertPlayer(t *testing.T) {
	t.Parallel()

	t.Run("正常系_追加", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("正常系_更新", func(t *testing.T) {
		t.Parallel()
	})
}

func TestConfig_RemoveAlertPlayer(t *testing.T) {
	t.Parallel()

	t.Run("正常系_対象IDあり", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("正常系_対象IDなし", func(t *testing.T) {
		t.Parallel()
	})
}

func createGameClientPath() error {
	if err := os.MkdirAll(validInstallPath, fs.ModePerm); err != nil {
		return err
	}

	gameExePath := filepath.Join(validInstallPath, GameExeName)

	return os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}
