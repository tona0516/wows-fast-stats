package service

import (
	"context"
	"testing"
	"time"

	"wfs/backend/data"
	"wfs/backend/domain"
	"wfs/backend/mock/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type ExpectedCanSubcribe struct {
	installPath string
	expected    bool
}

func TestBattlePublisher_CanSubcribe(t *testing.T) {
	t.Parallel()

	params := []ExpectedCanSubcribe{
		{
			installPath: "test",
			expected:    true,
		},
		{
			installPath: "",
			expected:    false,
		},
	}

	ctrl := gomock.NewController(t)
	mockUserConfig := repository.NewMockUserConfigInterface(ctrl)

	for _, v := range params {
		mockUserConfig.EXPECT().Load().Return(&domain.UserConfig{
			InstallPath: v.installPath,
		}, nil)

		bp := NewBattlePublisher(
			context.Background(),
			1,
			nil,
			mockUserConfig,
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
	testArena := data.TempArenaInfo{PlayerName: "testPlayer"}
	mockLocalFile := repository.NewMockLocalFileInterface(ctrl)
	mockLocalFile.EXPECT().TempArenaInfo("test").Return(testArena, nil).AnyTimes()

	mockUserConfig := repository.NewMockUserConfigInterface(ctrl)
	mockUserConfig.EXPECT().Load().Return(&domain.UserConfig{InstallPath: "test"}, nil).AnyTimes()

	// イベント発火履歴を記録するモック
	var events []string
	emitFunc := func(ctx context.Context, event string, args ...interface{}) {
		events = append(events, event)
	}

	bp := NewBattlePublisher(
		context.Background(),
		0, // intervalを0にして即時実行
		mockLocalFile,
		mockUserConfig,
		nil,
		emitFunc,
	)

	// CanSubcribeでinstallPathをセット
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
