package usecase

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"time"
	"wfs/internal/config"
	"wfs/internal/data"
	"wfs/internal/gateway"

	"github.com/samber/do/v2"
)

type PollMatch struct {
	pollingInterval time.Duration
	wails           gateway.Wails
	configStore     gateway.ConfigStore
	replayReader    gateway.ReplayReader
}

func NewPollMatch(i do.Injector) (*PollMatch, error) {
	config := do.MustInvoke[config.Config](i)
	return &PollMatch{
		pollingInterval: config.Basic.PollingInterval,
		wails:           do.MustInvoke[gateway.Wails](i),
		configStore:     do.MustInvoke[gateway.ConfigStore](i),
		replayReader:    do.MustInvoke[gateway.ReplayReader](i),
	}, nil
}

func (pm *PollMatch) Invoke(
	ctx context.Context,
	cancelCtx context.Context,
	channel chan data.TempArenaInfo,
) {
	userConfig, err := pm.configStore.UserConfig()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			pm.emitNeedInitialSetting(ctx)
			return
		}

		pm.emitError(ctx, err)
		return
	}

	if userConfig.InstallPath == "" {
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

			tempArenaInfo, err := pm.replayReader.TempArenaInfo(userConfig.InstallPath)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
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
