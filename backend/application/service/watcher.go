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
	storage        repository.StorageInterface
	eventsEmitFunc vo.EventEmit
	userConfig     domain.UserConfig
}

func NewWatcher(
	interval time.Duration,
	localFile repository.LocalFileInterface,
	storage repository.StorageInterface,
	eventsEmitFunc vo.EventEmit,
) *Watcher {
	return &Watcher{
		interval:       interval,
		localFile:      localFile,
		storage:        storage,
		eventsEmitFunc: eventsEmitFunc,
	}
}

func (w *Watcher) Prepare() error {
	userConfig, err := w.storage.ReadUserConfig()
	if err != nil {
		return failure.Wrap(err)
	}

	if userConfig.InstallPath == "" {
		return failure.New(apperr.InvalidInstallPath)
	}

	w.userConfig = userConfig
	return nil
}

func (w *Watcher) Start(appCtx context.Context, cancelCtx context.Context) {
	var latestHash string

	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			time.Sleep(w.interval)

			tempArenaInfo, err := w.localFile.TempArenaInfo(w.userConfig.InstallPath)
			if err != nil {
				if failure.Is(err, apperr.FileNotExist) || failure.Is(err, apperr.ReplayDirNotFoundError) {
					w.eventsEmitFunc(appCtx, EventEnd)
					continue
				}

				logger.Error(err)
				w.eventsEmitFunc(appCtx, EventErr, apperr.Unwrap(err))
				return
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
