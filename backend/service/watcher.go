package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/morikuni/failure"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
	EventErr   = "BATTLE_ERR"
)

type Watcher struct {
	interval       time.Duration
	localFile      repository.LocalFile
	userConfig     repository.UserConfigStore
	logger         repository.Logger
	eventsEmitFunc eventEmitFunc
	config         model.UserConfigV2
}

func NewWatcher(
	interval time.Duration,
	localFile repository.LocalFile,
	userConfig repository.UserConfigStore,
	logger repository.Logger,
	eventsEmitFunc eventEmitFunc,
) *Watcher {
	return &Watcher{
		interval:       interval,
		localFile:      localFile,
		userConfig:     userConfig,
		logger:         logger,
		eventsEmitFunc: eventsEmitFunc,
	}
}

func (w *Watcher) Prepare() error {
	config, err := w.userConfig.GetV2()
	if err != nil {
		return err
	}

	if config.InstallPath == "" {
		return failure.New(apperr.InvalidInstallPath)
	}

	w.config = config

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

			tempArenaInfo, err := w.localFile.ReadTempArenaInfo(w.config.InstallPath)
			if err != nil {
				if failure.Is(err, apperr.FileNotExist) || failure.Is(err, apperr.ReplayDirNotFoundError) {
					w.eventsEmitFunc(appCtx, EventEnd)

					continue
				}

				w.logger.Error(err, nil)
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
