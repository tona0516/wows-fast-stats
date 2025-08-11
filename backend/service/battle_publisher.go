package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
)

type BattlePublisher struct {
	ctx            context.Context
	interval       time.Duration
	localFile      repository.LocalFileInterface
	storage        repository.StorageInterface
	logger         repository.LoggerInterface
	eventsEmitFunc eventEmitFunc

	userConfig *data.UserConfigV2
}

func NewBattlePublisher(
	ctx context.Context,
	interval time.Duration,
	localFile repository.LocalFileInterface,
	storage repository.StorageInterface,
	logger repository.LoggerInterface,
	eventsEmitFunc eventEmitFunc,
) *BattlePublisher {
	return &BattlePublisher{
		ctx:            ctx,
		interval:       interval,
		localFile:      localFile,
		storage:        storage,
		logger:         logger,
		eventsEmitFunc: eventsEmitFunc,
	}
}

func (bp *BattlePublisher) CanSubcribe() bool {
	userConfig, err := bp.storage.UserConfigV2()
	if err != nil {
		return false
	}

	if userConfig.InstallPath == "" {
		return false
	}

	bp.userConfig = &userConfig
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
			tempArenaInfo, err := bp.localFile.TempArenaInfo(bp.userConfig.InstallPath)
			if err != nil {
				if failure.Is(err, apperr.FileNotExist) || failure.Is(err, apperr.ReplayDirNotFoundError) {
					bp.eventsEmitFunc(bp.ctx, EventEnd)
					continue
				}

				bp.logger.Error(err, nil)
				bp.eventsEmitFunc(bp.ctx, EventErr, apperr.Unwrap(err))
				return
			}

			// tempArenaInfo.jsonのハッシュから前回の戦闘を同じかを判定
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x", tempArenaInfo))))
			if hash == latestHash {
				continue
			}

			// 必要に応じて保存
			if bp.userConfig.SaveTempArenaInfo {
				if err := bp.localFile.SaveTempArenaInfo(tempArenaInfo); err != nil {
					bp.logger.Warn(err, nil)
				}
			}

			// 新規戦闘開始を通知
			latestHash = hash
			bp.eventsEmitFunc(bp.ctx, EventStart)
			channel <- tempArenaInfo
		}
	}
}
