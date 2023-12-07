package service

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

func TestWatcher_Start_戦闘開始(t *testing.T) {
	t.Parallel()

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, nil)

	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)
	go watcher.Start(ctx, ctx)

	time.Sleep(100 * time.Millisecond)
	assert.Len(t, events, 1)
	assert.Contains(t, events, EventStart)
}

func TestWatcher_Start_戦闘終了(t *testing.T) {
	t.Parallel()

	ignoreErrs := []failure.StringCode{
		apperr.FileNotExist,
		apperr.ReplayDirNotFoundError,
	}

	for _, ie := range ignoreErrs {
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, failure.New(ie))
		mockStorage := &mocks.StorageInterface{}
		mockStorage.On("ReadUserConfig").Return(domain.UserConfig{
			InstallPath: "install_path_test",
			Appid:       "abc123",
			FontSize:    "medium",
		}, nil)

		var events []string
		emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
			events = append(events, eventName)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, emitFunc)
		err := watcher.Prepare()
		require.NoError(t, err)

		go watcher.Start(ctx, ctx)

		events = nil
		time.Sleep(100 * time.Millisecond)
		assert.Contains(t, events, EventEnd)
	}
}

func TestWatcher_Start_エラー発生(t *testing.T) {
	t.Parallel()

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(
		domain.TempArenaInfo{},
		failure.New(apperr.UnexpectedError),
	)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)

	go watcher.Start(ctx, ctx)

	time.Sleep(100 * time.Millisecond)
	assert.Len(t, events, 1)
	assert.Contains(t, events, EventErr)
}

func TestWatcher_Start_キャンセル(t *testing.T) {
	t.Parallel()

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", mock.Anything).Return(domain.TempArenaInfo{}, nil)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}
	ctx, cancel := context.WithCancel(context.Background())

	watcher := NewWatcher(10*time.Millisecond, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)

	go watcher.Start(ctx, ctx)
	cancel()

	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, events)
}
