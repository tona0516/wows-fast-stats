package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"changeme/backend/infra"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
)

type Watcher struct {
	configRepo     infra.ConfigInterface
	taiRepo        infra.TempArenaInfoInterface
	eventsEmitFunc func(ctx context.Context, eventName string, optionalData ...interface{})
}

func NewWatcher(
	configRepo infra.ConfigInterface,
	taiRepo infra.TempArenaInfoInterface,
	eventsEmitFunc func(ctx context.Context, eventName string, optionalData ...interface{}),
) *Watcher {
	return &Watcher{
		configRepo:     configRepo,
		taiRepo:        taiRepo,
		eventsEmitFunc: eventsEmitFunc,
	}
}

func (w *Watcher) Start(appCtx context.Context, ctx context.Context) {
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
			if userConfig.Appid == "" {
				continue
			}

			tempArenaInfo, err := w.taiRepo.Get(userConfig.InstallPath)
			if err != nil {
				w.eventsEmitFunc(appCtx, EventEnd)
				continue
			}

			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x", tempArenaInfo))))

			if hash == latestHash {
				continue
			}

			latestHash = hash
			w.eventsEmitFunc(appCtx, EventStart)
		}
	}
}
