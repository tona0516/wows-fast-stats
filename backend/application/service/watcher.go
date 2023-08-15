package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/logger"

	"github.com/morikuni/failure"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
	EventErr   = "BATTLE_ERR"
)

type Watcher struct {
	interval       time.Duration
	localFile      repository.LocalFileInterface
	eventsEmitFunc vo.EventEmit
	appCtx         context.Context
	userConfig     domain.UserConfig
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

func (w *Watcher) Prepare(appCtx context.Context) error {
	userConfig, err := w.localFile.User()
	if err != nil {
		return failure.Wrap(err)
	}

	if userConfig.InstallPath == "" {
		return failure.New(apperr.InvalidInstallPath)
	}

	w.appCtx = appCtx
	w.userConfig = userConfig
	return nil
}

func (w *Watcher) Start(ctx context.Context) {
	var latestHash string

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(w.interval)

			tempArenaInfo, err := w.localFile.TempArenaInfo(w.userConfig.InstallPath)
			if err != nil {
				if failure.Is(err, apperr.FileNotExist) {
					w.eventsEmitFunc(w.appCtx, EventEnd)
					continue
				}

				logger.Error(err)
				w.eventsEmitFunc(w.appCtx, EventErr, apperr.Unwrap(err))
				return
			}

			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x", tempArenaInfo))))

			if hash == latestHash {
				continue
			}

			latestHash = hash
			w.eventsEmitFunc(w.appCtx, EventStart)
		}
	}
}
