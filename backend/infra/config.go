package infra

import (
	"changeme/backend/apperr"
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
)

const (
	dirName  string = "config"
	userName string = "user.json"
	appName  string = "app.json"
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

	return read(userName, config)
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
	return update(userName, config)
}

func (c *Config) App() (vo.AppConfig, error) {
	return read(appName, vo.AppConfig{})
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
	return update(appName, config)
}

func read[T any](filename string, defaultValue T) (T, error) {
	errCode := apperr.CfgRead

	_ = os.Mkdir(dirName, 0o755)

	f, err := os.ReadFile(filepath.Join(dirName, filename))
	if err != nil {
		return defaultValue, failure.Translate(err, errCode)
	}
	if err := json.Unmarshal(f, &defaultValue); err != nil {
		return defaultValue, failure.Translate(err, errCode)
	}

	return defaultValue, nil
}

func update[T any](filename string, target T) error {
	errCode := apperr.CfgUpdate

	_ = os.Mkdir(dirName, 0o755)

	file, err := os.Create(filepath.Join(dirName, filename))
	if err != nil {
		return failure.Translate(err, errCode)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(target); err != nil {
		return failure.Translate(err, errCode)
	}

	return nil
}
