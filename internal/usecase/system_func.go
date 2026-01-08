package usecase

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventsEmitFunc func(ctx context.Context, eventName string, optionalData ...any)

func NewEventsEmitFunc(i do.Injector) (EventsEmitFunc, error) {
	return func(ctx context.Context, eventName string, optionalData ...any) {
		runtime.EventsEmit(ctx, eventName, optionalData...)
	}, nil
}

type OpenDirectoryDialogFunc func(ctx context.Context) (string, error)

func NewOpenDirectoryDialogFunc(i do.Injector) (OpenDirectoryDialogFunc, error) {
	return func(ctx context.Context) (string, error) {
		return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	}, nil
}
