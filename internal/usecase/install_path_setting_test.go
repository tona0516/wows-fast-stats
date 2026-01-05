package usecase

import (
	"context"
	"errors"
	"testing"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestInstallPathSetting_Invoke(t *testing.T) {
	t.Run("正常系_選択したディレクトリのパスが保存される", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		originalConfig := data.UserConfig{
			Version:     1,
			InstallPath: "/old/path",
		}
		updatedConfig := data.UserConfig{
			Version:     1,
			InstallPath: "/new/path",
		}

		gomock.InOrder(
			mockLocalStorage.EXPECT().UserConfig().Return(originalConfig, nil),
			mockLocalStorage.EXPECT().SetUserConfig(updatedConfig).Return(nil),
		)

		openDirCalled := false
		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			openDirCalled = true
			return "/new/path", nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.True(t, openDirCalled)
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

	t.Run("異常系_UserConfigでエラーが発生", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockLocalStorage := mock.NewMockLocalStorage(ctrl)
		expectedErr := errors.New("user config read error")
		mockLocalStorage.EXPECT().UserConfig().Return(data.UserConfig{}, expectedErr)
		mockLocalStorage.EXPECT().SetUserConfig(gomock.Any()).Times(0)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "/new/path", nil
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

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "/new/path", nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		require.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("正常系_複数フィールドを持つ設定が保持される", func(t *testing.T) {
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
		expectedUpdatedConfig := data.UserConfig{
			Version:      1,
			InstallPath:  "/new/path",
			ZoomRate:     100,
			StatsExtra:   "some_value",
			IsSendReport: true,
		}

		gomock.InOrder(
			mockLocalStorage.EXPECT().UserConfig().Return(originalConfig, nil),
			mockLocalStorage.EXPECT().SetUserConfig(expectedUpdatedConfig).Return(nil),
		)

		mockOpenDirectoryDialog := func(ctx context.Context) (string, error) {
			return "/new/path", nil
		}

		instance := NewInstallPathSetting(mockLocalStorage, mockOpenDirectoryDialog)
		ok, err := instance.Invoke(context.Background())

		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
