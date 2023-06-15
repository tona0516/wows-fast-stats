package main

import (
	"changeme/backend/apperr"
	"changeme/backend/infra"
	"changeme/backend/service"
	"changeme/backend/vo"
	"context"
	"fmt"
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PARALLELS = 5

type App struct {
	ctx                 context.Context
	version             vo.Version
	env                 vo.Env
	cancelReplayWatcher context.CancelFunc
	configService       service.Config
	screenshotService   service.Screenshot
	replayWatcher       service.ReplayWatcher
	battleService       service.Battle
	reportService       service.Report
	logger              infra.LoggerInterface
	excludePlayer       mapset.Set[int]
}

func NewApp(
	env vo.Env,
	version vo.Version,
	configService service.Config,
	screenshotService service.Screenshot,
	replayWatcher service.ReplayWatcher,
	battleService service.Battle,
	reportService service.Report,
	logger infra.LoggerInterface,
) *App {
	return &App{
		env:               env,
		version:           version,
		configService:     configService,
		screenshotService: screenshotService,
		replayWatcher:     replayWatcher,
		battleService:     battleService,
		reportService:     reportService,
		logger:            logger,
		excludePlayer:     mapset.NewSet[int](),
	}
}

func (a *App) onStartup(ctx context.Context) {
	a.logger.Debug("onStartup() called")
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
	a.logger.Debug("onShutdown() called")

	// Save windows size
	appConfig, _ := a.configService.App()
	width, height := runtime.WindowGetSize(ctx)
	appConfig.Window.Width = width
	appConfig.Window.Height = height
	if err := a.configService.UpdateApp(appConfig); err != nil {
		a.logger.Warn("Failed to update AppConfig", err)
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
		a.logger.Error("Failed to get UserConfig", err)
		return result, err
	}

	result, err = a.battleService.Battle(userConfig)
	if err != nil {
		a.logger.Error("Failed to fetch Battle", err)

		if err := a.reportService.Send(err); err != nil {
			a.logger.Warn("Failed to send Report", err)
		}

		return result, err
	}

	return result, nil
}

func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		a.logger.Warn("Failed to get directory, path:"+path, err)
	}

	return path, err
}

func (a *App) UserConfig() (vo.UserConfig, error) {
	config, err := a.configService.User()
	if err != nil {
		a.logger.Warn("Failed to get UserConfig", err)
	}

	return config, err
}

func (a *App) ApplyUserConfig(config vo.UserConfig) error {
	err := a.configService.UpdateUser(config)
	if err != nil {
		a.logger.Warn("Failed to update UserConfig", err)
	}

	return err
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil && !errors.Is(err, apperr.SrvSs.Canceled) {
		a.logger.Warn("Failed to save screenshot, filename:"+filename+" base64Data:"+base64Data, err)
	}

	return err
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Warn("Failed to save screenshot, filename:"+filename+" base64Data:"+base64Data, err)
	}

	return err
}

func (a *App) AppVersion() vo.Version {
	return a.version
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		wraped := errors.WithStack(apperr.App.OpenDir.WithRaw(err))
		a.logger.Warn("Failed to open directory, path:"+path, wraped)
		return wraped
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
	players, err := a.configService.AlertPlayers()
	if err != nil {
		a.logger.Warn("Failed to get AlertPlayers", err)
	}

	return players, err
}

func (a *App) UpdateAlertPlayer(player vo.AlertPlayer) error {
	err := a.configService.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Warn("Failed to update AlertPlayer, player.Name:"+player.Name, err)
	}

	return err
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Warn("Failed to remove AlertPlayer, accountID:"+strconv.Itoa(accountID), err)
	}

	return err
}

func (a *App) SearchPlayer(prefix string) (vo.WGAccountList, error) {
	accountList, err := a.configService.SearchPlayer(prefix)
	if err != nil {
		a.logger.Warn("Failed to search player, prefix:"+prefix, err)
	}

	return accountList, err
}

func (a *App) AlertPatterns() []string {
	return vo.AlertPatterns
}

func (a *App) LogErrorForFrontend(err string) {
	//nolint:goerr113
	a.logger.Warn("Error has occurred in frontend", fmt.Errorf("%s", err))
}

func (a *App) FontSizes() []string {
	return vo.FontSizes
}

func (a *App) StatsPatterns() []string {
	return vo.StatsPatterns
}
