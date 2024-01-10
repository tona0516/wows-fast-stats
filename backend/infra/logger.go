package infra

import (
	"context"
	"os"
	"time"
	"wfs/backend/domain"
	"wfs/backend/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger struct {
	zlog   zerolog.Logger
	env    domain.Env
	report repository.ReportInterface
}

func NewLogger(
	env domain.Env,
	report repository.ReportInterface,
) *Logger {
	return &Logger{
		env:    env,
		report: report,
	}
}

func (l *Logger) Init(appCtx context.Context) {
	//nolint:reassign
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if l.env.IsDev {
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
		Out:        &frontendWriter{appCtx: appCtx},
		NoColor:    true,
	}
	logFile, _ := os.OpenFile(
		l.env.AppName+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o664,
	)

	multi := zerolog.MultiLevelWriter(consoleWriter, frontendWriter, logFile)

	l.zlog = zerolog.New(multi).
		With().
		Timestamp().
		Stack().
		Str("semver", l.env.Semver).
		Logger()
}

func (l *Logger) Debug(message string, contexts map[string]string) {
	e := l.zlog.Debug()
	addContext(e, contexts)
	e.Msg(message)
}

func (l *Logger) Info(message string, contexts map[string]string) {
	e := l.zlog.Info()
	addContext(e, contexts)
	e.Msg(message)
}

func (l *Logger) Warn(err error, contexts map[string]string) {
	e := l.zlog.Warn().Err(err)
	addContext(e, contexts)
	e.Send()

	if l.report != nil {
		l.report.Send("warn has occurred!", err)
	}
}

func (l *Logger) Error(err error, contexts map[string]string) {
	e := l.zlog.Error().Err(err)
	addContext(e, contexts)
	e.Send()

	if l.report != nil {
		l.report.Send("error has occurred!", err)
	}
}

func addContext(e *zerolog.Event, contexts map[string]string) {
	if contexts == nil {
		return
	}

	for key, value := range contexts {
		e = e.Str(key, value)
	}
}

//nolint:containedctx
type frontendWriter struct {
	appCtx context.Context
}

func (w *frontendWriter) Write(p []byte) (int, error) {
	runtime.EventsEmit(w.appCtx, "LOG", string(p))
	return len(p), nil
}
