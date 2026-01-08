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
	configStore     gateway.ConfigStore
	replayReader    gateway.ReplayReader
	eventsEmitFunc  EventsEmitFunc
}

func NewPollMatch(i do.Injector) (*PollMatch, error) {
	config := do.MustInvoke[config.Config](i)
	return &PollMatch{
		pollingInterval: config.Basic.PollingInterval,
		configStore:     do.MustInvoke[gateway.ConfigStore](i),
		replayReader:    do.MustInvoke[gateway.ReplayReader](i),
		eventsEmitFunc:  do.MustInvoke[EventsEmitFunc](i),
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
	pm.eventsEmitFunc(ctx, EventNeedInitialSetting)
}

func (pm *PollMatch) emitPollingStart(ctx context.Context) {
	pm.eventsEmitFunc(ctx, EventPollingStart)
}

func (pm *PollMatch) emitBattleStart(ctx context.Context) {
	pm.eventsEmitFunc(ctx, EventBattleStart)
}

func (pm *PollMatch) emitError(ctx context.Context, err error) {
	pm.eventsEmitFunc(ctx, EventErr, err)
}
