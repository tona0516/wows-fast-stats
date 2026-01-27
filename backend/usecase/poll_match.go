package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"
	"wfs/backend/usecase/service"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
)

const (
	eventStartPolling               = "START_POLLING"
	eventStartBattle                = "START_BATTLE"
	eventEmptyGameClientPathError   = "EMPTY_GAME_CLIENT_PATH_ERROR"
	eventInvalidGameClientPathError = "INVALID_GAME_CLIENT_PATH_ERROR"
	eventUnexpectedError            = "UNEXPECTED_ERROR"
)

type PollMatch struct {
	pollingInterval time.Duration
	replayReader    adapter.ReplayReader
	prefStore       adapter.PrefStore
	wails           adapter.Wails
	validator       *service.GameClientPathValidator
}

func NewPollMatch(i do.Injector) (*PollMatch, error) {
	config := do.MustInvoke[config.Config](i)
	return &PollMatch{
		pollingInterval: config.Basic.PollingInterval,
		replayReader:    do.MustInvoke[adapter.ReplayReader](i),
		prefStore:       do.MustInvoke[adapter.PrefStore](i),
		wails:           do.MustInvoke[adapter.Wails](i),
		validator:       do.MustInvoke[*service.GameClientPathValidator](i),
	}, nil
}

func (pm *PollMatch) Invoke(
	ctx context.Context,
	cancelCtx context.Context,
) {
	gameClientPath, err := pm.prefStore.GameClientPath()
	if err != nil {
		if failure.Is(err, core.ErrStringNotFound) {
			pm.wails.EmitEvent(ctx, eventEmptyGameClientPathError)
			return
		}

		pm.wails.EmitEvent(ctx, eventUnexpectedError, core.ErrorForDisplay(err))
		return
	}

	if err := pm.validator.Validate(gameClientPath); err != nil {
		pm.wails.EmitEvent(ctx, eventInvalidGameClientPathError)
		return
	}

	pm.wails.EmitEvent(ctx, eventStartPolling)

	var latestHash string
	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			time.Sleep(pm.pollingInterval)

			tempArenaInfo, err := pm.replayReader.TempArenaInfo(gameClientPath)
			if err != nil {
				if failure.Is(err, core.ErrTempArenaInfoNotFound) {
					continue
				}

				pm.wails.EmitEvent(ctx, eventUnexpectedError, core.ErrorForDisplay(err))
				return
			}

			hash := fmt.Sprintf("%x", sha256.Sum256(fmt.Append(nil, tempArenaInfo)))
			if hash == latestHash {
				continue
			}

			latestHash = hash
			pm.wails.EmitEvent(ctx, eventStartBattle, tempArenaInfo)
		}
	}
}
