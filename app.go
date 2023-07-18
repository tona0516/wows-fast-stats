package main

import (
	"context"
	"fmt"
	"strconv"
	"wfs/backend/apperr"
	"wfs/backend/application/repository"
	"wfs/backend/application/service"
	"wfs/backend/application/vo"
	"wfs/backend/domain"
	"wfs/backend/infra"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const PARALLELS = 5

type App struct {
	ctx               context.Context
	version           vo.Version
	env               vo.Env
	cancelWatcher     context.CancelFunc
	configService     service.Config
	screenshotService service.Screenshot
	watcherService    service.Watcher
	battleService     service.Battle
	reportService     service.Report
	updaterService    service.Updater
	logger            repository.LoggerInterface
	excludePlayer     map[int]bool
}

func NewApp(
	env vo.Env,
	version vo.Version,
	configService service.Config,
	screenshotService service.Screenshot,
	watcherService service.Watcher,
	battleService service.Battle,
	reportService service.Report,
	updaterService service.Updater,
	logger repository.LoggerInterface,
) *App {
	return &App{
		env:               env,
		version:           version,
		configService:     configService,
		screenshotService: screenshotService,
		watcherService:    watcherService,
		battleService:     battleService,
		reportService:     reportService,
		updaterService:    updaterService,
		logger:            logger,
		excludePlayer:     map[int]bool{},
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
	if a.cancelWatcher != nil {
		a.cancelWatcher()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelWatcher = cancel

	go a.watcherService.Start(a.ctx, ctx)
}

func (a *App) Battle() (domain.Battle, error) {
	var result domain.Battle

	userConfig, err := a.configService.User()
	if err != nil {
		a.logger.Error("Failed to get UserConfig", err)
		return result, apperr.ToFrontendError(err)
	}

	result, err = a.battleService.Battle(userConfig)
	if err != nil {
		a.logger.Error("Failed to fetch Battle", err)

		if err := a.reportService.Send(a.version, err); err != nil {
			a.logger.Warn("Failed to send Report", err)
		}

		return result, apperr.ToFrontendError(err)
	}

	return result, nil
}

func (a *App) SampleTeams() []domain.Team {
	players := make([]domain.Player, 8)

	tiers := []uint{
		11,
		10,
		9,
		8,
		7,
		6,
		5,
		4,
	}

	shipTypes := []domain.ShipType{
		domain.CV,
		domain.BB,
		domain.BB,
		domain.CL,
		domain.CL,
		domain.DD,
		domain.DD,
		domain.SS,
	}

	prs := []float64{
		0,
		750,
		1100,
		1350,
		1550,
		1750,
		2100,
		2450,
	}

	damageRatios := []float64{
		0,
		0.6,
		0.8,
		1.0,
		1.2,
		1.4,
		1.5,
		1.6,
	}

	winRates := []float64{
		0,
		47,
		50,
		52,
		54,
		56,
		60,
		65,
	}

	for i := range players {
		playerInfo := domain.PlayerInfo{
			ID:   1,
			Name: fmt.Sprintf("player_name%d", i+1),
			Clan: domain.Clan{
				Tag: "TEST",
			},
		}
		shipInfo := domain.ShipInfo{
			Name:      "Test Ship",
			Nation:    "japan",
			Tier:      tiers[i],
			Type:      shipTypes[i],
			AvgDamage: 10000,
		}
		shipStats := domain.ShipStats{
			Battles:            10,
			Damage:             10000 * damageRatios[i],
			WinRate:            winRates[i],
			WinSurvivedRate:    50,
			LoseSurvivedRate:   50,
			KdRate:             1,
			Kill:               1,
			Exp:                1000,
			MainBatteryHitRate: 50,
			TorpedoesHitRate:   5,
			PlanesKilled:       5,
			PR:                 prs[i],
		}
		overallStats := domain.OverallStats{
			Battles:          10,
			Damage:           10000 * damageRatios[i],
			WinRate:          winRates[i],
			WinSurvivedRate:  50,
			LoseSurvivedRate: 50,
			KdRate:           1,
			Kill:             1,
			Exp:              1000,
			AvgTier:          5,
			UsingShipTypeRate: domain.ShipTypeGroup{
				SS: 20,
				DD: 20,
				CL: 20,
				BB: 20,
				CV: 20,
			},
			UsingTierRate: domain.TierGroup{
				Low:    33.3,
				Middle: 33.3,
				High:   33.4,
			},
		}
		players[i] = domain.Player{
			PlayerInfo: playerInfo,
			ShipInfo:   shipInfo,
			PvPSolo: domain.PlayerStats{
				ShipStats:    shipStats,
				OverallStats: overallStats,
			},
			PvPAll: domain.PlayerStats{
				ShipStats:    shipStats,
				OverallStats: overallStats,
			},
		}
	}

	return []domain.Team{
		{
			Players: players,
		},
	}
}

func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		a.logger.Warn("Failed to get directory, path:"+path, err)
	}

	return path, apperr.ToFrontendError(err)
}

func (a *App) DefaultUserConfig() domain.UserConfig {
	return infra.DefaultUserConfig
}

func (a *App) UserConfig() (domain.UserConfig, error) {
	config, err := a.configService.User()
	if err != nil {
		a.logger.Warn("Failed to get UserConfig", err)
	}

	return config, apperr.ToFrontendError(err)
}

func (a *App) ApplyUserConfig(config domain.UserConfig) error {
	err := a.configService.UpdateOptional(config)
	if err != nil {
		a.logger.Warn("Failed to update UserConfig", err)
	}

	return apperr.ToFrontendError(err)
}

func (a *App) ApplyRequiredUserConfig(
	installPath string,
	appid string,
) (vo.ValidatedResult, error) {
	validatedResult, err := a.configService.UpdateRequired(installPath, appid)
	if err != nil {
		a.logger.Warn("Failed to update UserConfig for required", err)
	}

	return validatedResult, apperr.ToFrontendError(err)
}

func (a *App) ManualScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveWithDialog(a.ctx, filename, base64Data)
	if err != nil {
		a.logger.Warn("Failed to save screenshot, filename:"+filename+" base64Data:"+base64Data, err)
	}

	return apperr.ToFrontendError(err)
}

