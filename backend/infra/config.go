package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	ConfigDirName  string = "config"
	ConfigUserName string = "user.json"
	ConfigAppName  string = "app.json"
)

type Config struct{}

func (c *Config) User() (vo.UserConfig, error) {
	// note: set default value
	config := vo.UserConfig{
		FontSize: "medium",
		Displays: vo.Displays{
			Basic: vo.Basic{
				IsInAvg:    true,
				PlayerName: true,
				ShipInfo:   true,
			},
		},
	}

	return read(ConfigUserName, config)
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
	return update(ConfigUserName, config)
}

func (c *Config) App() (vo.AppConfig, error) {
	return read(ConfigAppName, vo.AppConfig{})
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
	return update(ConfigAppName, config)
}

func read[T any](filename string, defaultValue T) (T, error) {
	errDetail := apperr.Cfg.Read

	_ = os.Mkdir(ConfigDirName, 0o755)

	f, err := os.ReadFile(filepath.Join(ConfigDirName, filename))
	if err != nil {
		return defaultValue, errors.WithStack(errDetail.WithRaw(err))
	}
	if err := json.Unmarshal(f, &defaultValue); err != nil {
		return defaultValue, errors.WithStack(errDetail.WithRaw(err))
	}

	return defaultValue, nil
}

func update[T any](filename string, target T) error {
	errDetail := apperr.Cfg.Update

	_ = os.Mkdir(ConfigDirName, 0o755)

	file, err := os.Create(filepath.Join(ConfigDirName, filename))
	if err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(target); err != nil {
		return errors.WithStack(errDetail.WithRaw(err))
	}

	return nil
}
