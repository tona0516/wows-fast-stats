package usecase

import (
	"context"
	"errors"
	"io/fs"
	"slices"
	"sync"
	"testing"
	"time"
	"wfs/internal/data"
	"wfs/internal/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPollMatch_Invoke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		setupMock        func(*mock.MockConfigStore)
		expectEventNames []string
	}{
		{
			name: "UserConfigが存在しない場合",
			setupMock: func(m *mock.MockConfigStore) {
				m.EXPECT().
					UserConfig().
					Return(data.UserConfig{}, fs.ErrNotExist)
			},
			expectEventNames: []string{EventNeedInitialSetting},
		},
		{
			name: "UserConfig取得でエラーが発生した場合",
			setupMock: func(m *mock.MockConfigStore) {
				m.EXPECT().
					UserConfig().
					Return(data.UserConfig{}, errors.New("read error"))
			},
			expectEventNames: []string{EventErr},
		},
		{
			name: "InstallPathが空文字の場合",
			setupMock: func(m *mock.MockConfigStore) {
				userConfig := data.UserConfig{InstallPath: ""}
				m.EXPECT().
					UserConfig().
					Return(userConfig, nil)
			},
			expectEventNames: []string{EventNeedInitialSetting},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			configStore := mock.NewMockConfigStore(ctrl)
			mockReplayReader := mock.NewMockReplayReader(ctrl)
			tt.setupMock(configStore)

			var emittedEvents []string
			var eventsMutex sync.Mutex
			emitFunc := func(ctx context.Context, eventName string, optionalData ...any) {
				eventsMutex.Lock()
				defer eventsMutex.Unlock()
				emittedEvents = append(emittedEvents, eventName)
			}

			pm := NewPollMatch(50*time.Millisecond, configStore, mockReplayReader, emitFunc)

			ctx := context.Background()
			cancelCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			channel := make(chan data.TempArenaInfo, 10)
			go pm.Invoke(ctx, cancelCtx, channel)

			// チャンネルのデータを受け取る（タイムアウトで終了）
			<-cancelCtx.Done()
			close(channel)

			// イベント名の確認
			for _, expectedEvent := range tt.expectEventNames {
				found := slices.Contains(emittedEvents, expectedEvent)
				assert.True(t, found, "expected event %s not found in emitted events: %v", expectedEvent, emittedEvents)
			}
		})
	}
}

func TestPollMatch_InvokeWithDataChange(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	MockConfigStore := mock.NewMockConfigStore(ctrl)
	mockReplayReader := mock.NewMockReplayReader(ctrl)

	userConfig := data.UserConfig{InstallPath: "/path/to/install"}
	MockConfigStore.EXPECT().
		UserConfig().
		Return(userConfig, nil)

	tempArenaInfo1 := data.TempArenaInfo{
		Vehicles: []data.Vehicle{
			{ShipID: 1, Relation: 0, ID: 100, Name: "player1"},
		},
		DateTime:   "22.05.2023 12:34:56",
		MapID:      10,
		MatchGroup: "pvp",
		PlayerName: "player1",
	}

	// TempArenaInfoを複数回呼び出して、最初は同じデータ、その後異なるデータを返す
	callCount := 0
	mockReplayReader.EXPECT().
		TempArenaInfo("/path/to/install").
		DoAndReturn(func(path string) (data.TempArenaInfo, error) {
			callCount++
			if callCount > 2 {
				// 3回目以降は異なるデータを返す
				return data.TempArenaInfo{
					Vehicles: []data.Vehicle{
						{ShipID: 2, Relation: 0, ID: 100, Name: "player1"},
					},
					DateTime:   "22.05.2023 12:35:00",
					MapID:      10,
					MatchGroup: "pvp",
					PlayerName: "player1",
				}, nil
			}
			return tempArenaInfo1, nil
		}).
		AnyTimes()

	var emittedEvents []string
	var eventsMutex sync.Mutex
	emitFunc := func(ctx context.Context, eventName string, optionalData ...any) {
		eventsMutex.Lock()
		defer eventsMutex.Unlock()
		emittedEvents = append(emittedEvents, eventName)
	}

	pm := NewPollMatch(30*time.Millisecond, MockConfigStore, mockReplayReader, emitFunc)

	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	channel := make(chan data.TempArenaInfo, 10)
	go pm.Invoke(ctx, cancelCtx, channel)

	// チャンネルのデータを受け取る
	receivedData := []data.TempArenaInfo{}
	<-cancelCtx.Done()
	close(channel)
	for data := range channel {
		receivedData = append(receivedData, data)
	}

	// ポーリング開始とバトル開始イベントが発火されている
	eventsMutex.Lock()
	hasPollingStart := false
	hasBattleStart := false
	for _, event := range emittedEvents {
		if event == EventPollingStart {
			hasPollingStart = true
		}
		if event == EventBattleStart {
			hasBattleStart = true
		}
	}
	eventsMutex.Unlock()

	assert.True(t, hasPollingStart)
	assert.True(t, hasBattleStart)
}
