package infra

import (
	"os"
	"path/filepath"
	"wfs/backend/apperr"
	"wfs/backend/domain"

	"github.com/morikuni/failure"
)

const (
	// directory.
	ConfigDir string = "config"

	// file.
	UserConfigFile  string = "user.json"
	AlertPlayerFile string = "alert_player.json"
)

type ConfigV0 struct {
	userConfigPath  string
	alertPlayerPath string
}

func NewConfigV0() *ConfigV0 {
	return &ConfigV0{
		userConfigPath:  filepath.Join(ConfigDir, UserConfigFile),
		alertPlayerPath: filepath.Join(ConfigDir, AlertPlayerFile),
	}
}

func (c *ConfigV0) User() (domain.UserConfig, error) {
	config, err := readJSON(c.userConfigPath, domain.DefaultUserConfig)
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return domain.DefaultUserConfig, nil
	}

	return config, err
}

func (c *ConfigV0) IsExistUser() bool {
	_, err := os.Stat(c.userConfigPath)
	return err == nil
}

func (c *ConfigV0) DeleteUser() error {
	return os.RemoveAll(c.userConfigPath)
}

func (c *ConfigV0) AlertPlayers() ([]domain.AlertPlayer, error) {
	players, err := readJSON(c.alertPlayerPath, []domain.AlertPlayer{})
	if err != nil && failure.Is(err, apperr.FileNotExist) {
		return []domain.AlertPlayer{}, nil
	}

	return players, err
}

func (c *ConfigV0) IsExistAlertPlayers() bool {
	_, err := os.Stat(c.alertPlayerPath)
	return err == nil
}

func (c *ConfigV0) DeleteAlertPlayers() error {
	return os.RemoveAll(c.alertPlayerPath)
}
