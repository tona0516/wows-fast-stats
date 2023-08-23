package apperr

import (
	"errors"

	"github.com/morikuni/failure"
)

const (
	FileNotExist                failure.StringCode = "ファイルが存在しません。"
	ExpectedStatsParseError     failure.StringCode = "ExpectedStatsパースエラー"
	HTTPRequestError            failure.StringCode = "HTTPリクエストエラー"
	WGAPITemporaryUnavaillalble failure.StringCode = "WG APIが一時的に利用できません。リロードしてください。"
	WGAPIError                  failure.StringCode = "WG APIが利用できません。再起動するか、設定を見直してください。"
	NumbersAPIError             failure.StringCode = "NumbersAPIエラー"
	DiscordAPIError             failure.StringCode = "DiscordAPIエラー"
	GithubAPIError              failure.StringCode = "GithubAPIエラー"
	InvalidInstallPath          failure.StringCode = "選択したフォルダに「WorldOfWarships.exe」が存在しません。"
	InvalidAppID                failure.StringCode = "WG APIと通信できません。アプリケーションIDが間違っている可能性があります。"
	UserCanceled                failure.StringCode = "ユーザキャンセル"
	OpenDirectoryError          failure.StringCode = "フォルダが開けません。存在しない可能性があります。"
	FrontendError               failure.StringCode = "フロントエンドエラー"
	WailsError                  failure.StringCode = "Wailsエラー"
	UnexpectedError             failure.StringCode = "予期しないエラーが発生しました。"

	FailSafeProccess failure.StringCode = "フェイルセーフ用の処理が実行されました。"
)

func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	code, ok := failure.CodeOf(err)
	if !ok {
		return err
	}

	//nolint:goerr113
	return errors.New(code.ErrorCode())
}
