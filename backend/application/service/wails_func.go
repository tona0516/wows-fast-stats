package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SaveFileDialogFunc func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error)
type WindowGetSizeFunc func(ctx context.Context) (int, int)
type WindowSetSizeFunc func(ctx context.Context, width int, height int)
type OpenDirectoryDialogFunc func(ctx context.Context, dialogOptions runtime.OpenDialogOptions) (string, error)
