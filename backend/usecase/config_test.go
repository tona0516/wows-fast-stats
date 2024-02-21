package usecase

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock_repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	validInstallPath = "install_path_test"
	validAppID       = "abc123"
)

var errWargaming = failure.New(apperr.WGAPIError)

//nolint:paralleltest
func TestConfig_UpdateRequired(t *testing.T) {
	ctrl := gomock.NewController(t)

	// 準備
	err := createGameClientPath()
	require.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		config := createInputConfig()

		mockWargaming := mock_repository.NewMockWargamingInterface(ctrl)
		mockWargaming.EXPECT().Test(config.Appid).Return(true, nil)

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().UserConfigV2().Return(model.DefaultUserConfigV2, nil)
		mockStorage.EXPECT().WriteUserConfigV2(config).Return(nil)

		// テスト
		c := NewConfig(nil, mockWargaming, mockStorage, nil)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, model.RequiredConfigError{Valid: true}, actual)
		require.NoError(t, err)
	})

	t.Run("異常系_不正なインストールパス", func(t *testing.T) {
		config := model.DefaultUserConfigV2
		config.InstallPath = "invalid/path" // Note: 不正なパス
		config.Appid = "abc123"

		mockWargaming := mock_repository.NewMockWargamingInterface(ctrl)
		mockWargaming.EXPECT().Test(config.Appid).Return(true, nil)

		// テスト
		c := NewConfig(nil, mockWargaming, nil, nil)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, model.RequiredConfigError{InstallPath: apperr.InvalidInstallPath.ErrorCode()}, actual)
		require.NoError(t, err)
	})

	t.Run("異常系_不正なAppID", func(t *testing.T) {
		config := createInputConfig()

		mockWargaming := mock_repository.NewMockWargamingInterface(ctrl)
		mockWargaming.EXPECT().Test(config.Appid).Return(false, errWargaming) // Note: WG APIでエラー

		// テスト
		c := NewConfig(nil, mockWargaming, nil, nil)
		actual, err := c.UpdateRequired(config.InstallPath, config.Appid)

		// アサーション
		assert.Equal(t, model.RequiredConfigError{AppID: apperr.InvalidAppID.ErrorCode()}, actual)
		require.NoError(t, err)
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
		actualWritten := model.DefaultUserConfigV2
		actualWritten.FontSize = "small"

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().UserConfigV2().Return(model.DefaultUserConfigV2, nil)
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
		expected := []model.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
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

	ctrl := gomock.NewController(t)

	existingPlayers := []model.AlertPlayer{
		{AccountID: 1, Name: "Player1"},
		{AccountID: 2, Name: "Player2"},
	}

	t.Run("正常系_追加", func(t *testing.T) {
		t.Parallel()

		// 準備
		newPlayer := model.AlertPlayer{AccountID: 3, Name: "Player3"}

		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(existingPlayers, nil)
		mockStorage.EXPECT().WriteAlertPlayers(append(existingPlayers, newPlayer)).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		err := config.UpdateAlertPlayer(newPlayer)

		// アサーション
		require.NoError(t, err)
	})

	t.Run("正常系_更新", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return(existingPlayers, nil)
		mockStorage.EXPECT().WriteAlertPlayers([]model.AlertPlayer{
			{AccountID: 1, Name: "UpdatedPlayer"},
			{AccountID: 2, Name: "Player2"},
		}).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		err := config.UpdateAlertPlayer(model.AlertPlayer{AccountID: 1, Name: "UpdatedPlayer"})

		// アサーション
		require.NoError(t, err)
	})
}

func TestConfig_RemoveAlertPlayer(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_対象IDあり", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return([]model.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}, nil)
		mockStorage.EXPECT().WriteAlertPlayers([]model.AlertPlayer{
			{AccountID: 2, Name: "Player2"},
		}).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		err := config.RemoveAlertPlayer(1)

		// アサーション
		require.NoError(t, err)
	})

	t.Run("正常系_対象IDなし", func(t *testing.T) {
		t.Parallel()

		// 準備
		mockStorage := mock_repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().AlertPlayers().Return([]model.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}, nil)

		// テスト
		config := NewConfig(nil, nil, mockStorage, nil)
		err := config.RemoveAlertPlayer(3)

		// アサーション
		require.NoError(t, err)
	})
}

func createInputConfig() model.UserConfigV2 {
	config := model.DefaultUserConfigV2
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
