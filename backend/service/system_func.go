package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	eventEmitFunc           func(ctx context.Context, eventName string, optionalData ...interface{})
	openDirectoryDialogFunc func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	openWithDefaultAppFunc  func(input string) error
)
