package main

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"time"
	"wfs/backend/apperr"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"
	"wfs/backend/infra"
	"wfs/backend/service"

	"github.com/dgraph-io/badger/v4"
	"github.com/imroc/req/v3"
	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/ratelimit"
)

const (
	eventUpdateConfig       = "CONFIG_UPDATE"
	eventUpdateAlertPlayers = "ALERT_PLAYERS_UPDATE"
)

//nolint:containedctx
type App struct {
	config Config

	ctx                   context.Context
	logger                repository.Logger
	versionFetcher        repository.VersionFetcher
	configService         *service.Config
	screenshotService     *service.Screenshot
	watcherService        *service.Watcher
	battleService         *service.Battle
	configMigratorService *service.ConfigMigrator

	cancelWacthFunc context.CancelFunc
}

func NewApp(config Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogSetLogLevel(ctx, logger.INFO)

	if err := a.inject(a.config); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) MigrateIfNeeded() error {
	if err := a.configMigratorService.ExecuteIfNeeded(); err != nil {
		a.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	return nil
}

func (a *App) StartWatching() error {
	if err := a.watcherService.Prepare(); err != nil {
		a.logger.Error(err, nil)
		return apperr.Unwrap(err)
	}

	if a.cancelWacthFunc != nil {
		a.cancelWacthFunc()
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	a.cancelWacthFunc = cancel

	go a.watcherService.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (model.Battle, error) {
	result := model.Battle{}

	userConfig, err := a.configService.User()
	if err != nil {
		a.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	result, err = a.battleService.Get(a.ctx, userConfig)
	if err != nil {
		a.logger.Error(err, nil)
		return result, apperr.Unwrap(err)
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := a.configService.SelectDirectory(a.ctx)
	if err != nil {
		a.logger.Error(err, nil)
	}

	return path, apperr.Unwrap(err)
}

func (a *App) OpenDirectory(path string) error {
	err := a.configService.OpenDirectory(path)
	if err != nil {
		a.logger.Warn(err, nil)
	}

	return apperr.Unwrap(err)
}

func (a *App) DefaultUserConfig() model.UserConfigV2 {
	return model.DefaultUserConfigV2()
}

func (a *App) UserConfig() (model.UserConfigV2, error) {
	config, err := a.configService.User()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) UpdateUserConfig(config model.UserConfigV2) error {
	err := a.configService.UpdateOptional(config)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateInstallPath(path string) string {
	err := a.configService.ValidateInstallPath(path)
	if err != nil {
		a.logger.Error(err, nil)
	}

	if err := apperr.Unwrap(err); err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) UpdateInstallPath(path string) error {
	config, err := a.configService.UpdateInstallPath(path)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateConfig, config)
	}

	return apperr.Unwrap(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) (bool, error) {
	saved, err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error(err, nil)
	}
	return saved, apperr.Unwrap(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error(err, nil)
	}
	return apperr.Unwrap(err)
}

func (a *App) Semver() string {
	return a.config.App.Semver
}

func (a *App) AlertPlayers() ([]model.AlertPlayer, error) {
	players, err := a.configService.AlertPlayers()
	if err != nil {
		a.logger.Error(err, nil)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player model.AlertPlayer) error {
	players, err := a.configService.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	players, err := a.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Error(err, nil)
	} else {
		runtime.EventsEmit(a.ctx, eventUpdateAlertPlayers, players)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) map[string]int {
	return a.configService.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return model.AlertPatterns()
}

func (a *App) LogError(errString string, contexts map[string]string) {
	err := failure.New(apperr.FrontendError, failure.Message(errString))
	a.logger.Error(err, contexts)
}

func (a *App) LogInfo(message string, contexts map[string]string) {
	a.logger.Info(message, contexts)
}

func (a *App) LatestRelease() (model.LatestRelease, error) {
	latestRelease, err := a.versionFetcher.Fetch()
	return latestRelease, apperr.Unwrap(err)
}

func (a *App) inject(config Config) error {
	db, err := badger.Open(badger.DefaultOptions("./persistent_data"))
	if err != nil {
		return err
	}

	storage := infra.NewStorage(db)

	a.logger = infra.NewLogger(
		req.C().
			SetBaseURL(a.config.Discord.AlertURL).
			SetCommonRetryCount(a.config.Discord.MaxRetry).
			SetTimeout(time.Duration(a.config.Discord.TimeoutSec)*time.Second),
		req.C().
			SetBaseURL(a.config.Discord.InfoURL).
			SetCommonRetryCount(a.config.Discord.MaxRetry).
			SetTimeout(time.Duration(a.config.Discord.TimeoutSec)*time.Second),
		*storage,
		a.config.App.Name,
		a.config.App.Semver,
		a.config.Logger.ZerologLogLevel,
	)

	rateLimiter := ratelimit.New(a.config.Wargaming.RateLimitRPS)
	wargamingClient := req.C().
		SetBaseURL(a.config.Wargaming.URL).
		AddCommonQueryParam("application_id", a.config.Wargaming.AppID).
		SetCommonRetryCount(a.config.Wargaming.MaxRetry).
		SetTimeout(time.Duration(a.config.Wargaming.TimeoutSec) * time.Second).
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {
			rateLimiter.Take()
			return nil
		}).
		AddCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}

			var rb infra.WGResponseCommon[any]
			if err := json.Unmarshal(resp.Bytes(), &rb); err != nil {
				return true
			}

			if rb.GetStatus() == "error" {
				// Note:
				// https://developers.wargaming.net/documentation/guide/getting-started/#common-errors
				message := rb.GetError().Message
				if slices.Contains([]string{"REQUEST_LIMIT_EXCEEDED", "SOURCE_NOT_AVAILABLE"}, message) {
					return true
				}
			}

			return false
		})

	localFile := infra.NewLocalFile()
	warshipStore := infra.NewWarshipFetcher(
		db,
		wargamingClient,
		req.C().
			SetBaseURL(a.config.Numbers.URL).
			SetCommonRetryCount(a.config.Numbers.MaxRetry).
			SetTimeout(time.Duration(a.config.Numbers.TimeoutSec)*time.Second).
			EnableInsecureSkipVerify(),
	)
	clanFercher := infra.NewClanFetcher(
		wargamingClient,
		req.C().
			SetBaseURL(a.config.UnofficialWargaming.URL).
			SetCommonRetryCount(a.config.UnofficialWargaming.MaxRetry).
			SetTimeout(time.Duration(a.config.UnofficialWargaming.TimeoutSec)*time.Second),
	)
	rawStatFetcher := infra.NewRawStatFetcher(wargamingClient)
	battleMetaFetcher := infra.NewBattleMetaFetcher(wargamingClient)
	accountFetcher := infra.NewAccountFetcher(wargamingClient)
	userConfig := infra.NewUserConfigStore(db)
	alertPlayer := infra.NewAlertPlayerStore(db)
	a.versionFetcher = infra.NewVersionFetcher(
		req.C().
			SetBaseURL(a.config.Github.URL).
			SetCommonRetryCount(a.config.Github.MaxRetry).
			SetTimeout(time.Duration(a.config.Github.TimeoutSec)*time.Second),
		config.App.Semver,
	)

	// service
	a.configService = service.NewConfig(accountFetcher, userConfig, alertPlayer)
	a.screenshotService = service.NewScreenshot(localFile)
	a.battleService = service.NewBattle(
		localFile,
		warshipStore,
		clanFercher,
		rawStatFetcher,
		battleMetaFetcher,
		accountFetcher,
		a.logger,
	)
	a.watcherService = service.NewWatcher(1*time.Second, localFile, userConfig, a.logger, runtime.EventsEmit)
	a.configMigratorService = service.NewConfigMigrator(storage, userConfig, alertPlayer)

	return nil
}
