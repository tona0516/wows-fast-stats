package apperr

import (
	"errors"

	"github.com/morikuni/failure"
)

const (
	FileNotExist                failure.StringCode = "ファイルが存在しません。"
	ParseExpectedStatsError     failure.StringCode = "ユーザの予測成績データのパースに失敗しました。"
	HTTPRequestError            failure.StringCode = "HTTPリクエストエラー"
	WGAPITemporaryUnavaillalble failure.StringCode = "WG APIが一時的に利用できません。リロードしてください。"
	WGAPIError                  failure.StringCode = "WG APIが利用できません。再起動するか、設定を見直してください。"
	UWGAPIError                 failure.StringCode = "WG APIが利用できません。リロードしてください。"
	NumbersAPIError             failure.StringCode = "Numbers APIが利用できません。"
	DiscordAPISendLogError      failure.StringCode = "ログ送信に失敗しました。"
	GithubAPICheckUpdateError   failure.StringCode = "アプリ更新チェックに失敗しました。"
	InvalidInstallPath          failure.StringCode = "選択したフォルダに「WorldOfWarships.exe」が存在しません。"
	EmptyInstallPath            failure.StringCode = "インストールフォルダが空です。"
	OpenDirectoryError          failure.StringCode = "フォルダが開けません。存在しない可能性があります。"
	FrontendError               failure.StringCode = "フロントエンドエラー"
	WailsError                  failure.StringCode = "Wailsエラー"
	MigrationError              failure.StringCode = "データマイグレーションに失敗しました。configフォルダを削除して再起動してください。"
	ReplayDirNotFoundError      failure.StringCode = "replayフォルダが存在しません。"
	ExpectedStatsUnavaillalble  failure.StringCode = "PRを算出できませんでした。予測成績の取得に失敗しました。"
	InvalidTempArenaInfo        failure.StringCode = "不正なtempArenaInfo.jsonです。"
	UnexpectedError             failure.StringCode = "予期しないエラーが発生しました。"
)

func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	code, ok := failure.CodeOf(err)
	if !ok {
		return err
	}

	//nolint:err113
	return errors.New(code.ErrorCode())
}
