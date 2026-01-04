package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

//go:generate mockgen -source=$GOFILE -destination ../mock/$GOFILE -package mock
type Logger interface {
	SetOwnIGN(ownIGN string)
	Debug(message string, contexts map[string]string)
	Info(message string, contexts map[string]string)
	Error(err error, contexts map[string]string)
}

type logger struct {
	zlog   zerolog.Logger
	ownIGN string
}

func NewLogger(
	appName string,
	semver string,
	userDataDir string,
	logLevel zerolog.Level,
	alertDiscord DiscordApiClient,
	infoDiscord DiscordApiClient,
) Logger {
	zerolog.TimeFieldFormat = time.DateTime
	zerolog.SetGlobalLevel(logLevel)

	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	remoteWriter := remoteWriter{
		alertDiscord: alertDiscord,
		infoDiscord:  infoDiscord,
	}
	logFile, _ := os.OpenFile(
		filepath.Join(userDataDir, appName+".log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModePerm,
	)
	multiLevelWriter := zerolog.MultiLevelWriter(consoleWriter, &remoteWriter, logFile)

	zlog := zerolog.New(multiLevelWriter).
		With().
		Timestamp().
		Str("semver", semver).
		Logger()

	return &logger{zlog: zlog}
}

func (l *logger) SetOwnIGN(ownIGN string) {
	l.ownIGN = ownIGN
}

func (l *logger) Debug(message string, contexts map[string]string) {
	e := l.zlog.Debug().
		Str("ign", l.ownIGN).
		Str("message", message)
	l.addContext(e, contexts)
	e.Send()
}

func (l *logger) Info(message string, contexts map[string]string) {
	e := l.zlog.Info().
		Str("ign", l.ownIGN).
		Str("message", message)
	l.addContext(e, contexts)
	e.Send()
}

func (l *logger) Error(err error, contexts map[string]string) {
	e := l.zlog.Error().
		Str("ign", l.ownIGN).
		Str("error", fmt.Sprintf("%+v", err))
	l.addContext(e, contexts)
	e.Send()
}

func (l *logger) addContext(e *zerolog.Event, contexts map[string]string) {
	if len(contexts) == 0 {
		return
	}

	for key, value := range contexts {
		e = e.Str(key, value)
	}
}
