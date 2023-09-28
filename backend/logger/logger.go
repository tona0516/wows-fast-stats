package logger

import (
	"context"
	"os"
	"time"
	"wfs/backend/application/vo"
	"wfs/backend/logger/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
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
	env vo.Env,
	report repository.ReportInterface,
) {
	// zerolog setting
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if env.IsDev {
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
		Str("semver", env.Semver).
		Logger()

	instance.report = report
}

func Debug(message string, contexts ...vo.Pair) {
	e := instance.zlog.Debug()
	addContext(e, contexts...)
	e.Msg(message)
}

func Info(message string, contexts ...vo.Pair) {
	e := instance.zlog.Info()
	addContext(e, contexts...)
	e.Msg(message)
}

func Warn(err error, contexts ...vo.Pair) {
	e := instance.zlog.Warn().Err(err)
	addContext(e, contexts...)
	e.Send()

	if instance.report != nil {
		instance.report.Send("warn has occurred!", err)
	}
}

func Error(err error, contexts ...vo.Pair) {
	e := instance.zlog.Error().Err(err)
	addContext(e, contexts...)
	e.Send()

	if instance.report != nil {
		instance.report.Send("error has occurred!", err)
	}
}

func addContext(e *zerolog.Event, contexts ...vo.Pair) {
	if len(contexts) != 0 {
		for _, c := range contexts {
			e = e.Str(c.Key, c.Value)
		}
	}
}

type FrontendWriter struct {
	appCtx context.Context
}

func (w *FrontendWriter) Write(p []byte) (int, error) {
	runtime.EventsEmit(w.appCtx, "LOG", string(p))
	return len(p), nil
}
