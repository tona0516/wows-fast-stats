package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	WindowGetSize       func(ctx context.Context) (int, int)
	WindowSetSize       func(ctx context.Context, width int, height int)
	WindowGetPosition   func(ctx context.Context) (x int, y int)
	WindowSetPosition   func(ctx context.Context, x int, y int)
	SaveFileDialog      func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	OpenDirectoryDialog func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	OpenWithDefaultApp  func(input string) error
)
