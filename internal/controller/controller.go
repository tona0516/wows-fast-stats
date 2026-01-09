package controller

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"wfs/internal/adapter"
	"wfs/internal/config"
	"wfs/internal/data"
	"wfs/internal/usecase"

	"github.com/mitchellh/go-ps"
	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Controller struct {
	appConfig                 config.Config
	configStore               adapter.ConfigStore
	fetchBattleUsecase        *usecase.FetchBattle
	pollMatchUsecase          *usecase.PollMatch
	installPathSettingUsecase *usecase.InstallPathSetting
	updateCheckUsecase        *usecase.UpdateCheck

	ctx                 context.Context
	pollMatchCancelFunc context.CancelFunc
}

func NewController(i do.Injector) (*Controller, error) {
	return &Controller{
		appConfig:                 do.MustInvoke[config.Config](i),
		configStore:               do.MustInvoke[adapter.ConfigStore](i),
		fetchBattleUsecase:        do.MustInvoke[*usecase.FetchBattle](i),
		pollMatchUsecase:          do.MustInvoke[*usecase.PollMatch](i),
		installPathSettingUsecase: do.MustInvoke[*usecase.InstallPathSetting](i),
		updateCheckUsecase:        do.MustInvoke[*usecase.UpdateCheck](i),
	}, nil
}

func (c *Controller) StartPollingMatch() {
	if c.pollMatchCancelFunc != nil {
		c.pollMatchCancelFunc()
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	c.pollMatchCancelFunc = cancelFunc
	channel := make(chan data.TempArenaInfo)

	go c.pollMatchUsecase.Invoke(c.ctx, cancelCtx, channel)
	for tempArenaInfo := range channel {
		c.fetchBattleUsecase.Invoke(c.ctx, tempArenaInfo)
	}
}

func (c *Controller) GetUserConfig() (data.UserConfig, error) {
	config, err := c.configStore.UserConfig()
	if err == nil {
		return config, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return data.DefaultUserConfig(), nil
	}

	return data.UserConfig{}, failure.Wrap(err)
}

func (c *Controller) SaveUserConfig(config data.UserConfig) error {
	if err := c.configStore.SetUserConfig(config); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (c *Controller) TrySaveInstallPath() (bool, error) {
	return c.installPathSettingUsecase.Invoke(c.ctx)
}

func (c *Controller) CurrentVersion() string {
	return c.appConfig.Basic.Version
}

func (c *Controller) NewVersion() *data.NewVersion {
	return c.updateCheckUsecase.Invoke()
}

func (c *Controller) ShowMessageDialog(message string) {
	_, _ = runtime.MessageDialog(c.ctx, runtime.MessageDialogOptions{
		Title:   c.appConfig.Basic.Name,
		Message: message,
	})
}

// 構造体のバインド用のメソッド.
func (c *Controller) EmptyBattle() data.Battle {
	return data.Battle{}
}

func (c *Controller) OnStartup(ctx context.Context) {
	c.ctx = ctx

	if isAlreadyRunning() {
		c.ShowMessageDialog("すでに起動しています。")
		os.Exit(1)
		return
	}
}

func isAlreadyRunning() bool {
	ownPid := os.Getpid()
	ownPidInfo, err := ps.FindProcess(ownPid)
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	processes, err := ps.Processes()
	if err != nil {
		// Note: 可用性のためfalseを返す
		return false
	}

	isRunning := false
	for _, p := range processes {
		if p.Pid() != ownPid && p.Executable() == ownPidInfo.Executable() {
			isRunning = true
			break
		}
	}

	return isRunning
}
