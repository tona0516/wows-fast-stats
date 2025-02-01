package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock/repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

const validInstallPath = "install_path_test"

//nolint:paralleltest
func TestConfig_UpdateInstallPath(t *testing.T) {
	ctrl := gomock.NewController(t)

	// 準備
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockUserConfig.EXPECT().GetV2().Return(model.DefaultUserConfigV2(), nil)
		mockUserConfig.EXPECT().SaveV2(gomock.Any()).Return(nil)

		// テスト
		c := NewConfig(nil, mockUserConfig, nil)
		actual, err := c.UpdateInstallPath(validInstallPath)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, validInstallPath, actual.InstallPath)
	})

	t.Run("異常系", func(t *testing.T) {
		params := map[string]failure.StringCode{
			"":             apperr.EmptyInstallPath,   // 空文字
			"invalid/path": apperr.InvalidInstallPath, // 配下にWorldOfWarships.exeが存在しないパス
		}

		for path, expected := range params {
			// テスト
			c := NewConfig(nil, nil, nil)
			_, err := c.UpdateInstallPath(path)

			// アサーション
			assert.EqualError(t, apperr.Unwrap(err), expected.ErrorCode())
		}
	})
}

//nolint:paralleltest
func TestConfig_UpdateOptional(t *testing.T) {
	ctrl := gomock.NewController(t)

	// 準備
	err := createGameClientPath()
	assert.NoError(t, err)
	defer os.RemoveAll(validInstallPath)

	t.Run("正常系", func(t *testing.T) {
		config := createInputConfig()
		config.FontSize = "small"
		// Note: requiredな値を与えてもこれらの値はWriteUserConfigでは含まれない
		actualWritten := model.DefaultUserConfigV2()
		actualWritten.FontSize = "small"

		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockUserConfig.EXPECT().GetV2().Return(model.DefaultUserConfigV2(), nil)
		mockUserConfig.EXPECT().SaveV2(actualWritten).Return(nil)

		// テスト実行
		c := NewConfig(nil, mockUserConfig, nil)
		err = c.UpdateOptional(config)

		// アサーション
		assert.NoError(t, err)
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

		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockAlertPlayer.EXPECT().GetV1().Return(expected, nil)

		// テスト
		config := NewConfig(nil, nil, mockAlertPlayer)
		actual, err := config.AlertPlayers()

		// アサーション
		assert.NoError(t, err)
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
		expected := make([]model.AlertPlayer, 0)
		expected = append(expected, existingPlayers...)
		expected = append(expected, newPlayer)

		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockAlertPlayer.EXPECT().GetV1().Return(existingPlayers, nil)
		mockAlertPlayer.EXPECT().SaveV1(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockAlertPlayer)
		actual, err := config.UpdateAlertPlayer(newPlayer)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("正常系_更新", func(t *testing.T) {
		t.Parallel()

		expected := []model.AlertPlayer{
			{AccountID: 1, Name: "UpdatedPlayer"},
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockAlertPlayer.EXPECT().GetV1().Return(existingPlayers, nil)
		mockAlertPlayer.EXPECT().SaveV1(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockAlertPlayer)
		actual, err := config.UpdateAlertPlayer(model.AlertPlayer{AccountID: 1, Name: "UpdatedPlayer"})

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestConfig_RemoveAlertPlayer(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_対象IDあり", func(t *testing.T) {
		t.Parallel()

		expected := []model.AlertPlayer{
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockAlertPlayer.EXPECT().GetV1().Return([]model.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}, nil)
		mockAlertPlayer.EXPECT().SaveV1(expected).Return(nil)

		// テスト
		config := NewConfig(nil, nil, mockAlertPlayer)
		actual, err := config.RemoveAlertPlayer(1)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("正常系_対象IDなし", func(t *testing.T) {
		t.Parallel()

		expected := []model.AlertPlayer{
			{AccountID: 1, Name: "Player1"},
			{AccountID: 2, Name: "Player2"},
		}

		// 準備
		mockAlertPlayer := repository.NewMockAlertPlayerStore(ctrl)
		mockAlertPlayer.EXPECT().GetV1().Return(expected, nil)

		// テスト
		config := NewConfig(nil, nil, mockAlertPlayer)
		actual, err := config.RemoveAlertPlayer(3)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func createInputConfig() model.UserConfigV2 {
	config := model.DefaultUserConfigV2()
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
