package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
)

type Watcher struct {
	interval       time.Duration
	localFile      repository.LocalFileInterface
	eventsEmitFunc vo.EventEmit
}

func NewWatcher(
	interval time.Duration,
	localFile repository.LocalFileInterface,
	eventsEmitFunc vo.EventEmit,
) *Watcher {
	return &Watcher{
		interval:       interval,
		localFile:      localFile,
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
			time.Sleep(w.interval)

			userConfig, err := w.localFile.User()
			if err != nil {
				continue
			}
			if userConfig.Appid == "" {
				continue
			}

			tempArenaInfo, err := w.localFile.TempArenaInfo(userConfig.InstallPath)
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
