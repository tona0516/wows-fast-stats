package main

import (
	"context"
	"wfs/backend/apperr"
	"wfs/backend/application/service"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/infra"
	"wfs/backend/logger"
	"wfs/backend/logger/repository"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const EventOnload = "ONLOAD"

//nolint:containedctx
type App struct {
	ctx               context.Context
	env               vo.Env
	cancelWatcher     context.CancelFunc
	reportRepo        repository.ReportInterface
	configService     service.Config
	screenshotService service.Screenshot
	watcherService    service.Watcher
	battleService     service.Battle
	updaterService    service.Updater
	excludePlayers    domain.ExcludedPlayers
}

func NewApp(
	env vo.Env,
	reportRepo repository.ReportInterface,
	configService service.Config,
	screenshotService service.Screenshot,
	watcherService service.Watcher,
	battleService service.Battle,
	updaterService service.Updater,
) *App {
	return &App{
		env:               env,
		reportRepo:        reportRepo,
		configService:     configService,
		screenshotService: screenshotService,
		watcherService:    watcherService,
		battleService:     battleService,
		updaterService:    updaterService,
		excludePlayers:    domain.ExcludedPlayers{},
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	logger.Init(ctx, a.env, a.reportRepo)

	runtime.EventsOn(ctx, EventOnload, func(optionalData ...interface{}) {
		logger.Info("application started")
	})
}

func (a *App) StartWatching() error {
	if err := a.watcherService.Prepare(); err != nil {
		logger.Error(err)
		return apperr.Unwrap(err)
	}

	if a.cancelWatcher != nil {
		a.cancelWatcher()
	}
	cancelCtx, cancel := context.WithCancel(context.Background())
	a.cancelWatcher = cancel

	go a.watcherService.Start(a.ctx, cancelCtx)

	return nil
}

func (a *App) Battle() (domain.Battle, error) {
	result := domain.Battle{}

	userConfig, err := a.configService.User()
	if err != nil {
		logger.Error(err)
		return result, apperr.Unwrap(err)
	}

	result, err = a.battleService.Battle(userConfig)
	if err != nil {
		logger.Error(err)
		return result, apperr.Unwrap(err)
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := a.configService.SelectDirectory(a.ctx)
	if err != nil {
		logger.Error(err)
	}

	return path, apperr.Unwrap(err)
}

func (a *App) OpenDirectory(path string) error {
	err := a.configService.OpenDirectory(path)
	if err != nil {
		logger.Warn(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) DefaultUserConfig() domain.UserConfig {
	return infra.DefaultUserConfig
}

func (a *App) UserConfig() (domain.UserConfig, error) {
	config, err := a.configService.User()
	if err != nil {
		logger.Error(err)
	}

	return config, apperr.Unwrap(err)
}

func (a *App) ApplyUserConfig(config domain.UserConfig) error {
	err := a.configService.UpdateOptional(config)
	if err != nil {
		logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) ValidateRequiredConfig(
	installPath string,
	appid string,
) vo.RequiredConfigError {
	return a.configService.ValidateRequired(installPath, appid)
}

func (a *App) ApplyRequiredUserConfig(
	installPath string,
	appid string,
) (vo.RequiredConfigError, error) {
	validatedResult, err := a.configService.UpdateRequired(installPath, appid)
	if err != nil {
		logger.Error(err)
	}

	return validatedResult, apperr.Unwrap(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) (bool, error) {
	saved, err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		logger.Error(err)
	}
	return saved, apperr.Unwrap(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		logger.Error(err)
	}
	return apperr.Unwrap(err)
}

func (a *App) Semver() string {
	return a.env.Semver
}

func (a *App) ExcludePlayerIDs() []int {
	return a.excludePlayers.PlayerIDs()
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.excludePlayers.Add(playerID)
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	a.excludePlayers.Remove(playerID)
}

func (a *App) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := a.configService.AlertPlayers()
	if err != nil {
		logger.Error(err)
	}

	return players, apperr.Unwrap(err)
}

func (a *App) UpdateAlertPlayer(player domain.AlertPlayer) error {
	err := a.configService.UpdateAlertPlayer(player)
	if err != nil {
		logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		logger.Error(err)
	}

	return apperr.Unwrap(err)
}

func (a *App) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	accountList, err := a.configService.SearchPlayer(prefix)
	if err != nil {
		logger.Error(err)
	}

	return accountList, apperr.Unwrap(err)
}

func (a *App) AlertPatterns() []string {
	return domain.AlertPatterns
}

func (a *App) LogError(errString string) {
	err := failure.New(apperr.FrontendError, failure.Messagef("%s", errString))
	logger.Error(err)
}

func (a *App) LatestRelease() (domain.GHLatestRelease, error) {
	latestRelease, err := a.updaterService.Updatable()
	return latestRelease, apperr.Unwrap(err)
}
