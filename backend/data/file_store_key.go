package data

import "strconv"

type FileStoreKey string

func (k FileStoreKey) ToAlertPlayerKey(accountID int) FileStoreKey {
	return AlertPlayerKeyPrefix + FileStoreKey(strconv.Itoa(accountID))
}

func (k FileStoreKey) ToString() string {
	return string(k)
}

const (
	InstallPathKey       FileStoreKey = "install_path"
	SendReportKey        FileStoreKey = "send_report"
	ExpectedStatsKey     FileStoreKey = "expected_stats"
	OwnIGNKey            FileStoreKey = "own_ign"
	AlertPlayerKeyPrefix FileStoreKey = "alert_player_"
)
