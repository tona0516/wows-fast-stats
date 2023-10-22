package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	WindowGetSizeFunc       func(ctx context.Context) (int, int)
	WindowSetSizeFunc       func(ctx context.Context, width int, height int)
	WindowGetPosition       func(ctx context.Context) (x int, y int)
	WindowSetPosition       func(ctx context.Context, x int, y int)
	SaveFileDialogFunc      func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
	OpenDirectoryDialogFunc func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
	OpenWithDefaultAppFunc  func(input string) error
)
