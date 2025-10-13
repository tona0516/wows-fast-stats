package service

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/data"
	"wfs/backend/domain"
	"wfs/backend/repository"

	"github.com/morikuni/failure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GameExeName = "WorldOfWarships.exe"

type Setting struct {
	localFile           repository.LocalFileInterface
	userConfig          repository.UserConfigInterface
	wargaming           repository.WargamingInterface
	logger              repository.LoggerInterface
	OpenDirectoryDialog openDirectoryDialogFunc
	OpenWithDefaultApp  openWithDefaultAppFunc
}

func NewSetting(
	localFile repository.LocalFileInterface,
	userConfig repository.UserConfigInterface,
	wargaming repository.WargamingInterface,
	logger repository.LoggerInterface,
) *Setting {
	return &Setting{
		localFile:           localFile,
		userConfig:          userConfig,
		wargaming:           wargaming,
		logger:              logger,
		OpenDirectoryDialog: runtime.OpenDirectoryDialog,
		OpenWithDefaultApp: func(input string) error {
			return exec.Command("explorer", input).Start()
		},
	}
}

func (s *Setting) RequiredSetting() (data.RequiredSetting, error) {
	config, err := s.loadUserConfig()
	if err != nil {
		return data.RequiredSetting{}, err
	}

	return data.RequiredSetting{
		Version:     config.Version,
		InstallPath: config.InstallPath,
	}, nil
}

func (s *Setting) UpdateRequiredSetting(setting data.RequiredSetting) error {
	config, err := s.loadUserConfig()
	if err != nil {
		return err
	}

	config.Version = setting.Version
	config.InstallPath = setting.InstallPath

	return s.saveUserConfig(config)
}

func (s *Setting) OptionalSetting() (data.OptionalSetting, error) {
	config, err := s.loadUserConfig()
	if err != nil {
		return data.OptionalSetting{}, err
	}

	return data.OptionalSetting{
		Version:      config.Version,
		ZoomRate:     config.ZoomRate,
		StatsExtra:   config.StatsExtra,
		IsSendReport: config.IsSendReport,
	}, nil
}

func (s *Setting) UpdateOptionalSetting(setting data.OptionalSetting) error {
	config, err := s.loadUserConfig()
	if err != nil {
		return err
	}

	config.Version = setting.Version
	config.ZoomRate = setting.ZoomRate
	config.StatsExtra = setting.StatsExtra
	config.IsSendReport = setting.IsSendReport

	return s.saveUserConfig(config)
}

func (s *Setting) BasicColumnSetting() (data.BasicColumnSetting, error) {
	config, err := s.loadUserConfig()
	if err != nil {
		return data.BasicColumnSetting{}, err
	}

	return data.BasicColumnSetting{
		Version: config.Version,
		Player: data.PlayerColumnSetting{
			EnableNationFlag: config.Column.Player.EnableNationFlag,
			ColorPattern:     config.Column.Player.ColorPattern,
		},
		Ship: data.ShipColumnSetting{
			EnableNationFlag: config.Column.Ship.EnableNationFlag,
			IsColored:        config.Column.Ship.IsColored,
		},
	}, nil
}

func (s *Setting) UpdateBasicColumnSetting(setting data.BasicColumnSetting) error {
	config, err := s.loadUserConfig()
	if err != nil {
		return err
	}

	config.Version = setting.Version
	config.Column.Player = domain.PlayerColumnConfig{
		EnableNationFlag: setting.Player.EnableNationFlag,
		ColorPattern:     setting.Player.ColorPattern,
	}
	config.Column.Ship = domain.ShipColumnConfig{
		EnableNationFlag: setting.Ship.EnableNationFlag,
		IsColored:        setting.Ship.IsColored,
	}

	return s.saveUserConfig(config)
}

func (s *Setting) StatsColumnSettings() (data.StatsColumnSettings, error) {
	config, err := s.loadUserConfig()
	if err != nil {
		return data.StatsColumnSettings{}, err
	}

	stats := statsColumnToData(config.Column.Stats)
	stats.Version = config.Version

	return stats, nil
}

func (s *Setting) UpdateStatsColumnSettings(setting data.StatsColumnSettings) error {
	config, err := s.loadUserConfig()
	if err != nil {
		return err
	}

	config.Version = setting.Version
	config.Column.Stats = dataToStatsColumn(setting)

	return s.saveUserConfig(config)
}

func (s *Setting) ValidateInstallPath(path string) error {
	if path == "" {
		return failure.New(apperr.EmptyInstallPath)
	}

	if _, err := os.Stat(filepath.Join(path, GameExeName)); err != nil {
		return failure.New(apperr.InvalidInstallPath)
	}

	return nil
}

func (s *Setting) SearchPlayer(prefix string) (data.WGAccountList, error) {
	return s.wargaming.AccountListForSearch(prefix)
}

func (s *Setting) SelectDirectory(appCtx context.Context) (string, error) {
	selected, err := s.OpenDirectoryDialog(appCtx, runtime.OpenDialogOptions{})
	if err != nil {
		return selected, failure.New(apperr.WailsError, failure.Messagef("%s", err.Error()))
	}

	return selected, nil
}

func (s *Setting) OpenDirectory(path string) error {
	err := s.OpenWithDefaultApp(path)
	if err != nil {
		return failure.New(apperr.OpenDirectoryError, failure.Context{"path": path}, failure.Messagef("%s", err.Error()))
	}

	return nil
}

