package infra

import (
	"context"
	"fmt"
	"os"
	"time"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger struct {
	zlog    zerolog.Logger
	env     model.Env
	discord repository.DiscordInterface
}

func NewLogger(
	env model.Env,
	discord repository.DiscordInterface,
) *Logger {
	return &Logger{
		env:     env,
		discord: discord,
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
	reportWriter := zerolog.ConsoleWriter{
		TimeFormat: time.DateTime,
		Out:        &reportWriter{discord: l.discord},
		NoColor:    true,
	}
	logFile, _ := os.OpenFile(
		l.env.AppName+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o664,
	)

	multi := zerolog.MultiLevelWriter(consoleWriter, frontendWriter, reportWriter, logFile)

	l.zlog = zerolog.New(multi).
		With().
		Timestamp().
		Stack().
		Str("semver", l.env.Semver).
		Logger()
}

func (l *Logger) Debug(message string, contexts map[string]string) {
	e := l.zlog.Debug().Str("message", message)
	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Info(message string, contexts map[string]string) {
	e := l.zlog.Info().Str("message", message)
	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Warn(err error, contexts map[string]string) {
	e := l.zlog.Warn().Str("error", fmt.Sprintf("%+v", err))
	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Error(err error, contexts map[string]string) {
	e := l.zlog.Error().Str("error", fmt.Sprintf("%+v", err))
	addContext(e, contexts)
	e.Send()
}

func addContext(e *zerolog.Event, contexts map[string]string) {
	if len(contexts) == 0 {
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

type reportWriter struct {
	discord repository.DiscordInterface
}

func (w *reportWriter) Write(p []byte) (int, error) {
	_ = w.discord.Comment(string(p))
	return len(p), nil
}
