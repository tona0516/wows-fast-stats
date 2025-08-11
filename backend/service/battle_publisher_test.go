package service

import (
	"context"
	"testing"
	"time"

	"wfs/backend/data"
	"wfs/backend/mock/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type ExpectedCanSubcribe struct {
	userConfig data.UserConfigV2
	expected   bool
}

func TestBattlePublisher_CanSubcribe(t *testing.T) {
	t.Parallel()

	params := []ExpectedCanSubcribe{
		{
			userConfig: data.UserConfigV2{InstallPath: "test"},
			expected:   true,
		},
		{
			userConfig: data.UserConfigV2{},
			expected:   false,
		},
	}

	for _, v := range params {
		ctrl := gomock.NewController(t)
		mockStorage := repository.NewMockStorageInterface(ctrl)
		mockStorage.EXPECT().UserConfigV2().Return(v.userConfig, nil)

		bp := NewBattlePublisher(
			context.Background(),
			1,
			nil,
			mockStorage,
			nil,
			nil,
		)

		assert.Equal(t, v.expected, bp.CanSubcribe())
	}
}

func TestBattlePublisher_Subcribe(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	// モックの作成
	mockLocalFile := repository.NewMockLocalFileInterface(ctrl)
	mockStorage := repository.NewMockStorageInterface(ctrl)

	// テスト用データ
	testConfig := data.UserConfigV2{InstallPath: "test", SaveTempArenaInfo: false}
	testArena := data.TempArenaInfo{PlayerName: "testPlayer"}

	// UserConfigV2の返却値を設定
	mockStorage.EXPECT().UserConfigV2().Return(testConfig, nil)
	// TempArenaInfoの返却値を設定
	mockLocalFile.EXPECT().TempArenaInfo(testConfig.InstallPath).Return(testArena, nil).AnyTimes()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	bp := NewBattlePublisher(
		context.Background(),
		0, // intervalを0にして即時実行
		mockLocalFile,
		mockStorage,
		nil,
		emitFunc,
	)

	// CanSubcribeでuserConfigをセット
	assert.True(t, bp.CanSubcribe())

	// channelを用意
	ch := make(chan data.TempArenaInfo, 1)
	// キャンセル用context
	ctx, cancel := context.WithCancel(context.Background())

	// goroutineでSubcribeを実行
	go bp.Subcribe(ctx, ch)

	// channelから値を受信できるか
	select {
	case arena := <-ch:
		assert.Equal(t, testArena.PlayerName, arena.PlayerName)
		// EventStartが発火しているか
		assert.Contains(t, events, EventStart)
	case <-time.After(time.Second):
		t.Fatal("channelに値が送信されませんでした")
	}

	// goroutine停止
	cancel()
}