func (s *Setting) loadUserConfig() (domain.UserConfig, error) {
	cfg, err := s.userConfig.Load()
	if err != nil {
		if isNotFound(err) {
			return defaultUserConfig(), nil
		}

		return domain.UserConfig{}, failure.Wrap(err)
	}

	if cfg == nil {
		return defaultUserConfig(), nil
	}

	result := *cfg

	if result.Version == 0 {
		result.Version = 1
	}
	if result.ZoomRate == 0 {
		result.ZoomRate = 100
	}
	if result.StatsExtra == "" {
		result.StatsExtra = "pvp_all"
	}
	if result.Column.Player.ColorPattern == "" {
		result.Column.Player.ColorPattern = "none"
	}

	if result.Column.Stats == (domain.StatsColumnConfig{}) {
		result.Column.Stats = defaultUserConfig().Column.Stats
	}

	return result, nil
}

func (s *Setting) saveUserConfig(config domain.UserConfig) error {
	normalized := normalizeUserConfig(config)
	if err := s.userConfig.Save(normalized); err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func statsColumnToData(stats domain.StatsColumnConfig) data.StatsColumnSettings {
	return data.StatsColumnSettings{
		Battles:           toDataDetail(stats.Battles),
		Damage:            toDataDetail(stats.Damage),
		MaxDamage:         toDataDetail(stats.MaxDamage),
		WinRate:           toDataDetail(stats.WinRate),
		SurvivedRate:      toDataDetail(stats.SurvivedRate),
		KdRate:            toDataDetail(stats.KdRate),
		Kill:              toDataDetail(stats.Kill),
		Exp:               toDataDetail(stats.Exp),
		PR:                toDataDetail(stats.PR),
		HitRate:           toDataDetail(stats.HitRate),
		PlanesKilled:      toDataDetail(stats.PlanesKilled),
		PlatoonRate:       toDataDetail(stats.PlatoonRate),
		EfficiencyBadge:   toDataDetail(stats.EfficiencyBadge),
		ThreatLevel:       toDataDetail(stats.ThreatLevel),
		AvgTier:           toDataDetail(stats.AvgTier),
		UsingShipTypeRate: toDataDetail(stats.UsingShipTypeRate),
		UsingTierRate:     toDataDetail(stats.UsingTierRate),
	}
}

func dataToStatsColumn(settings data.StatsColumnSettings) domain.StatsColumnConfig {
	return domain.StatsColumnConfig{
		Battles:           toDomainDetail(settings.Battles),
		Damage:            toDomainDetail(settings.Damage),
		MaxDamage:         toDomainDetail(settings.MaxDamage),
		WinRate:           toDomainDetail(settings.WinRate),
		SurvivedRate:      toDomainDetail(settings.SurvivedRate),
		KdRate:            toDomainDetail(settings.KdRate),
		Kill:              toDomainDetail(settings.Kill),
		Exp:               toDomainDetail(settings.Exp),
		PR:                toDomainDetail(settings.PR),
		HitRate:           toDomainDetail(settings.HitRate),
		PlanesKilled:      toDomainDetail(settings.PlanesKilled),
		PlatoonRate:       toDomainDetail(settings.PlatoonRate),
		EfficiencyBadge:   toDomainDetail(settings.EfficiencyBadge),
		ThreatLevel:       toDomainDetail(settings.ThreatLevel),
		AvgTier:           toDomainDetail(settings.AvgTier),
		UsingShipTypeRate: toDomainDetail(settings.UsingShipTypeRate),
		UsingTierRate:     toDomainDetail(settings.UsingTierRate),
	}
}

func toDataDetail(detail domain.DetailStatsColumnConfig) data.StatsColumnSetting {
	return data.StatsColumnSetting{
		IsShowShip:    detail.IsShowShip,
		IsShowOverall: detail.IsShowOverall,
		Digit:         detail.Digit,
	}
}

func toDomainDetail(detail data.StatsColumnSetting) domain.DetailStatsColumnConfig {
	return domain.DetailStatsColumnConfig{
		IsShowShip:    detail.IsShowShip,
		IsShowOverall: detail.IsShowOverall,
		Digit:         detail.Digit,
	}
}

func normalizeUserConfig(config domain.UserConfig) domain.UserConfig {
	def := defaultUserConfig()

	if config.Version == 0 {
		config.Version = def.Version
	}

	if config.ZoomRate == 0 {
		config.ZoomRate = def.ZoomRate
	}

	if config.StatsExtra == "" {
		config.StatsExtra = def.StatsExtra
	}

	if config.Column.Player.ColorPattern == "" {
		config.Column.Player.ColorPattern = def.Column.Player.ColorPattern
	}

	if config.Column.Stats == (domain.StatsColumnConfig{}) {
		config.Column.Stats = def.Column.Stats
	}

	return config
}

func defaultUserConfig() domain.UserConfig {
	return domain.UserConfig{
		Version:      1,
		InstallPath:  "",
		ZoomRate:     100,
		StatsExtra:   "pvp_all",
		IsSendReport: true,
		Column: domain.ColumnConfig{
			Player: domain.PlayerColumnConfig{
				EnableNationFlag: false,
				ColorPattern:     "none",
			},
			Ship: domain.ShipColumnConfig{
				EnableNationFlag: true,
				IsColored:        false,
			},
			Stats: domain.StatsColumnConfig{
				Battles: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				Damage: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				MaxDamage: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				WinRate: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         1,
				},
				SurvivedRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				KdRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         2,
				},
				Kill: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         2,
				},
				Exp: domain.DetailStatsColumnConfig{
					IsShowShip:    true,
					IsShowOverall: true,
					Digit:         0,
				},
				PR: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				HitRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				PlanesKilled: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				PlatoonRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				EfficiencyBadge: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         0,
				},
				ThreatLevel: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				AvgTier: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingShipTypeRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
				UsingTierRate: domain.DetailStatsColumnConfig{
					IsShowShip:    false,
					IsShowOverall: false,
					Digit:         1,
				},
			},
		},
	}
}

func isNotFound(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}
