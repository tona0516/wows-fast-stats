package infra

import "changeme/backend/vo"

type ConfigInterface interface {
	User() (vo.UserConfig, error)
	UpdateUser(config vo.UserConfig) error
	App() (vo.AppConfig, error)
	UpdateApp(config vo.AppConfig) error
}
