package infra

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/imroc/req/v3"
	"github.com/rs/zerolog"
)

type Logger struct {
	zlog    zerolog.Logger
	storage Storage
	ign     string
}

func NewLogger(
	alertDiscordClient *req.Client,
	infoDiscordClient *req.Client,
	storage Storage,
	appName string,
	appSemver string,
	logLevel string,
) *Logger {
	zerolog.TimeFieldFormat = time.DateTime
	level, err := zerolog.ParseLevel(logLevel)

	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	reportWriter := reportWriter{
		alertDiscordClient: alertDiscordClient,
		infoDiscordClient:  infoDiscordClient,
	}

	var logFile *os.File
	if len(appName) > 0 {
		logFile, _ = os.Open(path.Clean(appName + ".log"))
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, &reportWriter, logFile)

	zlog := zerolog.New(multi).
		With().
		Timestamp().
		Str("semver", appSemver).
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
	alertDiscordClient *req.Client
	infoDiscordClient  *req.Client
}

func (w *reportWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	formatted := fmt.Sprintf("```%s```", prettyJSON(string(p)))

	var client *req.Client
	if level >= zerolog.WarnLevel {
		client = w.alertDiscordClient
	} else {
		client = w.infoDiscordClient
	}

	_, err := client.R().
		SetBody(map[string]string{"content": formatted}).
		Post("")
	if err != nil {
		log.Println("Failed to send to discord: ", err.Error())
	}

	return len(p), nil
}
