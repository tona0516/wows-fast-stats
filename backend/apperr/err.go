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
		return err
	}

	//nolint:exhaustive
	switch appError.Name {
	case WargamingAPITemporaryUnavaillalble:
		return ErrWargamingAPITemporaryUnavaillalble
	case ValidateInvalidInstallPath:
		return ErrInvalidInstallPath
	case ValidateInvalidAppID:
		return ErrInvalidAppID
	case ValidateInvalidFontSize:
		return ErrInvalidFontSize
	case OpenDirectory:
		return ErrOpenDirectory
	case UserCanceled:
		return nil
	default:
		return err
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
	ValidateInvalidInstallPath
	ValidateInvalidAppID
	ValidateInvalidFontSize
	ShowDialog
	UserCanceled
	OpenDirectory
	FrontendError
)

var (
	ErrWargamingAPITemporaryUnavaillalble = errors.New("WG APIが一時的に利用できません。リロードしてください。")
	ErrInvalidInstallPath                 = errors.New("選択したフォルダに「WorldOfWarships.exe」が存在しません。")
	ErrInvalidAppID                       = errors.New("WG APIと通信できません。アプリケーションIDが間違っている可能性があります。")
	ErrInvalidFontSize                    = errors.New("不正な文字サイズです。")
	ErrOpenDirectory                      = errors.New("フォルダが開けません。存在しない可能性があります。")
	ErrUnexpected                         = errors.New("予期しないエラーが発生しました。再起動してください。")
)