func (a *App) AutoScreenshot(filename string, base64Data string) error {
	err := a.screenshotService.SaveForAuto(filename, base64Data)
	if err != nil {
		a.logger.Warn("Failed to save screenshot, filename:"+filename+" base64Data:"+base64Data, err)
	}

	return apperr.ToFrontendError(err)
}

func (a *App) AppVersion() vo.Version {
	return a.version
}

func (a *App) OpenDirectory(path string) error {
	err := open.Run(path)
	if err != nil {
		wraped := apperr.New(apperr.OpenDirectory, err)
		a.logger.Warn("Failed to open directory, path:"+path, wraped)
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
		a.logger.Warn("Failed to get AlertPlayers", err)
	}

	return players, apperr.ToFrontendError(err)
}

func (a *App) UpdateAlertPlayer(player domain.AlertPlayer) error {
	err := a.configService.UpdateAlertPlayer(player)
	if err != nil {
		a.logger.Warn("Failed to update AlertPlayer, player.Name:"+player.Name, err)
	}

	return apperr.ToFrontendError(err)
}

func (a *App) RemoveAlertPlayer(accountID int) error {
	err := a.configService.RemoveAlertPlayer(accountID)
	if err != nil {
		a.logger.Warn("Failed to remove AlertPlayer, accountID:"+strconv.Itoa(accountID), err)
	}

	return apperr.ToFrontendError(err)
}

func (a *App) SearchPlayer(prefix string) (domain.WGAccountList, error) {
	accountList, err := a.configService.SearchPlayer(prefix)
	if err != nil {
		a.logger.Warn("Failed to search player, prefix:"+prefix, err)
	}

	return accountList, apperr.ToFrontendError(err)
}

func (a *App) AlertPatterns() []string {
	return domain.AlertPatterns
}

func (a *App) LogErrorForFrontend(err string) {
	a.logger.Warn("Error has occurred in frontend", apperr.New(apperr.FrontendError, errors.New(err)))
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
