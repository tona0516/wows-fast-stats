package main

import (
	"context"
	"wfs/backend/apperr"
	"wfs/backend/application/service"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/infra"
	"wfs/backend/logger"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const EventOnload = "ONLOAD"

type App struct {
	ctx               context.Context
	env               vo.Env
	cancelWatcher     context.CancelFunc
	configService     service.Config
	screenshotService service.Screenshot
	watcherService    service.Watcher
	battleService     service.Battle
	reportService     service.Report
	updaterService    service.Updater
	excludePlayer     map[int]bool
}

func NewApp(
	env vo.Env,
	configService service.Config,
	screenshotService service.Screenshot,
	watcherService service.Watcher,
	battleService service.Battle,
	reportService service.Report,
	updaterService service.Updater,
) *App {
	return &App{
		env:               env,
		configService:     configService,
		screenshotService: screenshotService,
		watcherService:    watcherService,
		battleService:     battleService,
		reportService:     reportService,
		updaterService:    updaterService,
		excludePlayer:     map[int]bool{},
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.ctx = ctx
	logger.Init(ctx, a.env)

	runtime.EventsOn(ctx, EventOnload, func(optionalData ...interface{}) {
		logger.Zerolog().Info().Msg("application started")
	})

	appConfig, err := a.configService.App()
	if err == nil {
		// Set window size
		window := appConfig.Window
		if window.Width > 0 && window.Height > 0 {
			runtime.WindowSetSize(ctx, window.Width, window.Height)
		}
	}
}

func (a *App) onShutdown(ctx context.Context) {
	logger.Zerolog().Info().Msg("application will shutdown...")

	// Save windows size
	appConfig, _ := a.configService.App()
	width, height := runtime.WindowGetSize(ctx)
	appConfig.Window.Width = width
	appConfig.Window.Height = height
	if err := a.configService.UpdateApp(appConfig); err != nil {
		logger.Zerolog().Error().Err(err).Send()
	}
}

func (a *App) StartWatching() error {
	userConfig, err := a.configService.User()
	if err != nil {
		return err
	}
	if userConfig.InstallPath == "" {
		return apperr.ErrInvalidInstallPath
	}

	if a.cancelWatcher != nil {
		a.cancelWatcher()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelWatcher = cancel
	go a.watcherService.Start(a.ctx, ctx, userConfig)

	return nil
}

func (a *App) Battle() (domain.Battle, error) {
	var result domain.Battle

	userConfig, err := a.configService.User()
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()
		return result, apperr.ToFrontendError(err)
	}

	result, err = a.battleService.Battle(userConfig)
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()

		if err := a.reportService.Send(err); err != nil {
			logger.Zerolog().Error().Err(err).Send()
		}

		return result, apperr.ToFrontendError(err)
	}

	return result, nil
}

func (a *App) SampleTeams() []domain.Team {
	return domain.SampleTeams()
}

func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		logger.Zerolog().Error().Err(err).Str("path", path).Send()
	}

	return path, apperr.ToFrontendError(err)
}

func (a *App) DefaultUserConfig() domain.UserConfig {
	return infra.DefaultUserConfig
}

func (a *App) UserConfig() (domain.UserConfig, error) {
	config, err := a.configService.User()
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()
	}

	return config, apperr.ToFrontendError(err)
}

func (a *App) ApplyUserConfig(config domain.UserConfig) error {
	err := a.configService.UpdateOptional(config)
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()
	}

	return apperr.ToFrontendError(err)
}

func (a *App) ApplyRequiredUserConfig(
	installPath string,
	appid string,
) (vo.ValidatedResult, error) {
	validatedResult, err := a.configService.UpdateRequired(installPath, appid)
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()
	}

	return validatedResult, apperr.ToFrontendError(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		logger.Zerolog().Error().Err(err).Str("filename", filename).Send()
	}

	return apperr.ToFrontendError(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		logger.Zerolog().Error().Err(err).Str("filename", filename).Send()
	}

	return apperr.ToFrontendError(err)
}

func (a *App) Semver() string {
	return a.env.Semver
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		wraped := apperr.New(apperr.OpenDirectory, err)
		logger.Zerolog().Error().Err(wraped).Str("path", path).Send()
		return apperr.ToFrontendError(wraped)
	}

	return nil
}

func (a *App) ExcludePlayerIDs() []int {
	ids := make([]int, 0, len(a.excludePlayer))
	for id := range a.excludePlayer {
		ids = append(ids, id)
	}

	return ids
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.excludePlayer[playerID] = true
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	delete(a.excludePlayer, playerID)
}

func (a *App) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := a.configService.AlertPlayers()
	if err != nil {
		logger.Zerolog().Error().Err(err).Send()
	}

	return players, apperr.ToFrontendError(err)
}

func (a *App) UpdateAlertPlayer(player domain.AlertPlayer) error {
	err := a.configService.UpdateAlertPlayer(player)
	if err != nil {
		logger.Zerolog().Error().Err(err).Str("player", player.Name).Send()
	}

	return apperr.ToFrontendError(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		logger.Zerolog().Error().Err(err).Int("account id", accountID).Send()
	}

	return apperr.ToFrontendError(err)
}

func (a *App) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	accountList, err := a.configService.SearchPlayer(prefix)
	if err != nil {
		logger.Zerolog().Error().Err(err).Str("prefix", prefix).Send()
	}

	return accountList, apperr.ToFrontendError(err)
}

func (a *App) AlertPatterns() []string {
	return domain.AlertPatterns
}

func (a *App) LogErrorForFrontend(errString string) {
	err := apperr.New(apperr.FrontendError, errors.New(errString))
	logger.Zerolog().Error().Err(err).Send()
}

func (a *App) FontSizes() []string {
	return vo.FontSizes
}

func (a *App) StatsPatterns() []string {
	return domain.StatsPatterns
}

func (a *App) LatestRelease() (domain.GHLatestRelease, error) {
	return a.updaterService.Updatable()
}
