package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"time"
	"wfs/internal/apperr"
	"wfs/internal/data"
	"wfs/internal/infra"

	"github.com/morikuni/failure"
)

type BattlePublisher struct {
	ctx            context.Context
	interval       time.Duration
	localStorage   infra.LocalStorage
	logger         infra.Logger
	eventsEmitFunc eventEmitFunc

	installPath string
}

func NewBattlePublisher(
	ctx context.Context,
	interval time.Duration,
	localStorage infra.LocalStorage,
	logger infra.Logger,
	eventsEmitFunc eventEmitFunc,
) *BattlePublisher {
	return &BattlePublisher{
		ctx:            ctx,
		interval:       interval,
		localStorage:   localStorage,
		logger:         logger,
		eventsEmitFunc: eventsEmitFunc,
	}
}

func (bp *BattlePublisher) CanSubcribe() bool {
	config, err := bp.localStorage.UserConfig()
	if err != nil {
		return false
	}

	if !filepath.IsAbs(config.InstallPath) {
		return false
	}

	bp.installPath = config.InstallPath
	return true
}

func (bp *BattlePublisher) Subcribe(cancelCtx context.Context, channel chan data.TempArenaInfo) {
	var latestHash string

	for {
		select {
		case <-cancelCtx.Done():
			return
		default:
			// interval秒待機
			time.Sleep(bp.interval)

			// tempArenaInfo.jsonを取得
			// 取得できない場合は、イベントを発行して次のループ
			tempArenaInfo, err := bp.localStorage.TempArenaInfo(bp.installPath)
			if err != nil {
				if failure.Is(err, apperr.FileNotExist) || failure.Is(err, apperr.ReplayDirNotFoundError) {
					bp.eventsEmitFunc(bp.ctx, EventEnd)
					continue
				}

				bp.logger.Error(err, nil)
				bp.eventsEmitFunc(bp.ctx, EventErr, apperr.ToStringCode(err))
				return
			}

			// tempArenaInfo.jsonのハッシュから前回の戦闘を同じかを判定
			hash := fmt.Sprintf("%x", sha256.Sum256(fmt.Appendf(nil, "%x", tempArenaInfo)))
			if hash == latestHash {
				continue
			}

			// 新規戦闘開始を通知
			latestHash = hash
			bp.eventsEmitFunc(bp.ctx, EventStart)
			channel <- tempArenaInfo
		}
	}
}
