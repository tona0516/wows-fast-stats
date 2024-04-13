package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	eventEmitFunc           func(ctx context.Context, eventName string, optionalData ...interface{})
	saveFileDialogFunc      func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	openDirectoryDialogFunc func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	openWithDefaultAppFunc  func(input string) error
)
