package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/mock/repository"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const validInstallPath = "install_path_test"

//nolint:paralleltest
func TestConfig_UpdateInstallPath(t *testing.T) {
	ctrl := gomock.NewController(t)

	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().UserConfigV2().Return(data.DefaultUserConfigV2(), nil)
		mockStorage.EXPECT().WriteUserConfigV2(gomock.Any()).Return(nil)

		// テスト
		c := NewConfig(nil, nil, mockStorage, nil)
		actual, err := c.UpdateInstallPath(validInstallPath)

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, validInstallPath, actual.InstallPath)
	})

	t.Run("異常系", func(t *testing.T) {
		params := map[string]failure.StringCode{
			"":             apperr.EmptyInstallPath,   // 空文字
			"invalid/path": apperr.InvalidInstallPath, // 配下にWorldOfWarships.exeが存在しないパス
		}

		for path, expected := range params {
			// テスト
			c := NewConfig(nil, nil, nil, nil)
			_, err := c.UpdateInstallPath(path)

			// アサーション
			require.EqualError(t, apperr.Unwrap(err), expected.ErrorCode())
		}
	})
}

//nolint:paralleltest
func TestConfig_UpdateOptional(t *testing.T) {
	ctrl := gomock.NewController(t)

	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		config := createInputConfig()
		config.FontSize = "small"
		// Note: requiredな値を与えてもこれらの値はWriteUserConfigでは含まれない
		actualWritten := data.DefaultUserConfigV2()
		actualWritten.FontSize = "small"

		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().UserConfigV2().Return(data.DefaultUserConfigV2(), nil)
		mockStorage.EXPECT().WriteUserConfigV2(actualWritten).Return(nil)

		// テスト実行
		c := NewConfig(nil, nil, mockStorage, nil)
		err = c.UpdateOptional(config)

		// アサーション
		require.NoError(t, err)
	})
}

func TestConfig_AlertPlayers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		expected := []data.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(expected, nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		actual, err := config.AlertPlayers()

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestConfig_UpdateAlertPlayer(t *testing.T) {
	t.Parallel()

	t.Run("正常系_追加", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		existingPlayers := []data.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		newPlayer := data.AlertPlayer{AccountID: 3, Name: "Player3"}
		expected := make([]data.AlertPlayer, 0)
		expected = append(expected, existingPlayers...)
		expected = append(expected, newPlayer)

		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(existingPlayers, nil)
		mockStorage.EXPECT().WriteAlertPlayers(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		actual, err := config.UpdateAlertPlayer(newPlayer)

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("正常系_更新", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)

		existingPlayers := []data.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		expected := []data.AlertPlayer{
			{AccountID: 1, Name: "UpdatedPlayer"},
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(existingPlayers, nil)
		mockStorage.EXPECT().WriteAlertPlayers(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		actual, err := config.UpdateAlertPlayer(data.AlertPlayer{AccountID: 1, Name: "UpdatedPlayer"})

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestConfig_RemoveAlertPlayer(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_対象IDあり", func(t *testing.T) {
		t.Parallel()

		expected := []data.AlertPlayer{
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return([]data.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}, nil)
		mockStorage.EXPECT().WriteAlertPlayers(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		actual, err := config.RemoveAlertPlayer(1)

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("正常系_対象IDなし", func(t *testing.T) {
		t.Parallel()

		expected := []data.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(expected, nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		actual, err := config.RemoveAlertPlayer(3)

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func createInputConfig() data.UserConfigV2 {
	config := data.DefaultUserConfigV2()
	config.InstallPath = validInstallPath

	return config
}

func createGameClientPath() error {
	if err := os.MkdirAll(validInstallPath, fs.ModePerm); err != nil {
		return err
	}

	gameExePath := filepath.Join(validInstallPath, GameExeName)

	return os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}
