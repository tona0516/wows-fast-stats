package usecase

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"wfs/backend/adapter"
	"wfs/backend/data"
	"wfs/backend/mock"

	"github.com/samber/do/v2"
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
		mockConfigStore := mock.NewMockConfigStore(ctrl)
		mockWails := mock.NewMockWails(ctrl)
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
			mockConfigStore.EXPECT().UserConfig().Return(originalConfig, nil),
			mockConfigStore.EXPECT().SetUserConfig(expectedUpdatedConfig).Return(nil),
		)

		mockWails.EXPECT().OpenDirectoryDialog(gomock.Any()).Return(tempDir, nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
			return mockConfigStore, nil
		})
		do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
			return mockWails, nil
		})
		instance, err := NewInstallPathSetting(injector)
		require.NoError(t, err)
		ok, err := instance.Invoke(context.Background())

		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("異常系_openDirectoryDialogでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockConfigStore(ctrl)
		mockWails := mock.NewMockWails(ctrl)
		mockConfigStore.EXPECT().UserConfig().Times(0)
		mockConfigStore.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockWails.EXPECT().OpenDirectoryDialog(gomock.Any()).Return("", errors.New("dialog error"))

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
			return mockConfigStore, nil
		})
		do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
			return mockWails, nil
		})
		instance, err := NewInstallPathSetting(injector)
		require.NoError(t, err)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_WorldOfWarships.exeが見つからない", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockConfigStore(ctrl)
		mockWails := mock.NewMockWails(ctrl)
		mockConfigStore.EXPECT().UserConfig().Times(0)
		mockConfigStore.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockWails.EXPECT().OpenDirectoryDialog(gomock.Any()).Return("/invalid/path", nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
			return mockConfigStore, nil
		})
		do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
			return mockWails, nil
		})
		instance, err := NewInstallPathSetting(injector)
		require.NoError(t, err)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_UserConfigでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockConfigStore(ctrl)
		mockWails := mock.NewMockWails(ctrl)
		expectedErr := errors.New("user config read error")
		mockConfigStore.EXPECT().UserConfig().Return(data.UserConfig{}, expectedErr)
		mockConfigStore.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		tempDir := createWorldOfWarshipsExeDir(t)

		mockWails.EXPECT().OpenDirectoryDialog(gomock.Any()).Return(tempDir, nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
			return mockConfigStore, nil
		})
		do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
			return mockWails, nil
		})
		instance, err := NewInstallPathSetting(injector)
		require.NoError(t, err)
		ok, err := instance.Invoke(context.Background())

		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("異常系_SetUserConfigでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockConfigStore(ctrl)
		mockWails := mock.NewMockWails(ctrl)
		originalConfig := data.UserConfig{
			Version:     1,
			InstallPath: "/old/path",
		}

		gomock.InOrder(
			mockConfigStore.EXPECT().UserConfig().Return(originalConfig, nil),
			mockConfigStore.EXPECT().SetUserConfig(gomock.Any()).Return(errors.New("user config write error")),
		)

		tempDir := createWorldOfWarshipsExeDir(t)

		mockWails.EXPECT().OpenDirectoryDialog(gomock.Any()).Return(tempDir, nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.ConfigStore, error) {
			return mockConfigStore, nil
		})
		do.Provide(injector, func(i do.Injector) (adapter.Wails, error) {
			return mockWails, nil
		})
		instance, err := NewInstallPathSetting(injector)
		require.NoError(t, err)
		ok, err := instance.Invoke(context.Background())

		require.Error(t, err)
		assert.False(t, ok)
	})
}
