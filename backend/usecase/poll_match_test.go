package usecase

import (
	"context"
	"testing"
	"time"
	"wfs/backend/core"
	"wfs/backend/mock"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPollMatch_GetInstallPath(t *testing.T) {
	t.Parallel()

	t.Run("正常系_保存済みのインストールパスを返す", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)
		expectedPath := "/path/to/wows"

		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{InstallPath: expectedPath}, nil)

		pm := &PollMatch{prefStore: mockPrefStore}

		path, err := pm.GetInstallPath()

		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
	})

	t.Run("異常系_Prefファイルが存在しない場合は初期設定エラーを返す", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{}, failure.New(core.ErrJSONNotFound))

		pm := &PollMatch{prefStore: mockPrefStore}

		path, err := pm.GetInstallPath()

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.True(t, failure.Is(err, core.ErrInitialSettingRequired))
	})

	t.Run("異常系_InstallPathが未設定の場合は初期設定エラーを返す", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{InstallPath: ""}, nil)

		pm := &PollMatch{prefStore: mockPrefStore}

		path, err := pm.GetInstallPath()

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.True(t, failure.Is(err, core.ErrInitialSettingRequired))
	})

	t.Run("異常系_Pref読み込みエラーをそのまま返す", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)
		expectedErr := failure.New(core.ErrJSONRead)

		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{}, expectedErr)

		pm := &PollMatch{prefStore: mockPrefStore}

		path, err := pm.GetInstallPath()

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Equal(t, expectedErr, err)
	})
}

func TestPollMatch_Invoke(t *testing.T) {
	t.Parallel()

	t.Run("正常系_TempArenaInfoの変更を検知して通知する", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockReplayReader := mock.NewMockReplayReader(ctrl)

		installPath := "/path/to/replays"
		info1 := core.TempArenaInfo{PlayerName: "player1", DateTime: "01.01.2024 00:00:00"}
		info2 := core.TempArenaInfo{PlayerName: "player2", DateTime: "02.01.2024 00:00:00"}

		callCount := 0
		mockReplayReader.EXPECT().
			TempArenaInfo(installPath).
			AnyTimes().
			DoAndReturn(func(string) (core.TempArenaInfo, error) {
				switch callCount {
				case 0:
					callCount++
					return info1, nil
				case 1:
					callCount++
					return info1, nil
				case 2:
					callCount++
					return info2, nil
				default:
					return info2, nil
				}
			})

		pm := &PollMatch{
			pollingInterval: 5 * time.Millisecond,
			replayReader:    mockReplayReader,
		}

		resultCh := make(chan PollingResult, 4)
		ctx := context.Background()
		cancelCtx, cancel := context.WithCancel(ctx)
		done := make(chan struct{})
		go func() {
			pm.Invoke(ctx, cancelCtx, installPath, resultCh)
			close(done)
		}()

		results := make([]PollingResult, 0, 2)
		timeout := time.After(200 * time.Millisecond)
		for len(results) < 2 {
			select {
			case res := <-resultCh:
				results = append(results, res)
				if len(results) == 2 {
					cancel()
				}
			case <-timeout:
				t.Fatalf("通知を受信できませんでした")
			}
		}

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("処理が終了しませんでした")
		}

		assert.NotNil(t, results[0].TempArenaInfo)
		assert.Equal(t, info1, *results[0].TempArenaInfo)
		assert.NoError(t, results[0].Error)

		assert.NotNil(t, results[1].TempArenaInfo)
		assert.Equal(t, info2, *results[1].TempArenaInfo)
		assert.NoError(t, results[1].Error)

		time.Sleep(20 * time.Millisecond)
		assert.Empty(t, resultCh)

		cancel()
	})

	t.Run("異常系_TempArenaInfoが見つからない場合はスキップする", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockReplayReader := mock.NewMockReplayReader(ctrl)

		installPath := "/path/to/replays"
		expectedInfo := core.TempArenaInfo{PlayerName: "player1", DateTime: "01.01.2024 00:00:00"}

		callCount := 0
		mockReplayReader.EXPECT().
			TempArenaInfo(installPath).
			AnyTimes().
			DoAndReturn(func(string) (core.TempArenaInfo, error) {
				if callCount == 0 {
					callCount++
					return core.TempArenaInfo{}, failure.New(core.ErrTempArenaInfoNotFound)
				}
				callCount++
				return expectedInfo, nil
			})

		pm := &PollMatch{
			pollingInterval: 5 * time.Millisecond,
			replayReader:    mockReplayReader,
		}

		resultCh := make(chan PollingResult, 2)
		ctx := context.Background()
		cancelCtx, cancel := context.WithCancel(ctx)
		done := make(chan struct{})
		go func() {
			pm.Invoke(ctx, cancelCtx, installPath, resultCh)
			close(done)
		}()

		var result PollingResult
		select {
		case result = <-resultCh:
			cancel()
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("TempArenaInfoの通知を受け取れませんでした")
		}

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("処理が終了しませんでした")
		}

		assert.NoError(t, result.Error)
		assert.NotNil(t, result.TempArenaInfo)
		assert.Equal(t, expectedInfo, *result.TempArenaInfo)
		assert.Empty(t, resultCh)
	})

	t.Run("異常系_TempArenaInfo取得でエラーが発生した場合は通知して終了する", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockReplayReader := mock.NewMockReplayReader(ctrl)

		installPath := "/path/to/replays"
		expectedErr := failure.New(core.ErrTempArenaInfoSearch)

		mockReplayReader.EXPECT().
			TempArenaInfo(installPath).
			Return(core.TempArenaInfo{}, expectedErr)

		pm := &PollMatch{
			pollingInterval: 5 * time.Millisecond,
			replayReader:    mockReplayReader,
		}

		resultCh := make(chan PollingResult, 1)
		ctx := context.Background()
		cancelCtx, cancel := context.WithCancel(ctx)
		done := make(chan struct{})
		go func() {
			pm.Invoke(ctx, cancelCtx, installPath, resultCh)
			close(done)
		}()

		var result PollingResult
		select {
		case result = <-resultCh:
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("エラー通知を受け取れませんでした")
		}

		cancel()
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("処理が終了しませんでした")
		}

		assert.Nil(t, result.TempArenaInfo)
		assert.Equal(t, expectedErr, result.Error)
		assert.Empty(t, resultCh)
	})
}
