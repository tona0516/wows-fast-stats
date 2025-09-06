package data

import "strconv"

type FileStorePath string

func (p FileStorePath) ToAlertPlayerPath(accountID int) FileStorePath {
	return AlertPlayerKeyPrefix + FileStorePath(strconv.Itoa(accountID))
}

func (p FileStorePath) ToString() string {
	return string(p)
}

const (
	InstallPathKey       FileStorePath = "install_path"
	SendReportKey        FileStorePath = "send_report"
	ExpectedStatsKey     FileStorePath = "expected_stats"
	OwnIGNKey            FileStorePath = "own_ign"
	AlertPlayerKeyPrefix FileStorePath = "alert_player_"
)
