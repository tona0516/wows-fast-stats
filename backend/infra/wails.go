package infra

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Wails struct{}

func NewWails(i do.Injector) (*Wails, error) {
	return &Wails{}, nil
}

func (w *Wails) EmitEvent(ctx context.Context, eventName string, optionalData ...any) {
	runtime.EventsEmit(ctx, eventName, optionalData...)
}

func (w *Wails) OpenDirectoryDialog(ctx context.Context) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
}
