package data

import "strconv"

type FileStorePath string

func (p FileStorePath) ToAlertPlayerFileName(accountID int) FileStorePath {
	return FilePathAlertPlayer + FileStorePath(strconv.Itoa(accountID)) + ".json"
}

func (p FileStorePath) ToString() string {
	return string(p)
}

const (
	// user_data.
	FileNameRequiredSetting    FileStorePath = "required_setting.json"
	FileNameOptionalSetting    FileStorePath = "optional_setting.json"
	FileNameBasicColumnSetting FileStorePath = "basic_column.json"
	FileNameStatsColumnSetting FileStorePath = "stats_column.json"
	// cache.
	FileNameExpectedStats FileStorePath = "expected_stats.json"
	FileNameOwnIGN        FileStorePath = "own_ign.txt"

	FilePathAlertPlayer FileStorePath = "alert_player/"
)
