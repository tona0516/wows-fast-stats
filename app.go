package main

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PARALLELS = 5

type App struct {
	version             vo.Version
	env                 vo.Env
	cancelReplayWatcher context.CancelFunc
	configService       service.Config
	screenshotService   service.Screenshot
	replayWatcher       service.ReplayWatcher
	battleService       service.Battle
	logger              infra.Logger
	ctx                 context.Context
	excludePlayer       mapset.Set[int]
}

func NewApp(
	env vo.Env,
	version vo.Version,
	configService service.Config,
	screenshotService service.Screenshot,
	replayWatcher service.ReplayWatcher,
	battleService service.Battle,
	logger infra.Logger,
) *App {
	return &App{
		env:               env,
		version:           version,
		configService:     configService,
		screenshotService: screenshotService,
		replayWatcher:     replayWatcher,
		battleService:     battleService,
		logger:            logger,
		excludePlayer:     mapset.NewSet[int](),
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.logger.Info("start app.")
	a.ctx = ctx

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
	a.logger.Info("shutdown app.")

	appConfig, err := a.configService.App()
	if err != nil {
		a.logger.Info("No app config.")
	}

	// Save windows size
	width, height := runtime.WindowGetSize(ctx)
	appConfig.Window.Width = width
	appConfig.Window.Height = height
	if err := a.configService.UpdateApp(appConfig); err != nil {
		a.logger.Warn("Failed to update app config.", err)
	}
}

func (a *App) Ready() {
	if a.cancelReplayWatcher != nil {
		a.cancelReplayWatcher()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelReplayWatcher = cancel

	go a.replayWatcher.Start(a.ctx, ctx)
}

func (a *App) Battle() (vo.Battle, error) {
	var result vo.Battle

	userConfig, err := a.configService.User()
	if err != nil {
		return result, err
	}

	result, err = a.battleService.Battle(userConfig)
	if err != nil {
		a.logger.Error("Failed to get battle.", err)
		return result, err
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		a.logger.Error("Failed to get client installing path.", err)
	}

	return path, err
}

func (a *App) UserConfig() (vo.UserConfig, error) {
	return a.configService.User()
}

func (a *App) ApplyUserConfig(config vo.UserConfig) error {
	err := a.configService.UpdateUser(config)
	if err != nil {
		a.logger.Error("Failed to apply user config.", err)
	}

	return err
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to save screenshot by manual.", err)
	}

	return err
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Error("Failed to save screenshot by auto.", err)
	}

	return err
}

func (a *App) AppVersion() vo.Version {
	return a.version
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		a.logger.Warn("Failed to open directory -> "+path, err)
		return errors.WithStack(apperr.App.OpenDir.WithRaw(err))
	}

	return nil
}

func (a *App) ExcludePlayerIDs() []int {
	return a.excludePlayer.ToSlice()
}

func (a *App) AddExcludePlayerID(playerID int) {
	a.excludePlayer.Add(playerID)
}

func (a *App) RemoveExcludePlayerID(playerID int) {
	a.excludePlayer.Remove(playerID)
}

func (a *App) AlertPlayers() ([]vo.AlertPlayer, error) {
	return a.configService.AlertPlayers()
}

func (a *App) UpdateAlertPlayer(player vo.AlertPlayer) error {
	return a.configService.UpdateAlertPlayer(player)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	return a.configService.RemoveAlertPlayer(accountID)
}

func (a *App) SearchPlayer(prefix string) (vo.WGAccountList, error) {
	return a.configService.SearchPlayer(prefix)
}

func (a *App) AlertPatterns() []string {
	return vo.AlertPatterns
}

func (a *App) LogError(err string) {
	//nolint:goerr113
	a.logger.Error("Error occurred in frontend.", fmt.Errorf("%s", err))
}
