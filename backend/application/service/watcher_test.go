package service

import (
	"context"
	"testing"
	"time"
	"wfs/backend/domain"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, emitFunc)
	go watcher.Start(ctx, ctx, config)

	// 20ms待ってEventStartが発行されたことを検証する
	time.Sleep(20 * time.Millisecond)
	assert.Contains(t, events, EventStart)

	// さらに100ms待ってEventStartが発行されなかったことを検証する
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

	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, errors.New("not exists"))

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}

	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, emitFunc)
	go watcher.Start(ctx, ctx, config)

	// 20ms待ってEventEndが発行されたことを検証する
	time.Sleep(100 * time.Millisecond)
	assert.Contains(t, events, EventEnd)

	// さらに100秒待ってEventEndが発行されなかったことを検証する
	events = nil
	time.Sleep(20 * time.Millisecond)
	assert.Contains(t, events, EventEnd)
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

	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("TempArenaInfo", config.InstallPath).Return(domain.TempArenaInfo{}, nil)

	var events []string
	emitFunc := func(ctx context.Context, eventName string, optionalData ...interface{}) {
		events = append(events, eventName)
	}
	interval := 10 * time.Millisecond

	watcher := NewWatcher(interval, mockLocalFile, emitFunc)

	// Startメソッドをゴルーチンで非同期に実行する
	go watcher.Start(ctx, ctx, config)

	// キャンセルしてイベントが発行されなかったことを検証する
	cancel()
	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, events)
}
