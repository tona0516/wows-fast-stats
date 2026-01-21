package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

type PollMatch struct {
	pollingInterval time.Duration
	wails           adapter.Wails
	prefStore       adapter.PrefStore
	replayReader    adapter.ReplayReader
}

func NewPollMatch(i do.Injector) (*PollMatch, error) {
	config := do.MustInvoke[config.Config](i)
	return &PollMatch{
		pollingInterval: config.Basic.PollingInterval,
		wails:           do.MustInvoke[adapter.Wails](i),
		prefStore:       do.MustInvoke[adapter.PrefStore](i),
		replayReader:    do.MustInvoke[adapter.ReplayReader](i),
	}, nil
}

func (pm *PollMatch) Invoke(
	ctx context.Context,
	cancelCtx context.Context,
	channel chan core.TempArenaInfo,
) {
	pref, err := pm.prefStore.Pref()
	if err != nil {
		if failure.Is(err, core.ErrJSONNotFound) {
			pm.emitNeedInitialSetting(ctx)
			return
		}

		pm.emitError(ctx, err)
		return
	}

	if pref.InstallPath == "" {
		pm.emitNeedInitialSetting(ctx)
		return
	}

	pm.emitPollingStart(ctx)

	var latestHash string
	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			time.Sleep(pm.pollingInterval)

			tempArenaInfo, err := pm.replayReader.TempArenaInfo(pref.InstallPath)
			if err != nil {
				if failure.Is(err, core.ErrTempArenaInfoNotFound) {
					continue
				}

				pm.emitError(ctx, err)
				continue
			}

			hash := fmt.Sprintf("%x", sha256.Sum256(fmt.Append(nil, tempArenaInfo)))
			if hash == latestHash {
				continue
			}

			latestHash = hash
			pm.emitBattleStart(ctx)
			channel <- tempArenaInfo
		}
	}
}

func (pm *PollMatch) emitNeedInitialSetting(ctx context.Context) {
	pm.wails.EmitEvent(ctx, EventNeedInitialSetting)
}

func (pm *PollMatch) emitPollingStart(ctx context.Context) {
	pm.wails.EmitEvent(ctx, EventPollingStart)
}

func (pm *PollMatch) emitBattleStart(ctx context.Context) {
	pm.wails.EmitEvent(ctx, EventBattleStart)
}

func (pm *PollMatch) emitError(ctx context.Context, err error) {
	pm.wails.EmitEvent(ctx, EventErr, err)
}
