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
	prefStore       adapter.PrefStore
	replayReader    adapter.ReplayReader
}

func NewPollMatch(i do.Injector) (*PollMatch, error) {
	config := do.MustInvoke[config.Config](i)
	return &PollMatch{
		pollingInterval: config.Basic.PollingInterval,
		prefStore:       do.MustInvoke[adapter.PrefStore](i),
		replayReader:    do.MustInvoke[adapter.ReplayReader](i),
	}, nil
}

func (pm *PollMatch) GetInstallPath() (string, error) {
	pref, err := pm.prefStore.Pref()
	if err != nil {
		if failure.Is(err, core.ErrJSONNotFound) {
			return "", failure.Translate(err, core.ErrInitialSettingRequired)
		}

		return "", err
	}

	if pref.InstallPath == "" {
		return "", failure.New(core.ErrInitialSettingRequired)
	}

	return pref.InstallPath, nil
}

func (pm *PollMatch) Invoke(
	ctx context.Context,
	cancelCtx context.Context,
	installPath string,
	result chan PollingResult,
) {
	var latestHash string
	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			time.Sleep(pm.pollingInterval)

			tempArenaInfo, err := pm.replayReader.TempArenaInfo(installPath)
			if err != nil {
				if failure.Is(err, core.ErrTempArenaInfoNotFound) {
					continue
				}

				result <- PollingResult{
					TempArenaInfo: nil,
					Error:         err,
				}
				return
			}

			hash := fmt.Sprintf("%x", sha256.Sum256(fmt.Append(nil, tempArenaInfo)))
			if hash == latestHash {
				continue
			}

			latestHash = hash
			result <- PollingResult{
				TempArenaInfo: &tempArenaInfo,
				Error:         nil,
			}
		}
	}
}
