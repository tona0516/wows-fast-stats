package usecase

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	validInstallPath = "install_path_test"
	validAppID       = "abc123"
)

var errWargaming = failure.New(apperr.WGAPIError)

//nolint:paralleltest
func TestConfig_UpdateRequired(t *testing.T) {
	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		config := createInputConfig()

		mockWargaming := &mocks.WargamingInterface{}
		mockStorage := &mocks.StorageInterface{}
		mockWargaming.On("Test", config.Appid).Return(true, nil)
		mockStorage.On("ReadUserConfig").Return(domain.DefaultUserConfig, nil)
		mockStorage.On("WriteUserConfig", config).Return(nil)

		// テスト
		c := NewConfig(nil, mockWargaming, mockStorage)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, vo.RequiredConfigError{Valid: true}, actual)
		require.NoError(t, err)
		mockWargaming.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("異常系_不正なインストールパス", func(t *testing.T) {
		config := domain.DefaultUserConfig
		config.InstallPath = "invalid/path" // Note: 不正なパス
		config.Appid = "abc123"

		mockWargaming := &mocks.WargamingInterface{}
		mockWargaming.On("Test", config.Appid).Return(true, nil)

		// テスト
		c := NewConfig(nil, mockWargaming, nil)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, vo.RequiredConfigError{InstallPath: apperr.InvalidInstallPath.ErrorCode()}, actual)
		require.NoError(t, err)
		mockWargaming.AssertExpectations(t)
	})

	t.Run("異常系_不正なAppID", func(t *testing.T) {
		config := createInputConfig()

		mockWargaming := &mocks.WargamingInterface{}
		mockWargaming.On("Test", config.Appid).Return(false, errWargaming) // Note: WG APIでエラー

		// テスト
		c := NewConfig(nil, mockWargaming, nil)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, vo.RequiredConfigError{AppID: apperr.InvalidAppID.ErrorCode()}, actual)
		require.NoError(t, err)
		mockWargaming.AssertExpectations(t)
	})
}

//nolint:paralleltest
func TestConfig_UpdateOptional(t *testing.T) {
	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		config := createInputConfig()
		config.FontSize = "small"
		// Note: requiredな値を与えてもこれらの値はWriteUserConfigでは含まれない
		actualWritten := domain.DefaultUserConfig
		actualWritten.FontSize = "small"

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadUserConfig").Return(domain.DefaultUserConfig, nil)
		mockStorage.On("WriteUserConfig", actualWritten).Return(nil)

		// テスト実行
		c := NewConfig(nil, nil, mockStorage)
		err = c.UpdateOptional(config)

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestConfig_AlertPlayers(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		expected := []domain.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadAlertPlayers").Return(expected, nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage)
		actual, err := config.AlertPlayers()

		// アサーション
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
		mockStorage.AssertExpectations(t)
	})
}

func TestConfig_UpdateAlertPlayer(t *testing.T) {
	t.Parallel()

	existingPlayers := []domain.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}

	t.Run("正常系_追加", func(t *testing.T) {
		t.Parallel()

		// 準備
		newPlayer := domain.AlertPlayer{AccountID: 3, Name: "Player3"}

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
		mockStorage.On("WriteAlertPlayers", append(existingPlayers, newPlayer)).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage)
		err := config.UpdateAlertPlayer(newPlayer)

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("正常系_更新", func(t *testing.T) {
		t.Parallel()

		// 準備
		playerToUpdate := domain.AlertPlayer{AccountID: 1, Name: "UpdatedPlayer"}

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
		mockStorage.On("WriteAlertPlayers", []domain.AlertPlayer{
			{AccountID: 1, Name: "UpdatedPlayer"},
			{AccountID: 2, Name: "Player2"},
		}).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage)
		err := config.UpdateAlertPlayer(playerToUpdate)

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestConfig_RemoveAlertPlayer(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		accountIDToRemove := 1
		existingPlayers := []domain.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadAlertPlayers").Return(existingPlayers, nil)
		mockStorage.On("WriteAlertPlayers", mock.Anything).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage)
		err := config.RemoveAlertPlayer(accountIDToRemove)

		// アサーション
		require.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})
}

func createInputConfig() domain.UserConfig {
	config := domain.DefaultUserConfig
	config.InstallPath = validInstallPath
	config.Appid = validAppID

	return config
}

func createGameClientPath() error {
	if err := os.MkdirAll(validInstallPath, fs.ModePerm); err != nil {
		return err
	}

	gameExePath := filepath.Join(validInstallPath, GameExeName)

	return os.WriteFile(gameExePath, []byte{}, fs.ModePerm)
}
