package infra

import (
	"context"
	"wfs/backend/data"

	"github.com/morikuni/failure"
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
	selected, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", failure.Translate(err, data.ErrWailsOpenDirectoryDialog)
	}

	return selected, nil
}
