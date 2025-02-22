package infra

import (
	"fmt"
	"os"
	"time"
	"wfs/backend/data"
	"wfs/backend/infra/webapi"

	"github.com/rs/zerolog"
)

type Logger struct {
	zlog    zerolog.Logger
	storage Storage
	ign     string
}

func NewLogger(
	env data.Env,
	alertDiscord webapi.Discord,
	infoDiscord webapi.Discord,
	storage Storage,
) *Logger {
	zerolog.TimeFieldFormat = time.DateTime

	if env.IsDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	reportWriter := reportWriter{
		alertDiscord: alertDiscord,
		infoDiscord:  infoDiscord,
	}

	var logFile *os.File
	if len(env.AppName) > 0 {
		logFile, _ = os.OpenFile(
			env.AppName+".log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0o664,
		)
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, &reportWriter, logFile)

	zlog := zerolog.New(multi).
		With().
		Timestamp().
		Str("semver", env.Semver).
		Logger()

	ign, _ := storage.OwnIGN()

	return &Logger{
		zlog:    zlog,
		storage: storage,
		ign:     ign,
	}
}

func (l *Logger) SetOwnIGN(ign string) {
	_ = l.storage.WriteOwnIGN(ign)
	l.ign = ign
}

func (l *Logger) Debug(message string, contexts map[string]string) {
	e := l.zlog.Debug().
		Str("ign", l.ign).
		Str("message", message)

	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Info(message string, contexts map[string]string) {
	e := l.zlog.Info().
		Str("ign", l.ign).
		Str("message", message)

	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Warn(err error, contexts map[string]string) {
	e := l.zlog.Warn().
		Str("ign", l.ign).
		Str("error", fmt.Sprintf("%+v", err))

	addContext(e, contexts)
	e.Send()
}

func (l *Logger) Error(err error, contexts map[string]string) {
	e := l.zlog.Error().
		Str("ign", l.ign).
		Str("error", fmt.Sprintf("%+v", err))

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

type reportWriter struct {
	zerolog.FilteredLevelWriter
	alertDiscord webapi.Discord
	infoDiscord  webapi.Discord
}

func (w *reportWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	formatted := fmt.Sprintf("```%s```", prettyJSON(string(p)))

	var discord webapi.Discord
	if level >= zerolog.WarnLevel {
		discord = w.alertDiscord
	} else {
		discord = w.infoDiscord
	}

	err := discord.Comment(formatted)
	if err != nil {
		fmt.Printf("Failed to send to discord: %s", err.Error())
	}

	return len(p), nil
}
