package apperr

import (
	"fmt"

	"github.com/pkg/errors"
)

type AppError struct {
	Name AppErrorName
	Raw  error
}

func New(
	name AppErrorName,
	raw error,
) error {
	return errors.WithStack(AppError{
		Name: name,
		Raw:  raw,
	})
}

func (e AppError) Error() string {
	if e.Raw != nil {
		return fmt.Sprintf("%s: %s", e.Name.String(), e.Raw.Error())
	}

	return e.Name.String()
}

func ToFrontendError(err error) error {
	var appError AppError
	if !errors.As(err, &appError) {
		return errors.Wrap(err, ErrUnexpected.Error())
	}

	//nolint:exhaustive
	switch appError.Name {
	case WargamingAPITemporaryUnavaillalble:
		return ErrWargamingAPITemporaryUnavaillalble
	case WargamingAPIError:
		return ErrWargamingAPI
	case OpenDirectory:
		return ErrOpenDirectory
	case UserCanceled:
		return nil
	default:
		return errors.Wrap(err, ErrUnexpected.Error())
	}
}

//go:generate stringer -type=AppErrorName
type AppErrorName int

const (
	HTTPRequest AppErrorName = iota
	WargamingAPITemporaryUnavaillalble
	WargamingAPIError
	NumbersAPIParse
	ReadFile
	WriteFile
	DecodeBase64
	DiscordAPIError
	ShowDialog
	UserCanceled
	OpenDirectory
	FrontendError
)

var (
	ErrWargamingAPITemporaryUnavaillalble = errors.New("WG APIが一時的に利用できません。リロードしてください。")
	ErrWargamingAPI                       = errors.New("WG APIが利用できません。再起動するか、設定を見直してください。")
	ErrInvalidInstallPath                 = errors.New("選択したフォルダに「WorldOfWarships.exe」が存在しません。")
	ErrInvalidAppID                       = errors.New("WG APIと通信できません。アプリケーションIDが間違っている可能性があります。")
	ErrOpenDirectory                      = errors.New("フォルダが開けません。存在しない可能性があります。")
	ErrUnexpected                         = errors.New("予期しないエラーが発生しました。再起動してください。")
)
