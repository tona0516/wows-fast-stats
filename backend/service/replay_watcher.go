package service

import (
	"changeme/backend/infra"
	"context"
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
    EVENT_START = "BATTLE_START"
    EVENT_END = "BATTLE_END"
    EVENT_ERR = "BATTLE_ERROR"
)

type ReplayWatcher struct {
    appCtx context.Context
    configRepo infra.Config
    taiRepo infra.TempArenaInfo
}

func NewReplayWatcher(appCtx context.Context, configRepo infra.Config, taiRepo infra.TempArenaInfo) *ReplayWatcher {
    return &ReplayWatcher{
        appCtx: appCtx,
        configRepo: configRepo,
        taiRepo: taiRepo,
    }
}

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
            if !strings.HasSuffix(event.Name, "tempArenaInfo.json") {
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
            w.notifyError(err)
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

    md5 := md5.Sum([]byte(fmt.Sprintf("%x", tempArenaInfo)))
    result = fmt.Sprintf("%x", md5)
    return result, nil
}

func (w *ReplayWatcher) genWatcher() (*fsnotify.Watcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
		return watcher, err
	}

    userConfig, err := w.configRepo.User()
    if err != nil {
        return watcher, err
    }

	if err := watcher.Add(filepath.Join(userConfig.InstallPath, "replays/")); err != nil {
        return watcher, err
	}

    return watcher, nil
}

func (w *ReplayWatcher) notifyStart() {
    runtime.EventsEmit(w.appCtx, EVENT_START)
}

func (w *ReplayWatcher) notifyEnd() {
    runtime.EventsEmit(w.appCtx, EVENT_END)
}

func (w *ReplayWatcher) notifyError(err error) {
    runtime.EventsEmit(w.appCtx, EVENT_ERR, err)
}

