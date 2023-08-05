package apperr

import (
	"github.com/pkg/errors"
)

func New(
	appErr AppError,
	raw error,
) error {
	if raw != nil {
		return errors.WithStack(errors.Wrap(raw, appErr.Error()))
	}

	return errors.WithStack(appErr)
}

type AppError error

var (
	ErrHTTPRequest                 AppError = errors.New("HTTPRequest")
	ErrWGAPITemporaryUnavaillalble AppError = errors.New("WG APIが一時的に利用できません。リロードしてください。")
	ErrWGAPI                       AppError = errors.New("WG APIが利用できません。再起動するか、設定を見直してください。")
	ErrNumbersAPIParse             AppError = errors.New("NumbersAPIParse")
	ErrDiscordAPI                  AppError = errors.New("DiscordAPIError")
	ErrInvalidInstallPath          AppError = errors.New("選択したフォルダに「WorldOfWarships.exe」が存在しません。")
	ErrInvalidAppID                AppError = errors.New("WG APIと通信できません。アプリケーションIDが間違っている可能性があります。")
	ErrReadFile                    AppError = errors.New("ReadFile")
	ErrWriteFile                   AppError = errors.New("WriteFile")
	ErrShowDialog                  AppError = errors.New("ShowDialog")
	ErrUserCanceled                AppError = errors.New("UserCanceled")
	ErrSelectDirectory             AppError = errors.New("フォルダが選択できませんでした。")
	ErrOpenDirectory               AppError = errors.New("フォルダが開けません。存在しない可能性があります。")
	ErrFrontend                    AppError = errors.New("FrontendError")
	ErrUnexpected                  AppError = errors.New("予期しないエラーが発生しました。")
)
