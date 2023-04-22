package repo

import (
	"changeme/backend/vo"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
}

func (c *Config) User() (vo.UserConfig, error) {
    os.Mkdir("config", 0755)

    // note: set default value
    config := vo.UserConfig{
        FontSize: "medium",
    }
	file, err := os.ReadFile(filepath.Join("config", "user.json"))
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
    return config, err
}

func (c *Config) UpdateUser(config vo.UserConfig) error {
    os.Mkdir("config", 0755)

    file, err := os.Create(filepath.Join("config", "user.json"))
	if err != nil {
		return err
	}
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(config)
    return err
}

func (c *Config) App() (vo.AppConfig, error) {
    var config vo.AppConfig
    os.Mkdir("config", 0755)

	file, err := os.ReadFile(filepath.Join("config", "app.json"))
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
    return config, err
}

func (c *Config) UpdateApp(config vo.AppConfig) error {
    os.Mkdir("config", 0755)

    file, err := os.Create(filepath.Join("config", "app.json"))
	if err != nil {
		return err
	}
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(config)
    return err
}
