package vo

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	EventEmit           func(ctx context.Context, eventName string, optionalData ...interface{})
	SaveFileDialog      func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	OpenDirectoryDialog func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	OpenWithDefaultApp  func(input string) error
)
