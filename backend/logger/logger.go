package logger

import (
	"context"
	"os"
	"time"
	"wfs/backend/logger/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	wailsLogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//nolint:gochecknoglobals
var instance = &Logger{}

type Logger struct {
	zlog   zerolog.Logger
	report repository.ReportInterface
}

func Init(
	appCtx context.Context,
	appVersion string,
	isDev bool,
	report repository.ReportInterface,
) {
	// wails logger setting
	runtime.LogSetLogLevel(appCtx, wailsLogger.ERROR)

	// zerolog setting
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if isDev {
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
		Str("version", appVersion).
		Logger()

	instance.report = report
}

func Debug(message string) {
	instance.zlog.Debug().Msg(message)
}

func Info(message string) {
	instance.zlog.Info().Msg(message)
}

func Warn(err error) {
	instance.zlog.Warn().Err(err).Send()

	if instance.report != nil {
		instance.report.Send(err)
	}
}

func Error(err error) {
	instance.zlog.Error().Err(err).Send()

	if instance.report != nil {
		instance.report.Send(err)
	}
}

type FrontendWriter struct {
	appCtx context.Context
}

func (w *FrontendWriter) Write(p []byte) (int, error) {
	runtime.EventsEmit(w.appCtx, "LOG", string(p))
	return len(p), nil
}
