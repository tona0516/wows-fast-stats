package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"wfs/backend/domain/model"
	"wfs/backend/domain/repository"

	"github.com/rs/zerolog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger struct {
	zlog         zerolog.Logger
	env          model.Env
	alertDiscord repository.DiscordInterface
	infoDiscord  repository.DiscordInterface
}

func NewLogger(
	env model.Env,
	alertDiscord repository.DiscordInterface,
	infoDiscord repository.DiscordInterface,
) *Logger {
	return &Logger{
		env:          env,
		alertDiscord: alertDiscord,
		infoDiscord:  infoDiscord,
	}
}

func (l *Logger) Init(appCtx context.Context) {
	zerolog.TimeFieldFormat = time.DateTime

	if l.env.IsDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	frontendWriter := zerolog.ConsoleWriter{
		Out:     &frontendWriter{appCtx: appCtx},
		NoColor: true,
	}
	reportWriter := reportWriter{
		alertDiscord: l.alertDiscord,
		infoDiscord:  l.infoDiscord,
	}
	logFile, _ := os.OpenFile(
		l.env.AppName+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o664,
	)

	multi := zerolog.MultiLevelWriter(consoleWriter, frontendWriter, &reportWriter, logFile)

	l.zlog = zerolog.New(multi).
		With().
		Timestamp().
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
	zerolog.FilteredLevelWriter
	alertDiscord repository.DiscordInterface
	infoDiscord  repository.DiscordInterface
}

func (w *reportWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	formatted := fmt.Sprintf("```%s```", pretty(string(p)))

	var discord repository.DiscordInterface
	if level >= zerolog.WarnLevel {
		discord = w.alertDiscord
	} else {
		discord = w.infoDiscord
	}
	_ = discord.Comment(formatted)

	return len(p), nil
}

func pretty(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}
