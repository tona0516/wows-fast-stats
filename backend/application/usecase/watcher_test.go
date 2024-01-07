package usecase

import (
	"context"
	"testing"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWatcher_Start(t *testing.T) {
	t.Parallel()

	t.Run("正常系_戦闘開始", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := domain.UserConfig{
			InstallPath: "install_path_test",
			Appid:       "abc123",
			FontSize:    "medium",
		}
		mockStorage := &mocks.StorageInterface{}
		mockLocalFile := &mocks.LocalFileInterface{}

		mockStorage.On("UserConfig").Return(config, nil)
		mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, nil)

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// テスト
		watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, nil, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)

		// アサーション
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Len(t, events, 1)
		assert.Contains(t, events, EventStart)

		mockStorage.AssertExpectations(t)
		mockLocalFile.AssertExpectations(t)
	})
	t.Run("正常系_戦闘終了", func(t *testing.T) {
		t.Parallel()

		ignoreErrs := []failure.StringCode{
			apperr.FileNotExist,
			apperr.ReplayDirNotFoundError,
		}

		for _, ie := range ignoreErrs {
			// 準備
			config := domain.UserConfig{
				InstallPath: "install_path_test",
				Appid:       "abc123",
				FontSize:    "medium",
			}
			mockStorage := &mocks.StorageInterface{}
			mockLocalFile := &mocks.LocalFileInterface{}

			mockStorage.On("UserConfig").Return(config, nil)
			mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, failure.New(ie))

			var events []string
			emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
				events = append(events, eventName)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// テスト
			watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, nil, emitFunc)
			go watcher.Start(ctx, ctx)

			// アサーション
			err := watcher.Prepare()
			require.NoError(t, err)

			events = nil
			time.Sleep(100 * time.Millisecond)
			assert.Contains(t, events, EventEnd)
		}
	})
	t.Run("正常系_キャンセル", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := domain.UserConfig{
			InstallPath: "install_path_test",
			Appid:       "abc123",
			FontSize:    "medium",
		}
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("UserConfig").Return(config, nil)

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}
		ctx, cancel := context.WithCancel(context.Background())

		// テスト
		watcher := NewWatcher(10*time.Millisecond, nil, mockStorage, nil, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)
		cancel()

		// アサーション
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Empty(t, events)

		mockStorage.AssertExpectations(t)
	})
	t.Run("異常系_エラー発生", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := domain.UserConfig{
			InstallPath: "install_path_test",
			Appid:       "abc123",
			FontSize:    "medium",
		}
		mockStorage := &mocks.StorageInterface{}
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLogger := &mocks.LoggerInterface{}

		mockStorage.On("UserConfig").Return(config, nil)
		mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(
			domain.TempArenaInfo{},
			failure.New(apperr.UnexpectedError),
		)
		mockLogger.On("Error", mock.Anything).Return()

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// テスト
		watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, mockLogger, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)

		// アサーション
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Len(t, events, 1)
		assert.Contains(t, events, EventErr)

		mockStorage.AssertExpectations(t)
		mockLocalFile.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
