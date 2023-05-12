package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"changeme/backend/infra"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
)

type ReplayWatcher struct {
	appCtx     context.Context
	configRepo infra.Config
	taiRepo    infra.TempArenaInfo
}

func NewReplayWatcher(
	appCtx context.Context,
	configRepo infra.Config,
	taiRepo infra.TempArenaInfo,
) *ReplayWatcher {
	return &ReplayWatcher{
		appCtx:     appCtx,
		configRepo: configRepo,
		taiRepo:    taiRepo,
	}
}

func (w *ReplayWatcher) Start(ctx context.Context) {
	var latestHash string

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second)

			userConfig, err := w.configRepo.User()
			if err != nil {
				continue
			}

			tempArenaInfo, err := w.taiRepo.Get(userConfig.InstallPath)
			if err != nil {
				runtime.EventsEmit(w.appCtx, EventEnd)
				continue
			}

			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x", tempArenaInfo))))

			if hash == latestHash {
				continue
			}

			latestHash = hash
			runtime.EventsEmit(w.appCtx, EventStart)
		}
	}
}
