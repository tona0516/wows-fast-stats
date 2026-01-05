package usecase

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// createWorldOfWarshipsExeDir テンポラリディレクトリを作成し、WorldOfWarships.exeファイルを配置する
func createWorldOfWarshipsExeDir(t *testing.T) string {
	tempDir := t.TempDir()
	exePath := filepath.Join(tempDir, "WorldOfWarships.exe")
	_, err := os.Create(exePath)
	require.NoError(t, err)
	return tempDir
}

func TestInstallPathSetting_Invoke(t *testing.T) {
	t.Run("正常系_InstallPathが更新され、他のフィールドの値は更新されない", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		originalConfig := data.UserConfig{
			Version:      1,
			InstallPath:  "/old/path",
			ZoomRate:     100,
			StatsExtra:   "some_value",
			IsSendReport: true,
		}

		tempDir := createWorldOfWarshipsExeDir(t)

		expectedUpdatedConfig := data.UserConfig{
			Version:      1,
			InstallPath:  tempDir,
			ZoomRate:     100,
			StatsExtra:   "some_value",
			IsSendReport: true,
		}

		gomock.InOrder(
			mockLocalStorage.EXPECT().UserConfig().Return(originalConfig, nil),
			mockLocalStorage.EXPECT().SetUserConfig(expectedUpdatedConfig).Return(nil),
		)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return tempDir, nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("正常系_キャンセルされた場合に保存されずエラーとしない", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		mockLocalStorage.EXPECT().UserConfig().Times(0)
		mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "", nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_openDirectoryDialogでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		mockLocalStorage.EXPECT().UserConfig().Times(0)
		mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "", errors.New("dialog error")
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_WorldOfWarships.exeが見つからない", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		mockLocalStorage.EXPECT().UserConfig().Times(0)
		mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "/invalid/path", nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_UserConfigでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		expectedErr := errors.New("user config read error")
		mockLocalStorage.EXPECT().UserConfig().Return(data.UserConfig{}, expectedErr)
		mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		tempDir := createWorldOfWarshipsExeDir(t)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return tempDir, nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_SetUserConfigでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		originalConfig := data.UserConfig{
			Version:     1,
			InstallPath: "/old/path",
		}

		gomock.InOrder(
			mockLocalStorage.EXPECT().UserConfig().Return(originalConfig, nil),
			mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Return(errors.New("user config write error")),
		)

		tempDir := createWorldOfWarshipsExeDir(t)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return tempDir, nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		require.Error(t, err)
		assert.False(t, ok)
	})
}
