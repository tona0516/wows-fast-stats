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
	"github.com/stretchr/testify/require"
)

func TestWatcher_Start_戦闘開始(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, nil)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(config, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)
	go watcher.Start(ctx, ctx)

	time.Sleep(20 * time.Millisecond)
	assert.Contains(t, events, EventStart)

	events = nil
	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, events)
}

func TestWatcher_Start_戦闘終了(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, failure.New(apperr.FileNotExist))
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(config, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)
	go watcher.Start(ctx, ctx)

	time.Sleep(20 * time.Millisecond)
	assert.Contains(t, events, EventEnd)

	events = nil
	time.Sleep(100 * time.Millisecond)
	assert.Contains(t, events, EventEnd)
}

func TestWatcher_Start_エラー発生(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(
		domain.TempArenaInfo{},
		failure.New(apperr.UnexpectedError),
	)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(config, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, mockStorage, emitFunc)
	err := watcher.Prepare()
	require.NoError(t, err)
	go watcher.Start(ctx, ctx)

	time.Sleep(20 * time.Millisecond)
	assert.Contains(t, events, EventErr)

	events = nil
	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, events)
}

func TestWatcher_Start_キャンセル(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := domain.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, nil)
	mockStorage := &mocks.StorageInterface{}
	mockStorage.On("ReadUserConfig").Return(config, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}
	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, mockStorage, emitFunc)

	err := watcher.Prepare()
	require.NoError(t, err)
	go watcher.Start(ctx, ctx)
	cancel()

	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, events)
}
