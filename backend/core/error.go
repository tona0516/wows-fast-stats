package core

import (
	"errors"

	"github.com/morikuni/failure"
)

const (
	ErrWGAPITemporaryUnavailable failure.StringCode = "[A101] WGサーバとの通信に失敗しました"
	ErrWGAPI                     failure.StringCode = "[A102] WGサーバとの通信に失敗しました"
	ErrClanAPI                   failure.StringCode = "[A201] WGサーバとの通信に失敗しました"
	ErrNumbersAPI                failure.StringCode = "[A301] Numbersサーバとの通信に失敗しました"
	ErrInvalidExpectedStats      failure.StringCode = "[A302] サーバ平均成績の取得に失敗しました"
	ErrGithubAPI                 failure.StringCode = "[A401] GitHubサーバとの通信に失敗しました"
	ErrDiscordAPI                failure.StringCode = "[A501] Discordサーバとの通信に失敗しました"
	ErrTempArenaInfoNotFound     failure.StringCode = "[B101] リプレイファイルが見つかりません"
	ErrTempArenaInfoSearch       failure.StringCode = "[B102] リプレイファイルの検索に失敗しました"
	ErrJSONNotFound              failure.StringCode = "[B201] ファイルが見つかりません"
	ErrJSONRead                  failure.StringCode = "[B202] ファイル読み込みに失敗しました"
	ErrJSONWrite                 failure.StringCode = "[B203] ファイル書き込みに失敗しました"
	ErrStringRead                failure.StringCode = "[B304] ファイル読み込みに失敗しました"
	ErrStringWrite               failure.StringCode = "[B305] ファイル書き込みに失敗しました"
	ErrEmptyGameClientPath       failure.StringCode = "[C101] ゲームクライアントパスが設定されていません"
	ErrNotGameClientPath         failure.StringCode = "[C102] ゲームクライアントパスが正しくありません"
	ErrSelectFolderCancelled     failure.StringCode = "[C103] フォルダの選択がキャンセルされました"
	ErrWailsOpenDirectoryDialog  failure.StringCode = "[E101] フォルダ選択ダイアログの表示に失敗しました"
	ErrUnexpected                failure.StringCode = "[Z999] 想定外のエラーが発生しました"
)

func ErrorForDisplay(err error) error {
	code, ok := failure.CodeOf(err)
	if !ok {
		//nolint:err113
		return errors.New(ErrUnexpected.ErrorCode())
	}
	//nolint:err113
	return errors.New(code.ErrorCode())
}
