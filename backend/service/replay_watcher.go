package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"

	"changeme/backend/apperr"
	"changeme/backend/infra"

	"github.com/fsnotify/fsnotify"
	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	EventStart = "BATTLE_START"
	EventEnd   = "BATTLE_END"
	EventErr   = "BATTLE_ERROR"
)

type ReplayWatcher struct {
	appCtx     context.Context
	configRepo infra.Config
	taiRepo    infra.TempArenaInfo
}

func NewReplayWatcher(appCtx context.Context, configRepo infra.Config, taiRepo infra.TempArenaInfo) *ReplayWatcher {
	return &ReplayWatcher{
		appCtx:     appCtx,
		configRepo: configRepo,
		taiRepo:    taiRepo,
	}
}

//nolint:cyclop
func (w *ReplayWatcher) Start(ctx context.Context) {
	latestHash, err := w.tempArenaInfoHash()
	if err != nil {
		w.notifyEnd()
	} else {
		w.notifyStart()
	}

	watcher, err := w.genWatcher()
	if err != nil {
		w.notifyError(err)

		return
	}
	defer watcher.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-watcher.Events:
			if !strings.HasSuffix(event.Name, infra.TempArenaInfoName) {
				continue
			}

			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				hash, err := w.tempArenaInfoHash()
				if err != nil {
					continue
				}
				if hash == latestHash {
					continue
				}
				latestHash = hash
				w.notifyStart()
			}

			if event.Has(fsnotify.Remove) {
				w.notifyEnd()
			}
		case err := <-watcher.Errors:
			w.notifyError(failure.Translate(err, apperr.ReplayWatchSvWatcherChan))

			return
		}
	}
}

func (w *ReplayWatcher) tempArenaInfoHash() (string, error) {
	var result string

	userConfig, err := w.configRepo.User()
	if err != nil {
		return result, err
	}

	tempArenaInfo, err := w.taiRepo.Get(userConfig.InstallPath)
	if err != nil {
		return result, err
	}

	md5 := sha256.Sum256([]byte(fmt.Sprintf("%x", tempArenaInfo)))
	result = fmt.Sprintf("%x", md5)

	return result, nil
}

func (w *ReplayWatcher) genWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, failure.Translate(err, apperr.ReplayWatchSvNewWatcher)
	}

	userConfig, err := w.configRepo.User()
	if err != nil {
		return watcher, err
	}

	if err := watcher.Add(filepath.Join(userConfig.InstallPath, infra.ReplayDir)); err != nil {
		return watcher, failure.Translate(err, apperr.ReplayWatchSvWatcherAdd)
	}

	return watcher, nil
}

func (w *ReplayWatcher) notifyStart() {
	runtime.EventsEmit(w.appCtx, EventStart)
}

func (w *ReplayWatcher) notifyEnd() {
	runtime.EventsEmit(w.appCtx, EventEnd)
}

func (w *ReplayWatcher) notifyError(err error) {
	runtime.EventsEmit(w.appCtx, EventErr, err)
}
