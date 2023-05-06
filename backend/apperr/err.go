package apperr

import (
	"github.com/morikuni/failure"
)

const (
	WGAccountInfo       failure.StringCode = "WGAccountInfo"
	WGAccountList       failure.StringCode = "WGAccountList"
	WGClansAccountInfo  failure.StringCode = "WGClansAccountInfo"
	WGClansInfo         failure.StringCode = "WGClansInfo"
	WGEncyclopediaShips failure.StringCode = "WGEncyclopediaShips"
	WGShipsStats        failure.StringCode = "WGShipsStats"
	WGEncyclopediaInfo  failure.StringCode = "WGEncyclopediaInfo"
	WGBattleArenas      failure.StringCode = "WGBattleArenas"
	WGBattleTypes       failure.StringCode = "WGBattleTypes"

	NSExpectedStatsReq   failure.StringCode = "NSExpectedStatsReq"
	NSExpectedStatsParse failure.StringCode = "NSExpectedStatsParse"

	CacheSerialize   failure.StringCode = "CacheSerialize"
	CacheDeserialize failure.StringCode = "CacheDeserialize"

	CfgRead   failure.StringCode = "CfgRead"
	CfgUpdate failure.StringCode = "CfgUpdate"

	ScreenshotSave failure.StringCode = "ScreenshotSave"

	TempArenaInfoGet  failure.StringCode = "TempArenaInfoGet"
	TempArenaInfoSave failure.StringCode = "TempArenaInfoSave"

	UnregisteredWarship failure.StringCode = "UnregisteredWarship"

	CfgSvInvalidInstallPath failure.StringCode = "CfgSvInvalidInstallPath"
	CfgSvInvalidAppID       failure.StringCode = "CfgSvInvalidAppID"
	CfgSvInvalidFontSize    failure.StringCode = "CfgSvInvalidFontSize"

	PrepareSvDeleteCache failure.StringCode = "PrepareSvDeleteCache"

	ReplayWatchSvNewWatcher  failure.StringCode = "ReplayWatchSvNewWatcher"
	ReplayWatchSvWatcherAdd  failure.StringCode = "ReplayWatchSvWatcherAdd"
	ReplayWatchSvWatcherChan failure.StringCode = "ReplayWatchSvWatcherChan"

	ScreenshotSvSaveDialog failure.StringCode = "ScreenshotSvSaveDialog"

	AppCwd     failure.StringCode = "AppCwd"
	AppOpenDir failure.StringCode = "AppOpenDir"
)
