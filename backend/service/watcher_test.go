package service

import (
	"context"
	"testing"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock/repository"
	"wfs/backend/domain/model"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestWatcher_Start(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系_戦闘開始", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := model.UserConfigV2{
			InstallPath: "install_path_test",
			FontSize:    "medium",
		}
		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockUserConfig.EXPECT().GetV2().Return(config, nil)

		mockLocalFile := repository.NewMockLocalFile(ctrl)
		mockLocalFile.EXPECT().ReadTempArenaInfo(config.InstallPath).Return(model.TempArenaInfo{}, nil).AnyTimes()

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// テスト
		watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockUserConfig, nil, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)

		// アサーション
		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Len(t, events, 1)
		assert.Contains(t, events, EventStart)
	})

	t.Run("正常系_戦闘終了", func(t *testing.T) {
		t.Parallel()

		ignoreErrs := []failure.StringCode{
			apperr.FileNotExist,
			apperr.ReplayDirNotFoundError,
		}

		for _, ie := range ignoreErrs {
			// 準備
			config := model.UserConfigV2{
				InstallPath: "install_path_test",
				FontSize:    "medium",
			}

			mockUserConfig := repository.NewMockUserConfigStore(ctrl)
			mockUserConfig.EXPECT().GetV2().Return(config, nil)

			mockLocalFile := repository.NewMockLocalFile(ctrl)
			mockLocalFile.EXPECT().ReadTempArenaInfo(config.InstallPath).Return(
				model.TempArenaInfo{},
				failure.New(ie),
			).AnyTimes()

			var events []string
			emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
				events = append(events, eventName)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// テスト
			watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockUserConfig, nil, emitFunc)
			go watcher.Start(ctx, ctx)

			// アサーション
			err := watcher.Prepare()
			assert.NoError(t, err)

			events = nil
			time.Sleep(100 * time.Millisecond)
			assert.Contains(t, events, EventEnd)
		}
	})

	t.Run("正常系_キャンセル", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := model.UserConfigV2{
			InstallPath: "install_path_test",
			FontSize:    "medium",
		}

		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockUserConfig.EXPECT().GetV2().Return(config, nil)

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}
		ctx, cancel := context.WithCancel(context.Background())

		// テスト
		watcher := NewWatcher(10*time.Millisecond, nil, mockUserConfig, nil, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)
		cancel()

		// アサーション
		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Empty(t, events)
	})

	t.Run("異常系_エラー発生", func(t *testing.T) {
		t.Parallel()

		// 準備
		config := model.UserConfigV2{
			InstallPath: "install_path_test",
			FontSize:    "medium",
		}

		mockUserConfig := repository.NewMockUserConfigStore(ctrl)
		mockUserConfig.EXPECT().GetV2().Return(config, nil)

		mockLocalFile := repository.NewMockLocalFile(ctrl)
		mockLocalFile.EXPECT().ReadTempArenaInfo(config.InstallPath).Return(
			model.TempArenaInfo{},
			failure.New(apperr.UnexpectedError),
		).AnyTimes()

		mockLogger := repository.NewMockLogger(ctrl)
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any())

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// テスト
		watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockUserConfig, mockLogger, emitFunc)
		err := watcher.Prepare()
		go watcher.Start(ctx, ctx)

		// アサーション
		assert.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		assert.Len(t, events, 1)
		assert.Contains(t, events, EventErr)
	})
}
