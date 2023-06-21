package service

import (
	"changeme/backend/vo"
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWatcher_Start_戦闘開始(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockConfigRepo := &mockConfigRepo{}
	mockConfigRepo.On("User").Return(config, nil)
	mockTaiRepo := &mockTempArenaInfoRepo{}
	mockTaiRepo.On("Get", config.InstallPath).Return(vo.TempArenaInfo{}, nil)

	// イベントが発行されたかどうかを検証するための変数
	var events []string

	// イベントを発行する関数
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	watcher := NewWatcher(mockConfigRepo, mockTaiRepo, emitFunc)

	// Startメソッドをゴルーチンで非同期に実行する
	go watcher.Start(ctx, ctx)

	// 2秒待ってEventStartが発行されたことを検証する
	time.Sleep(2 * time.Second)
	assert.Contains(t, events, EventStart)

	// 2秒待ってEventStartが発行されなかったことを検証する
	events = nil
	time.Sleep(2 * time.Second)
	assert.Empty(t, events)
}

func TestWatcher_Start_戦闘終了(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockConfigRepo := &mockConfigRepo{}
	mockConfigRepo.On("User").Return(config, nil)
	mockTaiRepo := &mockTempArenaInfoRepo{}

	mockTaiRepo.On("Get", config.InstallPath).Return(vo.TempArenaInfo{}, errors.New("not exists"))

	// イベントが発行されたかどうかを検証するための変数
	var events []string

	// イベントを発行する関数
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	watcher := NewWatcher(mockConfigRepo, mockTaiRepo, emitFunc)

	// Startメソッドをゴルーチンで非同期に実行する
	go watcher.Start(ctx, ctx)

	// 2秒待ってEventEndが発行されたことを検証する
	time.Sleep(2 * time.Second)
	assert.Contains(t, events, EventEnd)

	// 2秒待ってEventEndが発行されなかったことを検証する
	events = nil
	time.Sleep(2 * time.Second)
	assert.Contains(t, events, EventEnd)
}

func TestWatcher_Start_キャンセル(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := vo.UserConfig{
		InstallPath: "install_path_test",
		Appid:       "abc123",
		FontSize:    "medium",
	}

	mockConfigRepo := &mockConfigRepo{}
	mockConfigRepo.On("User").Return(config, nil)
	mockTaiRepo := &mockTempArenaInfoRepo{}
	mockTaiRepo.On("Get", config.InstallPath).Return(vo.TempArenaInfo{}, nil)

	// イベントが発行されたかどうかを検証するための変数
	var events []string

	// イベントを発行する関数
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	watcher := NewWatcher(mockConfigRepo, mockTaiRepo, emitFunc)

	// Startメソッドをゴルーチンで非同期に実行する
	go watcher.Start(ctx, ctx)

	// キャンセルしてイベントが発行されなかったことを検証する
	cancel()
	time.Sleep(2 * time.Second)
	assert.Empty(t, events)
}
