package main

type Config struct {
	App struct {
		Name   *string `yaml:"name"`
		Semver *string `yaml:"semver"`
		Width  *int    `yaml:"width"`
		Height *int    `yaml:"height"`
	} `yaml:"app"`
	Watcher struct {
		IntervalSec *int `yaml:"interval_sec"`
	} `yaml:"watcher"`
	Wargaming struct {
		URL             *string `yaml:"url"`
		MaxRetry        *int    `yaml:"max_retry"`
		TimeoutSec      *int    `yaml:"timeout_sec"`
		RetryIntervalMs *int    `yaml:"retry_interval_ms"`
		RateLimitRPS    *int    `yaml:"rate_limit_rps"`
		AppID           *string `yaml:"app_id"`
	} `yaml:"wargaming"`
	UnofficialWargaming struct {
		URL        *string `yaml:"url"`
		MaxRetry   *int    `yaml:"max_retry"`
		TimeoutSec *int    `yaml:"timeout_sec"`
	} `yaml:"unofficial_wargaming"`
	Numbers struct {
		URL        *string `yaml:"url"`
		MaxRetry   *int    `yaml:"max_retry"`
		TimeoutSec *int    `yaml:"timeout_sec"`
	} `yaml:"numbers"`
	Github struct {
		URL        *string `yaml:"url"`
		MaxRetry   *int    `yaml:"max_retry"`
		TimeoutSec *int    `yaml:"timeout_sec"`
	} `yaml:"github"`
	Discord struct {
		AlertURL   *string `yaml:"alert_url"`
		InfoURL    *string `yaml:"info_url"`
		MaxRetry   *int    `yaml:"max_retry"`
		TimeoutSec *int    `yaml:"timeout_sec"`
	} `yaml:"discord"`
	Local struct {
		StoragePath *string `yaml:"storage_path"`
	} `yaml:"local"`
	Logger struct {
		ZerologLogLevel *string `yaml:"zerolog_log_level"`
	} `yaml:"logger"`
}
