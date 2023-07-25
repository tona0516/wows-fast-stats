package logger

import (
	"context"
	"os"
	"time"
	"wfs/backend/application/vo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	wailsLogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//nolint:gochecknoglobals
var instance = &Logger{}

type Logger struct {
	zlog zerolog.Logger
}

func Init(appCtx context.Context, env vo.Env) {
	// wails logger setting
	runtime.LogSetLogLevel(appCtx, wailsLogger.ERROR)

	// zerolog setting
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if env.IsDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{
		TimeFormat: time.DateTime,
		Out:        os.Stdout,
	}
	frontendWriter := zerolog.ConsoleWriter{
		TimeFormat: time.DateTime,
		Out:        &FrontendWriter{appCtx: appCtx},
		NoColor:    true,
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, frontendWriter)

	instance.zlog = zerolog.New(multi).
		With().
		Timestamp().
		Stack().
		Str("version", env.Semver).
		Logger()
}

func Zerolog() *zerolog.Logger {
	return &instance.zlog
}

type FrontendWriter struct {
	appCtx context.Context
}

func (w *FrontendWriter) Write(p []byte) (int, error) {
	runtime.EventsEmit(w.appCtx, "LOG", string(p))
	return len(p), nil
}
